import csv

fields = ['reference_no', 'tracking_code_1', 'tracking_code_2']
rows = []
jnaTrackingCodes = {}

with open('jna-trackingcode-cancel-order-2021-03-05.csv', newline='') as csvfile:
    spamreader = csv.reader(csvfile, delimiter=',', quotechar='|')
    next(spamreader)
    for row in spamreader:
        jnaTrackingCodes[row[0]] = row[1]

with open('duplicate.csv', newline='', encoding='utf-8') as csvfile:
    spamreader = csv.reader(csvfile, delimiter=',', quotechar='|')
    next(spamreader)
    for row in spamreader:
        referenceNo = row[5]
        trackingCode = row[6]
        rows.append([referenceNo, trackingCode, jnaTrackingCodes[referenceNo]])
        print([referenceNo, trackingCode, jnaTrackingCodes[referenceNo]])

with open("merge-tracking-code-2021-03-05.csv", 'w', newline='') as csvfile:  
    csvwriter = csv.writer(csvfile)  
    csvwriter.writerow(fields)  
    csvwriter.writerows(rows)

print("Complete")