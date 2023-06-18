PROJECTNAME="wb_orders"
NATSPORT=9000

all: docker
	@ echo "  >  running nats server on port $(NATSPORT)"
	@ nats-server -js -p $(NATSPORT) &
	@ echo "  >  running project..."
	@ go run cmd/server/main.go

client:
	@ echo "  >  running client"
	@ go run cmd/client/main.go

docker:
	@ echo "  >  making docker container $(PROJECTNAME)..."
	@ sudo docker compose up -d
	@ make migration-up

# making only init migration (may be changed using ARGS="[version]")
migration-up:
	@ echo "  >  making migrations"
	@ sudo docker start orders_pg
	@ sleep 0.1
	@ cat schemas/0001_init.up.sql | sudo docker exec -i orders_pg  psql -U postgres -d postgres

# making only init migration (may be changed using ARGS="[version]")
migration-down:
	@ echo "  >  making migrations"
	@ sudo docker start orders_pg
	@ sleep 0.1
	@ cat schemas/0001_init.down.sql | sudo docker exec -i orders_pg  psql -U postgres -d postgres