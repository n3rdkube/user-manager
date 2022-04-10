package messanger

// MockMessagingSystem is just a mock
type MockMessagingSystem struct {
	Received        []Message
	MessageToReturn []Message
	ErrorToReturn   []error
	index           int
}

// GetMessageFromQueue here to implement interface
func (mm *MockMessagingSystem) GetMessageFromQueue() (*Message, error) {
	var err error
	if len(mm.ErrorToReturn) > mm.index {
		err = mm.ErrorToReturn[mm.index]
	}

	var mess *Message
	if len(mm.MessageToReturn) > mm.index {
		mess = &mm.MessageToReturn[mm.index]
	}

	mm.index++
	return mess, err
}

// GetMessageFromQueueWithCID here to implement interface
func (mm *MockMessagingSystem) GetMessageFromQueueWithCID(_ string) (*Message, error) {
	var err error
	if len(mm.ErrorToReturn) > mm.index {
		err = mm.ErrorToReturn[mm.index]
	}

	var mess *Message
	if len(mm.MessageToReturn) > mm.index {
		mess = &mm.MessageToReturn[mm.index]
	}

	mm.index++
	return mess, err
}

// PostMessageToQueue  here to implement interface
func (mm *MockMessagingSystem) PostMessageToQueue(message Message) error {
	var err error
	if len(mm.ErrorToReturn) > mm.index {
		err = mm.ErrorToReturn[mm.index]
	}
	mm.Received = append(mm.Received, message)

	mm.index++
	return err
}
