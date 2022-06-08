d = [
    {
        'trackingCode': 'SP177384685',
        'courierTrackingCode': 'TH190713STUH4A'
    },
    {
        "trackingCode": "SP177384676",
        "courierTrackingCode": "TH390113STUG9A",
    },
    {
        "trackingCode": "SP177384661",
        "courierTrackingCode": "TH310513STUG8A",

    },
    {
        "trackingCode": "SP177384746",
        "courierTrackingCode": "TH550113STUJ5A",

    },
    {
        "trackingCode": "SP177384731",
        "courierTrackingCode": "TH011713STUJ4B",

    },
    {
        "trackingCode": "SP177384727",
        "courierTrackingCode": "TH010313STUJ2C",
    },
    {
        "trackingCode": "SP177384712",
        "courierTrackingCode": "TH014413STUJ1A",
    },
    {
        "trackingCode": "SP177384708",
        "courierTrackingCode": "TH120413STUH8C",
    },
    {
        "trackingCode": "SP177384695",
        "courierTrackingCode": "TH200713STUH7D",
    },
    {
        "trackingCode": "SP177384638",
        "courierTrackingCode": "TH020113STUG2B",
    },
    {
        "trackingCode": "SP177384825",
        "courierTrackingCode": "TH014113STUN4C",
    },
    {
        "trackingCode": "SP177384816",
        "courierTrackingCode": "TH700813STUN1D",
    },
    {
        "trackingCode": "SP177384801",
        "courierTrackingCode": "TH681413STUM6A",
    },
    {
        "trackingCode": "SP177384799",
        "courierTrackingCode": "TH680113STUM0J",
    },
    {
        "trackingCode": "SP177384657",
        "courierTrackingCode": "TH013813STUG6A",
    },
    {
        "trackingCode": "SP177384784",
        "courierTrackingCode": "TH680113STUK9B",
    },
    {
        "trackingCode": "SP177384770",
        "courierTrackingCode": "TH471113STUK4A"
    },
    {
        "trackingCode": "SP177384765",
        "courierTrackingCode": "TH011813STUK2A"
    },
    {
        "trackingCode": "SP177384755",
        "courierTrackingCode": "TH580113STUJ6P"
    },
    {
        "trackingCode": "SP177384642",
        "courierTrackingCode": "TH012513STUG4D"
    }
]

count = 0

file = open("insert_order_parcel.sql", "w+", encoding = "utf-8")


for i in range(1050, 1069 + 1):
    file.write("INSERT INTO public.order_parcel(shop_id, sender_id, origin_id, destination_id, provider_code, price, payment_method, weight, width, length, height, tracking_code, status, cod_amount) VALUES ({}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, '{}', '{}', {});\n".format(
        2, 5, 5, i, 1, 0, 1, 1, 1, 1, 1, d[count]['courierTrackingCode'], 'booking', 0
    ))
    count += 1

file.close()