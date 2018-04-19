--建立用户表并生成100万用户数据,用户注册时间随机分布

CREATE TABLE public.users
(
    id serial not null unique,
    nickname text not null,
    mobileno text not null unique,
    createtime timestamp not null default now()
);

COMMENT ON TABLE public.users IS '用户表';
COMMENT ON COLUMN public.users.id IS 'id号';
COMMENT ON COLUMN public.users.nickname IS '呢称';
COMMENT ON COLUMN public.users.mobileno IS '手机号';
COMMENT ON COLUMN public.users.createtime IS '注册时间';

--生成用户数据

INSERT INTO public.users(nickname,mobileno,createtime) select
 t::text,t::text,'2015-12-01'::timestamp + (random()*31449600)::integer * 
interval '1 seconds'   FROM  generate_series(13800000001,13801000000) as t;

--建立索引

CREATE INDEX users_createtime_idx ON users USING BTREE(createtime);

--建立订单表并生成1000万订单数据,用户下单时间随机分布

CREATE TABLE orders
(
    id serial not null unique,
    createtime timestamp not null default now(),
    users_id integer not null,
    goods_name text not null
);

COMMENT ON TABLE public.orders IS '运力需求订单表';
COMMENT ON COLUMN public.orders.id IS 'id号';
COMMENT ON COLUMN public.orders.createtime IS '下单时间';
COMMENT ON COLUMN public.orders.users_id IS '用户id号';
COMMENT ON COLUMN public.orders.goods_name IS '货源名称';

--这是生成订单的函数，每天30000订单

CREATE OR REPLACE FUNCTION orders_make(a_date date) RETURNS TEXT AS
$$
BEGIN
    INSERT INTO orders(users_id,goods_name,createtime)
    SELECT
        users.id,
        md5(random()::Text),
        a_date::timestamp + (random()*86399)::integer * interval '1 seconds'
    FROM 
        users
    WHERE
        users.createtime < a_date
    ORDER BY 
        random()
    LIMIT 30000
    ;
    RETURN 'OK';
END;
$$
LANGUAGE PLPGSQL;

COMMENT ON FUNCTION orders_make(a_date date) IS '生成订单数据';

--调用函数生成数据，时间比较长，耐心等

SELECT orders_make('2015-12-01'::date + t) FROM generate_series(0,365) as t;

--把生成的数据序号重新生成

create sequence order_id_seq_tmp;
copy (select nextval('order_id_seq_tmp'),createtime,users_id,goods_name from orders 
order by createtime) to '/tmp/orders.txt';
truncate table orders;
copy orders from '/tmp/orders.txt';
SELECT SETVAL('orders_id_seq',(select id FROM orders ORDER BY id DESC LIMIT 1));
DROP sequence order_id_seq_tmp;

--建立索引

CREATE INDEX orders_createtime_idx ON orders USING BTREE(createtime);
CREATE INDEX orders_users_id_idx ON orders USING BTREE(users_id);
--数据生成后记得做一下统计

vacuum analyze