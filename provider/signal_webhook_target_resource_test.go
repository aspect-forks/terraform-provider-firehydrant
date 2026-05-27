package provider

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSignalWebhookTargetResource_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalWebhookTargetDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignalWebhookTargetConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSignalWebhookTargetExists("firehydrant_signal_webhook_target.test"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_webhook_target.test", "id"),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "name", fmt.Sprintf("test-webhook-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "url", "https://example.com/webhook"),
				),
			},
		},
	})
}

func TestAccSignalWebhookTargetResource_withDescription(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalWebhookTargetDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignalWebhookTargetConfig_withDescription(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSignalWebhookTargetExists("firehydrant_signal_webhook_target.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "name", fmt.Sprintf("test-webhook-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "url", "https://example.com/webhook"),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "description", "A test webhook target"),
				),
			},
		},
	})
}

func TestAccSignalWebhookTargetResource_update(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalWebhookTargetDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignalWebhookTargetConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSignalWebhookTargetExists("firehydrant_signal_webhook_target.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "name", fmt.Sprintf("test-webhook-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "url", "https://example.com/webhook"),
				),
			},
			{
				Config: testAccSignalWebhookTargetConfig_update(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSignalWebhookTargetExists("firehydrant_signal_webhook_target.test"),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "name", fmt.Sprintf("updated-webhook-%s", rName)),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "url", "https://example.com/updated-webhook"),
					resource.TestCheckResourceAttr("firehydrant_signal_webhook_target.test", "description", "An updated webhook target"),
				),
			},
		},
	})
}

func TestAccSignalWebhookTargetResourceImport_basic(t *testing.T) {
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalWebhookTargetDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignalWebhookTargetConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_signal_webhook_target.test",
				ImportState:       true,
				ImportStateVerify: true,
				// signing_key is write-only and never returned by the API
				ImportStateVerifyIgnore: []string{"signing_key"},
			},
		},
	})
}

func testAccCheckSignalWebhookTargetExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return fmt.Errorf("Error getting client: %s", err)
		}

		target, err := client.Sdk.Signals.GetSignalsWebhookTarget(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error fetching signal webhook target with ID %s: %s", rs.Primary.ID, err)
		}

		expected, got := rs.Primary.Attributes["name"], *target.Name
		if expected != got {
			return fmt.Errorf("Unexpected name. Expected: %s, got: %s", expected, got)
		}

		expected, got = rs.Primary.Attributes["url"], *target.URL
		if expected != got {
			return fmt.Errorf("Unexpected url. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

func testAccCheckSignalWebhookTargetDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return fmt.Errorf("Error getting client: %s", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "firehydrant_signal_webhook_target" {
				continue
			}

			_, err := client.Sdk.Signals.GetSignalsWebhookTarget(context.Background(), rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Signal webhook target %s still exists", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccSignalWebhookTargetConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_signal_webhook_target" "test" {
  name = "test-webhook-%s"
  url  = "https://example.com/webhook"
}
`, rName)
}

func testAccSignalWebhookTargetConfig_withDescription(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_signal_webhook_target" "test" {
  name        = "test-webhook-%s"
  url         = "https://example.com/webhook"
  description = "A test webhook target"
}
`, rName)
}

func testAccSignalWebhookTargetConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_signal_webhook_target" "test" {
  name        = "updated-webhook-%s"
  url         = "https://example.com/updated-webhook"
  description = "An updated webhook target"
}
`, rName)
}
