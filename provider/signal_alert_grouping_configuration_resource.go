package provider

import (
	"context"
	"fmt"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/firehydrant-go-sdk/models/sdkerrors"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceSignalAlertGroupingConfiguration() *schema.Resource {
	return &schema.Resource{
		Description:   "FireHydrant signal alert grouping configurations combine related alerts so that a group of alerts only notifies responders once.",
		CreateContext: createResourceFireHydrantSignalAlertGroupingConfiguration,
		UpdateContext: updateResourceFireHydrantSignalAlertGroupingConfiguration,
		ReadContext:   readResourceFireHydrantSignalAlertGroupingConfiguration,
		DeleteContext: deleteResourceFireHydrantSignalAlertGroupingConfiguration,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			// Required
			"reference_alert_time_period": {
				Type:     schema.TypeString,
				Required: true,
				Description: "How long an alert stays eligible for grouping, as an ISO 8601 duration. The window " +
					"starts when the first matching alert is received.",
			},
			"strategy": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "The strategy that decides which alerts get grouped together.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"substring": {
							Type:        schema.TypeList,
							Required:    true,
							MaxItems:    1,
							Description: "Group alerts whose field contains the given substrings.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Required
									"field_name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The alert field to match on. Valid values are `summary`, `body`, and `tags`.",
										ValidateFunc: validation.StringInSlice([]string{
											string(components.CreateSignalsAlertGroupingConfigurationFieldNameSummary),
											string(components.CreateSignalsAlertGroupingConfigurationFieldNameBody),
											string(components.CreateSignalsAlertGroupingConfigurationFieldNameTags),
										}, false),
									},
									"values": {
										Type:        schema.TypeSet,
										Required:    true,
										MinItems:    1,
										Description: "The substrings to look for in the field.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},

									// Optional
									"match_type": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										Description: "Whether an alert has to contain every value or just one of them. Valid " +
											"values are `and` and `or`.",
										ValidateFunc: validation.StringInSlice([]string{
											string(components.CreateSignalsAlertGroupingConfigurationMatchTypeAnd),
											string(components.CreateSignalsAlertGroupingConfigurationMatchTypeOr),
										}, false),
									},
								},
							},
						},
					},
				},
			},

			// Optional
			"action": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Description: "What to do with an alert that gets grouped. Exactly one of `link` or `fyi` may be set. " +
					"FireHydrant always resolves an action for a grouping configuration, so removing this block does " +
					"not clear the action that is already configured.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						// Optional
						"link": {
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true,
							Description:   "Link the grouped alert and do not notify anyone.",
							ConflictsWith: []string{"action.0.fyi"},
						},
						"fyi": {
							Type:          schema.TypeList,
							Optional:      true,
							MaxItems:      1,
							Description:   "Link the grouped alert and send an FYI notification to Slack.",
							ConflictsWith: []string{"action.0.link"},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									// Required
									"slack_channel_ids": {
										Type:        schema.TypeSet,
										Required:    true,
										MinItems:    1,
										Description: "The FireHydrant IDs of the Slack channels to notify.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func readResourceFireHydrantSignalAlertGroupingConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get the signal alert grouping configuration
	groupingID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Read signal alert grouping configuration: %s", groupingID), map[string]interface{}{
		"id": groupingID,
	})
	groupingResponse, err := client.Sdk.Signals.GetSignalsAlertGroupingConfiguration(ctx, groupingID)
	if err != nil {
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			tflog.Debug(ctx, fmt.Sprintf("Signal alert grouping configuration %s no longer exists", groupingID), map[string]interface{}{
				"id": groupingID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading signal alert grouping configuration %s: %v", groupingID, err)
	}

	// Set the resource attributes to the values we got from the API
	if err := setSignalAlertGroupingResourceData(d, groupingResponse); err != nil {
		return diag.Errorf("Error reading signal alert grouping configuration %s: %v", groupingID, err)
	}

	return diag.Diagnostics{}
}

// setSignalAlertGroupingResourceData writes an API response into state. Read goes
// through this rather than flattening inline so that a test can exercise the same
// path against the real schema: a flattened value whose shape does not match the
// schema only fails when d.Set walks it, which comparing flattened shapes to each
// other cannot catch.
func setSignalAlertGroupingResourceData(d *schema.ResourceData, groupingResponse *components.SignalsAPIGroupingEntity) error {
	// Process any data that could be nil. The SDK models every field as a
	// pointer, so an absent value becomes an empty string rather than a panic.
	// Absent values are still written to state so that a change made outside of
	// Terraform shows up as a diff instead of lingering in state.
	var referenceAlertTimePeriod string
	if groupingResponse.ReferenceAlertTimePeriod != nil {
		referenceAlertTimePeriod = *groupingResponse.ReferenceAlertTimePeriod
	}

	// Gather values from API response
	attributes := map[string]interface{}{
		"reference_alert_time_period": referenceAlertTimePeriod,
		"strategy":                    flattenSignalAlertGroupingStrategy(groupingResponse.Strategy),
		"action":                      flattenSignalAlertGroupingAction(groupingResponse.Action),
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("error setting %s: %w", key, err)
		}
	}

	return nil
}

