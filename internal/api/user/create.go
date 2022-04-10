package user

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/n3rdkube/user-manager/internal/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	sanitizedBytes, err := getUserData(r)
	if err != nil {
		logrus.Errorf("reading user data from body: %v", err)
		badRequest(w, r)
		return
	}

	err = h.Producer.PostMessageToQueue(ms.Message{Data: sanitizedBytes, MessageType: ms.CreateUser})
	if err != nil {
		logrus.Errorf("posting message: %v", err)
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func getUserData(r *http.Request) ([]byte, error) {
	var u models.User
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("reading full body: %w", err)
	}

	err = json.Unmarshal(b, &u)
	if err != nil {
		return nil, fmt.Errorf("unmashalling full body: %w", err)
	}

	err = u.ValidateUserInput()
	if err != nil {
		return nil, fmt.Errorf("validating user data: %w", err)
	}

	sanitizedBytes, err := json.Marshal(u)
	if err != nil {
		return nil, fmt.Errorf("marshalling: %w", err)
	}

	return sanitizedBytes, nil
}
