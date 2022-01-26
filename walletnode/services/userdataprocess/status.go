package userdataprocess

//
//import "github.com/pastelnetwork/gonode/walletnode/services/common"
//
//// List of task statuses.
//const (
//	StatusTaskStarted Status = iota
//	StatusConnected
//	StatusTaskFailed
//	StatusTaskCompleted
//	ErrorNotEnoughMasterNode
//	StatusErrorNotEnoughMasterNode
//)
//
//var statusNames = map[Status]string{
//	StatusTaskStarted:              "Task Started",
//	StatusConnected:                "Connected",
//	StatusTaskFailed:               "Task Failed",
//	StatusTaskCompleted:            "Task Completed",
//	StatusErrorNotEnoughMasterNode: "Error not enough Master Node to send data",
//}
//
//// Status represents status of the task
//type Status byte
//
//func (status Status) String() string {
//	if name, ok := statusNames[status]; ok {
//		return name
//	}
//	return ""
//}
//
//// IsFinal returns true if the status is the final.
//func (status Status) IsFinal() bool {
//	return status == StatusTaskCompleted || status == StatusTaskFailed
//}
//
//// IsFailure returns true if the task failed due to an error
//func (status Status) IsFailure() bool {
//	return status == common.StatusTaskFailed || status == common.StatusErrorNotEnoughSuperNode
//}
//
//// StatusNames returns a sorted list of status names.
//func StatusNames() []string {
//	list := make([]string, len(statusNames))
//	for i, name := range statusNames {
//		list[i] = name
//	}
//	return list
//}
