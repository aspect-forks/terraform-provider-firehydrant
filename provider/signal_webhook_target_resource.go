package provider

import (
	"context"
	"errors"

	"github.com/firehydrant/firehydrant-go-sdk/models/components"
	"github.com/firehydrant/terraform-provider-firehydrant/firehydrant"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSignalWebhookTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: createResourceFireHydrantSignalWebhookTarget,
		ReadContext:   readResourceFireHydrantSignalWebhookTarget,
		UpdateContext: updateResourceFireHydrantSignalWebhookTarget,
		DeleteContext: deleteResourceFireHydrantSignalWebhookTarget,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the webhook target",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL that the webhook target will notify",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional description of the webhook target",
			},
			"signing_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "An optional secret provided in the FH-Signature header when sending payloads. Not returned by the API once set.",
			},
		},
	}
}

func createResourceFireHydrantSignalWebhookTarget(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	name := d.Get("name").(string)
	tflog.Debug(ctx, "Create signal webhook target", map[string]interface{}{
		"name": name,
	})

	createReq := components.CreateSignalsWebhookTarget{
		Name: name,
		URL:  d.Get("url").(string),
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		createReq.Description = &desc
	}

	if v, ok := d.GetOk("signing_key"); ok {
		key := v.(string)
		createReq.SigningKey = &key
	}

	target, err := client.Sdk.Signals.CreateSignalsWebhookTarget(ctx, createReq)
	if err != nil {
		return diag.Errorf("Error creating signal webhook target %s: %v", name, err)
	}

	d.SetId(*target.ID)

	return readResourceFireHydrantSignalWebhookTarget(ctx, d, m)
}

func readResourceFireHydrantSignalWebhookTarget(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, "Read signal webhook target", map[string]interface{}{
		"id": id,
	})

	target, err := client.Sdk.Signals.GetSignalsWebhookTarget(ctx, id)
	if err != nil {
		if errors.Is(err, firehydrant.ErrorNotFound) {
			tflog.Debug(ctx, "Signal webhook target not found, removing from state", map[string]interface{}{
				"id": id,
			})
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading signal webhook target %s: %v", id, err)
	}

	attributes := map[string]interface{}{
		"name": *target.Name,
		"url":  *target.URL,
	}

	if target.Description != nil {
		attributes["description"] = *target.Description
	}

	for key, value := range attributes {
		if err := d.Set(key, value); err != nil {
			return diag.Errorf("Error setting %s for signal webhook target %s: %v", key, id, err)
		}
	}

	return diag.Diagnostics{}
}

func updateResourceFireHydrantSignalWebhookTarget(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, "Update signal webhook target", map[string]interface{}{
		"id": id,
	})

	name := d.Get("name").(string)
	url := d.Get("url").(string)
	updateReq := components.UpdateSignalsWebhookTarget{
		Name: &name,
		URL:  &url,
	}

	if v, ok := d.GetOk("description"); ok {
		desc := v.(string)
		updateReq.Description = &desc
	}

	if d.HasChange("signing_key") {
		if v, ok := d.GetOk("signing_key"); ok {
			key := v.(string)
			updateReq.SigningKey = &key
		}
	}

	_, err := client.Sdk.Signals.UpdateSignalsWebhookTarget(ctx, id, updateReq)
	if err != nil {
		return diag.Errorf("Error updating signal webhook target %s: %v", id, err)
	}

	return readResourceFireHydrantSignalWebhookTarget(ctx, d, m)
}

func deleteResourceFireHydrantSignalWebhookTarget(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*firehydrant.APIClient)

	id := d.Id()
	tflog.Debug(ctx, "Delete signal webhook target", map[string]interface{}{
		"id": id,
	})

	err := client.Sdk.Signals.DeleteSignalsWebhookTarget(ctx, id)
	if err != nil {
		return diag.Errorf("Error deleting signal webhook target %s: %v", id, err)
	}

	d.SetId("")
	return diag.Diagnostics{}
}