func createResourceFireHydrantSignalAlertGroupingConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Get attributes from config and construct the create request
	createRequest, err := buildSignalAlertGroupingCreateRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Create the new signal alert grouping configuration
	tflog.Debug(ctx, "Create signal alert grouping configuration", map[string]interface{}{
		"field_name": string(createRequest.Strategy.Substring.FieldName),
	})
	groupingResponse, err := client.Sdk.Signals.CreateSignalsAlertGroupingConfiguration(ctx, *createRequest)
	if err != nil {
		return diag.Errorf("Error creating signal alert grouping configuration: %v", err)
	}
	if groupingResponse.ID == nil {
		return diag.Errorf("Error creating signal alert grouping configuration: the API response did not include an ID")
	}

	// Set the new signal alert grouping configuration's ID in state
	d.SetId(*groupingResponse.ID)

	// Update state with the latest information from the API
	return readResourceFireHydrantSignalAlertGroupingConfiguration(ctx, d, m)
}

func updateResourceFireHydrantSignalAlertGroupingConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Construct the update request
	updateRequest, err := buildSignalAlertGroupingUpdateRequest(d)
	if err != nil {
		return diag.FromErr(err)
	}

	// Update the signal alert grouping configuration
	groupingID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Update signal alert grouping configuration: %s", groupingID), map[string]interface{}{
		"id": groupingID,
	})
	_, err = client.Sdk.Signals.UpdateSignalsAlertGroupingConfiguration(ctx, groupingID, *updateRequest)
	if err != nil {
		return diag.Errorf("Error updating signal alert grouping configuration %s: %v", groupingID, err)
	}

	// Update state with the latest information from the API
	return readResourceFireHydrantSignalAlertGroupingConfiguration(ctx, d, m)
}

