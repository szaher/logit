package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type httpServer struct {
	Log *Log
}

func newHttpServer() *httpServer{
	return &httpServer{
		Log: NewLog(),
	}
}

type ProduceRequest struct {
	Record Record `json:"record"`
}

type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

type ConsumeResponse struct {
	Record Record `json:"record"`
}

func NewHttpServer(address string) *http.Server{

	r := mux.NewRouter()
	httpsrv := newHttpServer()
	r.HandleFunc("/", httpsrv.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsrv.handleConsume).Methods("GET")

	return &http.Server{
		Addr:              address,
		Handler:           r,
		TLSConfig:         nil,
		ReadTimeout:       0,
		ReadHeaderTimeout: 0,
		WriteTimeout:      0,
		IdleTimeout:       0,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}
}


func (s *httpServer) handleProduce(w http.ResponseWriter, r *http.Request){
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	off, err := s.Log.Append(req.Record)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *httpServer) handleConsume(w http.ResponseWriter, r *http.Request){
	var req ConsumeRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	record, err := s.Log.Read(req.Offset)
	if err != nil{
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}


