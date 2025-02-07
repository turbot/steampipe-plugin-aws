package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/transfer"
	"github.com/aws/aws-sdk-go-v2/service/transfer/types"
	transferEndpoint "github.com/turbot/steampipe-plugin-aws/awsSupportedEndpoints"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

type TransferUserInfo = struct {
	types.DescribedUser
	UserName          *string
	ServerID          *string
	Arn               *string
	HomeDirectory     *string
	HomeDirectoryType *types.HomeDirectoryType
	Role              *string
	SshPublicKeyCount *int32
}

// // TABLE DEFINITION
func tableAwsTransferUser(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_transfer_user",
		Description: "AWS Transfer User",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"server_id", "user_name"}),
			Hydrate:    getTransferUser,
			Tags:       map[string]string{"service": "transfer", "action": "DescribeUser"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listTransferServers,
			Hydrate:       listTransferUsers,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "server_id", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "transfer", "action": "ListUsers"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getTransferUser,
				Tags: map[string]string{"service": "transfer", "action": "DescribeUser"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(transferEndpoint.AWS_TRANSFER_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the user.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_id",
				Description: "The ID of the server that the user is attached to.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerID"),
			},
			{
				Name:        "user_name",
				Description: "Specifies the name of the user whose ARN was specified. User names are used for authentication purposes.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "home_directory",
				Description: "Specifies the landing directory (folder) for a user when they log in to the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "home_directory_mappings",
				Description: "The landing directory (folder) for a user when they log in to the server using the client.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferUser,
			},
			{
				Name:        "home_directory_type",
				Description: "The type of landing directory (folder) you mapped for your users to see when they log in to the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "role",
				Description: "The Amazon Resource Name (ARN) of the AWS Identity and Access Management (IAM) role that controls your users' access to your Amazon S3 bucket or Amazon EFS file system.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "ssh_public_key_count",
				Description: "The number of SSH public keys stored for the user on the server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "ssh_public_keys",
				Description: "The public SSH keys stored for the user on the server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferUser,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferUser,
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("UserName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Arn").Transform(transform.EnsureStringArray),
			},
		}),
	}
}

// // HYDRATE FUNCTIONS
func listTransferUsers(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	id := h.Item.(types.ListedServer).ServerId

	// Create session
	svc, err := TransferClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_user.listTransferUsers", "client_error", err)
		return nil, err
	}

	maxItems := int32(100)
	params := &transfer.ListUsersInput{
		ServerId: id,
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			params.MaxResults = &limit
		}
	}

	paginator := transfer.NewListUsersPaginator(svc, params, func(o *transfer.ListUsersPaginatorOptions) {
		o.Limit = maxItems
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		d.WaitForListRateLimit(ctx)

		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_transfer_user.listTransferUsers", "api_error", err)
			return nil, err
		}

		for _, item := range output.Users {
			d.StreamListItem(ctx, TransferUserInfo{
				UserName:          item.UserName,
				ServerID:          id,
				Arn:               item.Arn,
				HomeDirectory:     item.HomeDirectory,
				HomeDirectoryType: &item.HomeDirectoryType,
				Role:              item.Role,
				SshPublicKeyCount: item.SshPublicKeyCount,
			})
		}
	}
	return nil, nil
}

func getTransferUser(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Session
	svc, err := TransferClient(ctx, d)
	if err != nil {
		return nil, err
	}

	var userName string
	var serverID string
	if h.Item != nil {
		userName = *h.Item.(TransferUserInfo).UserName
		serverID = *h.Item.(TransferUserInfo).ServerID
	} else {
		userName = d.EqualsQualString("user_name")
		serverID = d.EqualsQualString("server_id")
	}

	// check if userName or serverID is empty
	if userName == "" || serverID == "" {
		return nil, nil
	}

	// Build the params
	params := &transfer.DescribeUserInput{
		ServerId: &serverID,
		UserName: &userName,
	}

	// Get the service response
	op, err := svc.DescribeUser(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_user.getTransferUser", "api_error", err)
		return nil, err
	}
	sshPublicKeyCount := int32(len(op.User.SshPublicKeys))

	return TransferUserInfo{
		*op.User,
		&userName,
		&serverID,
		op.User.Arn,
		op.User.HomeDirectory,
		&op.User.HomeDirectoryType,
		op.User.Role,
		&sshPublicKeyCount,
	}, nil
}
