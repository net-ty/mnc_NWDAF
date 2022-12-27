import gdown
import os
import requests
import json
from AnLF.model_inference1 import *



def analytic_result() :

    if os.path.isfile("./AnLF/models/model.h5") :

        res_value = inference1()
        print("Inference completed")
        return (res_value)


    else: 
        root_nwdaf_server = f"http://192.168.221.130:8081/nnwdaf-mlmodelprovision/v1/subscriptions/sbd"
        req_body = {
            "notifUri": "http://192.168.221.130:8081/nnwdaf-mlmodelprovision/v1/subscriptions/sbd",
            "mLEventSubscs": [
                {
                    "mLEvent": "ML_model_provision_request_drop",       
                    "mLEventFilter": {
                        "not": {
                            "required": [
                                "anySlice",
                                "snssais"
                            ]
                        }
                    }
                }
            ]
        }
        headers = {'Content-type': 'application/json; charset=utf-8'}

        r = requests.put(root_nwdaf_server, headers = headers, data = json.dumps(req_body))
        res_body = json.loads(r.text)   

        MLmodel_url = res_body['m_l_file_addr'] # m_l_file_addr is in the response body)
        #event = res_body['event']
        #MLmodel_url = "https://drive.google.com/drive/folders/1zdpyXbnptSU4Yb0KdtJOrAUoXb3DN3EU?usp=sharing"

        fileaddr = format(MLmodel_url.get('m_l_model_url'))

#        output = "./AnLF/models"
#        gdown.download(fileaddr, output, quiet=False)

        gdown.download_folder(fileaddr, quiet=True, use_cookies=False, output="./AnLF/models")


        if os.path.isfile("./AnLF/models/model.h5") :

            res_value = inference1()
            print("Inference with provisioned ML model")
            return (res_value)