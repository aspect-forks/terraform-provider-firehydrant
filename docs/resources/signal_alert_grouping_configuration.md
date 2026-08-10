---
page_title: "FireHydrant Resource: firehydrant_signal_alert_grouping_configuration"
subcategory: "Signals"
---

# firehydrant_signal_alert_grouping_configuration Resource

FireHydrant signal alert grouping configurations combine related alerts so that a group of alerts only notifies
responders once. An alert that matches a configuration's strategy while an earlier matching alert is still open is
linked to that alert instead of paging on its own.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_signal_alert_grouping_configuration" "example" {
  reference_alert_time_period = "PT30M"

  strategy {
    substring {
      field_name = "summary"
      values     = ["disk pressure"]
    }
  }
}
```

Grouping on several substrings and linking the grouped alerts:
```hcl
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
```

Notifying Slack channels instead of staying silent:
```hcl
data "firehydrant_slack_channel" "example" {
  slack_channel_name = "#incidents"
}

resource "firehydrant_signal_alert_grouping_configuration" "example" {
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
```

## Argument Reference

The following arguments are supported:

* `reference_alert_time_period` - (Required) How long an alert stays eligible for grouping, as an ISO 8601
  duration, e.g. `PT30M`. The window starts when the first matching alert is received.
* `strategy` - (Required) The strategy that decides which alerts get grouped together.
  See [Strategy](#strategy) below for details.
* `action` - (Optional) What to do with an alert that gets grouped.
  See [Action](#action) below for details.

### Strategy

The `strategy` block supports:

* `substring` - (Required) Group alerts whose field contains the given substrings.
  See [Substring](#substring) below for details.

### Substring

The `substring` block supports:

* `field_name` - (Required) The alert field to match on. Valid values are `summary`, `body`, and `tags`.
* `values` - (Required) The substrings to look for in the field. At least one value is required.
* `match_type` - (Optional) Whether an alert has to contain every value or just one of them. Valid values are
  `and` and `or`.

### Action

The `action` block supports:

* `link` - (Optional) Link the grouped alert and do not notify anyone.
* `fyi` - (Optional) Link the grouped alert and send an FYI notification to Slack.
  See [FYI](#fyi) below for details.

Exactly one of `link` or `fyi` may be set. FireHydrant always resolves an action for a grouping configuration and
the API has no way to express a configuration with no action, so removing the `action` block from your
configuration does not clear the action that is already set in FireHydrant.

### FYI

The `fyi` block supports:

* `slack_channel_ids` - (Required) The FireHydrant IDs of the Slack channels to notify. Use the
  `firehydrant_slack_channel` data source to look these up. At least one channel is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the alert grouping configuration.

## Import

Signal alert grouping configuration resources can be imported using the resource ID, e.g.,

```
$ terraform import firehydrant_signal_alert_grouping_configuration.example 12345678-90ab-cdef-1234-567890abcdef
```
