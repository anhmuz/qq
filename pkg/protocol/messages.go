package protocol

type BaseMessage struct {
	Name string `json:"name"`
}

type BaseReplyMessage struct {
	Name string `json:"name"`
}

type AddMessage struct {
	BaseMessage
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RemoveMessage struct {
	BaseMessage
	Key string `json:"key"`
}

type GetMessage struct {
	BaseMessage
	Key string `json:"key"`
}

type GetAllMessage struct {
	BaseMessage
}

type AddReplyMessage struct {
	BaseReplyMessage
	Added bool `json:"added"`
}

type RemoveReplyMessage struct {
	BaseReplyMessage
	Removed bool `json:"removed"`
}

type GetReplyMessage struct {
	BaseReplyMessage
	Value *string `json:"value"`
}

type GetAllReplyMessage struct {
	BaseReplyMessage
	Entities []Entity `json:"entities"`
}
