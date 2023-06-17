package restserver

import natsserver "wb_internship/pkg/nats-server"

type Handler struct {
	*natsserver.NatsServer
}
