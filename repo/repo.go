package repo

import (
	"bufio"
	"encoding/json"
	"github.com/cephxdev/nero/internal/errors"
	"github.com/cephxdev/nero/repo/media"
	"github.com/cephxdev/nero/repo/media/meta"
	mime "github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"golang.org/x/exp/maps"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
)

// Repository is a media repository.
type Repository struct {
	id, path, lockPath string
	logger             *zap.Logger

	items map[uuid.UUID]*media.Media
	mu    sync.RWMutex
}

// NewMemory creates a Repository without a backing lock file and storage directory.
func NewMemory(id string, logger *zap.Logger) *Repository {
	return &Repository{
		id:     id,
		logger: logger,
	}
}

// NewFile creates a Repository persisted to a lock file.
// If lockPath exists, its content is loaded into the repository.
func NewFile(id, path, lockPath string, logger *zap.Logger) (*Repository, error) {
	var err error

	if !filepath.IsAbs(path) {
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, errors.Wrap(err, "failed to make repository path absolute")
		}
	}

	if err = os.MkdirAll(path, 0); err != nil {
		return nil, errors.Wrap(err, "failed to make repository directories")
	}

	var items map[uuid.UUID]*media.Media
	if _, err := os.Stat(lockPath); err == nil {
		f, err := os.Open(lockPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to open index file")
		}
		defer func() {
			if err0 := f.Close(); err0 != nil {
				err = multierr.Append(err, errors.Wrap(err0, "failed to close index file"))
			}
		}()

		items = make(map[uuid.UUID]*media.Media)

		s := bufio.NewScanner(f)
		for s.Scan() {
			if s.Text() == "" {
				continue // skip empty lines
			}

			var m media.Media
			if err := json.Unmarshal(s.Bytes(), &m); err != nil {
				return nil, errors.Wrap(err, "failed to read index file item")
			}

			if _, ok := items[m.ID]; ok {
				logger.Warn(
					"duplicate item in index",
					zap.String("repo", id),
					zap.String("id", m.ID.String()),
				)
				continue
			}

			absPath := m.Path
			if !filepath.IsAbs(absPath) {
				absPath = filepath.Join(path, m.Path)
			}

			if _, err := os.Stat(absPath); errors.Is(err, os.ErrNotExist) {
				logger.Warn(
					"missing item in index",
					zap.String("repo", id),
					zap.String("id", m.ID.String()),
				)
				continue
			}

			items[m.ID] = &media.Media{
				ID:     m.ID,
				Format: m.Format,
				Path:   absPath,
				Meta:   m.Meta,
			}
		}

		if err := s.Err(); err != nil {
			return nil, errors.Wrap(err, "failed to read index file")
		}
	}

	return &Repository{
		id:       id,
		path:     path,
		lockPath: lockPath,
		logger:   logger,
		items:    items,
	}, err
}

// ID returns the ID of the repository.
func (r *Repository) ID() string {
	return r.id
}

// Path returns the storage directory path of the repository.
// Returns an empty string if it is an in-memory repository (Memory).
func (r *Repository) Path() string {
	return r.path
}

// LockPath returns the lock file path of the repository.
// Returns an empty string if it is an in-memory repository (Memory).
func (r *Repository) LockPath() string {
	return r.lockPath
}

// Memory returns whether this repository is only in memory (without a backing lock file).
func (r *Repository) Memory() bool {
	return r.lockPath == ""
}

// Get tries to find media by its ID, returns nil if nothing was found.
func (r *Repository) Get(id uuid.UUID) *media.Media {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.items == nil {
		return nil
	}
	return r.items[id]
}

// Find tries to find media by a metadata query (meta.Matchable) and a format, returns nil if nothing was found.
// Supplying media.FormatUnknown means any format should be accepted.
func (r *Repository) Find(query string, format media.Format, amount int) []*media.Media {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.items == nil {
		return nil
	}

	var res []*media.Media
	for _, m := range r.items {
		if amount <= 0 {
			break
		}

		if format != media.FormatUnknown && m.Format != format { // format mismatch
			continue
		}

		if m.Meta == nil { // no meta to match against
			continue
		}

		if ma, ok := m.Meta.(meta.Matchable); !ok || !ma.Matches(query) { // meta didn't match
			continue
		}

		amount--
		res = append(res, m)
	}

	return res
}

