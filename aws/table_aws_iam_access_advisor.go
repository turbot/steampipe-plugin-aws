package aws

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

const maxRetries = 20
const retryIntervalMs = 500

type awsIamAccessAdvisorData struct {
	PrincipalArn               string
	Granularity                string
	LastAuthenticated          *time.Time
	LastAuthenticatedEntity    *string
	LastAuthenticatedRegion    *string
	ServiceName                *string
	ServiceNamespace           *string
	TotalAuthenticatedEntities *int64
	TrackedActionsLastAccessed interface{}
}

func tableAwsIamAccessAdvisor(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:             "aws_iam_access_advisor",
		Description:      "AWS IAM Access Advisor",
		DefaultTransform: transform.FromGo(),
		List: &plugin.ListConfig{
			KeyColumns: plugin.SingleColumn("principal_arn"),
			Hydrate:    listAccessAdvisor,
			Tags:    map[string]string{"service": "iam", "action": "GetServiceLastAccessedDetails"},
		},
		Columns: awsGlobalRegionColumns([]*plugin.Column{
			{
				Name:        "principal_arn",
				Description: "The ARN of the IAM resource (user, group, role, or managed policy) used to generate information about when the resource was last used in an attempt to access an AWS service.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_name",
				Description: "The name of the service in which access was attempted.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_namespace",
				Description: "The namespace of the service in which access was attempted.",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "last_authenticated",
				Description: "The date and time when an authenticated entity most recently attempted to access the service. AWS does not report unauthenticated requests.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "last_authenticated_entity",
				Description: "The ARN of the authenticated entity (user or role) that last attempted to access the service. AWS does not report unauthenticated requests.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_authenticated_region",
				Description: "The Region from which the authenticated entity (user or role) last attempted to access the service. AWS does not report unauthenticated requests.",
				Type:        proto.ColumnType_STRING,
			},

			{
				Name:        "total_authenticated_entities",
				Description: "The total number of authenticated principals (root user, IAM users, or IAM roles) that have attempted to access the service.",
				Type:        proto.ColumnType_INT,
			},
			{
				Name:        "tracked_actions_last_accessed",
				Description: "An array of objects that contains details about the most recent attempt to access a tracked action within the service.  Currently, only S3 supports action level tracking.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

//// LIST FUNCTION

func listAccessAdvisor(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// To simplify the table we always get ACTION_LEVEL.  ACTION_LEVEL is a superset of
	// SERVICE_LEVEL, and currently only s# supports action level actions anyway, so the
	// performance impact is minimal
	granularity := "ACTION_LEVEL"
	principalArn := d.EqualsQuals["principal_arn"].GetStringValue()

	commonData, err := getCommonColumns(ctx, d, h)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_access_advisor.listAccessAdvisor", "get_common_data_error", err)
		return nil, err
	}
	commonColumnData := commonData.(*awsCommonColumnData)

	// check if principalArn is empty or if the account id in the principalArn is not same with the account
	if principalArn == "" || len(strings.Split(principalArn, ":")) < 4 {
		return nil, nil
	} else if strings.Split(principalArn, ":")[4] != "aws" && strings.Split(principalArn, ":")[4] != commonColumnData.AccountId {
		return nil, nil
	}

	// Create Session
	svc, err := IAMClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_access_advisor.listAccessAdvisor", "client_error", err)
		return nil, err
	}

	// Generate the details.  We'll need the job id of this to get the details...
	generateResp, err := svc.GenerateServiceLastAccessedDetails(
		ctx,
		&iam.GenerateServiceLastAccessedDetailsInput{
			Arn:         &principalArn,
			Granularity: types.AccessAdvisorUsageGranularityType(granularity),
		})

	if err != nil {
		plugin.Logger(ctx).Error("aws_iam_access_advisor.listAccessAdvisor", "get_job_id_error", err)
		return nil, err
	}

	// Reduce the basic request limit down if the user has only requested a small number of rows
	maxItems := int32(1000)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxItems {
			if limit < 1 {
				maxItems = int32(1)
			} else {
				maxItems = int32(limit)
			}
		}
	}
	params := &iam.GetServiceLastAccessedDetailsInput{
		JobId:    generateResp.JobId,
		MaxItems: aws.Int32(maxItems),
	}

	retryNumber := 0
	for {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		resp, err := svc.GetServiceLastAccessedDetails(ctx, params)
		if err != nil {
			plugin.Logger(ctx).Error("aws_iam_access_advisor.listAccessAdvisor", "list_advisoer_details_error", err)
			return nil, err
		}

		// if job is still in progress, wait and retry
		if resp.JobStatus == "IN_PROGRESS" && retryNumber < maxRetries {
			retryNumber++
			plugin.Logger(ctx).Debug("GetServiceLastAccessedDetails in progress", "retryNumber", retryNumber)
			time.Sleep(retryIntervalMs * time.Millisecond)
			continue
		}

		// Stream results
		for _, serviceLastAccessed := range resp.ServicesLastAccessed {
			d.StreamListItem(ctx, &awsIamAccessAdvisorData{
				PrincipalArn:               principalArn,
				Granularity:                granularity,
				LastAuthenticated:          serviceLastAccessed.LastAuthenticated,
				LastAuthenticatedEntity:    serviceLastAccessed.LastAuthenticatedEntity,
				LastAuthenticatedRegion:    serviceLastAccessed.LastAuthenticatedRegion,
				ServiceName:                serviceLastAccessed.ServiceName,
				ServiceNamespace:           serviceLastAccessed.ServiceNamespace,
				TotalAuthenticatedEntities: aws.Int64(int64(*serviceLastAccessed.TotalAuthenticatedEntities)),
				TrackedActionsLastAccessed: serviceLastAccessed.TrackedActionsLastAccessed,
			})

			// Context may get cancelled due to manual cancellation or if the limit has been reached
			if d.RowsRemaining(ctx) == 0 {
				break
			}
		}
		if !resp.IsTruncated {
			break
		}
		params.Marker = resp.Marker
	}
	return nil, nil
}
