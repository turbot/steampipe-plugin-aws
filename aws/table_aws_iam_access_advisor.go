package aws

import (
	"context"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/turbot/steampipe-plugin-sdk/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/plugin"
	"github.com/turbot/steampipe-plugin-sdk/plugin/transform"
	"time"
)

type awsIamAccessAdvisorData struct {
	PrincipalArn               string
	Granularity                string
	LastAuthenticated          time.Time
	LastAuthenticatedEntity    string
	LastAuthenticatedRegion    string
	LastAccessedEntity         string
	LastAccessedRegion         string
	LastAccessedTime           time.Time
	ServiceName                string
	ServiceNamespace           string
	ActionName                 string
	TotalAuthenticatedEntities int64
}

func tableAwsIamAccessAdvisor(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_iam_access_advisor",
		Description: "AWS IAM Access Advisor",
		List: &plugin.ListConfig{
			KeyColumns: plugin.AllColumns([]string{"principal_arn", "granularity"}),
			Hydrate:    listAccessAdvisor,
		},
		Columns: awsColumns([]*plugin.Column{
			{
				Name:        "principal_arn",
				Description: "The ARN of the IAM resource (user, group, role, or managed policy) used to generate information about when the resource was last used in an attempt to access an AWS service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "granularity",
				Description: "Specifies the granularity of the report as either SERVICE_LEVEL or ACTION_LEVEL",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "last_authenticated",
				Description: "The date and time when an authenticated entity most recently attempted to access the service. AWS does not report unauthenticated requests.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "last_authenticated_entity",
				Description: "The ARN of the authenticated entity (user or role) that last attempted to access the service. AWS does not report unauthenticated requests.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "last_authenticated_region",
				Description: "The Region from which the authenticated entity (user or role) last attempted to access the service. AWS does not report unauthenticated requests.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "service_name",
				Description: "The name of the service in which access was attempted.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "service_namespace",
				Description: "The namespace of the service in which access was attempted.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "total_authenticated_entities",
				Description: "The total number of authenticated principals (root user, IAM users, or IAM roles) that have attempted to access the service.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "action_name",
				Description: "The name of the tracked action to which access was attempted. Tracked actions are actions that report activity to IAM.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "last_accessed_entity",
				Description: "The Amazon Resource Name (ARN).",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "last_accessed_region",
				Description: "The Region from which the authenticated entity (user or role) last attempted to access the tracked action. AWS does not report unauthenticated requests.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromGo(),
			},
			{
				Name:        "last_accessed_time",
				Description: "The date and time when an authenticated entity most recently attempted to access the tracked service. AWS does not report unauthenticated requests.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromGo(),
			},
		}),
	}
}

//// LIST FUNCTION

func listAccessAdvisor(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	logger.Trace("listAccessAdvisor")

	principalArn := d.KeyColumnQuals["principal_arn"].GetStringValue()
	granularity := d.KeyColumnQuals["granularity"].GetStringValue()

	logger.Info("Quals", "principal_arn", principalArn, "granularity", granularity)

	// Create Session
	svc, err := IAMService(ctx, d.ConnectionManager)
	if err != nil {
		return nil, err
	}

	generateResp, err := svc.GenerateServiceLastAccessedDetails(&iam.GenerateServiceLastAccessedDetailsInput{Arn: &principalArn, Granularity: &granularity})
	if err != nil {
		return nil, err
	}

	logger.Info("Generate", "jobId", &generateResp.JobId, "resp", &generateResp)

	time.Sleep(5 * time.Second)

	//jobId := "82bbb932-c67e-a664-8e78-cbc05b4e4ad7"

	resp, err := svc.GetServiceLastAccessedDetails(&iam.GetServiceLastAccessedDetailsInput{JobId: generateResp.JobId})
	if err != nil {
		return nil, err
	}
	logger.Info("Details", "jobId", &generateResp.JobId, "status", &resp.JobStatus, "resp", &resp)

	for _, serviceLastAccessed := range resp.ServicesLastAccessed {
		for _, trackedActionLastAccessed := range serviceLastAccessed.TrackedActionsLastAccessed {
			d.StreamListItem(ctx, &awsIamAccessAdvisorData{
				PrincipalArn:               principalArn,
				Granularity:                granularity,
				LastAuthenticated:          *serviceLastAccessed.LastAuthenticated,
				LastAuthenticatedEntity:    *serviceLastAccessed.LastAuthenticatedEntity,
				LastAuthenticatedRegion:    *serviceLastAccessed.LastAuthenticatedRegion,
				LastAccessedEntity:         *trackedActionLastAccessed.LastAccessedEntity,
				LastAccessedRegion:         *trackedActionLastAccessed.LastAccessedRegion,
				LastAccessedTime:           *trackedActionLastAccessed.LastAccessedTime,
				ServiceName:                *serviceLastAccessed.ServiceName,
				ServiceNamespace:           *serviceLastAccessed.ServiceNamespace,
				ActionName:                 *trackedActionLastAccessed.ActionName,
				TotalAuthenticatedEntities: *serviceLastAccessed.TotalAuthenticatedEntities,
			})
		}
	}

	return nil, nil
}
