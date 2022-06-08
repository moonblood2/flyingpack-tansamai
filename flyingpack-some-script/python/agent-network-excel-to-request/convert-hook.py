import csv
import json


records = []
count = 0
with open('./2021-03-01/rosegold-trackingcode-2021-03-01.csv', newline='') as csvfile:
    spamreader = csv.reader(csvfile, delimiter=',', quotechar='|')
    for row in spamreader:
        if "SHIPOP" in row[1]:
            data = {
                "data": {
                    "referenceNo": row[0],
                    "status": "complete"
                }
            }
            records.append(data)
        count += 1
        print(count)

with open("./2021-03-01/status.json", 'w', encoding='utf-8') as outfile:
    json.dump(records, outfile, ensure_ascii=False)

print("Complete " + str(count) + " rows")
