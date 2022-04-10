package messanger

const (
	// CreateUser should be used for create requests
	CreateUser = "createUser"
	// DeleteUser should be used for delete requests
	DeleteUser = "deleteUser"
	// UpdateUser should be used for update requests
	UpdateUser = "updateUser"
	// ListUser should be used for list requests
	ListUser = "listUser"
	// ListUserAnswer should be used for list callbacks
	ListUserAnswer = "listUserAnswer"
	// Notification should be used for notifications
	Notification = "notification"
)

// MessageQueueRead contains all method needed by consumers
type MessageQueueRead interface {
	GetMessageFromQueue() (*Message, error)
	GetMessageFromQueueWithCID(cID string) (*Message, error)
}

// MessageQueueWrite contains all method needed by producers
type MessageQueueWrite interface {
	PostMessageToQueue(message Message) error
}

// Message holds all data needed by manager and processor to read and send rabbitmq messages
type Message struct {
	Data        []byte
	MessageType string
	CID         string
}
