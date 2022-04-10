package user

import (
	"net/http"

	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/sirupsen/logrus"
)

func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	sanitizedBytes, err := getUserData(r)
	if err != nil {
		logrus.Errorf("reading user data from body: %v", err)
		badRequest(w, r)
		return
	}

	// we could check as well if urlID = userId in schema

	err = h.Producer.PostMessageToQueue(ms.Message{
		Data:        sanitizedBytes,
		MessageType: ms.UpdateUser,
	})
	if err != nil {
		logrus.Errorf("posting message: %v", err)
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
