import csv
import requests

url = 'https://rosegoldthailand.com/api/v1/delivery'
headers = {
    'api-key': 'HTTP_API_KEY'
}

count = 1

with open('refno-tracking.csv', encoding='utf-8') as csvfile:
    spamreader = csv.reader(csvfile, delimiter=',')
    next(spamreader)
    for row in spamreader:
        print(count)
        payload = {
            "data": {
                    "referenceNo": row[1],
                    "status": "fulfillment",
                    "trackingCode": row[2],
                    "items": [
                        {
                            "serialNumbers": []
                        }
                    ]
                }
        }
        print(payload)
        x = requests.put(url, json=payload, headers=headers)
        print(x.text)
        count += 1
print("Complete")
