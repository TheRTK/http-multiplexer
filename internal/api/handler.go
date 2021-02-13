package api

import (
	"net/http"
	"net/url"

	"github.com/TheRTK/http-multiplexer/internal/app"
)

type Handler struct {
	sem           chan struct{}
	getAppOptions app.OptionCreator
}

func NewHandler(getAppOptions app.OptionCreator, requestsLimit int) *Handler {
	h := &Handler{
		getAppOptions: getAppOptions,
		sem:           make(chan struct{}, requestsLimit),
	}

	return h
}

func (h Handler) Shutdown() {
	close(h.sem)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeResponseJSON(w, http.StatusMethodNotAllowed, ResponseWithMessage{Message: "Method not allowed"})
		return
	}

	if req.URL.Path == "/multiplexer" {
		/*
			Не понял по условию задания ограничение в 100 запросов должно быть для всех обработчиков или для конкретного.
			По идеи можно взять что-то типо netutil.LimitListener и тогда будет ограничение для всего сервера.
		*/
		select {
		case h.sem <- struct{}{}:
			h.PostURLHandler(app.NewContext(req.Context(), app.New(h.getAppOptions()...)), w, req)
			<-h.sem
		default:
			writeResponseJSON(w, http.StatusServiceUnavailable, ResponseWithMessage{Message: "Server is not unavailable"})
		}

		return
	}

	writeResponseJSON(w, http.StatusNotFound, ResponseWithMessage{Message: "Not found"})
}

type RequestBody struct {
	URL []string `json:"url"`
}

func (h Handler) PostURLHandler(c *app.Context, w http.ResponseWriter, req *http.Request) {
	var requestBody RequestBody

	if err := extractRequest(req, &requestBody); err != nil {
		writeResponseJSON(w, http.StatusBadRequest, ResponseWithMessage{Message: err.Error()})

		return
	}

	if len(requestBody.URL) == 0 || len(requestBody.URL) > 20 {
		writeResponseJSON(w, http.StatusBadRequest, ResponseWithMessage{Message: "Invalid count of urls"})

		return
	}

	urls, err := ParseURLArray(requestBody.URL)
	if err != nil {
		writeResponseJSON(w, http.StatusBadRequest, ResponseWithMessage{Message: err.Error()})

		return
	}

	requestService := c.App.GetRequestService()

	data, err := requestService.GetDataFromUrls(c.Ctx, urls)
	if err != nil {
		writeResponseJSON(w, http.StatusInternalServerError, ResponseWithMessage{Message: err.Error()})

		return
	}

	writeResponseJSON(w, http.StatusOK, data)
}

func ParseURLArray(data []string) ([]*url.URL, error) {
	urlSlice := make([]*url.URL, 0, len(data))

	for _, v := range data {
		urlItem, err := url.Parse(v)
		if err != nil {
			return nil, err
		}

		urlSlice = append(urlSlice, urlItem)
	}

	return urlSlice, nil
}
