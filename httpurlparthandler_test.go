package httpurlparthandler

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewHandler(t *testing.T) {
	tt := []struct {
		method string
		target string
		body io.Reader
		responseCode int
		responseBody string
	}{
		{
			method: http.MethodGet,
			target: "/",
			responseCode: http.StatusNotFound,
			responseBody: "404 page not found\n",
		},
		{
			method: http.MethodGet,
			target: "/api/",
			responseCode: http.StatusNotFound,
			responseBody: "404 page not found\n",
		},
		{
			method: http.MethodGet,
			target: "/api/tasks/",
			responseCode: http.StatusNotFound,
			responseBody: "",
		},
		{
			method: http.MethodGet,
			target: "/api/tasks/1000",
			responseCode: http.StatusOK,
			responseBody: "1000",
		},
		{
			method: http.MethodGet,
			target: "/api/notes/",
			responseCode: http.StatusNotFound,
			responseBody: "",
		},
		{
			method: http.MethodGet,
			target: "/api/notes/2000",
			responseCode: http.StatusOK,
			responseBody: "2000",
		},
	}

	handler := SetupHandler()

	for _, test := range tt {
		t.Run(test.target, func(t *testing.T) {
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, httptest.NewRequest(test.method, test.target, nil))
			if response.Code != test.responseCode {
				t.Errorf("Incorrect response code %v", response.Code)
			}
			if string(response.Body.Bytes()) != test.responseBody {
				t.Errorf("Incorrect response body %v", string(response.Body.Bytes()))
			}
		})
	}
}

func SetupHandler() http.Handler {
	handler := http.NewServeMux()

	tasksHandler := NewHandlerFunc("/api/", "tasks/", func(part string) http.HandlerFunc {
		return func(writer http.ResponseWriter, request *http.Request) {
			if part == "" {
				writer.WriteHeader(http.StatusNotFound)
				return
			}

			writer.WriteHeader(http.StatusOK)
			writer.Write([]byte(part))
		}
	})

	notesHandler := NewHandler("/api/", "notes/", func(noteID string) http.Handler {
		return &NoteHandler{ID: noteID}
	})

	handler.Handle(notesHandler.Route, notesHandler)
	handler.Handle(tasksHandler.Route, tasksHandler)

	return handler
}

type NoteHandler struct {
	ID string
}

func (handler *NoteHandler) ServeHTTP(writer http.ResponseWriter, r *http.Request) {
	if handler.ID == "" {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(handler.ID))
}
