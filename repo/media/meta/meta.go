package meta

// Type is a type of metadata.
type Type uint

const (
	// TypeGeneric is a generic, artist-attributed metadata type (GenericMetadata).
	TypeGeneric Type = iota
	// TypeAnime is an anime-attributed metadata type (AnimeMetadata).
	TypeAnime
)

// Metadata is a piece of media metadata.
type Metadata interface {
	// Type returns the type of the metadata.
	Type() Type
}

// Matchable is something that can be matched.
type Matchable interface {
	// Matches tries to match against a string query.
	Matches(query string) bool
}
