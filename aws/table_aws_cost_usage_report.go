package aws

import (
	"compress/gzip"
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	filehelpers "github.com/turbot/go-kit/files"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableAwsCostUsageReport(_ context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "aws_cost_usage_report",
		Description: "Cost & Usage reports available in local or S3 bucket.",
		List: &plugin.ListConfig{
			Hydrate: listCostUsageReports,
		},
		Columns: awsAccountColumns([]*plugin.Column{
			{
				Name:        "identity",
				Description: "Contains data about the entity that made or was targeted by the AWS request, including LineItemId, TimeInterval, and others.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "bill",
				Description: "Details related to AWS billing and associated entities.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "line_item",
				Description: "Specific line items pertaining to AWS services and their usage.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "product",
				Description: "Information about the AWS product or service in question.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("CUProduct"),
			},
			{
				Name:        "pricing",
				Description: "Details about the pricing associated with the AWS service or product.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "reservation",
				Description: "Details of the reservation, usually associated with AWS reserved instances.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "savings_plan",
				Description: "Information about the AWS savings plan, including duration, costs, and other related details.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "resource_tags",
				Description: "Tags associated with AWS resources, which aid in categorizing and managing them.",
				Type:        proto.ColumnType_JSON,
			},
		}),
	}
}

func listCostUsageReports(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
	records, err := getCostUsageReportContent(ctx, d)
	if err != nil {
		plugin.Logger(ctx).Error("listCostUsageReports", "getCostUsageReportContent", err)
		return nil, err
	}

	for _, record := range records {
		d.StreamListItem(ctx, record)

		// Context can be cancelled due to manual cancellation or the limit has been hit
		if d.RowsRemaining(ctx) == 0 {
			return nil, nil
		}
	}

	return nil, nil
}

// Get the report from the given files.
func getCostUsageReportContent(ctx context.Context, d *plugin.QueryData) ([]CSVRecord, error) {
	conn, err := costUsageReportContentCached(ctx, d, nil)
	if err != nil {
		return nil, err
	}

	return conn.([]CSVRecord), nil
}

// Cached form of the report.
var costUsageReportContentCached = plugin.HydrateFunc(costUsageReportContentUncached).Memoize()

