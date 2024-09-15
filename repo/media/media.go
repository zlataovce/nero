package media

import (
	"encoding/json"
	"fmt"
	"github.com/zlataovce/nero/repo/media/meta"
	"github.com/google/uuid"
)

// Format is a media format.
type Format uint

const (
	// FormatUnknown is an unknown media format.
	FormatUnknown Format = iota
	// FormatImage is a standard image media format, i.e. JPEG, PNG.
	FormatImage
	// FormatAnimatedImage is an animated image media format, i.e. GIF, APNG, WEBP.
	FormatAnimatedImage
)

// Media is a piece of media.
type Media struct {
	// ID is the media ID.
	ID uuid.UUID `json:"id"`
	// Format is the media format.
	Format Format `json:"format"`
	// Path is the media path.
	Path string `json:"path"`
	// Meta is the media metadata, may be nil.
	Meta meta.Metadata `json:"meta"`
}

// UnmarshalJSON reads data from a JSON representation.
func (m *Media) UnmarshalJSON(bytes []byte) error {
	var raw struct {
		ID     uuid.UUID       `json:"id"`
		Format Format          `json:"format"`
		Path   string          `json:"path"`
		Meta   json.RawMessage `json:"meta"`
	}
	if err := json.Unmarshal(bytes, &raw); err != nil {
		return err
	}

	m.ID = raw.ID
	m.Format = raw.Format
	m.Path = raw.Path

	var partialMeta struct {
		Type meta.Type `json:"type"`
	}
	if err := json.Unmarshal(raw.Meta, &partialMeta); err != nil {
		return err
	}

	switch partialMeta.Type {
	case meta.TypeGeneric:
		var meta0 meta.GenericMetadata
		if err := json.Unmarshal(raw.Meta, &meta0); err != nil {
			return err
		}

		m.Meta = &meta0
	case meta.TypeAnime:
		var meta0 meta.AnimeMetadata
		if err := json.Unmarshal(raw.Meta, &meta0); err != nil {
			return err
		}

		m.Meta = &meta0
	default:
		return fmt.Errorf("unexpected metadata type %d", partialMeta.Type)
	}

	return nil
}
