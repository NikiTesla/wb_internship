CREATE TABLE IF NOT EXISTS delivery_info (
    id serial PRIMARY KEY,
    name VARCHAR(255),
    phone VARCHAR(15),
    zip VARCHAR(10),
    city VARCHAR(55),
    address VARCHAR(100),
    region VARCHAR(55),
    email VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS payment_info (
    id serial PRIMARY KEY,
    transactions VARCHAR(255),
    request_id VARCHAR(55),
    currency VARCHAR(10),
    providerr VARCHAR(55),
    amount INT,
    payment_dt INT,
    bank VARCHAR(55),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE IF NOT EXISTS orders (
    id serial PRIMARY KEY,
    order_id VARCHAR(255),
    track_number VARCHAR(255),
    entr VARCHAR(55),
    delivery_id INT NOT NULL,
    payment_id INT NOT NULL,
    locale VARCHAR(5),
    internal_signaturex VARCHAR(55),
    customer_id VARCHAR(55),
    delivery_service VARCHAR(55),
    shardkey VARCHAR(55),
    sm_id INT,
    date_created TIMESTAMP,
    oof__shard VARCHAR(55),
    CONSTRAINT delivery FOREIGN KEY(delivery_id) REFERENCES delivery_info(id),
    CONSTRAINT payment FOREIGN KEY(payment_id) REFERENCES payment_info(id)
);

CREATE TABLE IF NOT EXISTS items (
    id serial PRIMARY KEY,
    order_id INT,
    chrt_id INT,
    track_number VARCHAR(55),
    price INT,
    rid VARCHAR(100),
    name VARCHAR(55),
    sale INT,
    size VARCHAR(55),
    total_price INT,
    nm_id INT,
    brand VARCHAR(55),
    stat INT,
    CONSTRAINT orders FOREIGN KEY(order_id) REFERENCES orders(id)
);