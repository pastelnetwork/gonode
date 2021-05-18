package state

// Status represents status of the state.
type Status interface {
	String() string
	IsFinal() bool
}
