package state

// List of job statuses.
const (
	StatusStarted Status = iota
	StatusAccepted
	StatusActivation
	StatusActivated
	StatusError
	StatusFinish
)

var statusNames = map[Status]string{
	StatusStarted:    "Registration Started",
	StatusAccepted:   "Artwork Accepted",
	StatusActivation: "Waiting Activation",
	StatusActivated:  "Activated",
	StatusError:      "Error",
	StatusFinish:     "Finish",
}

// Status represents status of the job
type Status byte

func (status Status) String() string {
	if name, ok := statusNames[status]; ok {
		return name
	}
	return ""
}
