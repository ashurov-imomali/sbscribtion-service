package api

import (
	"encoding/json"
	"github.com/ashurov-imomali/sbscribtion-service/internal/models"
	"github.com/ashurov-imomali/sbscribtion-service/internal/usecase"
	"github.com/google/uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

type Handler struct {
	s *usecase.Service
}

func New(s *usecase.Service) *Handler {
	return &Handler{s: s}
}

func (h *Handler) subscriptionHandle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.createSubscription(w, r)
	case http.MethodPut:
		h.updateSubscription(w, r)
	case http.MethodGet:
		h.getSubscriptionList(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) subscriptionHandleByID(w http.ResponseWriter, r *http.Request) {
	strId := strings.Trim(r.URL.Path, "/subscriptions/")
	id, err := uuid.Parse(strId)
	if err != nil {
		writeError(w, "invalid ID")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getSubscription(w, r, id)
	case http.MethodDelete:
		h.deleteSubscription(w, r, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createSubscription(w http.ResponseWriter, r *http.Request) {
	var request models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, "invalid json")
		log.Println(err)
		return
	}

	status, err := h.s.CreateSubscription(&request)
	if err != nil {
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, status, request)
}

func (h *Handler) getSubscription(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	subscribe, status, err := h.s.GetSubscribe(id)
	if err != nil {
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, status, subscribe)
}

func (h *Handler) getSubscriptionList(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	var filter models.SubscriptionFilter
	userID := query.Get("user_id")
	if userID != "" {
		if _, err := uuid.Parse(userID); err != nil {
			http.Error(w, "invalid user_id", http.StatusBadRequest)
			return
		}
		filter.UserID = &userID
	}

	if serviceName := query.Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if start := query.Get("start_date"); start != "" {
		t, err := time.Parse("2006-01-02", start)
		if err != nil {
			http.Error(w, "invalid start_date, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		filter.StartDate = &t
	}

	if end := query.Get("end_date"); end != "" {
		t, err := time.Parse("2006-01-02", end)
		if err != nil {
			http.Error(w, "invalid end_date, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		filter.EndDate = &t
	}

	subs, status, err := h.s.GetSubscriptions(filter)
	if err != nil {
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, status, subs)
}

func (h *Handler) updateSubscription(w http.ResponseWriter, r *http.Request) {
	var request models.Subscription
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeError(w, "INVALID JSON")
		return
	}
	updateSub, status, err := h.s.UpdateSubscription(request)
	if err != nil {
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, status, updateSub)
}

func (h *Handler) deleteSubscription(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	status, err := h.s.DeleteSubscription(id)
	if err != nil {
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}
	w.WriteHeader(status)
}

func (h *Handler) getTotalCost(w http.ResponseWriter, r *http.Request) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	userIDStr := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")

	if fromStr == "" || toStr == "" {
		http.Error(w, "from and to are required in YYYY-MM format", http.StatusUnprocessableEntity)
		return
	}

	fromDate, err := time.Parse("2006-01", fromStr)
	if err != nil {
		http.Error(w, "invalid from format, expected YYYY-MM", http.StatusUnprocessableEntity)
		return
	}

	toDate, err := time.Parse("2006-01", toStr)
	if err != nil {
		http.Error(w, "invalid to format, expected YYYY-MM", http.StatusUnprocessableEntity)
		return
	}

	var userID uuid.UUID
	if userIDStr != "" {
		userID, err = uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "invalid user_id format", http.StatusUnprocessableEntity)
			return
		}
	}

	total, status, err := h.s.GetTotalCost(fromDate, toDate, userID, serviceName)
	if err != nil {
		writeJSON(w, status, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, status, total)
}
