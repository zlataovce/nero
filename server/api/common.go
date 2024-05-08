package api

// MakeOptString converts a string to its pointer if it's not a zero value.
func MakeOptString(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

// MakeString converts a string pointer to a string or a zero value if it's nil.
func MakeString(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}
