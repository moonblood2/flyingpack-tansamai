import argparse
import json
import pandas as pd
from colorama import Fore, Back, Style 

parser = argparse.ArgumentParser(description="Compare json file and xslx file, recheck some field")

parser.add_argument("-j", "--json", help = "Input .json file")
parser.add_argument("-x", "--xslx", help = "Input .xslx file")

args = parser.parse_args()

if args.json == None:
    raise Exception("Input .json file is required.")
if args.xslx == None:
    raise Exception("Input .xslx file is required.")

f = open(args.json, 'r', encoding='utf-8')
records = json.load(f)
df = pd.read_excel(args.xslx)

rowCount = 0
errRecords = []

for data in records:
    rowCount += 1
    haveError = False
    if data["desAddress"] == "":
        haveError = True
        print("[desAddress] empty", data["desAddress"])
    if len(data["desPhoneNumber"]) > 10:
        haveError = True
        print("[desPhoneNumber] length more than 10 ditigs", data["desPhoneNumber"])

    if haveError == True:
        json_formatted_str = json.dumps(data, indent=4, ensure_ascii=False)
        print(Style.RESET_ALL)
        print("row: ", rowCount)
        print(Fore.RED + json_formatted_str)
        errRecords.append(rowCount)
        print(Style.RESET_ALL)
    
f.close()

print("There are " + str(len(errRecords)) + " error records.")
print("Error records: ", errRecords)
print("Completed")