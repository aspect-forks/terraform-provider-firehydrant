package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// This tests the resource with a configuration that only has the required
// attributes specified.
func TestAccSignalAlertGroupingConfigurationResource_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalAlertGroupingConfigurationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignalAlertGroupingConfigurationResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSignalAlertGroupingConfigurationResourceExistsWithAttributes_basic("firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "reference_alert_time_period", "PT30M"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.field_name", "summary"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.values.#", "1"),
					resource.TestCheckTypeSetElemAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.values.*", fmt.Sprintf("test-value-%s", rName)),
				),
			},
		},
	})
}

// This tests the resources ability to update and remove attributes
// with a configuration that has all attributes specified.
func TestAccSignalAlertGroupingConfigurationResource_update(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	rNameUpdated := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalAlertGroupingConfigurationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignalAlertGroupingConfigurationResourceConfig_basic(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSignalAlertGroupingConfigurationResourceExistsWithAttributes_basic("firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "reference_alert_time_period", "PT30M"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.field_name", "summary"),
				),
			},
			{
				Config: testAccSignalAlertGroupingConfigurationResourceConfig_update(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSignalAlertGroupingConfigurationResourceExistsWithAttributes_update("firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "reference_alert_time_period", "PT1H"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.field_name", "body"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.values.#", "2"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.match_type", "and"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "action.0.link", "true"),
				),
			},
			{
				Config: testAccSignalAlertGroupingConfigurationResourceConfig_basic(rNameUpdated),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckSignalAlertGroupingConfigurationResourceExistsWithAttributes_basic("firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration"),
					resource.TestCheckResourceAttrSet("firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "id"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "reference_alert_time_period", "PT30M"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.field_name", "summary"),
					resource.TestCheckResourceAttr(
						"firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration", "strategy.0.substring.0.values.#", "1"),
				),
			},
		},
	})
}

// This tests the resource's ability to import with a configuration that
// only has the required attributes specified.
func TestAccSignalAlertGroupingConfigurationResourceImport_basic(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalAlertGroupingConfigurationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignalAlertGroupingConfigurationResourceConfig_basic(rName),
			},
			{
				ResourceName:      "firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// This tests the resource's ability to import with a configuration that
// has all attributes specified.
func TestAccSignalAlertGroupingConfigurationResourceImport_allAttributes(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalAlertGroupingConfigurationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSignalAlertGroupingConfigurationResourceConfig_update(rName),
			},
			{
				ResourceName:      "firehydrant_signal_alert_grouping_configuration.test_signal_alert_grouping_configuration",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

// This tests that a configuration asking for both actions at once is rejected,
// because the API takes one action per grouping configuration.
func TestAccSignalAlertGroupingConfigurationResource_conflictingActions(t *testing.T) {
	t.Parallel()
	rName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testFireHydrantIsSetup(t) },
		ProviderFactories: sharedProviderFactories(),
		CheckDestroy:      testAccCheckSignalAlertGroupingConfigurationResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config:      testAccSignalAlertGroupingConfigurationResourceConfig_conflictingActions(rName),
				ExpectError: regexp.MustCompile("Conflicting configuration arguments"),
			},
		},
	})
}

