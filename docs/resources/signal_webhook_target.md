---
page_title: "FireHydrant Resource: firehydrant_signal_webhook_target"
subcategory: "Signals"
---

# firehydrant_signal_webhook_target Resource

FireHydrant signal webhook targets allow you to configure URLs that FireHydrant will notify when signals are received.

## Example Usage

Basic usage:
```hcl
resource "firehydrant_signal_webhook_target" "example" {
  name = "Example Webhook"
  url  = "https://example.com/webhook"
}
```

With an optional description and signing key:
```hcl
resource "firehydrant_signal_webhook_target" "example" {
  name        = "Example Webhook"
  url         = "https://example.com/webhook"
  description = "Forwards signals to the example service"
  signing_key = var.webhook_signing_key
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the webhook target.
* `url` - (Required) The URL that the webhook target will notify.
* `description` - (Optional) A description of the webhook target.
* `signing_key` - (Optional, Sensitive) A secret provided in the `FH-Signature` header when sending payloads. This value is write-only and will not be returned by the API once set.

## Import

Signal webhook targets can be imported using the webhook target ID. For example:

```shell
terraform import firehydrant_signal_webhook_target.example <id>
```
