import os
from os.path import join, dirname
from dotenv import load_dotenv
import sys
import argparse
import psycopg2

dotenv_path = join(dirname(__file__), '.env')
load_dotenv(dotenv_path)

if os.getenv('AN_POSTGRES_URI') == None:
    print("Not set env AN_POSTGRES_URI")
    sys.exit(0)

if os.getenv('SP_POSTGRES_URI') == None:
    print("Not set env SP_POSTGRES_URI")
    sys.exit(0)

def queryOrdersFromAN():
    orders = []
    conn = psycopg2.connect(dsn=os.getenv('AN_POSTGRES_URI'))
    cur = conn.cursor()
    query = "SELECT * FROM public.order WHERE tracking_codes = '{{}}' AND created_at >= '2021-09-20'".format()
    cur.execute(query)
    records = cur.fetchall()

    for r in records:
        orders.append({
            "anOrderId": r[0],
            "spOrderParcelId": r[2],
            "referenceNo": r[3],
            "trackingCode": r[15]
        })
    
    cur.close()
    conn.close()
    return orders

def queryOrderFromSP(orders):
    if len(orders) < 1:
        return []
    spOrderParcelIds = []
    for o in orders:
        spOrderParcelIds.append("'{}'".format(o['spOrderParcelId']))
    
    conn = psycopg2.connect(dsn=os.getenv('SP_POSTGRES_URI'))
    cur = conn.cursor()
    cur.execute("SELECT order_parcel_id, tracking_code FROM public.order_parcel WHERE order_parcel_id IN ({})".format(','.join(spOrderParcelIds)))
    records = cur.fetchall()

    orderParcelIdToTrackingCode = {}
    for r in records:
        orderParcelIdToTrackingCode[r[0]] = r[1]

    cur.close()
    conn.close()

    return orderParcelIdToTrackingCode

def updateTrackingCodeAN(orderParcelIdToTrackingCode):
    if len(orderParcelIdToTrackingCode) < 1:
        return
    conn = psycopg2.connect(dsn=os.getenv('AN_POSTGRES_URI'))
    cur = conn.cursor()
    query = ""
    for key in orderParcelIdToTrackingCode:
        query += "UPDATE public.order SET tracking_codes[1]='{}' WHERE sp_order_parcel_id='{}';".format(orderParcelIdToTrackingCode[key], key)
    
    cur.execute(query)
    conn.commit()
    
    cur.close()
    conn.close()

orders = queryOrdersFromAN()
print("OK 1 {}".format(len(orders)))

orderParcelIdToTrackingCode = queryOrderFromSP(orders)
print("OK 2 {}".format(len(orderParcelIdToTrackingCode)))
okCount = 0
for key in orderParcelIdToTrackingCode:
    if orderParcelIdToTrackingCode[key] != "":
        okCount += 1
print("{}, {}".format(okCount, len(orderParcelIdToTrackingCode) - okCount))

updateTrackingCodeAN(orderParcelIdToTrackingCode)
print("\rComplete")