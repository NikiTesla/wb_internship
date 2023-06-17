package restserver

import "github.com/gorilla/mux"

func (h Handler) InitRouter() *mux.Router {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/order", h.OrderGet)

	return rtr
}
