import os
from os.path import join, dirname
from dotenv import load_dotenv
import sys
import argparse
import psycopg2
import requests
import json

# Create .env file path. Load file from the path. Return list of spOrderParcelId.
dotenv_path = join(dirname(__file__), '.env')
load_dotenv(dotenv_path)

if os.getenv('AN_POSTGRES_URI') == None:
    print("Not set env AN_POSTGRES_URI")
    sys.exit(0)

if os.getenv('SP_POSTGRES_URI') == None:
    print("Not set env SP_POSTGRES_URI")
    sys.exit(0)

if os.getenv('RG_API_KEY') == None:
    print("Not set env RG_API_KEY")
    sys.exit(0)

# Args
parser = argparse.ArgumentParser(description='Update missing tracking code in Agent Network service.')
parser.add_argument('start', type=str,
                    help='Starat date in format "2021-03-17"')
parser.add_argument('end', type=str,
                    help='Starat date in format "2021-03-18"')

args = parser.parse_args()

# Find no tracking order in AgentNetwork DB.
def findNoHasTracking(startDate, endDate):
    conn = psycopg2.connect(dsn=os.getenv('AN_POSTGRES_URI'))
    cur = conn.cursor()
    cur.execute("SELECT * FROM public.order WHERE sp_order_parcel_id IS NOT NULL AND (tracking_code='' OR tracking_code IS NULL) AND created_at AT TIME ZONE 'asia/bangkok' BETWEEN '{}' AND '{}'".format(
        startDate, endDate
    ))
    records = cur.fetchall()

    orders = []
    for i in range(len(records)):
        orders.append({
            'anOrderId': records[i][0],
            'userId': records[i][1],
            'spOrderParcelId': records[i][2],
            'referenceNo': records[i][3]
        })
    
    cur.close()
    conn.close()

    return orders

def findOrders(startDate, endDate):
    conn = psycopg2.connect(dsn=os.getenv('AN_POSTGRES_URI'))
    cur = conn.cursor()
    cur.execute("SELECT * FROM public.order WHERE fulfillment_status=1 AND created_at AT TIME ZONE 'asia/bangkok' BETWEEN '{}' AND '{}'".format(
        startDate, endDate
    ))
    records = cur.fetchall()

    orders = []
    for i in range(len(records)):
        # print(records[i])
        orders.append({
            'anOrderId': records[i][0],
            'userId': records[i][1],
            'spOrderParcelId': records[i][2],
            'referenceNo': records[i][3],
            'trackingCode': records[i][15],
            'courierTrackingCode': records[i][15],
        })
    
    cur.close()
    conn.close()

    return orders


# Find Shippop's TrackingCode in Shipping DB. Return list of object {purchaseId, trackingCode}
def findShippopTrackingCode(orders):
    mapOrderParcelIdToOrderIndex = {}
    orderParcelIdsStr = []

    for i in range(len(orders)):
        mapOrderParcelIdToOrderIndex[orders[i]['spOrderParcelId']] = i
        orderParcelIdsStr.append("'{}'".format(orders[i]['spOrderParcelId']))
    
    conn = psycopg2.connect(dsn=os.getenv('SP_POSTGRES_URI'))
    cur = conn.cursor()
    cur.execute("SELECT * FROM public.order_parcel_shippop WHERE order_parcel_id IN ({})".format(
        ', '.join(orderParcelIdsStr)
    ))
    records = cur.fetchall()

    # populate trackingCode to order
    for i in range(len(records)):
        # 1=order_parcel_id, 2=purchase_id, 6=tracking_code
        orders[mapOrderParcelIdToOrderIndex[records[i][1]]]['purchaseId'] = records[i][2]
        orders[mapOrderParcelIdToOrderIndex[records[i][1]]]['trackingCode'] = records[i][6]
    
    cur.close()
    conn.close()

    return orders

def findMissingCourierTracingCode():
    orders = []

    conn = psycopg2.connect(dsn=os.getenv('SP_POSTGRES_URI'))
    cur = conn.cursor()
    cur.execute("SELECT * FROM public.order_parcel_shippop WHERE courier_tracking_code='' AND deleted_at=0")
    records = cur.fetchall()

    # populate trackingCode to order
    for i in range(len(records)):
        # 1=order_parcel_id, 2=purchase_id, 6=tracking_code
        orders.append({
            'spOrderParcelId': records[i][1],
            'purchaseId': records[i][2],
            'trackingCode': records[i][6]
        })

    cur.close()
    conn.close()

    return orders

def getCourierTrackingCodeFromSpDB(orders):
    conn = psycopg2.connect(dsn=os.getenv('SP_POSTGRES_URI'))
    cur = conn.cursor()

    for i in range(len(orders)):
        cur.execute("SELECT * FROM public.order_parcel WHERE order_parcel_id='{}'".format(
            orders[i]['spOrderParcelId']
        ))
        trackingCode = cur.fetchone()[13]
        orders[i]['courierTrackingCode'] = trackingCode
        print(i)
    cur.close()
    conn.close()
    
    return orders

