update public.order_parcel_shippop osp
set status=op.status, created_at=op.created_at
from public.order_parcel op
where osp.order_parcel_id = op.order_parcel_id
-- and op.created_at at time zone 'asia/bangkok' between '2021-03-08' and '2021-03-15'
-- and osp.status is null

