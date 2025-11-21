package mpesaexpress

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mraiyuu/M-Pesa/internal/response"
)

type handler struct {
	service Service
}

type InitiateMpesaExpressBody struct {
	PhoneNumber string `json:"phone_number"`
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) InitiateMpesaExpress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.Write(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	defer r.Body.Close()

	var req InitiateMpesaExpressBody
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		response.Write(w,http.StatusInternalServerError, "invalid body request",)
		return
	}

	if req.PhoneNumber == "" {
		http.Error(w, "phone number missing", http.StatusBadRequest)
		return 
	}

	stkPush, err := h.service.InitiateSTK(r.Context(), req.PhoneNumber)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Write(w, http.StatusOK, stkPush)
}
