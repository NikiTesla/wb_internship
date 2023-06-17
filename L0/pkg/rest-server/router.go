package restserver

import "github.com/gorilla/mux"

func (h Handler) InitRouter() *mux.Router {
	rtr := mux.NewRouter()

	rtr.HandleFunc("/order/{id:[0-9]+}", h.OrderGet)

	return rtr
}
