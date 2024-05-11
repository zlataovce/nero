package meta

import (
	"encoding/json"
	"strings"
)

// AnimeMetadata is a piece of anime-attributed metadata.
type AnimeMetadata struct {
	// Name is the anime name.
	Name string `json:"name"`

	lowerName string // transient cache for matching
}

// Type returns the type of the metadata (TypeAnime).
func (am *AnimeMetadata) Type() Type {
	return TypeAnime
}

// Matches tries to match against a string query.
func (am *AnimeMetadata) Matches(query string) bool {
	if am.lowerName == "" {
		am.lowerName = strings.ToLower(am.Name)
	}

	return strings.Contains(am.lowerName, strings.ToLower(query))
}

// MarshalJSON writes data into a JSON representation.
func (am *AnimeMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type Type   `json:"type"`
		Name string `json:"name"`
	}{
		Type: TypeAnime,
		Name: am.Name,
	})
}
