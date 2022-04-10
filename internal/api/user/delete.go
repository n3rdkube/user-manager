package user

import (
	"encoding/json"
	"net/http"
	"path"

	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/n3rdkube/user-manager/internal/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	u := models.User{
		ID: path.Base(r.URL.String()),
	}

	sanitizedBytes, err := json.Marshal(u)
	if err != nil {
		logrus.Errorf("mashalling: %v", err)
		internalServerError(w, r)
		return
	}

	err = h.Producer.PostMessageToQueue(ms.Message{
		Data:        sanitizedBytes,
		MessageType: ms.DeleteUser,
	})
	if err != nil {
		logrus.Errorf("posting message: %v", err)
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}
