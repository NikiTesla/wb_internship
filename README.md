# wb_internship
Repository for golang internship in wildberries

## L0
Service connects to nats and save orders (in proper model format) in database and in cache.
In case of mismatch between cache and database recreate cache from database.

-[] To run server you should just run Makefile with `make`
-[] To get order by id using http request you should send GET request to localhost:localport/order?id={id} 