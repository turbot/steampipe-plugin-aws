package aws

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice"
	"github.com/aws/aws-sdk-go-v2/service/directoryservice/types"
	"github.com/aws/smithy-go"

	directoryservicev1 "github.com/aws/aws-sdk-go/service/directoryservice"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsDirectoryServiceCertificate(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_directory_service_certificate",
		Description: "AWS Directory Service Certificate",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.AllColumns([]string{"directory_id", "certificate_id"}),
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"CertificateDoesNotExistException", "DirectoryDoesNotExistException", "InvalidParameterException"}),
			},
			Hydrate: getDirectoryServiceCertificate,
			Tags:    map[string]string{"service": "directoryservice", "action": "DescribeCertificate"},
		},
		List: &plugin.ListConfig{
			ParentHydrate: listDirectoryServiceDirectories,
			Hydrate:       listDirectoryServiceCertificates,
			Tags:          map[string]string{"service": "directoryservice", "action": "ListCertificates"},
			IgnoreConfig: &plugin.IgnoreConfig{
				ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{"DirectoryDoesNotExistException"}),
			},
			KeyColumns: []*plugin.KeyColumn{
				{
					Name:    "directory_id",
					Require: plugin.Optional,
				},
			},
		},
		HydrateConfig: []plugin.HydrateConfig{
			{
				Func:		getDirectoryServiceCertificate,
				Tags:    map[string]string{"service": "directoryservice", "action": "DescribeCertificate"},
			},
		},
		GetMatrixItemFunc: SupportedRegionMatrix(directoryservicev1.EndpointsID),
		Columns: awsRegionalColumns([]*plugin.Column{
			{
				Name:        "directory_id",
				Description: "The directory identifier.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "certificate_id",
				Description: "The identifier of the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "common_name",
				Description: "The common name for the certificate.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "type",
				Description: "The function that the registered certificate performs. Valid values include ClientLDAPS or ClientCertAuth. The default value is ClientLDAPS.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state",
				Description: "The state of the certificate. Valid values: Registering | Registered | RegisterFailed | Deregistering | Deregistered | DeregisterFailed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "expiry_date_time",
				Description: "The date and time when the certificate will expire.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "registered_date_time",
				Description: "The date and time that the certificate was registered.",
				Type:        proto.ColumnType_TIMESTAMP,
				Hydrate:     getDirectoryServiceCertificate,
			},
			{
				Name:        "state_reason",
				Description: "Describes a state change for the certificate.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getDirectoryServiceCertificate,
			},
			{
				Name:        "client_cert_auth_settings",
				Description: "A ClientCertAuthSettings object that contains client certificate authentication settings.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getDirectoryServiceCertificate,
			},

			// Steampipe Standard columns
			{
				Name:        "title",
				Description: resourceInterfaceDescription("title"),
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("CommonName"),
			},
		}),
	}
}

type CertInfo struct {
	DirectoryId            *string
	CertificateId          *string
	ClientCertAuthSettings *types.ClientCertAuthSettings
	CommonName             *string
	ExpiryDateTime         *time.Time
	RegisteredDateTime     *time.Time
	State                  types.CertificateState
	StateReason            *string
	Type                   types.CertificateType
}

//// LIST FUNCTION

func listDirectoryServiceCertificates(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	directory := h.Item.(types.DirectoryDescription)

	// Restrict the API call for other certificates.
	if d.EqualsQualString("directory_id") != "" && d.EqualsQualString("directory_id") != *directory.DirectoryId {
		return nil, nil
	}

	// Create Session
	svc, err := DirectoryServiceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_certificate.listDirectoryServiceCertificates", "connection_error", err)
		return nil, err
	}

	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	// Limiting the results
	maxLimit := int32(50)
	if d.QueryContext.Limit != nil {
		limit := int32(*d.QueryContext.Limit)
		if limit < maxLimit {
			maxLimit = limit
		}
	}

	// Build the params
	input := &directoryservice.ListCertificatesInput{
		DirectoryId: directory.DirectoryId,
		Limit:       aws.Int32(maxLimit),
	}

	pagesLeft := true
	// List call
	for pagesLeft {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		result, err := svc.ListCertificates(ctx, input)
		if err != nil {
			// In the case of parent hydrate the ignore config seems to not work fine. So we need to handle it manually
			// operation error Directory Service: ListCertificates, https response error StatusCode: 400, RequestID: 6238d084-f28d-42a7-876a-684b0ec0d999, UnsupportedOperationException: LDAPS operations are not supported for this Directory Type. : RequestId: 6238d084-f28d-42a7-876a-684b0ec0d999
			var ae smithy.APIError
			if errors.As(err, &ae) {
				if ae.ErrorCode() == "UnsupportedOperationException" {
					return nil, nil
				}
			}
			plugin.Logger(ctx).Error("aws_directory_service_certificate.listDirectoryServiceCertificates", "api_error", err)
			return nil, err
		}

		for _, cert := range result.CertificatesInfo {

			d.StreamListItem(ctx, CertInfo{
				DirectoryId:    directory.DirectoryId,
				CertificateId:  cert.CertificateId,
				CommonName:     cert.CommonName,
				ExpiryDateTime: cert.ExpiryDateTime,
				State:          cert.State,
				Type:           cert.Type,
			})

			// Context can be cancelled due to manual cancellation or the limit has been hit
			if d.RowsRemaining(ctx) == 0 {
				pagesLeft = false
			}
		}

		if result.NextToken != nil {
			input.NextToken = result.NextToken
		} else {
			pagesLeft = false
		}
	}

	return nil, err
}

//// HYDRATE FUNCTIONS

func getDirectoryServiceCertificate(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	var directoryID, certID string

	if h.Item != nil {
		data := h.Item.(CertInfo)
		directoryID = *data.DirectoryId
		certID = *data.CertificateId
	} else {
		directoryID = d.EqualsQualString("directory_id")
		certID = d.EqualsQualString("certificate_id")
	}

	if directoryID == "" || certID == "" {
		return nil, nil
	}

	// Create service
	svc, err := DirectoryServiceClient(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_certificate.getDirectoryServiceCertificate", "connection_error", err)
		return nil, err
	}
	if svc == nil {
		// Unsupported region, return no data
		return nil, nil
	}

	params := &directoryservice.DescribeCertificateInput{
		DirectoryId:   &directoryID,
		CertificateId: &certID,
	}

	op, err := svc.DescribeCertificate(ctx, params)
	if err != nil {
		plugin.Logger(ctx).Error("aws_directory_service_certificate.getDirectoryServiceCertificate", "api_error", err)
		return nil, err
	}

	if op.Certificate != nil {
		return CertInfo{
			DirectoryId:            &directoryID,
			CertificateId:          &certID,
			ClientCertAuthSettings: op.Certificate.ClientCertAuthSettings,
			CommonName:             op.Certificate.CommonName,
			ExpiryDateTime:         op.Certificate.ExpiryDateTime,
			RegisteredDateTime:     op.Certificate.RegisteredDateTime,
			State:                  op.Certificate.State,
			StateReason:            op.Certificate.StateReason,
			Type:                   op.Certificate.Type,
		}, nil
	}
	return nil, nil
}
