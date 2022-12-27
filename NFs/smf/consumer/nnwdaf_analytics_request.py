import requests


nwdaf_server = f"http://192.168.221.130:8080/nnwdaf-analyticsinfo/v1/analytics"
query_string = {
    'event-id': "REDUNDANT_TRANSMISSION",
    "ana-req":"",
    "event-filter":"",
    "tgt-ue":"imsi-208930000000005"
}
print("Redundant transmission experience related analytics requested")
r = requests.get(nwdaf_server, params=query_string)

print(r.status_code)
print(r.text)

