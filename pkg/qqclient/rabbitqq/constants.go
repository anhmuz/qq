package rabbitqq

const RpcQueue string = "rpc_queue"
const AmqpServerURL = "amqp://guest:guest@localhost:5672/"

const (
	RedisServerAddr     = "localhost:6379"
	RedisServerPassword = ""
)

const (
	AddMessageName    string = "add"
	RemoveMessageName string = "remove"
	GetMessageName    string = "get"
	GetAllMessageName string = "get all"
)

const ClientType = "rabbitmq"
