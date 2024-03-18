package gm

// Marshaler customizes the marshalling behavior for a type
type Marshaler interface {
	MarshalBCS() ([]byte, error)
}
