package azuread

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/graphrbac/1.6/graphrbac"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/graph"
	"github.com/terraform-providers/terraform-provider-azuread/azuread/helpers/validate"
)

func dataSourceUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceUserRead,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"object_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.UUID,
				ConflictsWith: []string{"user_principal_name"},
			},

			"user_principal_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validate.NoEmptyStrings,
				ConflictsWith: []string{"object_id"},
			},

			"account_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mail": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mail_nickname": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).usersClient
	ctx := meta.(*ArmClient).StopContext

	var user graphrbac.User

	if upn, ok := d.Get("user_principal_name").(string); ok && upn != "" {
		resp, err := client.Get(ctx, upn)
		if err != nil {
			return fmt.Errorf("Error making Read request on AzureAD User with ID %q: %+v", upn, err)
		}
		user = resp
	} else if oId, ok := d.Get("object_id").(string); ok && oId != "" {
		u, err := graph.UserGetByObjectId(&client, ctx, oId)
		if err != nil {
			return fmt.Errorf("Error finding Azure AD User with object ID %q: %+v", oId, err)
		}
		user = *u
	} else {
		return fmt.Errorf("one of `object_id` or `user_principal_name` must be supplied")
	}

	if user.ObjectID == nil {
		return fmt.Errorf("Azure AD User objectId is nil")
	}
	d.SetId(*user.ObjectID)

	d.Set("object_id", user.ObjectID)
	d.Set("user_principal_name", user.UserPrincipalName)
	d.Set("account_enabled", user.AccountEnabled)
	d.Set("display_name", user.DisplayName)
	d.Set("mail", user.Mail)
	d.Set("mail_nickname", user.MailNickname)

	return nil
}
