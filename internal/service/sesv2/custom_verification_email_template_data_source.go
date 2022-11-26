package sesv2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func DataSourceCustomVerificationEmailTemplate() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceCustomVerificationEmailTemplateRead,

		Schema: map[string]*schema.Schema{
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from_email_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_subject": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"template_content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"success_redirection_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failure_redirection_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

const (
	DSNameCustomVerificationEmailTemplate = "Custom Verification Email Template Data Source"
)

func dataSourceCustomVerificationEmailTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SESV2Client

	out, err := FindCustomVerificationEmailTemplateByID(ctx, conn, d.Get("template_name").(string))
	if err != nil {
		return create.DiagError(names.SESV2, create.ErrActionReading, DSNameCustomVerificationEmailTemplate, d.Get("template_name").(string), err)
	}
	templateName := aws.ToString(out.TemplateName)
	d.SetId(templateName)
	d.Set("from_email_address", out.FromEmailAddress)
	d.Set("template_subject", out.TemplateSubject)
	d.Set("template_content", out.TemplateContent)
	d.Set("success_redirection_url", out.SuccessRedirectionURL)
	d.Set("failure_redirection_url", out.FailureRedirectionURL)

	return nil
}
