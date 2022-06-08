import os
from os.path import join, dirname
from dotenv import load_dotenv
import sys
import argparse
import psycopg2
import pandas as pd

dotenv_path = join(dirname(__file__), '.env')
load_dotenv(dotenv_path)

if os.getenv('AN_POSTGRES_URI') == None:
    print("Not set env AN_POSTGRES_URI")
    sys.exit(0)

parser = argparse.ArgumentParser(description='Query bank account from AgentNetwork.')
parser.add_argument('i', type=str, help='COD report from Shippop eg. cod_xxxx-xx-xx.xlsx')
parser.add_argument('o', type=str, help='Outut file path eg. out_cod_xxxx-xx-xx.xlsx')

args = parser.parse_args()

if args.i == None:
    raise Exception("Input file is required.")
if args.o == None:
    raise Exception("Output file is required.")

def readTrackingCodeFromInputFile(filePath):
    trackingCodes = []
    df = pd.read_excel(filePath)
    for i in df.index:
        trackingCodes.append(df["เลขติดตามจากขนส่ง"][i])
    return trackingCodes

def queryBankAccountByTrackingCode(trackingCodes):
    _trackingCodes = []
    for i in range(len(trackingCodes)):
        _trackingCodes.append("'{}'".format(trackingCodes[i]))
    # print(','.join(_trackingCodes))
    conn = psycopg2.connect(dsn=os.getenv('AN_POSTGRES_URI'))
    cur = conn.cursor()
    cur.execute(
        "SELECT * FROM public.order o "+
        "INNER JOIN public.bank_account b "+
        "ON o.an_order_id = b.an_order_id "+
        "WHERE o.tracking_code IN ({})".format(','.join(_trackingCodes))
    )
    records = cur.fetchall()

    trackingCodeToBank = {}
    for i in range(len(records)):
        trackingCodeToBank[str(records[i][15])] = {
            'referenceNo': records[i][3],
            'bank': records[i][18 + 4],
            'accountNo': records[i][19 + 4],
            'accountName': records[i][20 + 4],
            'email': records[i][21 + 4]
        }
    
    cur.close()
    conn.close()

    return trackingCodeToBank

def insertBankAccountToFile(inputFilePath, outputFilePath, trackingCodeToBank):
    df = pd.read_excel(inputFilePath, dtype={"รหัส Tracking Code": str})
    df.insert(12, "referenceNo", "")
    df.insert(13, "bank", "")
    df.insert(14, "accountNo", "")
    df.insert(15, "accountName", "")
    df.insert(16, "email", "")

    # df['accountNo'].astype(str)
    df = df.astype(str)
    for i in df.index:
        trackingCode = df["เลขติดตามจากขนส่ง"][i]
        if trackingCode in trackingCodeToBank:
            # print(trackingCode)
            df.at[i, "referenceNo"] =  trackingCodeToBank[trackingCode]['referenceNo']
            df.at[i, "bank"] = trackingCodeToBank[trackingCode]['bank']
            df.at[i, "accountNo"] = trackingCodeToBank[trackingCode]['accountNo']
            df.at[i, "accountName"] = trackingCodeToBank[trackingCode]['accountName']
            df.at[i, "email"] = trackingCodeToBank[trackingCode]['email']
    
    df.to_excel(outputFilePath) 

inputFilePath = args.i
outputFilePath = args.o
trackingCodes = readTrackingCodeFromInputFile(inputFilePath)
print(trackingCodes)
trackingCodeToBank = queryBankAccountByTrackingCode(trackingCodes)
print(trackingCodeToBank)
insertBankAccountToFile(inputFilePath, outputFilePath, trackingCodeToBank)
print("complete")