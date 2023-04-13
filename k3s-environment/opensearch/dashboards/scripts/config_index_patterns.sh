#!/bin/sh

echo ""
echo "####################"
echo "## creating index ##"
echo "####################"
echo ""

curl -i \
    --fail \
    --insecure \
    -X PUT \
    -u 'admin:admin' \
    'https://192.168.1.194:9200/application-logs-00001' \
    -H 'Content-Type: application/json' \
    -d '{
    "settings": {
        "index": {
            "number_of_shards": 2,
            "number_of_replicas": 1
        }
    },
    "mappings": {
        "properties": {
            "@timestamp": {
                "type": "date"
            },
            "@version": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
            },
            "agent": {
                "properties": {
                    "ephemeral_id": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    },
                    "id": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    },
                    "name": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    },
                    "type": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    },
                    "version": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    }
                }
            },
            "container": {
                "properties": {
                    "id": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    },
                    "image": {
                        "properties": {
                            "name": {
                                "type": "text",
                                "fields": {
                                    "keyword": {
                                        "type": "keyword",
                                        "ignore_above": 256
                                    }
                                }
                            }
                        }
                    },
                    "name": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    }
                }
            },
            "docker": {
                "properties": {
                    "container": {
                        "properties": {
                            "labels": {
                                "type": "object"
                            }
                        }
                    }
                }
            },
            "ecs": {
                "properties": {
                    "version": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    }
                }
            },
            "event": {
                "properties": {
                    "original": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    }
                }
            },
            "host": {
                "properties": {
                    "name": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    }
                }
            },
            "input": {
                "properties": {
                    "type": {
                        "type": "text",
                        "fields": {
                            "keyword": {
                                "type": "keyword",
                                "ignore_above": 256
                            }
                        }
                    }
                }
            },
            "log": {
                "properties": {
                    "file": {
                        "properties": {
                            "path": {
                                "type": "text",
                                "fields": {
                                    "keyword": {
                                        "type": "keyword",
                                        "ignore_above": 256
                                    }
                                }
                            }
                        }
                    },
                    "offset": {
                        "type": "long"
                    }
                }
            },
            "message": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
            },
            "stream": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
            },
            "parsed": {
                "type": "object"
            },
            "tags": {
                "type": "text",
                "fields": {
                    "keyword": {
                        "type": "keyword",
                        "ignore_above": 256
                    }
                }
            }
        }
    },
    "aliases": {
        "application-logs": {}
    }
}' || echo "Failed to create logs index"

echo ""
echo ""
echo "############################"
echo "## creating index pattern ##"
echo "############################"
echo ""


curl -i \
    --fail \
    --insecure \
    -X POST \
    -u 'admin:admin' \
    'http://192.168.1.194:5601/api/saved_objects/index-pattern/application-logs' \
    -H 'osd-xsrf: true' \
    -H 'Content-Type: application/json' \
    -H 'securitytenant: global' \
    -d '{
        "attributes": {
            "title": "application-logs",
            "timeFieldName": "@timestamp"
        }
    }'

echo ""
echo ""

curl -i \
        --fail \
        --insecure \
        -X POST \
        -u 'admin:admin' \
        'http://192.168.1.194:5601/api/saved_objects/index-pattern/metrics-otel-v1-*' \
        -H 'osd-xsrf: true' \
        -H 'Content-Type: application/json' \
        -H 'securitytenant: global' \
        -d '{
            "attributes": {
                "title": "metrics-otel-v1-*",
                "timeFieldName": "@timestamp"
            }
        }'