// This test does a more in-depth check to test that what was in the
// configuration matches the information we get back from the API
// with a configuration that only has the required attributes specified.
func testAccCheckSignalAlertGroupingConfigurationResourceExistsWithAttributes_basic(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupingResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if groupingResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		groupingResponse, err := client.Sdk.Signals.GetSignalsAlertGroupingConfiguration(context.TODO(), groupingResource.Primary.ID)
		if err != nil {
			return err
		}

		if groupingResponse.ReferenceAlertTimePeriod == nil {
			return fmt.Errorf("Unexpected reference_alert_time_period. Expected a reference_alert_time_period, got: nil")
		}
		expected, got := groupingResource.Primary.Attributes["reference_alert_time_period"], *groupingResponse.ReferenceAlertTimePeriod
		if expected != got {
			return fmt.Errorf("Unexpected reference_alert_time_period. Expected: %s, got: %s", expected, got)
		}

		if groupingResponse.Strategy == nil || groupingResponse.Strategy.Substring == nil {
			return fmt.Errorf("Unexpected strategy. Expected a substring strategy, got: nil")
		}
		if groupingResponse.Strategy.Substring.FieldName == nil {
			return fmt.Errorf("Unexpected field_name. Expected a field_name, got: nil")
		}
		expected, got = groupingResource.Primary.Attributes["strategy.0.substring.0.field_name"], *groupingResponse.Strategy.Substring.FieldName
		if expected != got {
			return fmt.Errorf("Unexpected field_name. Expected: %s, got: %s", expected, got)
		}

		return nil
	}
}

// This test does a more in-depth check to test that what was in the
// configuration matches the information we get back from the API
// with a configuration that has all attributes specified.
func testAccCheckSignalAlertGroupingConfigurationResourceExistsWithAttributes_update(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		groupingResource, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if groupingResource.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		groupingResponse, err := client.Sdk.Signals.GetSignalsAlertGroupingConfiguration(context.TODO(), groupingResource.Primary.ID)
		if err != nil {
			return err
		}

		if groupingResponse.ReferenceAlertTimePeriod == nil {
			return fmt.Errorf("Unexpected reference_alert_time_period. Expected a reference_alert_time_period, got: nil")
		}
		expected, got := groupingResource.Primary.Attributes["reference_alert_time_period"], *groupingResponse.ReferenceAlertTimePeriod
		if expected != got {
			return fmt.Errorf("Unexpected reference_alert_time_period. Expected: %s, got: %s", expected, got)
		}

		if groupingResponse.Strategy == nil || groupingResponse.Strategy.Substring == nil {
			return fmt.Errorf("Unexpected strategy. Expected a substring strategy, got: nil")
		}
		substring := groupingResponse.Strategy.Substring

		if substring.FieldName == nil {
			return fmt.Errorf("Unexpected field_name. Expected a field_name, got: nil")
		}
		expected, got = groupingResource.Primary.Attributes["strategy.0.substring.0.field_name"], *substring.FieldName
		if expected != got {
			return fmt.Errorf("Unexpected field_name. Expected: %s, got: %s", expected, got)
		}

		if substring.MatchType == nil {
			return fmt.Errorf("Unexpected match_type. Expected a match_type, got: nil")
		}
		expected, got = groupingResource.Primary.Attributes["strategy.0.substring.0.match_type"], *substring.MatchType
		if expected != got {
			return fmt.Errorf("Unexpected match_type. Expected: %s, got: %s", expected, got)
		}

		if len(substring.Values) != 2 {
			return fmt.Errorf("Unexpected values. Expected 2 values, got: %d", len(substring.Values))
		}

		if groupingResponse.Action == nil {
			return fmt.Errorf("Unexpected action. Expected an action, got: nil")
		}
		if groupingResponse.Action.Link == nil || !*groupingResponse.Action.Link {
			return fmt.Errorf("Unexpected action.link. Expected: true, got: %v", groupingResponse.Action.Link)
		}

		return nil
	}
}

// This tests that the resource gets destroyed properly
func testAccCheckSignalAlertGroupingConfigurationResourceDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client, err := getAccTestClient()
		if err != nil {
			return err
		}

		for _, stateResource := range s.RootModule().Resources {
			if stateResource.Type != "firehydrant_signal_alert_grouping_configuration" {
				continue
			}

			if stateResource.Primary.ID == "" {
				return fmt.Errorf("No instance ID is set")
			}

			_, err := client.Sdk.Signals.GetSignalsAlertGroupingConfiguration(context.TODO(), stateResource.Primary.ID)
			if err == nil {
				return fmt.Errorf("Signal alert grouping configuration %s still exists", stateResource.Primary.ID)
			}
		}

		return nil
	}
}

