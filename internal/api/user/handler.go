package user

import (
	"net/http"
	"regexp"

	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/sirupsen/logrus"
)

var (
	listUserWithOptionsRe = regexp.MustCompile(`^/users[/]*$`)
	updateUserRe          = regexp.MustCompile(`^/users/[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$`)
	deleteUserRe          = regexp.MustCompile(`^/users/[{]?[0-9a-fA-F]{8}-([0-9a-fA-F]{4}-){3}[0-9a-fA-F]{12}[}]?$`)
	createUserRe          = regexp.MustCompile(`^/users[/]*$`)
)

// NewHandler creates a user handler
func NewHandler(Consumer ms.MessageQueueRead, Producer ms.MessageQueueWrite) *Handler {
	return &Handler{
		consumer: Consumer,
		Producer: Producer,
	}
}

// Handler manages all /users request
type Handler struct {
	consumer ms.MessageQueueRead
	Producer ms.MessageQueueWrite
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	switch {
	case r.Method == http.MethodGet && listUserWithOptionsRe.MatchString(r.URL.Path):
		logrus.Infof("received get request to list users")
		h.list(w, r)
		return
	case r.Method == http.MethodPut && updateUserRe.MatchString(r.URL.Path):
		logrus.Infof("received put request to update user")
		h.update(w, r)
		return
	case r.Method == http.MethodDelete && deleteUserRe.MatchString(r.URL.Path):
		logrus.Infof("received delete user request")
		h.delete(w, r)
		return
	case r.Method == http.MethodPost && createUserRe.MatchString(r.URL.Path):
		logrus.Infof("received create user request")
		h.create(w, r)
		return
	default:
		logrus.Infof("received request that cannot be handled %q %q", r.Method, r.URL.Path)
		notFound(w, r)
		return
	}
}

func badRequest(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write([]byte("bad request"))
}

func internalServerError(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write([]byte("internal server error"))
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	_, _ = w.Write([]byte("not found"))
}
