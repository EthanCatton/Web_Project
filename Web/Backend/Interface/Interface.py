import fastapi
from fastapi import FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
import subprocess
import uvicorn
import json
import sys
import os

query = "placeholder"
interface = FastAPI()

# code yoinked from gpt here because its basic config
interface.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@interface.post("/get_query")
async def get_query(request: Request):
    data = await request.json()
    query = data.get("searchbar")
    result = get_results(query)
    return result


# sends input from search into query script
def get_results(query):
    print("Query passed to subprocess", repr(query))
    proc = subprocess.run(
        ["python", subpath], input=query + "\n", text=True, capture_output=True
    )

    result = proc.stdout
    result = result.strip()

    err = proc.stderr
    err = err.strip()

    if err != "":
        return {"result": err}
    else:
        print("returning result")
        return {"result": result}


if __name__ == "__main__":
    print("running")

    # Gets correct relative path for subproccess
    filepath = sys.path[0]
    filedir = os.path.dirname(filepath)
    subpath = filedir + "\Query\Query.py"
    print(subpath)

    uvicorn.run(interface, host="127.0.0.1", port=8000)
