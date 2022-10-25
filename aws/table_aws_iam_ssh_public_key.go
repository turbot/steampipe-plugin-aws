package aws

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

func tableAwsIamSshPublicKey(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_ssh_public_key",
		Description: "AWS IAM SSH public key.",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"ssh_public_key_id", "user_name"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: isNotFoundErrorV2([]string{"ValidationError", "NoSuchEntity", "InvalidParameter"}),
			},
			Hydrate: getIamSshPublicKey,
		},
		List: &plugin.ListConfig{
			ParentHydrate: listIamUsers,
			Hydrate:       listIamSshPublicKeys,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "user_name", Require: plugin.Optional},
			},
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "ssh_public_key_id",
				Description: "The unique identifier for the SSH public key.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSHPublicKeyId"),
			},
			{
				Name:        "ssh_public_key_body_pem",
				Description: "The SSH public key, PEM encoded.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSHPublicKeyBody"),
				Hydrate:     getIamSshPublicKeyAsPem,
			},
			{
				Name:        "ssh_public_key_body_rsa",
				Description: "The SSH public key, SSH RSA encoded.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSHPublicKeyBody"),
				Hydrate:     getIamSshPublicKey,
			},
			{
				Name:        "Fingerprint",
				Description: "The MD5 message digest of the SSH public key.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getIamSshPublicKey,
			},
			{
				Name:        "user_name",
				Description: "The name of the IAM user associated with the SSH public key.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "status",
				Description: "The status of the SSH public key. Active means that the key can be used for authentication with an CodeCommit repository. Inactive means that the key cannot.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "upload_date",
				Description: "The date when the SSH public key was uploaded.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("SSHPublicKeyId"),
			},
		}),
	}
}

//// LIST FUNCTION

func listIamSshPublicKeys(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	user := h.Item.(types.User)

	// Minimize the API call with the given user_name
	equalQuals := d.KeyColumnQuals
	if equalQuals["user_name"] != nil {
		if equalQuals["user_name"].GetStringValue() != "" {
			if equalQuals["user_name"].GetStringValue() != "" && equalQuals["user_name"].GetStringValue() != *user.UserName {
				return nil, nil
			}
		} else if len(getListValues(equalQuals["user_name"].GetListValue())) > 0 {
			if !strings.Contains(fmt.Sprint(getListValues(equalQuals["user_name"].GetListValue())), *user.UserName) {
				return nil, nil
			}
		}
	}

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_ssh_public_key.listIamSshPublicKeys", "client_error", err)
		return nil, err
	}

	params := &iam.ListSSHPublicKeysInput{UserName: user.UserName}

	paginator := iam.NewListSSHPublicKeysPaginator(svc, params, func(o *iam.ListSSHPublicKeysPaginatorOptions) {
		o.Limit = 10
		o.StopOnDuplicateToken = true
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_ssh_public_key.listIamSshPublicKeys", "api_error", err)
			return nil, err
		}

		for _, ssh := range output.SSHPublicKeys {
			d.StreamListItem(ctx, ssh)

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.QueryStatus.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getIamSshPublicKey(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	userName := d.KeyColumnQuals["user_name"].GetStringValue()
	if userName == "" {
		userName = *h.Item.(types.SSHPublicKeyMetadata).UserName
	}

	sshPublicKeyId := d.KeyColumnQuals["ssh_public_key_id"].GetStringValue()
	if sshPublicKeyId == "" {
		sshPublicKeyId = *h.Item.(types.SSHPublicKeyMetadata).SSHPublicKeyId
	}

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_ssh_public_key.getIamSshPublicKey", "client_error", err)
		return nil, err
	}

	params := &iam.GetSSHPublicKeyInput{
		UserName:       &userName,
		SSHPublicKeyId: &sshPublicKeyId,
		Encoding:       "SSH",
	}

	op, err := svc.GetSSHPublicKey(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_ssh_public_key.getIamSshPublicKey", "api_error", err)
		return nil, err
	}

	return *op.SSHPublicKey, nil
}

func getIamSshPublicKeyAsPem(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	userName := d.KeyColumnQuals["user_name"].GetStringValue()
	if userName == "" {
		userName = *h.Item.(types.SSHPublicKeyMetadata).UserName
	}

	sshPublicKeyId := d.KeyColumnQuals["ssh_public_key_id"].GetStringValue()
	if sshPublicKeyId == "" {
		sshPublicKeyId = *h.Item.(types.SSHPublicKeyMetadata).SSHPublicKeyId
	}

	// Get client
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_ssh_public_key.getIamSshPublicKey", "client_error", err)
		return nil, err
	}

	params := &iam.GetSSHPublicKeyInput{
		UserName:       &userName,
		SSHPublicKeyId: &sshPublicKeyId,
		Encoding:       "PEM",
	}

	op, err := svc.GetSSHPublicKey(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_ssh_public_key.getIamSshPublicKey", "api_error", err)
		return nil, err
	}

	return *op.SSHPublicKey, nil
}
