
terraform {
  required_version = "~> 1.3"
  required_providers {
    elasticsearch = {
      source  = "phillbaker/elasticsearch"
      version = "2.0.7"
    }
  }
}

# provider "elasticsearch" {
#   # url        = "http://127.0.0.1:9200"
#   # kibana_url = "http://127.0.0.1:5601"
#   # url        = "http://elasticsearch:9200"
#   # kibana_url = "http://kibana:5601"
#   url         = "http://node-0.example.com:9200"
#   healthcheck = false
# }

# resource "elasticsearch_xpack_index_lifecycle_policy" "logger" {
#   name = "logger-lifecycle"
#   body = <<EOF
# {
#   "policy": {
#     "phases": {
#       "hot": {
#         "min_age": "0ms",
#         "actions": {
#           "rollover": {
#             "max_size": "50mb"
#           }
#         }
#       }
#     }
#   }
# }
# EOF
# }

# resource "elasticsearch_composable_index_template" "logger" {
#   name = "application-logs"
#   body = <<EOF
# {
#   "index_patterns": ["application-logs-*"],
#   "template": {
#     "settings": {
#       "index": {
#         "lifecycle": {
#             "name": "${elasticsearch_xpack_index_lifecycle_policy.logger.name}",
#             "rollover_alias": "application-logs"
#         }
#       }
#     },
#     "aliases": {
#       "application-logs": { }
#     },
#     "mappings": {
#       "dynamic":"true"
#     }
#   },
#   "priority": 200,
#   "version": 3
# }
# EOF
# }

# resource "elasticsearch_index" "logger" {
#   name               = "application-logs-001"
#   number_of_shards   = 1
#   number_of_replicas = 1
#   aliases = jsonencode({
#     "application-logs" = {
#       "is_write_index" = true
#     }
#   })
# }

# Configure the Elasticsearch provider
provider "elasticsearch" {
  url = "http://node-0.example.com:9200"
  healthcheck = false
}

# Create an index template
resource "elasticsearch_index_template" "template_1" {
  name = "template_1"
  body = <<EOF
{
  "template": "te*",
  "settings": {
    "number_of_shards": 1
  },
  "mappings": {
    "type1": {
      "_source": {
        "enabled": false
      },
      "properties": {
        "host_name": {
          "type": "keyword"
        },
        "created_at": {
          "type": "date",
          "format": "EEE MMM dd HH:mm:ss Z YYYY"
        }
      }
    }
  }
}
EOF
}