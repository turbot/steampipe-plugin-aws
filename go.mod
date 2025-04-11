module github.com/turbot/steampipe-plugin-aws

go 1.23.1

toolchain go1.23.2

require (
	github.com/aws/aws-sdk-go v1.51.19
	github.com/aws/aws-sdk-go-v2 v1.27.0
	github.com/aws/aws-sdk-go-v2/config v1.27.16
	github.com/aws/aws-sdk-go-v2/credentials v1.17.16
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.16.21
	github.com/aws/aws-sdk-go-v2/service/accessanalyzer v1.29.1
	github.com/aws/aws-sdk-go-v2/service/account v1.16.4
	github.com/aws/aws-sdk-go-v2/service/acm v1.25.4
	github.com/aws/aws-sdk-go-v2/service/acmpca v1.29.4
	github.com/aws/aws-sdk-go-v2/service/amplify v1.21.5
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.23.6
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.20.4
	github.com/aws/aws-sdk-go-v2/service/appconfig v1.29.2
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v1.27.4
	github.com/aws/aws-sdk-go-v2/service/apprunner v1.28.8
	github.com/aws/aws-sdk-go-v2/service/appstream v1.34.4
	github.com/aws/aws-sdk-go-v2/service/appsync v1.31.4
	github.com/aws/aws-sdk-go-v2/service/athena v1.40.4
	github.com/aws/aws-sdk-go-v2/service/auditmanager v1.32.4
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.40.5
	github.com/aws/aws-sdk-go-v2/service/backup v1.34.2
	github.com/aws/aws-sdk-go-v2/service/cloudcontrol v1.18.4
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.49.0
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.35.4
	github.com/aws/aws-sdk-go-v2/service/cloudsearch v1.22.4
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.39.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.37.0
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.35.1
	github.com/aws/aws-sdk-go-v2/service/codeartifact v1.26.1
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.34.0
	github.com/aws/aws-sdk-go-v2/service/codecommit v1.22.4
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.25.4
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.26.4
	github.com/aws/aws-sdk-go-v2/service/codestarnotifications v1.22.7
	github.com/aws/aws-sdk-go-v2/service/cognitoidentity v1.23.6
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.36.3
	github.com/aws/aws-sdk-go-v2/service/configservice v1.46.4
	github.com/aws/aws-sdk-go-v2/service/costexplorer v1.37.1
	github.com/aws/aws-sdk-go-v2/service/costoptimizationhub v1.4.7
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v1.38.4
	github.com/aws/aws-sdk-go-v2/service/dax v1.19.4
	github.com/aws/aws-sdk-go-v2/service/directoryservice v1.24.4
	github.com/aws/aws-sdk-go-v2/service/dlm v1.24.4
	github.com/aws/aws-sdk-go-v2/service/docdb v1.34.0
	github.com/aws/aws-sdk-go-v2/service/drs v1.25.3
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.31.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.156.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.27.4
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.23.4
	github.com/aws/aws-sdk-go-v2/service/ecs v1.41.7
	github.com/aws/aws-sdk-go-v2/service/efs v1.28.4
	github.com/aws/aws-sdk-go-v2/service/eks v1.42.1
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.38.1
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.23.4
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.24.4
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.30.5
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.28.4
	github.com/aws/aws-sdk-go-v2/service/emr v1.39.5
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.30.4
	github.com/aws/aws-sdk-go-v2/service/firehose v1.28.5
	github.com/aws/aws-sdk-go-v2/service/fms v1.31.4
	github.com/aws/aws-sdk-go-v2/service/fsx v1.43.4
	github.com/aws/aws-sdk-go-v2/service/glacier v1.22.4
	github.com/aws/aws-sdk-go-v2/service/globalaccelerator v1.23.1
	github.com/aws/aws-sdk-go-v2/service/glue v1.78.0
	github.com/aws/aws-sdk-go-v2/service/guardduty v1.41.1
	github.com/aws/aws-sdk-go-v2/service/health v1.24.4
	github.com/aws/aws-sdk-go-v2/service/iam v1.31.4
	github.com/aws/aws-sdk-go-v2/service/identitystore v1.23.5
	github.com/aws/aws-sdk-go-v2/service/inspector v1.21.4
	github.com/aws/aws-sdk-go-v2/service/inspector2 v1.24.4
	github.com/aws/aws-sdk-go-v2/service/iot v1.53.3
	github.com/aws/aws-sdk-go-v2/service/kafka v1.31.2
	github.com/aws/aws-sdk-go-v2/service/keyspaces v1.10.8
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.27.4
	github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2 v1.25.2
	github.com/aws/aws-sdk-go-v2/service/kinesisvideo v1.23.4
	github.com/aws/aws-sdk-go-v2/service/kms v1.31.0
	github.com/aws/aws-sdk-go-v2/service/lakeformation v1.31.5
	github.com/aws/aws-sdk-go-v2/service/lambda v1.54.0
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.37.0
	github.com/aws/aws-sdk-go-v2/service/macie2 v1.38.4
	github.com/aws/aws-sdk-go-v2/service/mediastore v1.20.4
	github.com/aws/aws-sdk-go-v2/service/memorydb v1.19.8
	github.com/aws/aws-sdk-go-v2/service/mgn v1.28.0
	github.com/aws/aws-sdk-go-v2/service/mq v1.22.4
	github.com/aws/aws-sdk-go-v2/service/neptune v1.31.6
	github.com/aws/aws-sdk-go-v2/service/networkfirewall v1.38.5
	github.com/aws/aws-sdk-go-v2/service/oam v1.10.1
	github.com/aws/aws-sdk-go-v2/service/opensearch v1.32.4
	github.com/aws/aws-sdk-go-v2/service/organizations v1.27.3
	github.com/aws/aws-sdk-go-v2/service/pinpoint v1.29.0
	github.com/aws/aws-sdk-go-v2/service/pipes v1.11.4
	github.com/aws/aws-sdk-go-v2/service/pricing v1.28.1
	github.com/aws/aws-sdk-go-v2/service/quicksight v1.64.2
	github.com/aws/aws-sdk-go-v2/service/ram v1.25.4
	github.com/aws/aws-sdk-go-v2/service/rds v1.77.0
	github.com/aws/aws-sdk-go-v2/service/redshift v1.43.5
	github.com/aws/aws-sdk-go-v2/service/redshiftserverless v1.17.4
	github.com/aws/aws-sdk-go-v2/service/resourceexplorer2 v1.10.5
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.21.4
	github.com/aws/aws-sdk-go-v2/service/route53 v1.40.4
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.23.4
	github.com/aws/aws-sdk-go-v2/service/route53resolver v1.27.4
	github.com/aws/aws-sdk-go-v2/service/s3 v1.54.3
	github.com/aws/aws-sdk-go-v2/service/s3control v1.44.4
	github.com/aws/aws-sdk-go-v2/service/sagemaker v1.135.0
	github.com/aws/aws-sdk-go-v2/service/scheduler v1.8.8
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.28.6
	github.com/aws/aws-sdk-go-v2/service/securityhub v1.47.2
	github.com/aws/aws-sdk-go-v2/service/securitylake v1.13.3
	github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository v1.20.4
	github.com/aws/aws-sdk-go-v2/service/servicecatalog v1.28.4
	github.com/aws/aws-sdk-go-v2/service/servicediscovery v1.29.5
	github.com/aws/aws-sdk-go-v2/service/servicequotas v1.21.4
	github.com/aws/aws-sdk-go-v2/service/ses v1.22.4
	github.com/aws/aws-sdk-go-v2/service/sfn v1.26.4
	github.com/aws/aws-sdk-go-v2/service/shield v1.25.7
	github.com/aws/aws-sdk-go-v2/service/simspaceweaver v1.10.4
	github.com/aws/aws-sdk-go-v2/service/sns v1.29.4
	github.com/aws/aws-sdk-go-v2/service/sqs v1.31.4
	github.com/aws/aws-sdk-go-v2/service/ssm v1.49.5
	github.com/aws/aws-sdk-go-v2/service/ssoadmin v1.25.5
	github.com/aws/aws-sdk-go-v2/service/sts v1.28.10
	github.com/aws/aws-sdk-go-v2/service/support v1.21.4
	github.com/aws/aws-sdk-go-v2/service/timestreamwrite v1.25.9
	github.com/aws/aws-sdk-go-v2/service/transfer v1.45.0
	github.com/aws/aws-sdk-go-v2/service/waf v1.20.4
	github.com/aws/aws-sdk-go-v2/service/wafregional v1.21.4
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.48.2
	github.com/aws/aws-sdk-go-v2/service/wellarchitected v1.29.4
	github.com/aws/aws-sdk-go-v2/service/workspaces v1.38.4
	github.com/aws/smithy-go v1.22.2
	github.com/gocarina/gocsv v0.0.0-20201208093247-67c824bc04d4
	github.com/goccy/go-yaml v1.11.3
	github.com/golang/protobuf v1.5.4
	github.com/hashicorp/go-hclog v1.6.3
	github.com/rs/dnscache v0.0.0-20230804202142-fc85eb664529
	github.com/turbot/go-kit v1.1.0
	github.com/turbot/steampipe-plugin-sdk/v5 v5.11.5
	golang.org/x/sync v0.11.0
	golang.org/x/text v0.22.0
)

require golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.49.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.49.0 // indirect
	golang.org/x/mod v0.19.0 // indirect
	golang.org/x/tools v0.23.0 // indirect
)

require (
	cloud.google.com/go v0.112.1 // indirect
	cloud.google.com/go/compute/metadata v0.3.0 // indirect
	cloud.google.com/go/iam v1.1.6 // indirect
	cloud.google.com/go/storage v1.38.0 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/agext/levenshtein v1.2.3 // indirect
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.2 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.7 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.3.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.9.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.17.7 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssmincidents v1.30.4
	github.com/aws/aws-sdk-go-v2/service/sso v1.20.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.24.3 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/dgraph-io/ristretto v0.2.0 // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/eko/gocache/lib/v4 v4.1.6 // indirect
	github.com/eko/gocache/store/bigcache/v4 v4.2.1 // indirect
	github.com/eko/gocache/store/ristretto/v4 v4.2.1 // indirect
	github.com/fatih/color v1.17.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.2 // indirect
	github.com/googleapis/gax-go/v2 v2.12.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.7.5 // indirect
	github.com/hashicorp/go-plugin v1.6.1 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.7.0 // indirect
	github.com/hashicorp/hcl/v2 v2.20.1 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.17.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0
	github.com/oklog/run v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/stevenle/topsort v0.2.0 // indirect
	github.com/tkrajina/go-reflector v0.5.6 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/zclconf/go-cty v1.14.4 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.24.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.24.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	go.opentelemetry.io/proto/otlp v1.2.0 // indirect
	golang.org/x/crypto v0.35.0 // indirect
	golang.org/x/exp v0.0.0-20240719175910-8a7402abbf56 // indirect
	golang.org/x/net v0.36.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	google.golang.org/api v0.171.0 // indirect
	google.golang.org/genproto v0.0.0-20240227224415-6ceb2ff114de // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240604185151-ef581f913117 // indirect
	google.golang.org/grpc v1.66.0 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
