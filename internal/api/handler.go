package api

import (
	"net/http"

	"github.com/TheRTK/http-multiplexer/internal/app"
)

type Handler struct {
	requestsLimit int

	getAppOptions app.OptionCreator
}

func NewHandler(getAppOptions app.OptionCreator, requestsLimit int) *Handler {
	return &Handler{
		getAppOptions: getAppOptions,
		requestsLimit: requestsLimit,
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeResponseJSON(w, http.StatusMethodNotAllowed, ResponseWithMessage{Message: "Method not allowed"})
		return
	}

	if req.URL.Path == "/url" {
		h.PostURLHandler(app.NewContext(req.Context(), app.New(h.getAppOptions()...)), w, req)
		return
	}

	writeResponseJSON(w, http.StatusNotFound, ResponseWithMessage{Message: "Not found"})
}

type UrlRequestBody struct {
	Url []string `json:"url"`
}

func (h Handler) PostURLHandler(c *app.Context, w http.ResponseWriter, req *http.Request) {
	var requestBody UrlRequestBody

	if err := extractRequest(req, &requestBody); err != nil {
		writeResponseJSON(w, http.StatusBadRequest, ResponseWithMessage{Message: err.Error()})

		return
	}

	if len(requestBody.Url) == 0 || len(requestBody.Url) > 20 {
		writeResponseJSON(w, http.StatusBadRequest, ResponseWithMessage{Message: "Invalid count of urls! Min is 1 and max is 20"})

		return
	}

	// c.App.GetRequestService()

	writeResponseJSON(w, http.StatusOK, ResponseWithMessage{Message: ""})
}
