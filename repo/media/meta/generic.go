package meta

import "encoding/json"

// GenericMetadata is a piece of artist-attributed metadata.
type GenericMetadata struct {
	// Source is the media source, i.e. a URL.
	Source string `json:"source"`
	// Artist is the identifier of the artist, i.e. their name.
	Artist string `json:"artist"`
	// ArtistLink is a reference to the artist, i.e. a URL.
	ArtistLink string `json:"artist_link"`
}

// Type returns the type of the metadata (TypeGeneric).
func (gm *GenericMetadata) Type() Type {
	return TypeGeneric
}

// MarshalJSON writes data into a JSON representation.
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
