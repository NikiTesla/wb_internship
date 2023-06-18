package restserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func (h *Handler) OrderGet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("incorrect format of id"))
		return
	}

	order, err := h.NatsServer.GetFromCache(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("No such order saved"))
		return
	}

	data, err := json.Marshal(&order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("cannot marshall order, error: %s", err)
		return
	}

	w.Write(data)
}
