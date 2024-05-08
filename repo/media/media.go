package media

import (
	"encoding/json"
	"fmt"
	"github.com/cephxdev/nero/repo/media/meta"
	"github.com/google/uuid"
)

type Format uint

const (
	FormatUnknown Format = iota
	FormatImage
	FormatAnimatedImage
)

type Media struct {
	ID     uuid.UUID     `json:"id"`
	Format Format        `json:"format"`
	Path   string        `json:"path"`
	Meta   meta.Metadata `json:"meta"`
}

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
