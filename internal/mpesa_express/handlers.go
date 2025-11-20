package mpesaexpress

import (
	"log"
	"net/http"

	"github.com/mraiyuu/M-Pesa/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) InitiateSTK(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		json.Write(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	token, err := h.service.GetAccessToken(r.Context())
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, token)
}
