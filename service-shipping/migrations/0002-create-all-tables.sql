CREATE TABLE public.contact (
    id smallint NOT NULL,
    name character varying(256),
    phone_number character varying(10),
    address character varying(256),
    district character varying(128),
    state character varying(128),
    province character varying(128),
    postcode character varying(8),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL,
    user_id uuid NOT NULL
);
CREATE TABLE public.destination (
    id bigint NOT NULL,
    name character varying(256) NOT NULL,
    phone_number character varying(10) NOT NULL,
    address character varying(256) NOT NULL,
    district character varying(128) NOT NULL,
    state character varying(128) NOT NULL,
    province character varying(128) NOT NULL,
    postcode character varying(8) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0
);
CREATE SEQUENCE public.destination_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.destination_id_seq OWNED BY public.destination.id;
CREATE TABLE public.order_parcel (
    order_parcel_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    sender_id bigint,
    origin_id bigint NOT NULL,
    destination_id bigint NOT NULL,
    provider_code smallint NOT NULL,
    price numeric NOT NULL,
    payment_method smallint NOT NULL,
    weight real NOT NULL,
    width real NOT NULL,
    length real NOT NULL,
    height real NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL,
    tracking_code character varying(32) DEFAULT ''::character varying NOT NULL,
    status character varying(32) DEFAULT ''::character varying NOT NULL,
    cod_amount numeric DEFAULT 0 NOT NULL,
    user_id uuid,
    cod_status character varying(20) DEFAULT ''::character varying NOT NULL,
    cod_transferred_date date,
    status_completed_date date
);
CREATE SEQUENCE public.order_parcel_flash_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TABLE public.order_parcel_flash (
    id bigint DEFAULT nextval('public.order_parcel_flash_id_seq'::regclass) NOT NULL,
    order_parcel_id uuid NOT NULL,
    pno character varying(20) NOT NULL,
    state smallint,
    state_text character varying(32),
    cod_amount numeric,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0
);
CREATE SEQUENCE public.order_parcel_shippop_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
CREATE TABLE public.order_parcel_shippop (
    id bigint DEFAULT nextval('public.order_parcel_shippop_id_seq'::regclass) NOT NULL,
    order_parcel_id uuid NOT NULL,
    purchase_id bigint NOT NULL,
    status character varying(32) DEFAULT ''::character varying NOT NULL,
    courier_code character varying(32) DEFAULT ''::character varying NOT NULL,
    courier_tracking_code character varying(32) DEFAULT ''::character varying NOT NULL,
    tracking_code character varying(32) DEFAULT ''::character varying NOT NULL,
    cod_amount numeric DEFAULT 0 NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL
);
CREATE TABLE public.order_parcel_shippop_flash (
    id bigint NOT NULL,
    order_parcel_shippop_id bigint NOT NULL,
    sort_code character varying(13),
    dst_code character varying(50),
    sorting_line_code character varying(5),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL
);
CREATE SEQUENCE public.order_parcel_shippop_flash_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.order_parcel_shippop_flash_id_seq OWNED BY public.order_parcel_shippop_flash.id;
CREATE TABLE public.order_product (
    id bigint NOT NULL,
    sender_id bigint NOT NULL,
    product_id integer NOT NULL,
    quantity integer NOT NULL,
    payment_method smallint NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL,
    user_id uuid
);
CREATE SEQUENCE public.order_product_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.order_product_id_seq OWNED BY public.order_product.id;
CREATE TABLE public.origin (
    id bigint NOT NULL,
    name character varying(256) NOT NULL,
    phone_number character varying(10) NOT NULL,
    address character varying(256) NOT NULL,
    district character varying(128) NOT NULL,
    state character varying(128) NOT NULL,
    province character varying(128) NOT NULL,
    postcode character varying(8) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0
);
CREATE SEQUENCE public.origin_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.origin_id_seq OWNED BY public.origin.id;
CREATE TABLE public.product (
    id integer NOT NULL,
    name character varying(64) NOT NULL,
    price numeric NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL,
    user_id uuid NOT NULL
);
CREATE SEQUENCE public.product_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.product_id_seq OWNED BY public.product.id;
CREATE TABLE public.sender (
    id bigint NOT NULL,
    sender_type smallint NOT NULL,
    national_id_number character varying(13),
    tax_id_number character varying(13),
    passport_number character varying(9),
    birth_date date,
    name character varying(256) NOT NULL,
    phone_number character varying(10) NOT NULL,
    address character varying(256) NOT NULL,
    district character varying(128) NOT NULL,
    state character varying(128) NOT NULL,
    province character varying(128) NOT NULL,
    postcode character varying(8) NOT NULL,
    created_at timestamp with time zone,
    deleted_at bigint DEFAULT 0
);
CREATE SEQUENCE public.sender_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.sender_id_seq OWNED BY public.sender.id;
CREATE SEQUENCE public.shop_id_seq
    AS smallint
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.shop_id_seq OWNED BY public.contact.id;
CREATE TABLE public."user" (
    user_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    email character varying(320) NOT NULL,
    name character varying(256) NOT NULL,
    role smallint NOT NULL,
    password character varying(60) NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL
);
ALTER TABLE ONLY public.contact ALTER COLUMN id SET DEFAULT nextval('public.shop_id_seq'::regclass);
ALTER TABLE ONLY public.destination ALTER COLUMN id SET DEFAULT nextval('public.destination_id_seq'::regclass);
ALTER TABLE ONLY public.order_parcel_shippop_flash ALTER COLUMN id SET DEFAULT nextval('public.order_parcel_shippop_flash_id_seq'::regclass);
ALTER TABLE ONLY public.order_product ALTER COLUMN id SET DEFAULT nextval('public.order_product_id_seq'::regclass);
ALTER TABLE ONLY public.origin ALTER COLUMN id SET DEFAULT nextval('public.origin_id_seq'::regclass);
ALTER TABLE ONLY public.product ALTER COLUMN id SET DEFAULT nextval('public.product_id_seq'::regclass);
ALTER TABLE ONLY public.sender ALTER COLUMN id SET DEFAULT nextval('public.sender_id_seq'::regclass);
ALTER TABLE ONLY public.destination
    ADD CONSTRAINT destination_phone_number_deleted_at_key UNIQUE (phone_number, deleted_at);
