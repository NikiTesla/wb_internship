# wb_internship
Repository for golang internship in wildberries

## L0
Service connects to nats and save orders (in proper model format) in database and in cache.
In case of mismatch between cache and database recreate cache from database.
