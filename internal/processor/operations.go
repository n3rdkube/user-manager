package processor

import (
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"

	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/n3rdkube/user-manager/internal/models"
)

func (p *Processor) addUser(data []byte) error {
	usr := models.User{}

	err := json.Unmarshal(data, &usr)
	if err != nil {
		return fmt.Errorf("unmrshalling message data to create user: %w", err)
	}

	err = p.st.AddUser(usr)
	if err != nil {
		return fmt.Errorf("adding user to db: %w", err)
	}

	logrus.Infof("sending notification, user was added!")
	err = p.notifications.PostMessageToQueue(ms.Message{
		Data:        []byte(fmt.Sprintf("user added %q", usr.ID)),
		MessageType: ms.Notification,
	})
	if err != nil {
		return fmt.Errorf("failed to notify systems: %w", err)
	}

	return nil
}

func (p *Processor) deleteUser(data []byte) error {
	usr := models.User{}

	err := json.Unmarshal(data, &usr)
	if err != nil {
		return fmt.Errorf("unmrshalling message data to delete user: %w", err)
	}

	err = p.st.DeleteUser(usr.ID)
	if err != nil {
		return fmt.Errorf("deleting user from db: %w", err)
	}

	logrus.Infof("sending notification, user was deleted!")
	err = p.notifications.PostMessageToQueue(ms.Message{
		Data:        []byte(fmt.Sprintf("user deleted %q", usr.ID)),
		MessageType: ms.Notification,
	})
	if err != nil {
		return fmt.Errorf("failed to notify systems: %w", err)
	}

	return nil
}

func (p *Processor) updateUser(data []byte) error {
	usr := models.User{}

	err := json.Unmarshal(data, &usr)
	if err != nil {
		return fmt.Errorf("unmrshalling message data to update user: %w", err)
	}

	err = p.st.UpdateUser(usr)
	if err != nil {
		return fmt.Errorf("updateing user in db: %w", err)
	}

	logrus.Infof("sending notification, user was updated!")
	err = p.notifications.PostMessageToQueue(ms.Message{
		Data:        []byte(fmt.Sprintf("user updated %q", usr.ID)),
		MessageType: ms.Notification,
	})
	if err != nil {
		return fmt.Errorf("failed to notify systems: %w", err)
	}

	return nil
}

func (p *Processor) listUser(data []byte, cid string) error {
	listOptions := models.ListOptions{}
	userList := models.ListUsers{}

	err := json.Unmarshal(data, &listOptions)
	if err != nil {
		return fmt.Errorf("unmrshalling message listOptions to list user: %w", err)
	}

	userList, err = p.st.ListUser(listOptions)
	if err != nil {
		return fmt.Errorf("listing users from db: %w", err)
	}

	callbackData, err := json.Marshal(userList)
	if err != nil {
		return fmt.Errorf("marshalling user list: %w", err)
	}

	logrus.Infof("sending data, user list was fetched!")
	answer := ms.Message{
		Data:        callbackData,
		MessageType: ms.ListUserAnswer,
		CID:         cid,
	}
	err = p.mqCallback.PostMessageToQueue(answer)
	if err != nil {
		return fmt.Errorf("publishing result of list: %w", err)
	}

	return nil
}
