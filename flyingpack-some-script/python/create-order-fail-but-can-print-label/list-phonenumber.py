import json

with open("sample-data.json", encoding="utf-8") as json_file:
    data = json.load(json_file)
    an_parcels = data['an_parcels']
    print(len(an_parcels))
    for i in range(len(an_parcels)):
        print("'{}',".format(an_parcels[i]['destination']['phone_number']))