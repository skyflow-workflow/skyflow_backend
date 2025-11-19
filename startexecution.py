#！/usr/bin/env python3
import sys
import argparse
import requests
import urllib
import json

default_host = "http://localhost:8080/"

def read_file(file_path):
    with open(file_path, 'r') as file:
        data = file.read()
    return data


def main(statemachine_definition:str, statemachine_uri:str, inputdata:str, title:str):

    body = {
        "definition": statemachine_definition,
        "statemachine_uri": statemachine_uri,
        "input": inputdata,
        "title": title
    }

    urlpath = urllib.parse.urljoin (host, "/api/v1/StartExecution")
    hearers = {
        "Content-Type": "application/json"
    }
    print("StartExecution API URL: ", urlpath)
    print("StartExecution API Body: \n", json.dumps(body))
    print("\n")
    response = requests.post(urlpath, data=json.dumps(body), headers=hearers)
    if response.status_code == 200:
        respJson = response.json()
        if respJson["success"] == True:
            print("StartExecution started successfully")
            print("Response: ", response.json())
        else:
            print("Execution failed with error: ", respJson)
            return
    else:
        print("Execution failed with status code: ", response.status_code)
        print(response.text)



if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Skyflow Cmdline StartExecution Tool")
    parser.add_argument("--statemachine_file", help="state machine defintion file path")
    parser.add_argument("--statemachine_uri", help="statemachine uri", default="")
    parser.add_argument("--inputfile", help="input file path")
    parser.add_argument("--input", help="input data")
    parser.add_argument("--title", help="title of the execution", default="test execution")
    parser.add_argument("--host", help="cloudflow host", default=default_host)

    args = parser.parse_args()

    host = args.host

    if args.statemachine_file:
        statemachine_definition = read_file(args.statemachine_file)

    inputdata = ""
    if args.input:
        inputdata = args.input
    elif args.inputfile:
        inputdata = read_file(args.inputfile)

    main(statemachine_definition, args.statemachine_uri, inputdata, args.title)
