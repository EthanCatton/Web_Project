import pymongo as pm
import json
import nltk
from nltk.corpus import stopwords as nltkstopwords
from nltk.stem import WordNetLemmatizer
from nltk.tokenize import word_tokenize
import sys
import string
 
#queries the mongo database and returns the results to the frontend via interface.py

url_dict={}

def main():
 query=getquery()
 query=proccesquery(query)
 

 global c
 c=pm.MongoClient("mongodb://localhost:27018")
 
 getmongo(query)
 
 
def proccesquery(data): 
 #just a copy of the lematising code from processor
 
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

    #Needs to be list for pipeline 
    #data=" ".join(data)
    return data


def getquery():
 query=sys.stdin.readline()
 query=query.strip()
 return query

def getmongo(query):
 db=c["database"]
 col=db["page_information"]  

 #AI ASSISTED TO MAKE CAUSE GOD HELP ME IN FORMATTING MATCH SCORE
 aggpipeline = [
        {
            "$project": {
                "_id": 0,
                "URL": {
                    "$ifNull": ["$URL", "$page_information.URL"]
                },
                "Name": {
                    "$ifNull": ["$URL", "$page_information.Name"]
                },
                "match_score": {
                    "$sum": [
                        {
                            "$ifNull": [
                                {
                                    "$toInt": { 
                                        "$getField": {
                                            "field": kw,
                                            "input": {
                                                "$ifNull": [ 
                                                    "$word_frequency",
                                                    {"$ifNull": ["$page_information.word_frequency", {}]}
                                                ]
                                            }
                                        }
                                    }
                                },
                                0
                            ]
                        }
                        for kw in query
                    ]
                }
            }
        },
        {"$match": {"match_score": {"$gt": 0}}}, 
        {"$sort": {"match_score": -1}}
    ]
 results = list(col.aggregate(aggpipeline))

 results=json.dumps(results)

 
 print(results) 
  








main()



# dictionary of all urls containing keywords
# add up total of keywords found per url for a score in the dict
# order results by that score


