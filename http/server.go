package http

import (
	"fmt"
	"net/http"

	"github.com/barkbay/zac/rate"
	"github.com/gorilla/mux"
)

type HttpServer struct {
	warningRates *rate.WarningRates
}

func NewHttpServer(wr *rate.WarningRates) *HttpServer {
	return &HttpServer{
		warningRates: wr,
	}
}

func (s *HttpServer) Listen() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", s.index)
	router.HandleFunc("/{namespace}", s.rate)
	http.ListenAndServe(":8080", router)
}

func (s *HttpServer) rate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if namespace, ok := vars["namespace"]; ok {
		rate, exists := s.warningRates.GetWarningRate(namespace)
		if exists && rate.RateCounter.Rate() > 0 {
			fmt.Fprintf(w, "%+v", rate)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404- No anormal rate"))
}

func (s *HttpServer) index(w http.ResponseWriter, r *http.Request) {
	s.warningRates.Dump(w)
}