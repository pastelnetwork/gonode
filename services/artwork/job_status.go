package artwork

// Constants of the job statuss.
const (
	JobStatusStarted JobStatus = iota
	JobStatusAccepted
	JobStatusActivation
	JobStatusActivated
	JobStatusError
)

var jobStatusNames = map[JobStatus]string{
	JobStatusStarted:    "Registration Started",
	JobStatusAccepted:   "Artwork Accepted",
	JobStatusActivation: "Waiting Activation",
	JobStatusActivated:  "Activated",
	JobStatusError:      "Error",
}

// JobStatus represents status of the job
type JobStatus byte

func (status JobStatus) String() string {
	if name, ok := jobStatusNames[status]; ok {
		return name
	}
	return ""
}
