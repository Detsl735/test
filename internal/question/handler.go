package question

import (
	"errors"
	"net/http"
	"strconv"
	"testTask/internal/handlers"
	"testTask/pkg/logging"
)

type handler struct {
	logger  *logging.Logger
	service Service
}

func NewHandler(logger *logging.Logger, service Service) handlers.Handler {
	return &handler{
		logger:  logger,
		service: service,
	}
}

func (h *handler) Register(router *http.ServeMux) {
	router.HandleFunc("GET /questions/", h.GetAll)
	router.HandleFunc("GET /questions/{id}", h.GetById)
	router.HandleFunc("POST /questions/", h.Create)
	router.HandleFunc("DELETE /questions/{id}", h.Delete)
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.GetAll(r.Context())
	if err != nil {
		h.logger.Errorf("list questions error: %v", err)
		handlers.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	handlers.WriteJSON(w, http.StatusOK, list)
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	q, err := h.service.GetByID(r.Context(), uint(idUint))
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			handlers.WriteError(w, http.StatusNotFound, err.Error())
		default:
			h.logger.Errorf("get question error: %v", err)
			handlers.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	handlers.WriteJSON(w, http.StatusOK, q)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateQuestionRequest
	if err := handlers.ReadJSON(r, &req); err != nil {
		h.logger.Errorf("failed to decode request: %v", err)
		handlers.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Warnf("failed to close request body: %v", err)
		}
	}()

	q, err := h.service.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyText):
			handlers.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			h.logger.Errorf("create question error: %v", err)
			handlers.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, q)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		handlers.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(r.Context(), uint(idUint)); err != nil {
		h.logger.Errorf("delete question error: %v", err)
		handlers.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
