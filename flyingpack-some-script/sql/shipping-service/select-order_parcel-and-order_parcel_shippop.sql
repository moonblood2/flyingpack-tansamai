select * from public.order_parcel op
left join public.order_parcel_shippop ops
on op.order_parcel_id = ops.order_parcel_id