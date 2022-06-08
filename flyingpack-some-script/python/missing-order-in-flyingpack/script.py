import pandas as pd

df = pd.read_excel('address.xlsx', dtype={"phone_number": str, "postcode": str})

file = open("insert.sql","w+", encoding='utf-8') 

for i in df.index:
    print(df['address'][i])
    file.write("INSERT INTO public.destination(name, phone_number, address, district, state, province, postcode) VALUES ('{}', '{}', '{}', '{}', '{}', '{}', '{}');\n".format(
        df['name'][i].strip(), df['phone_number'][i].strip(), df['address'][i].strip(), df['district'][i].strip(), df['state'][i].strip(), df['province'][i].strip(), df['postcode'][i].strip()
    ))


file.close()