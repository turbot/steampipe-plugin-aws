package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	ec2v1 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/turbot/steampipe-plugin-sdk/v4/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v4/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsEc2KeyPair(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_key_pair",
		Description: "AWS EC2 Key Pair",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("key_name"),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"InvalidKeyPair.NotFound", "InvalidKeyPair.Unavailable", "InvalidKeyPair.Malformed"}),
			},
			Hydrate: getEc2KeyPair,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2KeyPairs,
			KeyColumns: []*plugin.KeyColumn{
				{Name: "key_pair_id", Require: plugin.Optional},
				{Name: "key_fingerprint", Require: plugin.Optional},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(ec2v1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "key_name",
				Description: "The name of the key pair",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_pair_id",
				Description: "The ID of the key pair",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "key_fingerprint",
				Description: "If key pair was created using CreateKeyPair, this is the SHA-1 digest of the DER encoded private key. If key pair was created using ImportKeyPair to provide AWS the public key, this is the MD5 public key fingerprint as specified in section 4 of RFC4716",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "tags_src",
				Description: "A list of tags assigned to the key pair",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Tags"),
			},
			{
				Name:        "tags",
				Description: resourceInterfaceDescription("tags"),
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(getEc2KeyPairTurbotTags),
			},
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("KeyName"),
			},
			{
				Name:        "akas",
				Description: resourceInterfaceDescription("akas"),
				Type:        proto.ColumnType_JSON,
				Hydrate:     getAwsEc2KeyPairAkas,
				Transform:   transform.FromValue(),
			},
		}),
	}
}

//// LIST FUNCTION

func listEc2KeyPairs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {

	// Create Session
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_key_pair.listEc2KeyPairs", "connection_error", err)
		return nil, err
	}

	input := &ec2.DescribeKeyPairsInput{}

	filters := buildEc2KeyPairFilter(d.Quals)

	if len(filters) > 0 {
		input.Filters = filters
	}

	resp, err := svc.DescribeKeyPairs(ctx, input)

	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_key_pair.listEc2KeyPairs", "api_error", err)
		return nil, err
	}

	for _, keyPair := range resp.KeyPairs {
		d.StreamListItem(ctx, keyPair)

		// Context may get cancelled due to manual cancellation or if the limit has been reached
		if d.QueryStatus.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2KeyPair(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	keyName := d.KeyColumnQuals["key_name"].GetStringValue()

	// create service
	svc, err := EC2Client(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_key_pair.getEc2KeyPair", "connection_error", err)
		return nil, err
	}

	params := &ec2.DescribeKeyPairsInput{
		KeyNames: []string{keyName},
	}

	op, err := svc.DescribeKeyPairs(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_key_pair.getEc2KeyPair", "api_error", err)
		return nil, err
	}

	if op.KeyPairs != nil && len(op.KeyPairs) > 0 {
		return op.KeyPairs[0], nil
	}
	return nil, nil
}

func getAwsEc2KeyPairAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	region := d.KeyColumnQualString(matrixKeyRegion)
	keyPair := h.Item.(types.KeyPairInfo)

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_ec2_key_pair.getAwsEc2KeyPairAkas", "common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + region + ":" + commonColumnData.AccountId + ":key-pair/" + *keyPair.KeyName}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getEc2KeyPairTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	keyPair := d.HydrateItem.(types.KeyPairInfo)
	var turbotTagsMap map[string]string
	if keyPair.Tags == nil {
		return nil, nil
	}

	turbotTagsMap = map[string]string{}
	for _, i := range keyPair.Tags {
		turbotTagsMap[*i.Key] = *i.Value
	}

	return &turbotTagsMap, nil
}

//// UTILITY FUNCTIONS

// Build ec2 key-pair list call input filter
func buildEc2KeyPairFilter(quals plugin.KeyColumnQualMap) []types.Filter {
	filters := make([]types.Filter, 0)

	filterQuals := map[string]string{
		"key_pair_id":     "key-pair-id",
		"key_fingerprint": "fingerprint",
	}

	for columnName, filterName := range filterQuals {
		if quals[columnName] != nil {
			filter := types.Filter{
				Name: aws.String(filterName),
			}
			value := getQualsValueByColumn(quals, columnName, "string")
			val, ok := value.(string)
			if ok {
				filter.Values = []string{val}
			}
			filters = append(filters, filter)
		}
	}
	return filters
}