func costUsageReportContentUncached(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (any, error) {

	// Read the config
	resolvedPaths, err := resolveCostUsageReportPaths(ctx, d)
	if err != nil {
		return nil, err
	}

	var recordLogs []CSVRecord
	for _, path := range resolvedPaths {
		var rec []CSVRecord
		var csvFile io.Reader

		content, err := os.Open(path)
		if err != nil {
			plugin.Logger(ctx).Error("costUsageReportContentUncached", "failed to read file", err, "path", path)
			return nil, err
		}

		// check if the file is compressed
		if strings.HasSuffix(path, ".gz") {
			// Create a Gzip reader for the file
			gzipReader, err := gzip.NewReader(content)
			if err != nil {
				plugin.Logger(ctx).Error("costUsageReportContentUncached", "failed to uncompressed", err)
				return nil, err
			}
			csvFile = gzipReader
		} else {
			csvFile = content
		}

		if err := gocsv.Unmarshal(csvFile, &rec); err != nil {
			plugin.Logger(ctx).Error("costUsageReportContentUncached", "UnmarshalFile", err)
		}
		recordLogs = append(recordLogs, rec...)
	}

	return recordLogs, nil
}

func resolveCostUsageReportPaths(ctx context.Context, d *plugin.QueryData) ([]string, error) {
	config := GetConfig(d.Connection)

	if config.CostUsageReportPaths == nil {
		return nil, errors.New("cost_usage_report_paths must be set in the config")
	}

	// Gather file path matches for the glob
	var matches, resolvedPaths []string
	paths := config.CostUsageReportPaths
	for _, i := range paths {
		// List the files in the given source directory
		files, err := d.GetSourceFiles(i)
		if err != nil {
			return nil, err
		}
		matches = append(matches, files...)
	}

	// Sanitize the matches to ignore the directories
	for _, i := range matches {

		// Ignore directories
		if filehelpers.DirectoryExists(i) {
			continue
		}
		resolvedPaths = append(resolvedPaths, i)
	}

	return resolvedPaths, nil
}

type CSVRecord struct {
	Identity
	Bill
	LineItem
	CUProduct
	Pricing
	Reservation
	SavingsPlan
	ResourceTags
}

type Identity struct {
	LineItemId   string `csv:"identity/LineItemId"`
	TimeInterval string `csv:"identity/TimeInterval"`
}

type Bill struct {
	InvoiceId              string `csv:"bill/InvoiceId"`
	InvoicingEntity        string `csv:"bill/InvoicingEntity"`
	BillingEntity          string `csv:"bill/BillingEntity"`
	BillType               string `csv:"bill/BillType"`
	PayerAccountId         string `csv:"bill/PayerAccountId"`
	BillingPeriodStartDate string `csv:"bill/BillingPeriodStartDate"`
	BillingPeriodEndDate   string `csv:"bill/BillingPeriodEndDate"`
}

type LineItem struct {
	UsageAccountId        string `csv:"lineItem/UsageAccountId"`
	LineItemType          string `csv:"lineItem/LineItemType"`
	UsageStartDate        string `csv:"lineItem/UsageStartDate"`
	UsageEndDate          string `csv:"lineItem/UsageEndDate"`
	ProductCode           string `csv:"lineItem/ProductCode"`
	UsageType             string `csv:"lineItem/UsageType"`
	Operation             string `csv:"lineItem/Operation"`
	AvailabilityZone      string `csv:"lineItem/AvailabilityZone"`
	ResourceId            string `csv:"lineItem/ResourceId"`
	UsageAmount           string `csv:"lineItem/UsageAmount"`
	NormalizationFactor   string `csv:"lineItem/NormalizationFactor"`
	NormalizedUsageAmount string `csv:"lineItem/NormalizedUsageAmount"`
	CurrencyCode          string `csv:"lineItem/CurrencyCode"`
	UnblendedRate         string `csv:"lineItem/UnblendedRate"`
	UnblendedCost         string `csv:"lineItem/UnblendedCost"`
	BlendedRate           string `csv:"lineItem/BlendedRate"`
	BlendedCost           string `csv:"lineItem/BlendedCost"`
	LineItemDescription   string `csv:"lineItem/LineItemDescription"`
	TaxType               string `csv:"lineItem/TaxType"`
	LegalEntity           string `csv:"lineItem/LegalEntity"`
}

type CUProduct struct {
	ProductName                      string `csv:"product/ProductName"`
	SizeFlex                         string `csv:"product/SizeFlex"`
	AbdInstanceClass                 string `csv:"product/abdInstanceClass"`
	AlarmType                        string `csv:"product/alarmType"`
	AttachmentType                   string `csv:"product/attachmentType"`
	Availability                     string `csv:"product/availability"`
	AvailabilityZone                 string `csv:"product/availabilityZone"`
	BackupService                    string `csv:"product/backupservice"`
	BrokerEngine                     string `csv:"product/brokerEngine"`
	Bundle                           string `csv:"product/bundle"`
	BundleDescription                string `csv:"product/bundleDescription"`
	BundleGroup                      string `csv:"product/bundleGroup"`
	CacheEngine                      string `csv:"product/cacheEngine"`
	CacheType                        string `csv:"product/cacheType"`
	CacheMemorySize                  string `csv:"product/cachememorysize"`
	CapacityStatus                   string `csv:"product/capacitystatus"`
	Category                         string `csv:"product/category"`
	CiType                           string `csv:"product/ciType"`
	ClassicNetworkingSupport         string `csv:"product/classicnetworkingsupport"`
	ClockSpeed                       string `csv:"product/clockSpeed"`
	Component                        string `csv:"product/component"`
	ComputeFamily                    string `csv:"product/computeFamily"`
	ComputeType                      string `csv:"product/computeType"`
	ContentType                      string `csv:"product/contentType"`
	CpuType                          string `csv:"product/cputype"`
	CurrentGeneration                string `csv:"product/currentGeneration"`
	DatabaseEdition                  string `csv:"product/databaseEdition"`
	DatabaseEngine                   string `csv:"product/databaseEngine"`
	DataTransferOut                  string `csv:"product/datatransferout"`
	DedicatedEbsThroughput           string `csv:"product/dedicatedEbsThroughput"`
	DeploymentOption                 string `csv:"product/deploymentOption"`
	Description                      string `csv:"product/description"`
	DestinationCountryIsoCode        string `csv:"product/destinationCountryIsoCode"`
	DirectorySize                    string `csv:"product/directorySize"`
	DirectoryType                    string `csv:"product/directoryType"`
	DirectoryTypeDescription         string `csv:"product/directoryTypeDescription"`
	Durability                       string `csv:"product/durability"`
	Ecu                              string `csv:"product/ecu"`
	Edition                          string `csv:"product/edition"`
	EndpointType                     string `csv:"product/endpointType"`
	EngineCode                       string `csv:"product/engineCode"`
	EnhancedNetworkingSupport        string `csv:"product/enhancedNetworkingSupport"`
	EnhancedNetworkingSupported      string `csv:"product/enhancedNetworkingSupported"`
	EquivalentOnDemandSku            string `csv:"product/equivalentondemandsku"`
	EventType                        string `csv:"product/eventType"`
	FeeCode                          string `csv:"product/feeCode"`
	FeeDescription                   string `csv:"product/feeDescription"`
	FileSystemType                   string `csv:"product/fileSystemType"`
	FindingGroup                     string `csv:"product/findingGroup"`
	FindingSource                    string `csv:"product/findingSource"`
	FindingStorage                   string `csv:"product/findingStorage"`
	FreeUsageIncluded                string `csv:"product/freeUsageIncluded"`
	FromLocation                     string `csv:"product/fromLocation"`
	FromLocationType                 string `csv:"product/fromLocationType"`
	FromRegionCode                   string `csv:"product/fromRegionCode"`
	Gpu                              string `csv:"product/gpu"`
	GpuMemory                        string `csv:"product/gpuMemory"`
	GraphqlOperation                 string `csv:"product/graphqloperation"`
	Group                            string `csv:"product/group"`
	GroupDescription                 string `csv:"product/groupDescription"`
	InsightsType                     string `csv:"product/insightstype"`
	Instance                         string `csv:"product/instance"`
	InstanceFamily                   string `csv:"product/instanceFamily"`
	InstanceName                     string `csv:"product/instanceName"`
	InstanceType                     string `csv:"product/instanceType"`
	InstanceTypeFamily               string `csv:"product/instanceTypeFamily"`
	IntelAvx2Available               string `csv:"product/intelAvx2Available"`
	IntelAvxAvailable                string `csv:"product/intelAvxAvailable"`
	IntelTurboAvailable              string `csv:"product/intelTurboAvailable"`
	Invocation                       string `csv:"product/invocation"`
	Io                               string `csv:"product/io"`
	IoRequestType                    string `csv:"product/ioRequestType"`
	License                          string `csv:"product/license"`
	LicenseModel                     string `csv:"product/licenseModel"`
	Location                         string `csv:"product/location"`
	LocationType                     string `csv:"product/locationType"`
	LogsDestination                  string `csv:"product/logsDestination"`
	MarketOption                     string `csv:"product/marketoption"`
	MaxIopsBurstPerformance          string `csv:"product/maxIopsBurstPerformance"`
	MaxIopsVolume                    string `csv:"product/maxIopsvolume"`
	MaxThroughputVolume              string `csv:"product/maxThroughputvolume"`
	MaxVolumeSize                    string `csv:"product/maxVolumeSize"`
	MaximumExtendedStorage           string `csv:"product/maximumExtendedStorage"`
	Memory                           string `csv:"product/memory"`
	MemoryGib                        string `csv:"product/memoryGib"`
	MemoryType                       string `csv:"product/memorytype"`
	MessageCountFee                  string `csv:"product/messageCountfee"`
	MessageDeliveryFrequency         string `csv:"product/messageDeliveryFrequency"`
	MessageDeliveryOrder             string `csv:"product/messageDeliveryOrder"`
	MessageType                      string `csv:"product/messageType"`
	MeteringType                     string `csv:"product/meteringType"`
	MinVolumeSize                    string `csv:"product/minVolumeSize"`
	NetworkPerformance               string `csv:"product/networkPerformance"`
	NormalizationSizeFactor          string `csv:"product/normalizationSizeFactor"`
	OperatingSystem                  string `csv:"product/operatingSystem"`
	Operation                        string `csv:"product/operation"`
	Origin                           string `csv:"product/origin"`
	OriginationIdType                string `csv:"product/originationIdType"`
	Overhead                         string `csv:"product/overhead"`
	PackSize                         string `csv:"product/packSize"`
	ParameterType                    string `csv:"product/parameterType"`
	PhysicalCpu                      string `csv:"product/physicalCpu"`
	PhysicalGpu                      string `csv:"product/physicalGpu"`
	PhysicalProcessor                string `csv:"product/physicalProcessor"`
	PlatoClassificationType          string `csv:"product/platoclassificationtype"`
	PlatoInstanceName                string `csv:"product/platoinstancename"`
	PlatoInstanceType                string `csv:"product/platoinstancetype"`
	PlatoPricingType                 string `csv:"product/platopricingtype"`
	PlatoStorageName                 string `csv:"product/platostoragename"`
	PlatoStorageType                 string `csv:"product/platostoragetype"`
	PlatoUsageType                   string `csv:"product/platousagetype"`
	PlatoVolumeType                  string `csv:"product/platovolumetype"`
	PreInstalledSw                   string `csv:"product/preInstalledSw"`
	PricingUnit                      string `csv:"product/pricingUnit"`
	ProcessorArchitecture            string `csv:"product/processorArchitecture"`
	ProcessorFeatures                string `csv:"product/processorFeatures"`
	ProductFamily                    string `csv:"product/productFamily"`
	Protocol                         string `csv:"product/protocol"`
	Provisioned                      string `csv:"product/provisioned"`
	QPresent                         string `csv:"product/qPresent"`
	QueueType                        string `csv:"product/queueType"`
	Recipient                        string `csv:"product/recipient"`
	Region                           string `csv:"product/region"`
	RegionCode                       string `csv:"product/regionCode"`
	RequestDescription               string `csv:"product/requestDescription"`
	RequestType                      string `csv:"product/requestType"`
	ResourceAssessment               string `csv:"product/resourceAssessment"`
	ResourceEndpoint                 string `csv:"product/resourceEndpoint"`
	ResourcePriceGroup               string `csv:"product/resourcePriceGroup"`
	ResourceType                     string `csv:"product/resourceType"`
	RootVolume                       string `csv:"product/rootvolume"`
	RouteType                        string `csv:"product/routeType"`
	RoutingTarget                    string `csv:"product/routingTarget"`
	RoutingType                      string `csv:"product/routingType"`
	RunningMode                      string `csv:"product/runningMode"`
	ScanType                         string `csv:"product/scanType"`
	ServiceCode                      string `csv:"product/servicecode"`
	ServiceName                      string `csv:"product/servicename"`
	Sku                              string `csv:"product/sku"`
	SoftwareIncluded                 string `csv:"product/softwareIncluded"`
	StandardGroup                    string `csv:"product/standardGroup"`
	StandardStorage                  string `csv:"product/standardStorage"`
	StandardStorageRetentionIncluded string `csv:"product/standardStorageRetentionIncluded"`
	Steps                            string `csv:"product/steps"`
	Storage                          string `csv:"product/storage"`
	StorageClass                     string `csv:"product/storageClass"`
	StorageFamily                    string `csv:"product/storageFamily"`
	StorageMedia                     string `csv:"product/storageMedia"`
	StorageType                      string `csv:"product/storageType"`
	SubscriptionType                 string `csv:"product/subscriptionType"`
	Tenancy                          string `csv:"product/tenancy"`
	Throughput                       string `csv:"product/throughput"`
	ThroughputCapacity               string `csv:"product/throughputCapacity"`
	ThroughputClass                  string `csv:"product/throughputClass"`
	TierType                         string `csv:"product/tiertype"`
	TimeWindow                       string `csv:"product/timeWindow"`
	ToLocation                       string `csv:"product/toLocation"`
	ToLocationType                   string `csv:"product/toLocationType"`
	ToRegionCode                     string `csv:"product/toRegionCode"`
	TransferType                     string `csv:"product/transferType"`
	Type                             string `csv:"product/type"`
	UsageFamily                      string `csv:"product/usageFamily"`
	UsageGroup                       string `csv:"product/usageGroup"`
	UsageVolume                      string `csv:"product/usageVolume"`
	UsageType                        string `csv:"product/usagetype"`
	UserVolume                       string `csv:"product/uservolume"`
	Vcpu                             string `csv:"product/vcpu"`
	Version                          string `csv:"product/version"`
	VolumeApiName                    string `csv:"product/volumeApiName"`
	VolumeType                       string `csv:"product/volumeType"`
	VpcNetworkingSupport             string `csv:"product/vpcnetworkingsupport"`
	WithActiveUsers                  string `csv:"product/withActiveUsers"`
}

type Pricing struct {
	RateCode           string `csv:"pricing/RateCode"`
	RateId             string `csv:"pricing/RateId"`
	Currency           string `csv:"pricing/currency"`
	PublicOnDemandCost string `csv:"pricing/publicOnDemandCost"`
	PublicOnDemandRate string `csv:"pricing/publicOnDemandRate"`
	Term               string `csv:"pricing/term"`
	Unit               string `csv:"pricing/unit"`
}

type Reservation struct {
	AmortizedUpfrontCostForUsage              string `csv:"reservation/AmortizedUpfrontCostForUsage"`
	AmortizedUpfrontFeeForBillingPeriod       string `csv:"reservation/AmortizedUpfrontFeeForBillingPeriod"`
	EffectiveCost                             string `csv:"reservation/EffectiveCost"`
	EndTime                                   string `csv:"reservation/EndTime"`
	ModificationStatus                        string `csv:"reservation/ModificationStatus"`
	NormalizedUnitsPerReservation             string `csv:"reservation/NormalizedUnitsPerReservation"`
	NumberOfReservations                      string `csv:"reservation/NumberOfReservations"`
	RecurringFeeForUsage                      string `csv:"reservation/RecurringFeeForUsage"`
	StartTime                                 string `csv:"reservation/StartTime"`
	SubscriptionId                            string `csv:"reservation/SubscriptionId"`
	TotalReservedNormalizedUnits              string `csv:"reservation/TotalReservedNormalizedUnits"`
	TotalReservedUnits                        string `csv:"reservation/TotalReservedUnits"`
	UnitsPerReservation                       string `csv:"reservation/UnitsPerReservation"`
	UnusedAmortizedUpfrontFeeForBillingPeriod string `csv:"reservation/UnusedAmortizedUpfrontFeeForBillingPeriod"`
	UnusedNormalizedUnitQuantity              string `csv:"reservation/UnusedNormalizedUnitQuantity"`
	UnusedQuantity                            string `csv:"reservation/UnusedQuantity"`
	UnusedRecurringFee                        string `csv:"reservation/UnusedRecurringFee"`
	UpfrontValue                              string `csv:"reservation/UpfrontValue"`
}

type SavingsPlan struct {
	TotalCommitmentToDate                      string `csv:"savingsPlan/TotalCommitmentToDate"`
	SavingsPlanARN                             string `csv:"savingsPlan/SavingsPlanARN"`
	SavingsPlanRate                            string `csv:"savingsPlan/SavingsPlanRate"`
	UsedCommitment                             string `csv:"savingsPlan/UsedCommitment"`
	SavingsPlanEffectiveCost                   string `csv:"savingsPlan/SavingsPlanEffectiveCost"`
	AmortizedUpfrontCommitmentForBillingPeriod string `csv:"savingsPlan/AmortizedUpfrontCommitmentForBillingPeriod"`
	RecurringCommitmentForBillingPeriod        string `csv:"savingsPlan/RecurringCommitmentForBillingPeriod"`
}

type ResourceTags struct {
	CloudFormationStackName string `csv:"resourceTags/aws:cloudformation:stack-name"`
	CreatedBy               string `csv:"resourceTags/aws:createdBy"`
	CostCenter              string `csv:"resourceTags/user:Cost Center"`
	CreatedByActor          string `csv:"resourceTags/user:CreatedByActor"`
	Department              string `csv:"resourceTags/user:Department"`
	Name                    string `csv:"resourceTags/user:Name"`
	Owner                   string `csv:"resourceTags/user:Owner"`
}
