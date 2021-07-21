

package userdata

// List common code 
const (
	SUCCESS_PROCESS  Code = iota
	ERROR_NOT_ENOUGH_SUPERNODE
	ERROR_NOT_ENOUGH_SUPERNODE_RESPONSE
	ERROR_NOT_ENOUGH_SUPERNODE_CONFIRM
	ERROR_ON_CONTENT
)

var Description = map[Code]string{
	SUCCESS_PROCESS:               			"User specified data set successfully",
	ERROR_NOT_ENOUGH_SUPERNODE:         	"Not enough SuperNodes sending data to",
	ERROR_NOT_ENOUGH_SUPERNODE_RESPONSE:	"Not enough SuperNodes reply",
	ERROR_NOT_ENOUGH_SUPERNODE_CONFIRM:		"SuperNodes response mismatch"
	ERROR_ON_CONTENT:            			"There is error on field(s) in user specified data",
}