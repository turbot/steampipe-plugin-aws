package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/transfer"
	"github.com/aws/aws-sdk-go-v2/service/transfer/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsTransferServer(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_transfer_server",
		Description: "AWS Transfer Server",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("server_id"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"ResourceNotFoundException"}),
			},
			Hydrate: getTransferServer,
			Tags:    map[string]string{"service": "transfer", "action": "DescribeServer"},
		},
		List: &plugin.ListConfig{
			Hydrate: listTransferServers,
			Tags:    map[string]string{"service": "transfer", "action": "ListServers"},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func: getTransferServer,
				Tags: map[string]string{"service": "transfer", "action": "DescribeServer"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(AWS_TRANSFER_SERVICE_ID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "arn",
				Description: "The Amazon Resource Name (ARN) for the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "server_id",
				Description: "The system-assigned unique identifier for the server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "domain",
				Description: "Specifies the domain of the storage system that is used for file transfers.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "identity_provider_type",
				Description: "The mode of authentication for a server.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "endpoint_type",
				Description: "Specifies the type of VPC endpoint that your server is connected to.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "logging_role",
				Description: "The Amazon Resource Name (ARN) of the Identity and Access Management (IAM) role that allows a server to turn on Amazon CloudWatch logging for Amazon S3 or Amazon EFSevents. When set, you can view user activity in your CloudWatch logs.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The condition of the server that was described.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "user_count",
				Description: "Specifies the number of users that are assigned to a server.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "certificate",
				Description: "Specifies the ARN of the Amazon Web ServicesCertificate Manager (ACM) certificate.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "host_key_fingerprint",
				Description: "Specifies the Base64-encoded SHA256 fingerprint of the server's host key.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "pre_authentication_login_banner",
				Description: "Specifies a string to display when users connect to a server. This string is displayed before the user authenticates.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "post_authentication_login_banner",
				Description: "Specifies a string to display when users connect to a server. This string is displayed after the user authenticates.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "security_policy_name",
				Description: "Specifies the name of the security policy that is attached to the server.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "identity_provider_details",
				Description: "Specifies information to call a customer-supplied authentication API.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "protocol_details",
				Description: "The protocol settings that are configured for your server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "endpoint_details",
				Description: "The virtual private cloud (VPC) endpoint settings that are configured for your server.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "protocols",
				Description: "Specifies the file transfer protocol or protocols over which your file transfer protocol client can connect to your server's endpoint.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "workflow_details",
				Description: "Specifies the workflow ID for the workflow to assign and the execution role that's used for executing the workflow.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferServer,
			},
			{
				Name:        "structured_log_destinations",
				Description: "Specifies the log groups to which your server logs are sent.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getTransferServer,
			},

			// Steampipe standard columns
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("ServerId"),
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

//// LIST FUNCTION

func listTransferServers(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	// Create session
	svc, err := TransferClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_server.listTransferServers", "client_error", err)
		return nil, err
	}
	if svc == nil {
		return nil, nil
	}
	// Limiting the results
	maxLimit := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	input := &transfer.ListServersInput{
		MaxResults: &maxLimit,
	}

	paginator := transfer.NewListServersPaginator(svc, input, func(o *transfer.ListServersPaginatorOptions) {
		o.Limit = maxLimit
		o.StopOnDuplicateToken = true
	})

	// List call
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			plugin.Logger(ctx).Error("aws_transfer_server.listTransferServers", "api_error", err)
			return nil, err
		}

		for _, items := range output.Servers {
			d.StreamListItem(ctx, items)

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getTransferServer(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// get service
	svc, err := TransferClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_server.getTransferServer", "connection_error", err)
		return nil, err
	}

	var serverID string
	if h.Item != nil {
		serverID = *h.Item.(types.ListedServer).ServerId
	} else {
		serverID = d.EqualsQuals["server_id"].GetStringValue()
	}

	// Empty Check
	if serverID == "" {
		return nil, nil
	}

	// Build the params
	params := &transfer.DescribeServerInput{
		ServerId: &serverID,
	}

	// Get call
	op, err := svc.DescribeServer(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_transfer_server.getTransferServer", "api_error", err)
		return nil, err
	}

	if op.Server != nil {
		return op.Server, nil
	}
	return nil, nil
}
