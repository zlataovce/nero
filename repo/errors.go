package repo

import "fmt"

// ErrDuplicateID is an error about a duplicate media ID in a repository.
type ErrDuplicateID struct {
	// ID is the offending ID.
	ID string
	// Repo is the repository ID.
	Repo string
}

// Error returns the string representation of the error.
func (edi *ErrDuplicateID) Error() string {
	return fmt.Sprintf("duplicate media ID %s in repository %s", edi.ID, edi.Repo)
}
