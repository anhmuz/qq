package rabbitqq

type AddMessage struct {
	Name  string `json:"name"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RemoveMessage struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type GetMessage struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type GetAllMessage struct {
	Name string `json:"name"`
}

type AddReplyMessage struct {
	Added bool `json:"added"`
}

type RemoveReplyMessage struct {
	Removed bool `json:"removed"`
}

type GetReplyMessage struct {
	Value *string `json:"value"`
}

type GetAllReplyMessage struct {
	Entities map[string]string `json:"entities"`
}
