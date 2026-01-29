import redis
import pymongo as pm
from pymongo.errors import BulkWriteError
import time
import json


# manages communication between the two databases

# rlimit controls how frequently redis is uploaded to pymongo, every x entries
rlimit = 20


def init_db(c):
    db = c["database"]
    col = db["page_information"]
    col.create_index(
        [("page_information.URL", 1), ("page_information.Name", 1)], unique=True
    )


def monitor():
    global rlimit

    r = redis.Redis(host="192.168.57.3", port=6379, decode_responses=True)
    while True:

        rlen = r.llen("page_information")
        if rlen > rlimit:
            with r.lock("lock", timeout=60, blocking_timeout=10):
                # repeat incase incremented inbetween
                rlen = r.llen("page_information")
                data = r.lrange("page_information", 0, -1)
                r.ltrim("page_information", rlen, -1)
            manage_json(data)


def manage_json(data):
    db_arr = []
    for entry in data:

        esc_data = json.loads(entry)
        # just the og page data from before, with stuff like url,title and content
        page_data = esc_data[0]
        page_data = json.loads(page_data)
        page_data = json.loads(page_data)
        # a dictionary of words and their frequency within page content
        freq_data = esc_data[1]
        freq_data = json.loads(freq_data)

        combined_data = dict(page_information=page_data, word_frequency=freq_data)
        db_arr.append(combined_data)
    print("Appending new data for Mongo")
    send_to_db(db_arr)


def send_to_db(db_arr):
    db = c["database"]
    col = db["page_information"]
    try:
        sent_data = col.insert_many(db_arr)
    except BulkWriteError as err:
        print("Duplicates found")

    le = col.count_documents({})
    print("Mongo db size:", le)


def main():
    print("running")
    global c
    c = pm.MongoClient("192.168.57.30")
    init_db(c)
    monitor()


main()
