package server

import (
	"github.com/ashurov-imomali/sbscribtion-service/internal/api"
	"net/http"
	"time"
)

func NewServer(addr string, h *api.Handler) *http.Server {
	mux := http.NewServeMux()
	h.RegisterRouters(mux)

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return srv
}
