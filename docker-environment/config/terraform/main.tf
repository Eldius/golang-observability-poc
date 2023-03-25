
terraform {
  required_version = "~> 1.3"
  required_providers {
    elasticsearch = {
      source  = "phillbaker/elasticsearch"
      version = "2.0.7"
    }
  }
}

provider "elasticsearch" {
  # url        = "http://127.0.0.1:9200"
  # kibana_url = "http://127.0.0.1:5601"
  url        = "http://elasticsearch:9200"
  kibana_url = "http://kibana:5601"
}

resource "elasticsearch_xpack_index_lifecycle_policy" "logger" {
  name = "logger-lifecycle"
  body = <<EOF
{
  "policy": {
    "phases": {
      "hot": {
        "min_age": "0ms",
        "actions": {
          "rollover": {
            "max_size": "50mb"
          }
        }
      }
    }
  }
}
EOF
}

resource "elasticsearch_composable_index_template" "logger" {
  name = "filebeat"
  body = <<EOF
{
  "index_patterns": ["filebeat-*"],
  "template": {
    "settings": {
      "index": {
        "lifecycle": {
            "name": "${elasticsearch_xpack_index_lifecycle_policy.logger.name}"
        }
      }
    },
    "aliases": {
      "filebeat": { }
    },
    "mappings": {
      "dynamic":"true"
    }
  },
  "priority": 200,
  "version": 3
}
EOF
}

resource "elasticsearch_kibana_object" "logging_index_pattern" {
  body = <<EOF
[
  {
    "_id": "index-pattern:app-logging",
    "_type": "doc",
    "_source": {
      "type": "index-pattern",
      "index-pattern": {
        "title": "filebeat-7.15.2-*",
        "timeFieldName": "@timestamp"
      }
    }
  }
]
EOF
}

resource "elasticsearch_index" "logger" {
  name               = "filebeat-7.15.2-001"
  number_of_shards   = 1
  number_of_replicas = 1
  aliases = jsonencode({
    "application-logs" = {
      "is_write_index" = true
    }
  })
}

resource "elasticsearch_kibana_alert" "test" {
  name = "terraform-alert"
  schedule {
    interval = "1m"
  }
  conditions {
    aggregation_type     = "count"
    term_size            = 6
    threshold_comparator = "="
    time_window_size     = 5
    time_window_unit     = "m"
    group_by             = "top"
    threshold            = [1000]
    index                = ["application-logs"]
    time_field           = "@timestamp"
    # aggregation_field    = "message"
    term_field           = "message.keyword"
  }
  actions {
    id             = "c87f0dc6-c301-4988-aee9-95d391359a39"
    action_type_id = ".index"
    params = {
      level   = "info"
      message = "alert '{{alertName}}' is active for group '{{context.group}}':\n\n- Value: {{context.value}}\n- Conditions Met: {{context.conditions}} over {{params.timeWindowSize}}{{params.timeWindowUnit}}\n- Timestamp: {{context.date}}"
    }
  }
}