package messaging

// CommonProtoMsg override proto message
type CommonProtoMsg struct{}

func (*CommonProtoMsg) String() string { return "domain msg" }
func (*CommonProtoMsg) Reset()         {}
func (*CommonProtoMsg) ProtoMessage()  {}