#getCourierTrackingCode get courier tracking code from Shoppop's Label API.
def getCourierTrackingCodeFromShippop(orders):
    url = "http://mkpservice.shippop.com/label/"

    for i in range(len(orders)):
        payload = {
            "api_key": "e7672f8b5c520d05c71cd4e0cc6d4de1cb9d429d",
            "purchase_id": orders[i]['purchaseId'],
            "tracking_code": orders[i]['trackingCode'],
            "type": "json"
        }

        res = requests.post(url=url, json=payload)
        label = res.json()['json']['labels'][0]
        # populate courierTrackingCode to order
        orders[i]['courierTrackingCode'] = label['courierTrackingCode']

    return orders

# updateSpTrackingCode update tracking code in Shpping DB.
def updateSpTrackingCodeOrderParcel(orders):
    conn = psycopg2.connect(dsn=os.getenv('SP_POSTGRES_URI'))
    cur = conn.cursor()

    for i in range(len(orders)):       
        # updatae tracking_code in order_parcel   
        print(orders[i]['courierTrackingCode'] + ", " + orders[i]['spOrderParcelId'])   
        cur.execute("UPDATE public.order_parcel SET tracking_code='{}' WHERE order_parcel_id='{}'".format(
            orders[i]['courierTrackingCode'],
            orders[i]['spOrderParcelId']
        ))

        conn.commit()

    cur.close()
    conn.close()

# updateSpTrackingCode2 update tracking code in Shpping DB, order_parcel_shippop.
def updateSpTrackingCodeOrderParcelShippop(orders):
    conn = psycopg2.connect(dsn=os.getenv('SP_POSTGRES_URI'))
    cur = conn.cursor()

    for i in range(len(orders)):       
        # updatae tracking_code in order_parcel  
        print(orders[i]['courierTrackingCode'] + ", " + orders[i]['spOrderParcelId'])   
        cur.execute("UPDATE public.order_parcel_shippop SET courier_tracking_code='{}' WHERE order_parcel_id='{}'".format(
            orders[i]['courierTrackingCode'],
            orders[i]['spOrderParcelId']
        ))
        conn.commit()

    cur.close()
    conn.close()

def updateAnOrderTrackingCode(orders):
    conn = psycopg2.connect(dsn=os.getenv('AN_POSTGRES_URI'))
    cur = conn.cursor()

    for i in range(len(orders)):       
        cur.execute("UPDATE public.order SET tracking_code='{}' WHERE an_order_id='{}'".format(
            orders[i]['courierTrackingCode'],
            orders[i]['anOrderId']
        ))
        conn.commit()
        print(i)
    cur.close()
    conn.close()

# hookAnTrackingCode hook fullfilment to AgentNetwork service for update trakcingCode.
def hookAnTrackingCode(orders):
    url = "https://service-agent-network.herokuapp.com/closed/hook/order-status"

    for i in range(len(orders)):
        payload = {
            "sp_order_parcel_id": orders[i]["spOrderParcelId"],
            #"status": "shipping",
            "tracking_code": orders[i]["courierTrackingCode"]
        }
        res = requests.post(url=url, json=payload)
        print(i)
        print(res.json())

# printOrders
def printOrders(orders):
    for i in range(len(orders)):
        print(json.dumps(orders[i], indent = 1))

#Define 
startDate = args.start
endDate = args.end

'''
orders = findNoHasTracking(startDate, endDate)
print("findNoHasTracking: OK")
orders = findShippopTrackingCode(orders)
print("findShippopTrackingCode: OK")
orders = getCourierTrackingCodeFromShippop(orders)
print("getCourierTrackingCodeFromShippop: OK")
updateSpTrackingCodeOrderParcelShippop(orders)
print("updateSpTrackingCodeOrderParcelShippop: OK")
updateAnOrderTrackingCode(orders)
print("updateAnOrderTrackingCode: OK")
hookAnTrackingCode(orders)
print("hook: OK")
print("There are {} orders.".format(len(orders)))
printOrders(orders)
print("Completed.")
'''

'''
updateSpTrackingCode(orders)
print("updateSpTrackingCode: OK")
updateSpTrackingCode2(orders)
print("updateSpTrackingCode2: OK")
hookAnTrackingCode(orders)
print("hookAnTrackingCode: OK")
'''

'''
orders = findMissingCourierTracingCode()
print("findMissingCourierTracingCode: OK")
orders = getCourierTrackingCode(orders)
print("getCourierTrackingCode: OK")
updateSpTrackingCode(orders)
print("updateSpTrackingCode: OK")
updateSpTrackingCode2(orders)
print("updateSpTrackingCode2: OK")
'''

'''
orders = findNoHasTracking(startDate, endDate)
print("findNoHasTracking: OK")
orders = getCourierTrackingCodeFromSpDB(orders)
print("getCourierTrackingCodeFromSpDB: OK")
updateAnOrderTrackingCode(orders)
print("updateAnOrderTrackingCode: OK")
hookAnTrackingCode(orders)
print("hook: OK")
'''

orders = findOrders('2021-05-14 00:00', '2021-05-17 18:00')
print("findOrders: OK")
hookAnTrackingCode(orders)
print("hook: OK")

printOrders(orders)
print("There are {} orders.".format(len(orders)))
print("Completed.")