package dupedetection

// DDServerStats represents the DD-Server stats result
type DDServerStats struct {
	MaxConcurrent  int32
	Executing      int32
	WaitingInQueue int32
	Version        string
}
