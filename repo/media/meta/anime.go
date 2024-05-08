package meta

import "encoding/json"

type AnimeMetadata struct {
	Name string `json:"name"`
}

func (am *AnimeMetadata) Type() Type {
	return TypeAnime
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
