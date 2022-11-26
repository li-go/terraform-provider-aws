package sesv2

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func ResourceCustomVerificationEmailTemplate() *schema.Resource {
	return &schema.Resource{
		CreateWithoutTimeout: resourceCustomVerificationEmailTemplateCreate,
		ReadWithoutTimeout:   resourceCustomVerificationEmailTemplateRead,
		UpdateWithoutTimeout: resourceCustomVerificationEmailTemplateUpdate,
		DeleteWithoutTimeout: resourceCustomVerificationEmailTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"template_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"from_email_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_subject": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"success_redirection_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"failure_redirection_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

const (
	ResNameCustomVerificationEmailTemplate = "Custom Verification Email Template"
)

func resourceCustomVerificationEmailTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SESV2Client

	in := &sesv2.CreateCustomVerificationEmailTemplateInput{
		TemplateName:          aws.String(d.Get("template_name").(string)),
		FromEmailAddress:      aws.String(d.Get("from_email_address").(string)),
		TemplateSubject:       aws.String(d.Get("template_subject").(string)),
		TemplateContent:       aws.String(d.Get("template_content").(string)),
		SuccessRedirectionURL: aws.String(d.Get("success_redirection_url").(string)),
		FailureRedirectionURL: aws.String(d.Get("failure_redirection_url").(string)),
	}

	out, err := conn.CreateCustomVerificationEmailTemplate(ctx, in)
	if err != nil {
		return create.DiagError(names.SESV2, create.ErrActionCreating, ResNameCustomVerificationEmailTemplate, d.Get("template_name").(string), err)
	}

	if out == nil {
		return create.DiagError(names.SESV2, create.ErrActionCreating, ResNameCustomVerificationEmailTemplate, d.Get("template_name").(string), errors.New("empty output"))
	}

	d.SetId(d.Get("template_name").(string))

	return resourceCustomVerificationEmailTemplateRead(ctx, d, meta)
}

func resourceCustomVerificationEmailTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SESV2Client

	out, err := FindCustomVerificationEmailTemplateByID(ctx, conn, d.Id())

	if !d.IsNewResource() && tfresource.NotFound(err) {
		log.Printf("[WARN] SESV2 CustomVerificationEmailTemplate (%s) not found, removing from state", d.Id())
		d.SetId("")
		return nil
	}

	if err != nil {
		return create.DiagError(names.SESV2, create.ErrActionReading, ResNameCustomVerificationEmailTemplate, d.Id(), err)
	}

	d.Set("template_name", d.Id())
	d.Set("from_email_address", out.FromEmailAddress)
	d.Set("template_subject", out.TemplateSubject)
	d.Set("template_content", out.TemplateContent)
	d.Set("success_redirection_url", out.SuccessRedirectionURL)
	d.Set("failure_redirection_url", out.FailureRedirectionURL)

	return nil
}

func resourceCustomVerificationEmailTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SESV2Client

	in := &sesv2.UpdateCustomVerificationEmailTemplateInput{
		TemplateName: aws.String(d.Id()),
	}

	if v, ok := d.GetOk("from_email_address"); ok {
		in.FromEmailAddress = aws.String(v.(string))
	}
	if v, ok := d.GetOk("template_subject"); ok {
		in.TemplateSubject = aws.String(v.(string))
	}
	if v, ok := d.GetOk("template_content"); ok {
		in.TemplateContent = aws.String(v.(string))
	}
	if v, ok := d.GetOk("success_redirection_url"); ok {
		in.SuccessRedirectionURL = aws.String(v.(string))
	}
	if v, ok := d.GetOk("failure_redirection_url"); ok {
		in.FailureRedirectionURL = aws.String(v.(string))
	}

	log.Printf("[DEBUG] Update SES custom verification email template: %#v", in)
	_, err := conn.UpdateCustomVerificationEmailTemplate(ctx, in)
	if err != nil {
		return create.DiagError(names.SESV2, create.ErrActionUpdating, ResNameCustomVerificationEmailTemplate, d.Id(), err)
	}

	return resourceCustomVerificationEmailTemplateRead(ctx, d, meta)
}

func resourceCustomVerificationEmailTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn := meta.(*conns.AWSClient).SESV2Client

	log.Printf("[INFO] Deleting SESV2 CustomVerificationEmailTemplate %s", d.Id())

	_, err := conn.DeleteCustomVerificationEmailTemplate(ctx, &sesv2.DeleteCustomVerificationEmailTemplateInput{
		TemplateName: aws.String(d.Id()),
	})

	if err != nil {
		var nfe *types.NotFoundException
		if errors.As(err, &nfe) {
			return nil
		}

		return create.DiagError(names.SESV2, create.ErrActionDeleting, ResNameCustomVerificationEmailTemplate, d.Id(), err)
	}

	return nil
}

func FindCustomVerificationEmailTemplateByID(ctx context.Context, conn *sesv2.Client, id string) (*sesv2.GetCustomVerificationEmailTemplateOutput, error) {
	in := &sesv2.GetCustomVerificationEmailTemplateInput{
		TemplateName: aws.String(id),
	}
	out, err := conn.GetCustomVerificationEmailTemplate(ctx, in)
	if err != nil {
		var nfe *types.NotFoundException
		if errors.As(err, &nfe) {
			return nil, &resource.NotFoundError{
				LastError:   err,
				LastRequest: in,
			}
		}

		return nil, err
	}

	if out == nil {
		return nil, tfresource.NewEmptyResultError(in)
	}

	return out, nil
}
