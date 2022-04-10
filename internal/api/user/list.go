package user

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/n3rdkube/user-manager/internal/models"
	"github.com/sirupsen/logrus"
)

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	cID := fmt.Sprintf("%f", rand.Float64())

	sanitizedBytes, err := getListOptionsData(r)
	if err != nil {
		logrus.Errorf("getting list options: %v", err)
		badRequest(w, r)
		return
	}

	logrus.Infof("posting message %q", sanitizedBytes)
	err = h.Producer.PostMessageToQueue(ms.Message{
		Data:        sanitizedBytes,
		MessageType: ms.ListUser,
		CID:         cID,
	})
	if err != nil {
		logrus.Errorf("posting message: %v", err)
		internalServerError(w, r)
		return
	}

	logrus.Infof("waiting for message")
	m, err := h.consumer.GetMessageFromQueueWithCID(cID)
	if err != nil {
		logrus.Errorf("getting rpc message: %v", err)
		internalServerError(w, r)
		return
	}

	var list models.ListUsers
	err = json.Unmarshal(m.Data, &list)
	if err != nil {
		logrus.Errorf("From the callback we did not received a list: %v", err)
		internalServerError(w, r)
		return
	}

	sanitizedBytes, err = json.Marshal(&list)
	if err != nil {
		logrus.Errorf("From the callback we received a list we cannot marshall: %v", err)
		internalServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(sanitizedBytes)
}

func getListOptionsData(r *http.Request) ([]byte, error) {

	var pn int
	if pnString := r.URL.Query().Get("page.number"); pnString != "" {
		n, err := strconv.Atoi(pnString)
		if err != nil {
			return nil, fmt.Errorf("parsing page number: %w", err)
		}
		pn = n
	}

	var rows int
	if rowsString := r.URL.Query().Get("page.rows"); rowsString != "" {
		r, err := strconv.Atoi(rowsString)
		if err != nil {
			return nil, fmt.Errorf("parsing page number: %w", err)
		}
		rows = r
	}

	f := models.ListOptions{
		Include: models.User{
			Country:   r.URL.Query().Get("include.country"),
			FirstName: r.URL.Query().Get("include.firstname"),
			LastName:  r.URL.Query().Get("include.lastname"),
		},
		PageNumber:  pn,
		RowsPerPage: rows,
	}

	sanitizedBytes, err := json.Marshal(f)
	if err != nil {
		return nil, fmt.Errorf("marshalling: %w", err)
	}

	return sanitizedBytes, nil
}
