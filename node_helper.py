#!/usr/bin/env python
# -*- coding: utf-8 -*-
# Author: gexiao
# Created on 2017-10-10 20:56

from pymongo import MongoClient

MONGODB_HOST = '127.0.0.1'
MONGODB_PORT = 27017
MONGODB_USER = 'v2exuser'
MONGODB_PASSWORD = 'readwrite'
MONGODB_DBNAME = 'v2ex'

if MONGODB_USER:
    client = MongoClient(MONGODB_HOST, MONGODB_PORT,
                         username=MONGODB_USER, password=MONGODB_PASSWORD,
                         authSource=MONGODB_DBNAME, authMechanism='SCRAM-SHA-1')
else:
    client = MongoClient(MONGODB_HOST, MONGODB_PORT)

db = client.v2ex
node_collection = db.node


def find_node(name):
    if name:
        return node_collection.find_one({'$or': [{'name': name},
                                                 {'title': name},
                                                 {'title_alternative': name}]})
    return None
