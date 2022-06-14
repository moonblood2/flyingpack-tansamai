ALTER TABLE public.order
ADD COLUMN sort_code VARCHAR(20),
ADD COLUMN line_code VARCHAR(5),
ADD COLUMN sorting_line_code VARCHAR(5),
ADD COLUMN dst_store_name VARCHAR(30)
