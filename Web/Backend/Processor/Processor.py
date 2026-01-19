import nltk
from nltk.corpus import stopwords as nltkstopwords
from nltk.stem import WordNetLemmatizer
from nltk.tokenize import word_tokenize
import string
import socket as sc
import asyncio as asy
import json
import threading
import redis
from collections import Counter


checker=set()
lock=threading.Lock()

upload_counter=0

# this file should take singular inputs from the queue and reduce them down to keywords to be stored in the database
# both the original input and keywords should be kept

async def data_input(reader,writer):
   
    while True:
     data=await reader.readline() 
     if not data:
         break
     data=data.decode()
     data=data.strip()
     if data:
      sequence(data)
     writer.write("received".encode())
     await writer.drain()
    writer.close()
    await writer.wait_closed()   
    print("recieved data")

def reformat(og_data):
    global checker 
    #Converts to json to gets fields 
    full_data=json.loads(og_data)
    url=full_data["URL"]
    page_data=full_data["Content"]

    #Duplicate mitigation attempt
    #the issue was actually print buffering but this has been kept
    #just in case of future weird issues
    with lock:
       if url in checker:
          flag=True
       else:
        flag=False
        checker.add(url)
        if len(checker)>50:
           checker.clear()


    return og_data,page_data,flag




def process(data): 
    print("processing data")
    
    
    if not isinstance(data,str):
       return "None"
    
    
    data=data.lower()
   
    
    data=nltk.word_tokenize(data)
    #remove punctuation
    data=[text for text in data if text not in string.punctuation]

    #remove stopwords
    stopwords=set(nltkstopwords.words("english"))
    data=[text for text in data if text not in stopwords]

    #custom removal
    removal=["\\n","''","``"]
    data=[text for text in data if text not in removal]


    #reduce down
    lemmatizer=WordNetLemmatizer()
    data=[lemmatizer.lemmatize(text) for text in data] 

    data=Counter(data)
    data=dict(data)
    

    return data


def data_output(og_data,lem_data):
   global upload_counter
   raw_json=json.dumps(og_data)
   lem_json=json.dumps(lem_data)
   stored_data=json.dumps([raw_json,lem_json]) 

   r=redis.Redis(host="192.168.57.3",port=6379,decode_responses=True) 
   upload_counter+=1
   print("Pushing data to Redis:", upload_counter,og_data[0:200])  
   r.lpush("page_information",stored_data)
   #rlen=r.llen("page_information")
  # print("page information length",rlen)
   r.close()
    

def sequence(og_data):
    global checker 
    data,page_data,flag=reformat(og_data) 
 
    if data=="None" or flag==True:
       return
    lem_data=process(page_data)
    data_output(data,lem_data) 

       

  
   

async def main():
    server= await asy.start_server(data_input,"192.168.57.15",5757,limit=5048000)  
    async with server:
        await server.serve_forever()  

#repeat safeguard 
if __name__ == "__main__":  
    asy.run(main())

