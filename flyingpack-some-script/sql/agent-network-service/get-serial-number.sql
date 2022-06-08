SELECT o.reference_no, p.name, osn.serial_number FROM public.order o
LEFT JOIN public.order_product op
ON o.an_order_id = op.an_order_id
LEFT JOIN public.product p
ON op.an_product_id = p.an_product_id
LEFT JOIN public.an_order_product_serial_number osn
ON op.an_order_product_id = osn.an_order_product_id
WHERE osn.serial_number != '' AND o.created_at >= '2021-06-07'