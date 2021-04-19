package state

// List of task statuses.
const (
	StatusStarted Status = iota
	StatusAccepted
	StatusActivation
	StatusActivated
	StatusError
)

var statusNames = map[Status]string{
	StatusStarted:    "Registration Started",
	StatusAccepted:   "Artwork Accepted",
	StatusActivation: "Waiting Activation",
	StatusActivated:  "Activated",
	StatusError:      "Error",
}

// Status represents status of the task
type Status byte

func (status Status) String() string {
	if name, ok := statusNames[status]; ok {
		return name
	}
	return ""
}

// StatusNames returns a sorted list of status names.
func StatusNames() []string {
	list := make([]string, len(statusNames))
	for i, name := range statusNames {
		list[i] = name
	}
	return list
}
