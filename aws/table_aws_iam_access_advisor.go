package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
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
		},
		Columns: awsColumns([]*plugin.Column{
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

func listAccessAdvisor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAccessAdvisor")

	// To simplify the table we always get ACTION_LEVEL.  ACTION_LEVEL is a superset of
	// SERVICE_LEVEL, and currenty only s# supports action level actions anyway, so the
	// performance impact is minimal
	granularity := "ACTION_LEVEL"
	principalArn := d.KeyColumnQuals["principal_arn"].GetStringValue()

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	// Generate the details.  We'll need the job id of this to get the details...
	generateResp, err := svc.GenerateServiceLastAccessedDetails(&iam.GenerateServiceLastAccessedDetailsInput{Arn: &principalArn, Granularity: &granularity})
	if err != nil {
		return nil, err
	}
	logger.Debug("listAccessAdvisor generateResp", "jobId", *generateResp.JobId, "resp", *generateResp)

	params := &iam.GetServiceLastAccessedDetailsInput{
		JobId: generateResp.JobId,
	}
	retryNumber := 0
	for true {
		resp, err := svc.GetServiceLastAccessedDetails(params)
		if err != nil {
			return nil, err
		}
		logger.Debug("listAccessAdvisor Details", "jobId", *generateResp.JobId, "status", *resp.JobStatus, "resp", *resp)

		// if job is still in progress, wait and retry
		if *resp.JobStatus == "IN_PROGRESS" && retryNumber < maxRetries {
			retryNumber++
			logger.Debug("GetServiceLastAccessedDetails in progress", "retryNumber", retryNumber)
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
				TotalAuthenticatedEntities: serviceLastAccessed.TotalAuthenticatedEntities,
				TrackedActionsLastAccessed: serviceLastAccessed.TrackedActionsLastAccessed,
			})
		}
		if !*resp.IsTruncated {
			break
		}
		params.Marker = resp.Marker
	}
	return nil, nil
}
