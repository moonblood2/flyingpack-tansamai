import json
import csv

mapPhoneNumberAnParcelIndex = {}

json_file = open("sample-data.json", encoding="utf-8")
data = json.load(json_file)
an_parcels = data['an_parcels']
# print(len(an_parcels))
for i in range(len(an_parcels)):
    mapPhoneNumberAnParcelIndex[an_parcels[i]['destination']['phone_number']] = i
    # print("'{}',".format(an_parcels[i]['destination']['phone_number']))

print(mapPhoneNumberAnParcelIndex)

csvfile = open("destinations.csv", encoding="utf-8")
spamreader = csv.reader(csvfile, delimiter=',')
for row in spamreader:
    # print(row)
    an_parcels[mapPhoneNumberAnParcelIndex[row[2]]]['destination'] = {
        "name": row[1],
        "address": row[3],
        "district": row[4],
        "state": row[5],
        "province": row[6],
        "postcode": row[7],
        "phone_number": row[2]
    }

json_file.close()
csvfile.close()

with open('sample-data-new.json', 'w', encoding="utf-8") as outfile:
    json.dump(data, outfile, ensure_ascii=False)