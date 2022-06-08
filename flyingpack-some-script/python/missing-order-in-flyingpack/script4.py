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
spOrderParcelId = [
"244ae3ea-b9cd-40e7-b7fc-fd7ba4247e5c",
"cd1b5d40-c250-447e-9758-d015fd4c10a0",
"27186077-9bc2-4c6f-a53f-07b1a0269bc0",
"b5864d5e-a27e-4937-b412-f14fa01d453d",
"a71eb57b-8623-43d6-b256-fa22ed2acde2",
"c0e840c3-65fc-4150-a799-914176185eee",
"30ec904a-9eed-471f-877a-d0ad228cbdd4",
"8ad4040a-b031-482f-9828-7b5d0de70a91",
"2d15a82a-62cb-4f05-b671-ff16d00d4efa",
"3df93470-dcbf-4cd3-97be-10643588ead3",
"a746b5d2-07be-44ab-aeb9-5bb03b1f7de1",
"6412e81d-91f2-4281-820f-4ccd9783b839",
"13297ef0-c535-4392-aea2-c702bd8bdf7f",
"65fd35f0-1a3a-4714-81bb-bc792c9acff7",
"561da5b8-72db-4c43-935e-c94c48585e9e",
"3dda4084-65bb-43e1-8483-f04f5ad6b97c",
"138d3372-a945-4883-b2ba-a1036a6cffd5",
"f2902a05-07e9-47af-963a-c3c281ef9d93",
"c16d39c8-3160-47b3-b27e-2b7f7223ca9f",
"6a677607-11eb-4f4b-b90f-cb4e6b036ea4"
]
anOrderId = [
"1e8cc4e7-b316-4567-ac96-fcc8a048ff6f",
"682ea6ec-1e5e-4839-bd47-c120e6c97977",
"af342176-da61-4bd3-b532-ca0551b0fcbb",
"5cea25e4-a83b-4d19-b6a8-83be92c98663",
"03c1d1da-6f49-41a6-a7af-d0da82ee2f82",
"842b130c-b110-4c4a-aecd-dcdb723de30b",
"5a0756fb-82ca-47e2-9ab4-fb3b5d36416b",
"0cf44b9f-7e4e-4f70-8fc2-c2d46e76f5d1",
"f25b5a39-1a50-4968-b7db-6025a67a7f34",
"913f193a-5044-484e-9975-b53e38e0cf41",
"a9ab3a97-cb93-41ea-9934-fc91f7b75d6f",
"421e801d-7778-4b9b-895a-570acd587a34",
"8f77b826-cf5b-4984-9de3-5f3e32014224",
"5b070e04-1b41-41c3-8fe7-a089e72d9e8f",
"9000fde4-231d-4401-97d9-70155181db37",
"3b23d0a8-8cc9-40c6-8c6f-60191209f949",
"52eeaff3-d327-4c92-b1f1-db3817def0d9",
"ab09a4f6-7efa-43b0-bb66-db573911da6b",
"cc9946a7-364a-4c24-8297-e1ef2bdd2822",
"dfb9eb25-9e86-47dc-9bcf-2acd266ce208"
]

file = open("an_update_order.sql", "w+", encoding="utf-8")

for i in range(len(spOrderParcelId)):
    file.write("UPDATE public.order SET sp_order_parcel_id='{}', tracking_code='{}' WHERE an_order_id='{}';\n".format(
        spOrderParcelId[i], d[i]["courierTrackingCode"], anOrderId[i]
    ))


file.close()