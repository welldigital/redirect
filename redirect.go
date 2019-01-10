package redirect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Code used to redirect.
type Code int

// CodeTemporary is a temporary redirect.
var CodeTemporary Code = http.StatusFound

// CodePermanent is a permanent redirect.
var CodePermanent Code = http.StatusMovedPermanently

// NewHandler creates a new handler.
func NewHandler(defaultTo string, defaultCode Code) *Handler {
	return &Handler{
		Redirects: make(map[string]redirect),
		Default: redirect{
			URL:  defaultTo,
			Code: defaultCode,
		},
	}
}

// Handler which redirects.
type Handler struct {
	Redirects map[string]redirect
	Default   redirect
}

type redirect struct {
	URL  string
	Code Code
}

// Add a redirect to the Handler.
func (h *Handler) Add(from, to string, code Code) *Handler {
	h.Redirects[normalise(from)] = redirect{
		URL:  to,
		Code: code,
	}
	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	to, ok := h.Redirects[normalise(r.URL.Path)]
	if !ok {
		to = h.Default
	}
	log(r.URL.Path, to.URL, to.Code)
	http.Redirect(w, r, to.URL, int(to.Code))
}

func normalise(path string) string {
	return strings.ToLower(strings.TrimSuffix(path, "/"))
}

func log(from, to string, code Code) {
	b, err := json.Marshal(map[string]interface{}{
		"time": time.Now().UTC(),
		"from": from,
		"to":   to,
		"code": code,
	})
	if err != nil {
		fmt.Printf("error marshalling JSON: %v\n", err)
		return
	}
	fmt.Println(string(b))
}