func deleteResourceFireHydrantSignalAlertGroupingConfiguration(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Get the API client
	client := m.(*firehydrant.APIClient)

	// Delete the signal alert grouping configuration
	groupingID := d.Id()
	tflog.Debug(ctx, fmt.Sprintf("Delete signal alert grouping configuration: %s", groupingID), map[string]interface{}{
		"id": groupingID,
	})
	err := client.Sdk.Signals.DeleteSignalsAlertGroupingConfiguration(ctx, groupingID)
	if err != nil {
		// A grouping configuration that has already been deleted is not an error,
		// it just needs to come out of state.
		if sdkErr, ok := err.(*sdkerrors.SDKError); ok && sdkErr.StatusCode == 404 {
			tflog.Debug(ctx, fmt.Sprintf("Signal alert grouping configuration %s no longer exists", groupingID), map[string]interface{}{
				"id": groupingID,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error deleting signal alert grouping configuration %s: %v", groupingID, err)
	}

	d.SetId("")

	return diag.Diagnostics{}
}

// buildSignalAlertGroupingCreateRequest turns a configuration into a create
// request. The strategy and the time period are always sent, because the schema
// requires both.
func buildSignalAlertGroupingCreateRequest(d *schema.ResourceData) (*components.CreateSignalsAlertGroupingConfiguration, error) {
	substring, err := expandSignalAlertGroupingSubstring(d)
	if err != nil {
		return nil, err
	}

	createRequest := &components.CreateSignalsAlertGroupingConfiguration{
		ReferenceAlertTimePeriod: d.Get("reference_alert_time_period").(string),
		Strategy: components.CreateSignalsAlertGroupingConfigurationStrategy{
			Substring: &components.CreateSignalsAlertGroupingConfigurationSubstring{
				FieldName: components.CreateSignalsAlertGroupingConfigurationFieldName(substring.fieldName),
				Values:    substring.values,
			},
		},
	}
	if substring.matchType != "" {
		matchType := components.CreateSignalsAlertGroupingConfigurationMatchType(substring.matchType)
		createRequest.Strategy.Substring.MatchType = &matchType
	}

	if action := expandSignalAlertGroupingAction(d); action != nil {
		createRequest.Action = &components.CreateSignalsAlertGroupingConfigurationAction{
			Link: action.link,
		}
		if action.slackChannelIDs != nil {
			createRequest.Action.Fyi = &components.CreateSignalsAlertGroupingConfigurationFyi{
				SlackChannelIds: action.slackChannelIDs,
			}
		}
	}

	return createRequest, nil
}

// buildSignalAlertGroupingUpdateRequest turns a configuration into an update
// request. Every field the API models as optional is still sent, because the API
// replaces the strategy wholesale and the schema requires everything the update
// needs, with the one exception of the action: the API has no way to express
// "this grouping configuration has no action", so removing the action block
// leaves the action that is already configured in place.
func buildSignalAlertGroupingUpdateRequest(d *schema.ResourceData) (*components.UpdateSignalsAlertGroupingConfiguration, error) {
	substring, err := expandSignalAlertGroupingSubstring(d)
	if err != nil {
		return nil, err
	}

	referenceAlertTimePeriod := d.Get("reference_alert_time_period").(string)
	updateRequest := &components.UpdateSignalsAlertGroupingConfiguration{
		ReferenceAlertTimePeriod: &referenceAlertTimePeriod,
		Strategy: &components.UpdateSignalsAlertGroupingConfigurationStrategy{
			Substring: &components.UpdateSignalsAlertGroupingConfigurationSubstring{
				FieldName: components.UpdateSignalsAlertGroupingConfigurationFieldName(substring.fieldName),
				Values:    substring.values,
			},
		},
	}
	if substring.matchType != "" {
		matchType := components.UpdateSignalsAlertGroupingConfigurationMatchType(substring.matchType)
		updateRequest.Strategy.Substring.MatchType = &matchType
	}

	if action := expandSignalAlertGroupingAction(d); action != nil {
		updateRequest.Action = &components.UpdateSignalsAlertGroupingConfigurationAction{
			Link: action.link,
		}
		if action.slackChannelIDs != nil {
			updateRequest.Action.Fyi = &components.UpdateSignalsAlertGroupingConfigurationFyi{
				SlackChannelIds: action.slackChannelIDs,
			}
		}
	}

	return updateRequest, nil
}

// signalAlertGroupingSubstring holds the substring strategy from a
// configuration, so that create and update can build their own request types
// from one set of values.
type signalAlertGroupingSubstring struct {
	fieldName string
	values    []string
	matchType string
}

// expandSignalAlertGroupingSubstring pulls the substring strategy out of a
// configuration. The schema requires strategy and strategy.substring, so a
// missing block means the configuration was not validated, which is worth an
// error rather than a panic.
func expandSignalAlertGroupingSubstring(d *schema.ResourceData) (*signalAlertGroupingSubstring, error) {
	strategyList := d.Get("strategy").([]interface{})
	if len(strategyList) == 0 || strategyList[0] == nil {
		return nil, fmt.Errorf("strategy is required for a signal alert grouping configuration")
	}
	strategy := strategyList[0].(map[string]interface{})

	substringList := strategy["substring"].([]interface{})
	if len(substringList) == 0 || substringList[0] == nil {
		return nil, fmt.Errorf("strategy.substring is required for a signal alert grouping configuration")
	}
	substring := substringList[0].(map[string]interface{})

	return &signalAlertGroupingSubstring{
		fieldName: substring["field_name"].(string),
		values:    expandStringSet(substring["values"].(*schema.Set)),
		matchType: substring["match_type"].(string),
	}, nil
}

// signalAlertGroupingAction holds the action from a configuration. A nil
// slackChannelIDs means the configuration has no fyi block, which is different
// from an empty one.
type signalAlertGroupingAction struct {
	link            *bool
	slackChannelIDs []string
}

// expandSignalAlertGroupingAction pulls the action out of a configuration,
// returning nil when there is no action block to send. The schema keeps link and
// fyi from being configured together, so this only has to decide which one the
// configuration asked for.
func expandSignalAlertGroupingAction(d *schema.ResourceData) *signalAlertGroupingAction {
	actionList := d.Get("action").([]interface{})
	if len(actionList) == 0 || actionList[0] == nil {
		return nil
	}
	action := actionList[0].(map[string]interface{})

	expanded := &signalAlertGroupingAction{}

	// An fyi action links the grouped alert as well as notifying Slack, so the
	// API sets link on its own and reads it back. Sending fyi on its own keeps
	// that from looking like a request for both actions when a configuration
	// moves from link to fyi, because link is still in state at that point.
	fyiList := action["fyi"].([]interface{})
	if len(fyiList) > 0 && fyiList[0] != nil {
		fyi := fyiList[0].(map[string]interface{})
		expanded.slackChannelIDs = expandStringSet(fyi["slack_channel_ids"].(*schema.Set))
		return expanded
	}

	if link := action["link"].(bool); link {
		expanded.link = &link
	}

	return expanded
}

// flattenSignalAlertGroupingStrategy converts the strategy the API returns into
// the nested blocks the schema uses.
func flattenSignalAlertGroupingStrategy(strategy *components.NullableSignalsAPIGroupingEntityStrategyEntity) []interface{} {
	if strategy == nil || strategy.Substring == nil {
		return []interface{}{}
	}
	substring := strategy.Substring

	var fieldName, matchType string
	if substring.FieldName != nil {
		fieldName = *substring.FieldName
	}
	if substring.MatchType != nil {
		matchType = *substring.MatchType
	}

	// The API returns both the current values list and a singular value that
	// predates it. Grouping configurations created before values existed only
	// have the singular field, so fall back to it rather than importing a
	// configuration with no values at all.
	values := substring.Values
	if len(values) == 0 && substring.Value != nil {
		values = []string{*substring.Value}
	}

	// The fields belong under a substring block, not directly under strategy.
	// Flattening them one level too shallow makes d.Set walk the value against the
	// schema and try to write strategy.0.field_name, which is not an address the
	// schema defines.
	return []interface{}{
		map[string]interface{}{
			"substring": []interface{}{
				map[string]interface{}{
					"field_name": fieldName,
					"values":     values,
					"match_type": matchType,
				},
			},
		},
	}
}

// flattenSignalAlertGroupingAction converts the action the API returns into the
// nested blocks the schema uses.
func flattenSignalAlertGroupingAction(action *components.NullableSignalsAPIGroupingEntityActionEntity) []interface{} {
	if action == nil {
		return []interface{}{}
	}

	flattened := map[string]interface{}{}
	if action.Link != nil {
		flattened["link"] = *action.Link
	}
	if action.Fyi != nil && len(action.Fyi.SlackChannels) > 0 {
		flattened["fyi"] = []interface{}{
			map[string]interface{}{
				"slack_channel_ids": flattenSignalAlertGroupingSlackChannelIDs(action.Fyi.SlackChannels),
			},
		}
	}
	if len(flattened) == 0 {
		return []interface{}{}
	}

	return []interface{}{flattened}
}

// flattenSignalAlertGroupingSlackChannelIDs maps the Slack channels the API
// returns back to the identifiers a configuration uses. Create and update take
// slack_channel_ids as FireHydrant's own IDs, which is what the
// firehydrant_slack_channel data source exposes, so the ID Slack assigned that
// read also returns is not what belongs in state.
func flattenSignalAlertGroupingSlackChannelIDs(channels []components.IntegrationsSlackSlackChannelEntity) []string {
	firehydrantIDs := make([]string, 0, len(channels))
	for _, channel := range channels {
		if channel.ID != nil {
			firehydrantIDs = append(firehydrantIDs, *channel.ID)
		}
	}

	return firehydrantIDs
}
