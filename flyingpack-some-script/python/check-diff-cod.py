import pandas as pd

mapTrackingCode = {}

df2 = pd.read_excel("cod_2021-05-10_SHP.xlsx")
df1 = pd.read_excel("cod_2021-05-10_FLY.xlsx")

for i in df1.index:
    mapTrackingCode[df1["tracking code"][i]] = True

for i in df2.index:
    if df2["เลขติดตามจากขนส่ง"][i] not in mapTrackingCode:
        print(df2["เลขติดตามจากขนส่ง"][i])

    print(df2["ยอดเก็บปลายทาง"][i])
    if df2["ยอดเก็บปลายทาง"][i] == 1:
        print("COD = 1")
        print(df2["เลขติดตามจากขนส่ง"][i])


# print(mapTrackingCode)


