package meta

import (
	"encoding/json"
	"strings"
)

type AnimeMetadata struct {
	Name string `json:"name"`
}

func (am *AnimeMetadata) Type() Type {
	return TypeAnime
}

func (am *AnimeMetadata) Matches(query string) bool {
	return strings.Contains(strings.ToLower(am.Name), strings.ToLower(query))
}

func (am *AnimeMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type Type   `json:"type"`
		Name string `json:"name"`
	}{
		Type: TypeAnime,
		Name: am.Name,
	})
}