// Random picks N random media out of the repository.
func (r *Repository) Random(n int) []*media.Media {
	if n <= 0 {
		return nil
	}

	v := r.Items()
	rand.Shuffle(len(v), func(i, j int) {
		v[i], v[j] = v[j], v[i]
	})

	if len(v) > n {
		v = v[:n]
	}
	return v
}

// Create creates and inserts new media into the repository.
// Returns errors.ErrUnsupported for repositories without a backing storage directory.
func (r *Repository) Create(b []byte, m meta.Metadata) (*media.Media, error) {
	if r.path == "" {
		return nil, errors.ErrUnsupported
	}

	var (
		err error

		id    = uuid.New()
		type_ = mime.Detect(b)
		path  = filepath.Join(r.path, id.String()+type_.Extension())
	)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open file")
	}
	defer func() {
		if err0 := f.Close(); err0 != nil {
			err = multierr.Append(err, errors.Wrap(err0, "failed to close file"))
		}
	}()

	if _, err = f.Write(b); err != nil {
		return nil, errors.Wrap(err, "failed to write file")
	}

	m0 := &media.Media{
		ID:     id,
		Format: media.FormatUnknown,
		Path:   path,
		Meta:   m,
	}
	switch type_.String() {
	case "image/jpeg", "image/png":
		m0.Format = media.FormatImage
	case "image/vnd.mozilla.apng", "image/gif", "image/webp":
		m0.Format = media.FormatAnimatedImage
	}

	err = r.Add(m0)
	return m0, err
}

// Add inserts new media into the repository.
func (r *Repository) Add(m *media.Media) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.items == nil {
		r.items = make(map[uuid.UUID]*media.Media, 1)
	} else if _, ok := r.items[m.ID]; ok {
		return &ErrDuplicateID{
			ID:   m.ID.String(),
			Repo: r.id,
		}
	}

	r.items[m.ID] = m
	return r.save()
}

// Remove removes media from the repository by its ID.
func (r *Repository) Remove(id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.items, id)
	return r.save()
}

// Items returns all pieces of media in the repository.
func (r *Repository) Items() []*media.Media {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return maps.Values(r.items)
}

// Close cleans up after the repository.
// The repository should not be used anymore after calling Close.
func (r *Repository) Close() error {
	return nil
}

func (r *Repository) save() (err error) {
	if r.lockPath == "" {
		return nil
	}

	if _, err := os.Stat(r.lockPath); err == nil {
		if err = os.Rename(r.lockPath, r.lockPath+".old"); err != nil {
			return errors.Wrap(err, "failed to move index file")
		}
	}

	f, err := os.OpenFile(r.lockPath, os.O_WRONLY|os.O_CREATE, 0)
	if err != nil {
		return errors.Wrap(err, "failed to open index file")
	}
	defer func() {
		if err0 := f.Close(); err0 != nil {
			err = multierr.Append(err, errors.Wrap(err0, "failed to close index file"))
		}
	}()

	for _, m := range r.items {
		if err = r.write(f, m); err != nil {
			return err
		}
	}

	return err
}

func (r *Repository) write(f *os.File, m *media.Media) error {
	path, err0 := filepath.Rel(r.path, m.Path)
	if err0 != nil {
		path = m.Path
	}

	b, err := json.Marshal(&media.Media{
		ID:     m.ID,
		Format: m.Format,
		Path:   path,
		Meta:   m.Meta,
	})
	if err != nil {
		return errors.Wrap(err, "failed to serialize index item")
	}

	if _, err = f.Write(append(b, []byte("\n")...)); err != nil {
		return errors.Wrap(err, "failed to write index item")
	}

	return nil
}
