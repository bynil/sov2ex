#!/usr/bin/env python
# -*- coding: utf-8 -*-
# Author: gexiao
# Created on 2017-09-19 22:12

from elasticsearch import Elasticsearch

ES_HOST = '127.0.0.1:9200'
TOPIC_ALIAS_NAME = 'topic'
TOPIC_TYPE_NAME = 'topic'

es = Elasticsearch([ES_HOST])


def generate_search_body(keyword, es_from, es_size):
    body = {
        "from": es_from,
        "size": es_size,
        "highlight": {
            "order": "score",
            "fragment_size": 80,
            "fields": {
                "title": {
                    "number_of_fragments": 1
                },
                "content": {
                    "number_of_fragments": 1
                },
                "postscript_list.content": {
                    "number_of_fragments": 1
                }
            }
        },
        "_source": ["title",
                    "content",
                    "created",
                    "id",
                    "node",
                    "replies",
                    "member"],
        "query": {
            "function_score": {
                "query": {
                    "bool": {
                        "must_not": [
                            {
                                "term": {
                                    "deleted": True
                                }
                            }
                        ],
                        "should": [
                            {
                                "match": {
                                    "title": {
                                        "query": keyword,
                                        "analyzer": "ik_smart",
                                        "boost": 3
                                    }
                                }
                            },
                            {
                                "bool": {
                                    "should": [
                                        {
                                            "match": {
                                                "content": {
                                                    "query": keyword,
                                                    "analyzer": "ik_smart",
                                                    "boost": 2
                                                }
                                            }
                                        },
                                        {
                                            "nested": {
                                                "path": "postscript_list",
                                                "score_mode": "max",
                                                "query": {
                                                    "match": {
                                                        "postscript_list.content": {
                                                            "query": keyword,
                                                            "analyzer": "ik_smart",
                                                            "boost": 2
                                                        }
                                                    }
                                                }
                                            }
                                        }
                                    ]
                                }
                            },
                            {
                                "match": {
                                    "all_reply": {
                                        "query": keyword,
                                        "analyzer": "ik_smart",
                                        "boost": 1.5
                                    }
                                }
                            }
                        ]
                    }
                },
                "functions": [
                    {
                        "filter": {"match_phrase": {
                            "all_content": {
                                "query": keyword,
                                "analyzer": "ik_max_word",
                                "slop": 0
                            }
                        }},
                        "weight": 50
                    },
                    {
                        "field_value_factor": {
                            "field": "bonus",
                            "missing": 0,
                            "modifier": "none",
                            "factor": 1
                        }
                    }
                ],
                "score_mode": "sum",
                "boost_mode": "sum"
            }
        }
    }
    return body


def es_search(keyword, es_from, es_size):
    return es.search(index=TOPIC_ALIAS_NAME, doc_type=TOPIC_TYPE_NAME,
                     body=generate_search_body(keyword, es_from, es_size))


def es_analyze(keyword):
    body = {
        'text': keyword,
        'analyzer': 'ik_smart'
    }
    return es.indices.analyze(index=TOPIC_ALIAS_NAME, body=body)


def es_clause_count(keyword):
    return len(es_analyze(keyword)['tokens'])