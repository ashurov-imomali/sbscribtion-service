package api

import "net/http"

func (h *Handler) RegisterRouters(mux *http.ServeMux) {
	//stat
	{
		mux.HandleFunc("/subscriptions/total-cost", h.getTotalCost)
	}
	//crud
	{
		mux.HandleFunc("/subscriptions", h.subscriptionHandle)
		mux.HandleFunc("/subscriptions/", h.subscriptionHandleByID)
	}
}
