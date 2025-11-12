package question

import (
	"fmt"
	"net/http"
	"testTask/internal/handlers"
	"testTask/pkg/logging"
)

type handler struct {
	logger *logging.Logger
}

func NewHandler(logger *logging.Logger) handlers.Handler {
	return &handler{
		logger: logger,
	}
}

func (h *handler) Register(router *http.ServeMux) {
	router.HandleFunc("GET /questions/", h.GetQuestions)
	router.HandleFunc("GET /questions/{id}", h.GetUserById)
	router.HandleFunc("POST /questions/", h.CreateUser)
	router.HandleFunc("DELETE /questions/{id}", h.DeleteUser)
}

func (h *handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "List of questions")
}

func (h *handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Question")
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create user")
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete user")
}
