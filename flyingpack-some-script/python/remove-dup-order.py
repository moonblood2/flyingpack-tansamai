import pandas as pd
import numpy as np

sheet1 = pd.read_excel('Upload_GELสินค้าหมด202103_29_31.xlsx', sheet_name=0, dtype={'เบอร์โทร': str})
sheet1_tmp = sheet1.copy()
sheet2 = pd.read_excel('Upload_GELสินค้าหมด202103_29_31.xlsx', sheet_name=1, dtype={'เบอร์โทร': str})
sheet3 = pd.read_excel('Upload_GELสินค้าหมด202103_29_31.xlsx', sheet_name=2)
print(sheet3)

map_phonnumber = {}
map_referenceNo = {}

for i in sheet2.index:
    map_phonnumber['0' + sheet2['เบอร์โทร'][i]] = True
for i in sheet3.index:
    map_referenceNo[sheet3['referenceNo'][i]] = True

# print(map_phonnumber)
# print(map_referenceNo)

count = 0

for i in sheet1.index:
    if sheet1['เบอร์โทร'][i] in map_phonnumber or sheet1['referenceNo'][i] in map_referenceNo:
        sheet1_tmp = sheet1_tmp.drop(i)
        count += 1

with pd.ExcelWriter('output.xlsx', mode="w") as writer:  
    sheet1_tmp.to_excel(writer, sheet_name='Sheet1')
print(count)