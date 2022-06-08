import csv
import argparse

parser = argparse.ArgumentParser(description="Convert .xlsx file to SQL query command. For update order parcel such as tracking_code and status.")

parser.add_argument("-i", "--input", help = "Input .csv file")
parser.add_argument("-o", "--output", help = "Output .sql file")

args = parser.parse_args()

if args.input == None:
    raise Exception("Input file is required.")
if args.output == None:
    raise Exception("Output file is required.")

mapStatus = {
    "รอยืนยัน": "wait",
    "เตรียมจัดส่ง": "booking",
    "ระหว่างจัดส่ง": "shipping",
    "สำเร็จ": "complete",
    "รายการโดนยกเลิก": "cancel",
    "รายการตีกลับ": "return",
}

def filterValue(value, defaultValue, valueType):
    if (value == "-"):
        return defaultValue
    return valueType(value)

outFile = open(args.output, "w+")

with open(args.input, "r", encoding="utf-8") as csvfile:
    spamreader = csv.reader(csvfile, delimiter=',')
    next(spamreader)
    for row in spamreader:
        if str(row[8]) in mapStatus:
            trackingCode = str(row[7])
            status = mapStatus[str(row[8])]
            weight = filterValue(row[9], 1, int)
            width = filterValue(row[10], 1, int)
            length = filterValue(row[11], 1, int)
            height = filterValue(row[12], 1, int)
            price =  filterValue(row[13], 0.0, float)

            print(trackingCode, status, weight, width, length, height, price)
            query = "UPDATE public.order_parcel "
            query += "SET status='{}', weight={}, width={}, length={}, height={}, price={} ".format(status, weight, width, length, height, price)
            query += "WHERE tracking_code='{}';".format(trackingCode)

            outFile.write(query + "\n")
        else:
            print(str(row[8]), "not in mapStatus.")

outFile.close()