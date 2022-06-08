import pandas as pd
import json
import math
import argparse

parser = argparse.ArgumentParser(description="Convert excel to JSON payload.")

parser.add_argument("-i", "--input", help = "Input .xlsx file")
parser.add_argument("-o", "--output", help = "Output .json file")

args = parser.parse_args()

if args.input == None:
    raise Exception("Input file is required.")
if args.output == None:
    raise Exception("Output file is required.")

df = pd.read_excel(args.input, dtype={"เบอร์ผู้รับ": str, "COD": int})

productCodes = ["CREAM", "GEL", "SAKANA", "ALC-FREE", "SPRAY20ml", "ALC_GEL1", "ALC_GEL2", "Spray50ml", "SAKANA_N", "CREAM_N", "GEL_N", "ACNE", "CREAM15ml"]

records = []
rowCount = 0

for i in df.index:
    rowCount += 1
    print(i + 1)
    data = {}
    #apiKey
    data["apiKey"] = "bKMf1t7pqvCnsF7hXy1W3ZUqSg6ZaNkRqc65gUlsmA5XAUqiSzqAl9dwuWdSLf1gmkIgn18ajzf6JJu4SUVwkKyj3KFS10TCbSfjwjQR9HuLppMVqsr1Ju65EasPM3VH"
    # data["apiKey"] = "CTSvnpNEMJElQfe5UwhdEGdXSXuoz0ANCOw6WNKU098S11y4L9qo9QeaTZUU62KImKsEzteVxf4UI4pFfyYbiz7b53XFJiILDT8a8Sb9724cxb8pBalYyyopjsgmqaFe"
    #referenceNo
    data["referenceNo"] = df["เลขที่การขาย"][i]
    #desName
    data["desName"] = df["ชื่อผู้รับ(ลูกค้า)"][i].strip()
    #desPhoneNumber
    data["desPhoneNumber"] = df["เบอร์ผู้รับ"][i].strip()

    #Address    
    addressFull = df["ที่อยู่ผู้รับ"][i]
    addressArr = addressFull.split()
    addressArrLen = len(addressArr)
    
    desProvince = addressArr[addressArrLen - 1].strip().replace("จ.", "")
    desDistrict = addressArr[addressArrLen - 2].strip().replace("ต.", "").replace("แขวง", "")
    desSubdistrict = addressArr[addressArrLen - 3].strip().replace("อ.", "")
    desAddress = '  '.join(addressArr[:addressArrLen - 3])

    data["desAddress"] = desAddress
    data["desSubdistrict"] = desSubdistrict
    data["desDistrict"] = desDistrict
    data["desProvince"] = desProvince
    data["desPostcode"] = str(df["รหัสไปรษณีย์ผู้รับ"][i])

    #Courier Code
    courierName = df["รูปแบบขนส่ง"][i]
    courierCode = 0
    if ("FLASH" in courierName) or ("Flash" in courierName):
        courierCode = 1
    elif "KERRY" in courierName:
        courierCode = 2
    elif "EMS" in courierName:
        courierCode = 4
    else:
        raise Exception('Something went wrong with "courier_code"')

    data["courierCode"] = courierCode

    #COD
    data["codAmount"] = int(df["COD"][i])

    #Items
    data["items"] = []
    for productCode in productCodes:
        if math.isnan(df[productCode][i]) == False:
            data["items"].append(
                {
                    "productCode": productCode,
                    "quantity": int(df[productCode][i])
                }
            )

    records.append(data)

with open(args.output, 'w', encoding='utf-8') as outfile:
    json.dump(records, outfile, ensure_ascii=False)

print("Completed, " + str(rowCount) + " rows")