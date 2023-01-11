package rabbitqq

const RpcQueue string = "rpc_queue"
const AmqpServerURL = "amqp://guest:guest@localhost:5672/"

const (
	AddMessageName    string = "add"
	RemoveMessageName string = "remove"
	GetMessageName    string = "get"
	GetAllMessageName string = "get all"
)

const DefaultUserIdValue string = ""
