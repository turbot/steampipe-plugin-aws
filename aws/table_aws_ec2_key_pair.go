package aws

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

//// TABLE DEFINITION

func tableAwsEc2KeyPair(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_ec2_key_pair",
		Description: "AWS EC2 Key Pair",
		Get: &plugin.GetConfig{
			KeyColumns:        plugin.SingleColumn("key_name"),
			ShouldIgnoreError: isNotFoundError([]string{"InvalidKeyPair.NotFound", "InvalidKeyPair.Unavailable", "InvalidKeyPair.Malformed"}),
			ItemFromKey:       keyDetailsFromKey,
			Hydrate:           getEc2KeyPair,
		},
		List: &plugin.ListConfig{
			Hydrate: listEc2KeyPairs,
		},
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

//// BUILD HYDRATE INPUT

func keyDetailsFromKey(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	quals := d.KeyColumnQuals
	keyName := quals["key_name"].GetStringValue()
	item := &ec2.KeyPairInfo{
		KeyName: &keyName,
	}
	return item, nil
}

//// LIST FUNCTION

func listEc2KeyPairs(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	plugin.Logger(ctx).Trace("listEc2KeyPairs", "AWS_REGION", defaultRegion)

	// Create Session
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	resp, err := svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})

	for _, keyPair := range resp.KeyPairs {
		d.StreamListItem(ctx, keyPair)
	}
	return nil, err
}

//// HYDRATE FUNCTIONS

func getEc2KeyPair(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	defaultRegion := GetDefaultRegion()
	keyPair := h.Item.(*ec2.KeyPairInfo)

	// create service
	svc, err := Ec2Service(ctx, d, defaultRegion)
	if err != nil {
		return nil, err
	}

	params := &ec2.DescribeKeyPairsInput{
		KeyNames: []*string{aws.String(*keyPair.KeyName)},
	}

	op, err := svc.DescribeKeyPairs(params)
	if err != nil {
		return nil, err
	}

	if op.KeyPairs != nil && len(op.KeyPairs) > 0 {
		return op.KeyPairs[0], nil
	}
	return nil, nil
}

func getAwsEc2KeyPairAkas(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("getAwsEc2KeyPairAkas")
	keyPair := h.Item.(*ec2.KeyPairInfo)
	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// Get data for turbot defined properties
	akas := []string{"arn:" + commonColumnData.Partition + ":ec2:" + commonColumnData.Region + ":" + commonColumnData.AccountId + ":key-pair/" + *keyPair.KeyName}

	return akas, nil
}

//// TRANSFORM FUNCTIONS

func getEc2KeyPairTurbotTags(_ context.Context, d *transform.TransformData) (interface{}, error) {
	keyPair := d.HydrateItem.(*ec2.KeyPairInfo)
	return ec2TagsToMap(keyPair.Tags)
}
