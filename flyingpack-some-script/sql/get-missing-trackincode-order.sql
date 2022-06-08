SELECT * 
FROM public.order 
WHERE tracking_code='' 
AND sp_order_parcel_id IS NOT NULL 
AND created_at AT TIME ZONE 'Asia/Bangkok' BETWEEN '2021-03-17' AND '2021-03-18'