ALTER TABLE ONLY public.destination
    ADD CONSTRAINT destination_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.order_parcel_flash
    ADD CONSTRAINT order_parcel_flash_order_parcel_id_deleted_at_key UNIQUE (order_parcel_id, deleted_at);
ALTER TABLE ONLY public.order_parcel_flash
    ADD CONSTRAINT order_parcel_flash_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.order_parcel
    ADD CONSTRAINT order_parcel_pkey PRIMARY KEY (order_parcel_id);
ALTER TABLE ONLY public.order_parcel_shippop_flash
    ADD CONSTRAINT order_parcel_shippop_flash_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.order_parcel_shippop
    ADD CONSTRAINT order_parcel_shippop_order_parcel_id_deleted_at_key UNIQUE (order_parcel_id, deleted_at);
ALTER TABLE ONLY public.order_parcel_shippop
    ADD CONSTRAINT order_parcel_shippop_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.order_product
    ADD CONSTRAINT order_product_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.origin
    ADD CONSTRAINT origin_phone_number_deleted_at_key UNIQUE (phone_number, deleted_at);
ALTER TABLE ONLY public.origin
    ADD CONSTRAINT origin_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.sender
    ADD CONSTRAINT sender_phone_number_deleted_at_key UNIQUE (phone_number, deleted_at);
ALTER TABLE ONLY public.sender
    ADD CONSTRAINT sender_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.contact
    ADD CONSTRAINT shop_phone_number_deleted_at_key UNIQUE (phone_number, deleted_at);
ALTER TABLE ONLY public.contact
    ADD CONSTRAINT shop_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_email_deleted_at_key UNIQUE (email, deleted_at);
ALTER TABLE ONLY public."user"
    ADD CONSTRAINT user_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.contact
    ADD CONSTRAINT contact_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(user_id) NOT VALID;
ALTER TABLE ONLY public.order_parcel
    ADD CONSTRAINT order_parcel_destination_id_fkey FOREIGN KEY (destination_id) REFERENCES public.destination(id);
ALTER TABLE ONLY public.order_parcel_flash
    ADD CONSTRAINT order_parcel_flash_order_parcel_id_fkey FOREIGN KEY (order_parcel_id) REFERENCES public.order_parcel(order_parcel_id);
ALTER TABLE ONLY public.order_parcel
    ADD CONSTRAINT order_parcel_origin_id_fkey FOREIGN KEY (origin_id) REFERENCES public.origin(id);
ALTER TABLE ONLY public.order_parcel
    ADD CONSTRAINT order_parcel_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public.sender(id);
ALTER TABLE ONLY public.order_parcel_shippop_flash
    ADD CONSTRAINT order_parcel_shippop_flash_order_parcel_shippop_id_fkey FOREIGN KEY (order_parcel_shippop_id) REFERENCES public.order_parcel_shippop(id);
ALTER TABLE ONLY public.order_parcel_shippop
    ADD CONSTRAINT order_parcel_shippop_order_parcel_id_fkey FOREIGN KEY (order_parcel_id) REFERENCES public.order_parcel(order_parcel_id);
ALTER TABLE ONLY public.order_parcel
    ADD CONSTRAINT order_parcel_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(user_id) NOT VALID;
ALTER TABLE ONLY public.order_product
    ADD CONSTRAINT order_product_product_id_fkey FOREIGN KEY (product_id) REFERENCES public.product(id);
ALTER TABLE ONLY public.order_product
    ADD CONSTRAINT order_product_sender_id_fkey FOREIGN KEY (sender_id) REFERENCES public.sender(id);
ALTER TABLE ONLY public.order_product
    ADD CONSTRAINT order_product_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(user_id) NOT VALID;
ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_user_id_fkey FOREIGN KEY (user_id) REFERENCES public."user"(user_id) NOT VALID;