func testAccSignalAlertGroupingConfigurationResourceConfig_basic(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_signal_alert_grouping_configuration" "test_signal_alert_grouping_configuration" {
  reference_alert_time_period = "PT30M"

  strategy {
    substring {
      field_name = "summary"
      values     = ["test-value-%s"]
    }
  }
}`, rName)
}

func testAccSignalAlertGroupingConfigurationResourceConfig_update(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_signal_alert_grouping_configuration" "test_signal_alert_grouping_configuration" {
  reference_alert_time_period = "PT1H"

  strategy {
    substring {
      field_name = "body"
      values     = ["test-value-%s", "test-other-value-%s"]
      match_type = "and"
    }
  }

  action {
    link = true
  }
}`, rName, rName)
}

func testAccSignalAlertGroupingConfigurationResourceConfig_conflictingActions(rName string) string {
	return fmt.Sprintf(`
resource "firehydrant_signal_alert_grouping_configuration" "test_signal_alert_grouping_configuration" {
  reference_alert_time_period = "PT30M"

  strategy {
    substring {
      field_name = "summary"
      values     = ["test-value-%s"]
    }
  }

  action {
    link = true

    fyi {
      slack_channel_ids = ["00000000-0000-4000-8000-000000000000"]
    }
  }
}`, rName)
}

func TestFlattenSignalAlertGroupingStrategy(t *testing.T) {
	fieldName, matchType, legacyValue := "tags", "or", "legacy-value"

	for _, tc := range []struct {
		name     string
		strategy *components.NullableSignalsAPIGroupingEntityStrategyEntity
		expected []interface{}
	}{
		{
			name:     "nil strategy",
			strategy: nil,
			expected: []interface{}{},
		},
		{
			name:     "strategy without a substring",
			strategy: &components.NullableSignalsAPIGroupingEntityStrategyEntity{},
			expected: []interface{}{},
		},
		{
			name: "substring with values",
			strategy: &components.NullableSignalsAPIGroupingEntityStrategyEntity{
				Substring: &components.NullableSignalsAPIGroupingEntityStrategyEntitySubstringEntity{
					FieldName: &fieldName,
					Values:    []string{"one", "two"},
					MatchType: &matchType,
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"field_name": "tags",
					"values":     []string{"one", "two"},
					"match_type": "or",
				},
			},
		},
		{
			// Grouping configurations created before the API had a values list
			// only carry the singular value, so it has to stand in for the list.
			name: "substring with only the legacy singular value",
			strategy: &components.NullableSignalsAPIGroupingEntityStrategyEntity{
				Substring: &components.NullableSignalsAPIGroupingEntityStrategyEntitySubstringEntity{
					FieldName: &fieldName,
					Value:     &legacyValue,
				},
			},
			expected: []interface{}{
				map[string]interface{}{
					"field_name": "tags",
					"values":     []string{"legacy-value"},
					"match_type": "",
				},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := flattenSignalAlertGroupingStrategy(tc.strategy)
			if !reflect.DeepEqual(tc.expected, got) {
				t.Fatalf("Unexpected strategy. Expected: %#v, got: %#v", tc.expected, got)
			}
		})
	}
}

// The API takes slack_channel_ids as FireHydrant's own IDs, so the ID Slack
// assigned, which read returns alongside it, has to stay out of state.
func TestFlattenSignalAlertGroupingSlackChannelIDs(t *testing.T) {
	firehydrantID, slackID, name := "76a9fc36-10b8-43f9-b180-51146a7e4d50", "C01010101Z", "#incidents"
	channels := []components.IntegrationsSlackSlackChannelEntity{
		{ID: &firehydrantID, SlackChannelID: &slackID, Name: &name},
	}

	expected := []string{firehydrantID}
	got := flattenSignalAlertGroupingSlackChannelIDs(channels)
	if !reflect.DeepEqual(expected, got) {
		t.Fatalf("Unexpected slack_channel_ids. Expected: %#v, got: %#v", expected, got)
	}
}

