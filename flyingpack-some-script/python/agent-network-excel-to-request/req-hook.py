import json
import requests

f = open("./2021-03-01/status.json", 'r', encoding='utf-8')
records = json.load(f) 
count = 0
for data in records:
    count += 1
    print(count)
    print(data)
    url = 'https://rosegoldthailand.com/api/v1/delivery'
    headers = {
        'content-type': 'application/json',
        'api-key': 'HTTP_API_KEY'
    }
    x = requests.put(url, json=data, headers=headers)
    print(x.text)
    print("-----------------------------------------------------------")
f.close()

print("Completed")