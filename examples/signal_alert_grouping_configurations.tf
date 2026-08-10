resource "firehydrant_signal_alert_grouping_configuration" "example" {
  reference_alert_time_period = "PT30M"

  strategy {
    substring {
      field_name = "summary"
      values     = ["disk pressure", "memory pressure"]
      match_type = "or"
    }
  }

  action {
    link = true
  }
}

data "firehydrant_slack_channel" "example" {
  slack_channel_name = "#incidents"
}

resource "firehydrant_signal_alert_grouping_configuration" "example_fyi" {
  reference_alert_time_period = "PT1H"

  strategy {
    substring {
      field_name = "tags"
      values     = ["service:checkout"]
    }
  }

  action {
    fyi {
      slack_channel_ids = [data.firehydrant_slack_channel.example.id]
    }
  }
}