// signalAlertGroupingUIPayloads are the request bodies the FireHydrant UI sends
// when creating an alert grouping configuration, one per action. The provider has
// to produce the same body for an equivalent configuration, so these pin the
// request building against payloads that are known to be accepted. Note that the
// fyi payload carries no link key at all, even though an fyi action links the
// grouped alert too.
var signalAlertGroupingUIPayloads = []struct {
	name    string
	action  []interface{}
	values  []interface{}
	field   string
	payload string
}{
	{
		name:   "link action",
		action: []interface{}{map[string]interface{}{"link": true}},
		values: []interface{}{"foobarbaz"},
		field:  "tags",
		payload: `{
		    "strategy": {
		        "substring": {
		            "field_name": "tags",
		            "values": [
		                "foobarbaz"
		            ],
		            "match_type": "and"
		        }
		    },
		    "reference_alert_time_period": "PT30M",
		    "action": {
		        "link": true
		    }
		}`,
	},
	{
		name: "fyi action",
		action: []interface{}{
			map[string]interface{}{
				"fyi": []interface{}{
					map[string]interface{}{
						"slack_channel_ids": []interface{}{"76a9fc36-10b8-43f9-b180-51146a7e4d50"},
					},
				},
			},
		},
		values: []interface{}{"foobarbaz2"},
		field:  "summary",
		payload: `{
		    "strategy": {
		        "substring": {
		            "field_name": "summary",
		            "values": [
		                "foobarbaz2"
		            ],
		            "match_type": "and"
		        }
		    },
		    "reference_alert_time_period": "PT30M",
		    "action": {
		        "fyi": {
		            "slack_channel_ids": [
		                "76a9fc36-10b8-43f9-b180-51146a7e4d50"
		            ]
		        }
		    }
		}`,
	},
}

func TestBuildSignalAlertGroupingRequestsMatchUIPayloads(t *testing.T) {
	for _, tc := range signalAlertGroupingUIPayloads {
		t.Run(tc.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, resourceSignalAlertGroupingConfiguration().Schema, map[string]interface{}{
				"reference_alert_time_period": "PT30M",
				"strategy": []interface{}{
					map[string]interface{}{
						"substring": []interface{}{
							map[string]interface{}{
								"field_name": tc.field,
								"values":     tc.values,
								"match_type": "and",
							},
						},
					},
				},
				"action": tc.action,
			})

			createRequest, err := buildSignalAlertGroupingCreateRequest(d)
			if err != nil {
				t.Fatalf("Error building the create request: %v", err)
			}
			assertMatchesSignalAlertGroupingUIPayload(t, createRequest, tc.payload)

			updateRequest, err := buildSignalAlertGroupingUpdateRequest(d)
			if err != nil {
				t.Fatalf("Error building the update request: %v", err)
			}
			assertMatchesSignalAlertGroupingUIPayload(t, updateRequest, tc.payload)
		})
	}
}

// assertMatchesSignalAlertGroupingUIPayload marshals a request and compares it to
// the payload the UI sends, ignoring key order and whitespace.
func assertMatchesSignalAlertGroupingUIPayload(t *testing.T, request interface{}, payload string) {
	t.Helper()

	got, err := json.Marshal(request)
	if err != nil {
		t.Fatalf("Error marshalling the request: %v", err)
	}

	var gotBody, expectedBody interface{}
	if err := json.Unmarshal(got, &gotBody); err != nil {
		t.Fatalf("Error unmarshalling the request: %v", err)
	}
	if err := json.Unmarshal([]byte(payload), &expectedBody); err != nil {
		t.Fatalf("Error unmarshalling the UI payload: %v", err)
	}

	if !reflect.DeepEqual(expectedBody, gotBody) {
		t.Fatalf("Unexpected request body.\nExpected: %s\nGot:      %s", payload, got)
	}
}
