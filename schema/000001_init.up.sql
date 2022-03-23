CREATE TABLE delivery
(
    id serial not null unique,
    name varchar(255) not null,
    phone varchar(12) not null,
    zip varchar(255) not null,
    city varchar(255) not null,
    address varchar(255) not null,
    region varchar(255) not null,
    email varchar(255) not null,
    primary key(phone, email)
);

CREATE TABLE item
(
    id serial not null unique,
    chrt_id serial not null unique,
    track_number varchar(255) not null,
    price int not null,
    rid varchar(255) not null,
    name varchar(255) not null,
    sale int not null,
    size varchar(255) not null,
    total_price int not null,
    nm_id int not null,
    brand varchar(255) not null,
    status int not null
);

CREATE TABLE payment
(
   id serial not null unique,
   transaction varchar(255) not null unique,
   request_id varchar(255),
   currency varchar(255) not null,
   provider varchar(255) not null,
   amount int not null,
   payment_dt int not null,
   bank varchar(255) not null,
   delivery_cost int not null,
   goods_total int not null,
   custom_fee int not null
);

CREATE TABLE orders
(
    id serial not null unique,
    order_uid varchar(255) not null unique,
    track_number varchar(255) not null unique,
    entry varchar(255) not null,
    payment_id int references  payment(id) on delete cascade not null,
    delivery_id int references delivery(id) on delete cascade not null,
    locale varchar(255) not null,
    internal_signature varchar(255),
    customer_id varchar(255) not null,
    delivery_service varchar(255) not null,
    shardkey varchar(255) not null,
    sm_id int not null,
    data_created date not null,
    oof_shard varchar(255) not null,
    primary key(order_uid, track_number, customer_id)

);

CREATE TABLE order_items
(
    id serial not null unique,
    order_id int references orders(id) on delete cascade not null,
    item_id int references item(id) on delete cascade not null
)
