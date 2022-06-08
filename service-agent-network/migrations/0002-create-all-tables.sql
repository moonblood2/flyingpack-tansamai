CREATE TABLE public.an_order_product_serial_number (
    an_order_product_serial_number_id uuid NOT NULL,
    an_order_product_id uuid NOT NULL,
    serial_number_start character varying(30) DEFAULT ''::character varying,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL,
    serial_number_end character varying(30) DEFAULT ''::character varying
);
CREATE TABLE public.bank_account (
    bank_account_id bigint NOT NULL,
    bank character varying(128) NOT NULL,
    account_no character varying(15) NOT NULL,
    account_name character varying(256) NOT NULL,
    email character varying(128) NOT NULL,
    an_order_id uuid NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL,
    fi_code character varying(4) DEFAULT ''::character varying NOT NULL
);
CREATE SEQUENCE public.bank_account_bank_account_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
ALTER SEQUENCE public.bank_account_bank_account_id_seq OWNED BY public.bank_account.bank_account_id;
CREATE TABLE public."order" (
    an_order_id uuid NOT NULL,
    user_id uuid NOT NULL,
    sp_order_parcel_id uuid,
    reference_no character varying(30) NOT NULL,
    des_name character varying(256),
    des_phone_number character varying(20),
    des_address character varying(256),
    des_subdistrict character varying(128),
    des_district character varying(128),
    des_province character varying(128),
    des_postcode character varying(8),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL,
    courier_code smallint,
    cod_amount numeric,
    fulfillment_status smallint DEFAULT 0 NOT NULL,
    shipping_status character varying(32) DEFAULT ''::character varying NOT NULL,
    cod_status character varying(20) DEFAULT ''::character varying NOT NULL,
    status_completed_date date,
    cod_transferred_date date,
    packaged_imageg_url character varying(256) DEFAULT ''::character varying,
    tracking_codes text[],
    jna_cod_transferred_date date
);
COMMENT ON COLUMN public."order".an_order_id IS 'an_order_id is an id of order table of AgenNetwork service';
COMMENT ON COLUMN public."order".sp_order_parcel_id IS 'sp_order_parcel_id is an id from order parcel id in the Shipping service';
CREATE TABLE public.order_product (
    an_order_product_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    an_order_id uuid NOT NULL,
    an_product_id uuid NOT NULL,
    quantity integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL
);
CREATE TABLE public.product (
    an_product_id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    product_code character varying(20) NOT NULL,
    name character varying(50),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL,
    img_url character varying(200) DEFAULT ''::character varying NOT NULL,
    serial_regex character varying(64) DEFAULT ''::character varying NOT NULL,
    robotic_sku character varying(20) DEFAULT ''::character varying
);
CREATE TABLE public.slip (
    user_id uuid NOT NULL,
    message character varying,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0 NOT NULL
);
CREATE TABLE public.user_api (
    user_id uuid NOT NULL,
    api_key_prefix character varying(9),
    api_key_hash character varying(128),
    hook_url_order character varying(256),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    deleted_at bigint DEFAULT 0
);
ALTER TABLE ONLY public.bank_account ALTER COLUMN bank_account_id SET DEFAULT nextval('public.bank_account_bank_account_id_seq'::regclass);
ALTER TABLE ONLY public.an_order_product_serial_number
    ADD CONSTRAINT an_order_product_serial_number_pkey PRIMARY KEY (an_order_product_serial_number_id);
ALTER TABLE ONLY public.user_api
    ADD CONSTRAINT api_info_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.bank_account
    ADD CONSTRAINT bank_account_pkey PRIMARY KEY (bank_account_id);
ALTER TABLE ONLY public."order"
    ADD CONSTRAINT delivery_pkey PRIMARY KEY (an_order_id);
ALTER TABLE ONLY public.order_product
    ADD CONSTRAINT order_product_pkey PRIMARY KEY (an_order_product_id);
ALTER TABLE ONLY public."order"
    ADD CONSTRAINT order_reference_no_user_id_key UNIQUE (reference_no, user_id);
ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_pkey PRIMARY KEY (an_product_id);
ALTER TABLE ONLY public.product
    ADD CONSTRAINT product_product_code_user_id_deleted_at_key UNIQUE (product_code, user_id, deleted_at);
ALTER TABLE ONLY public.slip
    ADD CONSTRAINT slip_pkey PRIMARY KEY (user_id);
ALTER TABLE ONLY public.an_order_product_serial_number
    ADD CONSTRAINT an_order_product_serial_number_an_order_product_id_fkey FOREIGN KEY (an_order_product_id) REFERENCES public.order_product(an_order_product_id);
ALTER TABLE ONLY public.bank_account
    ADD CONSTRAINT bank_account_an_order_id_fkey FOREIGN KEY (an_order_id) REFERENCES public."order"(an_order_id);
ALTER TABLE ONLY public.order_product
    ADD CONSTRAINT order_product_an_order_id_fkey FOREIGN KEY (an_order_id) REFERENCES public."order"(an_order_id) NOT VALID;
ALTER TABLE ONLY public.order_product
    ADD CONSTRAINT order_product_an_product_id_fkey FOREIGN KEY (an_product_id) REFERENCES public.product(an_product_id) NOT VALID;