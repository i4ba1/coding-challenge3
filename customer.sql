create database customer_db;

create table tbl_customer
(
    customer_id   varchar(64) not null
        constraint tbl_customer_pk
            primary key,
    customer_name varchar(80) not null,
    phone_number  varchar(20) not null,
    email         varchar(50) not null,
    dob           date,
    sex           boolean     not null,
    salt          bytea       not null,
    password      text        not null,
    created_date  timestamp with time zone
);

alter table tbl_customer
    owner to mni;

INSERT INTO public.tbl_customer (customer_id, customer_name, phone_number, email, dob, sex, salt, password, created_date) VALUES ('wwnzdardea', 'Muhamad Nizar Iqbal', '087882458829', 'muhamad.iqbal1983@gmail.com', '2020-10-16', true, E'\\xCF20CBD73F92A553FDA3C705E4F6E26C', 'OzcM2FbhSUFCz1dx9gTcLnFcLlPvKXMlDNAoc-KrYdHsJ26B4aV5xMt-gD06NY_JueO5SNxEryMfqDhLwpIdGw==', '2020-10-16 14:14:25.266025');
INSERT INTO public.tbl_customer (customer_id, customer_name, phone_number, email, dob, sex, salt, password, created_date) VALUES ('hkizilwxkx', 'Muhamad Nizar Iqbal', '087882458839', 'muhamad.iqbal1981@gmail.com', '2020-10-16', true, E'\\x2993B7130C08427E22AEF7CB45D44B51', 'YJUQNNpsLiX_NT17jc1GIOevMC8orLQiZgYcKOVdiPOc3ia3VTaNBiH0KUdYbUmAnNhnsybVDDCpKFXDZKa93A==', '2020-10-16 16:21:03.853353');

create table tbl_product
(
    product_id   varchar(64)              not null
        constraint pk_product_id
            primary key,
    product_name varchar(80)              not null,
    basic_price  money                    not null,
    created_date timestamp with time zone not null
);

alter table tbl_product
    owner to mni;

create table tbl_payment_method
(
    payment_method_id varchar(64)              not null
        constraint tbl_payment_method_pkey
            primary key,
    method_name       varchar(70)              not null,
    code              varchar(10)              not null,
    created_date      timestamp with time zone not null
);

alter table tbl_payment_method
    owner to mni;

create table tbl_order
(
    order_id          varchar(64) not null
        constraint pk_order_id
            primary key,
    customer_id       varchar(64) not null
        constraint tbl_order_tbl_customer_customer_id_fk
            references tbl_customer,
    order_number      integer     not null,
    payment_method_id varchar(64) not null
        constraint tbl_order_tbl_payment_method_payment_method_id_fk
            references tbl_payment_method
);

alter table tbl_order
    owner to mni;

create table tbl_order_detail
(
    order_detail_id varchar(64)              not null
        constraint tbl_order_detail_pk
            primary key,
    order_id        varchar(64)              not null
        constraint tbl_order_detail_tbl_order_order_id_fk
            references tbl_order,
    product_id      varchar(64)              not null
        constraint tbl_order_detail_tbl_product_product_id_fk
            references tbl_product,
    quantity        integer                  not null,
    created_date    timestamp with time zone not null
);

alter table tbl_order_detail
    owner to mni;

