package storage

const createDatabase = `
-- Create Users
CREATE SEQUENCE IF NOT EXISTS users_user_id_seq;
CREATE TABLE IF NOT EXISTS users
(
    user_id integer DEFAULT nextval('users_user_id_seq') PRIMARY KEY,
    login text NOT NULL UNIQUE,
    hash text NOT NULL
);
ALTER SEQUENCE users_user_id_seq OWNED BY users.user_id;

-- Create Orders
CREATE SEQUENCE IF NOT EXISTS  orders_order_id_seq;
CREATE TABLE IF NOT EXISTS orders
(
    order_id integer DEFAULT nextval('orders_order_id_seq') PRIMARY KEY,
    user_id integer,
    number text NOT NULL UNIQUE,
    status text,
    accrual integer ,
    uploaded_at timestamp NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(user_id)
);
ALTER SEQUENCE orders_order_id_seq OWNED BY orders.order_id;
CREATE INDEX IF NOT EXISTS idx_orders_user
    ON orders(user_id);

-- Create Withdrawals
CREATE SEQUENCE IF NOT EXISTS  withdrawals_draw_id_seq;
CREATE TABLE IF NOT EXISTS withdrawals
(
    draw_id integer DEFAULT nextval('withdrawals_draw_id_seq') PRIMARY KEY,
    user_id integer,
    order_num text,
    sum integer,
    processed_at timestamp,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(user_id)
);
ALTER SEQUENCE withdrawals_draw_id_seq OWNED BY withdrawals.draw_id;
CREATE INDEX IF NOT EXISTS idx_withdrawals_user
    ON withdrawals(user_id);
`
