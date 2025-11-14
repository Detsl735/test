package answer

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
	router.HandleFunc("GET /answers/{id}", h.GetById)
	router.HandleFunc("POST /questions/{id}/answers/", h.Create)
	router.HandleFunc("DELETE /answers/{id}", h.Delete)
}

func (h *handler) GetById(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || idUint == 0 {
		handlers.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	ans, err := h.service.GetByID(r.Context(), uint(idUint))
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			handlers.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		h.logger.Errorf("get answer error: %v", err)
		handlers.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	handlers.WriteJSON(w, http.StatusOK, ans)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	qidStr := r.PathValue("id")
	qidUint, err := strconv.ParseUint(qidStr, 10, 64)
	if err != nil || qidUint == 0 {
		handlers.WriteError(w, http.StatusBadRequest, "invalid question id")
		return
	}

	var req CreateAnswerRequest
	if err := handlers.ReadJSON(r, &req); err != nil {
		h.logger.Errorf("failed to decode create answer request: %v", err)
		handlers.WriteError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			h.logger.Warnf("failed to close request body: %v", err)
		}
	}()

	req.QuestionID = uint(qidUint)

	ans, err := h.service.Create(r.Context(), &req)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmptyText),
			errors.Is(err, ErrEmptyUserID),
			errors.Is(err, ErrInvalidQuestion):
			handlers.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			h.logger.Errorf("create answer error: %v", err)
			handlers.WriteError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}

	handlers.WriteJSON(w, http.StatusCreated, ans)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	idUint, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || idUint == 0 {
		handlers.WriteError(w, http.StatusBadRequest, "invalid id")
		return
	}

	if err := h.service.Delete(r.Context(), uint(idUint)); err != nil {
		h.logger.Errorf("delete answer error: %v", err)
		handlers.WriteError(w, http.StatusInternalServerError, "internal error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
