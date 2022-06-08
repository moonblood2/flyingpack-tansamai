from functools import cache
import os
from os.path import join, dirname
from dotenv import load_dotenv
import sys
import psycopg2
import requests
import traceback



dotenv_path = join(dirname(__file__), '.env')
load_dotenv(dotenv_path)

if sys.argv[1] == "":
    print("Not pass date")
    sys.exit(0)

if sys.argv[2] == "":
    print("Not pass anorder_id_file")
    sys.exit(0)

if sys.argv[3] == "":
    print("Not pass hook_log_file")
    sys.exit(0)

if os.getenv('AN_POSTGRES_URI') == None:
    print("Not set env AN_POSTGRES_URI")
    sys.exit(0)

if os.getenv('TOKEN') == None:
    print("Not set env TOKEN")
    sys.exit(0)

AN_POSTGRES_URI = os.getenv('AN_POSTGRES_URI')
TOKEN = os.getenv('TOKEN')
DATE = sys.argv[1]
AN_ORDER_ID_FILE = sys.argv[2]
HOOK_LOG_FILE = sys.argv[3]

def queryOrdersFromAN():
    print("queryOrdersFromAN: ", end="")

    orders = []
    if os.path.exists(AN_ORDER_ID_FILE):
        f = open(AN_ORDER_ID_FILE, "r")
        lines = f.readlines()
        f.close()
        for l in lines:
            orders.append({
                "anOrderId": l.strip("\n")
            })
    else:
        out = open(AN_ORDER_ID_FILE, "w+")
        
        conn = psycopg2.connect(dsn=AN_POSTGRES_URI)
        cur = conn.cursor()
        query = '''
            SELECT DISTINCT o.an_order_id FROM public."order" o
            LEFT JOIN "public".order_product op
            ON o.an_order_id=op.an_order_id
            LEFT JOIN "public".product p
            ON op.an_product_id=p.an_product_id
            LEFT JOIN "public".an_order_product_serial_number aops
            on op.an_order_product_id = aops.an_order_product_id
            WHERE o.created_at >= '{}';
        '''.format(DATE)
        cur.execute(query)
        records = cur.fetchall()

        for r in records:
            orders.append({
                "anOrderId": r[0],
            })
            out.write("{}\n".format(r[0]))
        
        cur.close()
        conn.close()
        out.close()

    print("OK")
    return orders

def sendSerialNumbers(orders):
    out = open(HOOK_LOG_FILE, "w+")
    url = "https://service-agent-network.herokuapp.com/closed-api/order/order-product/serial-number"
    headers = {
        'Authorization': "Bearer {}".format(TOKEN)
    }
    try:
        for i, o in enumerate(orders):
            data = {
                "anOrderId": o['anOrderId'],
                "serialNumbers": [],
            }
            res = requests.put(url=url, headers=headers, json=data)
            out.write("{} {}\n".format(res.status_code, res.json()))
            print("\rsendSerialNumbers: {:.2f}% ({}/{})".format(((i+1)/len(orders) * 100), i+1, len(orders)), end="")
            
            if i % 10 == 0:
                out.flush()
    except Exception:
        traceback.print_exc()
    finally:
        out.close()
    print("\nsendSerialNumbers: OK\n")

orders = queryOrdersFromAN()
print("len(orders): ", len(orders))
sendSerialNumbers(orders)