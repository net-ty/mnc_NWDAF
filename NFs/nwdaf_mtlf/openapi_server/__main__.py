#!/usr/bin/env python3

import connexion

from openapi_server import encoder

# from multiprocessing import Process
import subprocess


def main():
    app = connexion.App(__name__, specification_dir="./openapi/")
    app.app.json_encoder = encoder.JSONEncoder
    app.add_api("openapi.yaml", arguments={"title": "Nnwdaf_MLModelProvision"}, pythonic_params=True)

    cmd = "/usr/local/go/bin/go run nwdaf.go"
    subprocess.run([cmd], shell=True)
    app.run(port=8081)


if __name__ == "__main__":
    main()
