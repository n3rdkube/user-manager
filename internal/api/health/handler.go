package health

import (
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
)

var (
	healthRe = regexp.MustCompile(`^/health[/]*$`)
)

// NewHandler creates a health handler
func NewHandler() *Handler { return &Handler{} }

// Handler is the health handler of the manager
type Handler struct {
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet && healthRe.MatchString(r.URL.Path):
		logrus.Infof("received get request for health")
		h.health(w, r)
		return
	default:
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("not found"))
		return
	}
}
