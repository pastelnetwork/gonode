package types

import (
	"time"
)

type File struct {
	FileID                       string
	UploadTimestamp              time.Time
	Path                         string
	FileIndex                    string
	BaseFileID                   string
	TaskID                       string
	RegTxid                      string
	ActivationTxid               string
	ReqBurnTxnAmount             float64
	BurnTxnID                    string
	ReqAmount                    float64
	IsConcluded                  bool
	CascadeMetadataTicketID      string
	UUIDKey                      string
	HashOfOriginalBigFile        string
	NameOfOriginalBigFileWithExt string
	SizeOfOriginalBigFile        float64
	DataTypeOfOriginalBigFile    string
	StartBlock                   int32
	DoneBlock                    int
}

type Files []*File

type RegistrationAttempt struct {
	ID           int
	FileID       string
	RegStartedAt time.Time
	ProcessorSNS string
	FinishedAt   time.Time
	IsSuccessful bool
	ErrorMessage string
}

type ActivationAttempt struct {
	ID                  int
	FileID              string
	ActivationAttemptAt time.Time
	IsSuccessful        bool
	ErrorMessage        string
}

func (fs Files) GetBase() *File {
	for _, f := range fs {
		if f.FileIndex == "0" {
			return f
		}
	}

	return nil
}
