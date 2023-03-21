-- Create Users
CREATE SEQUENCE users_user_id_seq;
CREATE TABLE IF NOT EXISTS users
(
    user_id integer DEFAULT nextval('users_user_id_seq') PRIMARY KEY,
    login text NOT NULL UNIQUE,
    hash text NOT NULL
);
ALTER SEQUENCE users_user_id_seq OWNED BY users.user_id;

-- Create Orders
CREATE TABLE IF NOT EXISTS statuses
(
    status_id smallint PRIMARY KEY,
    status text NOT NULL
);
CREATE SEQUENCE orders_order_id_seq;
CREATE TABLE IF NOT EXISTS orders
(
    order_id integer DEFAULT nextval('orders_order_id_seq') PRIMARY KEY,
    user_id integer,
    number text NOT NULL,
    status_id smallint,
    accrual float ,
    uploaded_at timestamptz NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(user_id),
    CONSTRAINT fk_status
        FOREIGN KEY(status_id)
            REFERENCES statuses(status_id)

);
ALTER SEQUENCE orders_order_id_seq OWNED BY orders.order_id;
CREATE INDEX idx_orders_user
    ON orders(user_id);

-- Create Withdrawals
CREATE SEQUENCE withdrawals_draw_id_seq;
CREATE TABLE IF NOT EXISTS withdrawals
(
    draw_id integer DEFAULT nextval('withdrawals_draw_id_seq') PRIMARY KEY,
    user_id integer,
    order_id integer,
    sum float,
    processed_at timestamptz,
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(user_id),
    CONSTRAINT fk_order
        FOREIGN KEY(order_id)
            REFERENCES orders(order_id)
);
ALTER SEQUENCE withdrawals_draw_id_seq OWNED BY withdrawals.draw_id;
CREATE INDEX idx_withdrawals_user
    ON withdrawals(user_id);
