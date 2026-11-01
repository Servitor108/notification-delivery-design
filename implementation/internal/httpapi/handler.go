package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Servitor108/notification-delivery-design/implementation/internal/notify"
)

type Handler struct{ service *notify.Service }
func New(s *notify.Service)*Handler{return &Handler{service:s}}
func(h *Handler) Register(mux *http.ServeMux){
	mux.HandleFunc("GET /healthz",func(w http.ResponseWriter,r *http.Request){write(w,200,map[string]string{"status":"ok"})})
	mux.HandleFunc("POST /v1/notifications",h.create)
	mux.HandleFunc("GET /v1/notifications/{id}",h.get)
}
func(h *Handler)create(w http.ResponseWriter,r *http.Request){
	var in struct{Recipient string `json:"recipient_ref"`;Channel string `json:"channel"`;Template string `json:"template_id"`;Payload map[string]string `json:"payload"`}
	d:=json.NewDecoder(http.MaxBytesReader(w,r.Body,1<<20));d.DisallowUnknownFields()
	if err:=d.Decode(&in);err!=nil{write(w,400,map[string]string{"error":"invalid request body"});return}
	key:=strings.TrimSpace(r.Header.Get("Idempotency-Key"));n,created,err:=h.service.Enqueue(r.Context(),notify.Notification{IdempotencyKey:key,Recipient:in.Recipient,Channel:in.Channel,Template:in.Template,Payload:in.Payload})
	if err!=nil{status:=400;if errors.Is(err,notify.ErrDuplicate){status=409};write(w,status,map[string]string{"error":err.Error()});return}
	status:=202;if !created{status=200};write(w,status,n)
}
func(h *Handler)get(w http.ResponseWriter,r *http.Request){n,err:=h.service.Get(r.Context(),r.PathValue("id"));if err!=nil{write(w,404,map[string]string{"error":"notification not found"});return};write(w,200,n)}
func write(w http.ResponseWriter,status int,v any){w.Header().Set("Content-Type","application/json");w.WriteHeader(status);_ = json.NewEncoder(w).Encode(v)}
