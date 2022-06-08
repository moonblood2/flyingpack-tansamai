import requests

in_f = open("Tracking_Oct2_Oct3.txt", "r")
out_f = open("out.sql", "w+")

in_lines = in_f.readlines()

def gen_query():
    orders = []
    out_lines = []

    for l in in_lines:
        x = l.strip().split("\t")
        out_lines.append(
            "UPDATE public.order SET tracking_codes[1]='{}', fulfillment_status=1 WHERE reference_no='{}';\n".format(x[1], x[0])
        )
        orders.append({
            "referenceNo": x[0],
            "trackingCode": x[1],
        })
    return orders, out_lines


# print(orders)

def hook(orders):
    #Hook to SMITH
    url = "https://rosegoldthailand.com/api/v1/delivery"
    i = 421
    while i < len(orders):
        headers = {
            "api-key": "HTTP_API_KEY", 
            "Content-Type": "application/json",
        }
        payload = {
            "data": {
                "referenceNo": orders[i]["referenceNo"],
                "trackingCode": orders[i]["trackingCode"],
            },
        }
        res = requests.put(url=url, headers=headers, json=payload)
        print(i)
        print(res.json())
        i += 1

# hook()
orders, out_lines = gen_query()

#Write SQL to file
out_f.writelines(out_lines)

in_f.close()
out_f.close()