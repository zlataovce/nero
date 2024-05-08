package meta

import "encoding/json"

type GenericMetadata struct {
	Source     string `json:"source"`
	Artist     string `json:"artist"`
	ArtistLink string `json:"artist_link"`
}

func (gm *GenericMetadata) Type() Type {
	return TypeGeneric
}

func (gm *GenericMetadata) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type       Type   `json:"type"`
		Source     string `json:"source"`
		Artist     string `json:"artist"`
		ArtistLink string `json:"artist_link"`
	}{
		Type:       TypeGeneric,
		Source:     gm.Source,
		Artist:     gm.Artist,
		ArtistLink: gm.ArtistLink,
	})
}
