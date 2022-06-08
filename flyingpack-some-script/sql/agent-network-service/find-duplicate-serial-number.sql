WITH serial_number_cte AS (
	SELECT opsn.*, c.count FROM (
		SELECT serial_number_start, COUNT(serial_number_start) FROM "public".an_order_product_serial_number
		WHERE serial_number_start!=''
		GROUP BY serial_number_start
	) AS c
	LEFT JOIN an_order_product_serial_number opsn
	ON opsn.serial_number_start = c.serial_number_start
	WHERE c.count > 1
)

SELECT o.reference_no, p.name, osn.serial_number_start FROM public.order o
LEFT JOIN public.order_product op
ON o.an_order_id = op.an_order_id
LEFT JOIN public.product p
ON op.an_product_id = p.an_product_id
LEFT JOIN serial_number_cte osn
ON op.an_order_product_id = osn.an_order_product_id
WHERE serial_number_start!=''
ORDER BY serial_number_start
