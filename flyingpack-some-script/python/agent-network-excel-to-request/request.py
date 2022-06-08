import argparse
import json 
import requests
from colorama import Fore, Back, Style 

parser = argparse.ArgumentParser(description="Send http request to AgentNetwork service, create order.")

parser.add_argument("-i", "--input", help = "Input .json file")

args = parser.parse_args()

if args.input == None:
    raise Exception("Input file is required.")

f = open(args.input, 'r', encoding='utf-8')
records = json.load(f) 

recordCount = 0
errRecords = []

url = 'https://service-agent-network.herokuapp.com/open-api/order'
# url = 'http://localhost:3000/open-api/order'

for data in records:
    recordCount += 1
    print("row: ", recordCount)
    print("data: ", data)
    x = requests.post(url, json = data)
    response = json.loads(x.text)
    if response["code"] == 1:
        print(Fore.GREEN + json.dumps(response, indent=4))
    else:
        print(Fore.RED + json.dumps(response, indent=4))
        errRecords.append(recordCount)
    print(Style.RESET_ALL)
    print("-----------------------------------------------------------")
f.close()

print("There are " + str(len(errRecords)) + " error records.")
print("Error records: ", errRecords)
print("Completed")