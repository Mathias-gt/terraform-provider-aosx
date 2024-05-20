package aosx

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceRestconf() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRestconfCreate,
		ReadContext:   resourceRestconfRead,
		UpdateContext: resourceRestconfUpdate,
		DeleteContext: resourceRestconfDelete,

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
				id := d.Id()
				d.Set("path", id)

				diags := resourceRestconfRead(ctx, d, m)
				if diags.HasError() {
					return nil, fmt.Errorf("failed to import resource: %s", diags[0].Summary)
				}

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
                        "delete_path": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRestconfCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*Client)

	path := d.Get("path").(string)
	content := d.Get("content").(string)

	err := client.CreateRestconf(ctx, path, content)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(path)

	return resourceRestconfRead(ctx, d, meta)
}

func resourceRestconfRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	path := d.Get("path").(string)

	// Read the configuration from the device
	apiResponse, err := client.ReadRestconf(ctx, path)
	if err != nil {
		return diag.FromErr(err)
	}

	// Decode and re-encode the API response to ensure consistent formatting
	var apiResponseFormatted map[string]interface{}
	if err := json.Unmarshal([]byte(apiResponse), &apiResponseFormatted); err != nil {
		return diag.FromErr(err)
	}
	apiResponseJson, err := json.Marshal(apiResponseFormatted)
	if err != nil {
		return diag.FromErr(err)
	}

	// Read the stored configuration from the state
	content := d.Get("content").(string)

	// Decode and re-encode the stored configuration to ensure consistent formatting
	var contentFormatted map[string]interface{}
	if content == "" {
		contentFormatted = make(map[string]interface{})
	} else {
		if err := json.Unmarshal([]byte(content), &contentFormatted); err != nil {
			return diag.FromErr(err)
		}
	}

	contentJson, err := json.Marshal(contentFormatted)
	if err != nil {
		return diag.FromErr(err)
	}

	// Compare the formatted API response and the formatted stored configuration
	if string(apiResponseJson) != string(contentJson) {
		d.Set("content", string(apiResponseJson))
	}

	return nil
}

func resourceRestconfUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	path := d.Id()
	content := d.Get("content").(string)

	err := client.UpdateRestconf(ctx, path, content)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceRestconfRead(ctx, d, m)
}

func resourceRestconfDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	path := d.Get("delete_path").(string)

	err := client.DeleteRestconf(ctx, path)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
