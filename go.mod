module github.com/turbot/steampipe-plugin-aws

go 1.21

require (
	github.com/aws/aws-sdk-go v1.44.189
	github.com/aws/aws-sdk-go-v2 v1.24.0
	github.com/aws/aws-sdk-go-v2/config v1.18.10
	github.com/aws/aws-sdk-go-v2/credentials v1.13.10
	github.com/aws/aws-sdk-go-v2/service/accessanalyzer v1.19.1
	github.com/aws/aws-sdk-go-v2/service/account v1.8.0
	github.com/aws/aws-sdk-go-v2/service/acm v1.17.1
	github.com/aws/aws-sdk-go-v2/service/amplify v1.13.0
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.16.1
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.13.1
	github.com/aws/aws-sdk-go-v2/service/appconfig v1.15.1
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v1.17.1
	github.com/aws/aws-sdk-go-v2/service/appstream v1.20.5
	github.com/aws/aws-sdk-go-v2/service/athena v1.23.1
	github.com/aws/aws-sdk-go-v2/service/auditmanager v1.23.0
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.26.1
	github.com/aws/aws-sdk-go-v2/service/backup v1.19.1
	github.com/aws/aws-sdk-go-v2/service/cloudcontrol v1.11.1
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.25.1
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.24.0
	github.com/aws/aws-sdk-go-v2/service/cloudsearch v1.14.0
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.22.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.25.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.20.1
	github.com/aws/aws-sdk-go-v2/service/codeartifact v1.15.0
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.20.1
	github.com/aws/aws-sdk-go-v2/service/codecommit v1.14.0
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.16.1
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.14.0
	github.com/aws/aws-sdk-go-v2/service/cognitoidentity v1.15.14
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.23.0
	github.com/aws/aws-sdk-go-v2/service/configservice v1.29.1
	github.com/aws/aws-sdk-go-v2/service/costexplorer v1.25.0
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v1.23.1
	github.com/aws/aws-sdk-go-v2/service/dax v1.12.0
	github.com/aws/aws-sdk-go-v2/service/directoryservice v1.16.0
	github.com/aws/aws-sdk-go-v2/service/dlm v1.15.6
	github.com/aws/aws-sdk-go-v2/service/docdb v1.20.1
	github.com/aws/aws-sdk-go-v2/service/drs v1.10.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.18.1
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.81.0
	github.com/aws/aws-sdk-go-v2/service/ecr v1.18.1
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.15.0
	github.com/aws/aws-sdk-go-v2/service/ecs v1.23.1
	github.com/aws/aws-sdk-go-v2/service/efs v1.19.2
	github.com/aws/aws-sdk-go-v2/service/eks v1.27.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.26.1
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.15.0
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.15.1
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.19.1
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.18.1
	github.com/aws/aws-sdk-go-v2/service/emr v1.22.1
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.18.0
	github.com/aws/aws-sdk-go-v2/service/firehose v1.16.0
	github.com/aws/aws-sdk-go-v2/service/fms v1.24.3
	github.com/aws/aws-sdk-go-v2/service/fsx v1.28.1
	github.com/aws/aws-sdk-go-v2/service/glacier v1.14.0
	github.com/aws/aws-sdk-go-v2/service/globalaccelerator v1.16.0
	github.com/aws/aws-sdk-go-v2/service/glue v1.40.0
	github.com/aws/aws-sdk-go-v2/service/guardduty v1.18.0
	github.com/aws/aws-sdk-go-v2/service/health v1.16.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.19.0
	github.com/aws/aws-sdk-go-v2/service/identitystore v1.16.0
	github.com/aws/aws-sdk-go-v2/service/inspector v1.13.0
	github.com/aws/aws-sdk-go-v2/service/inspector2 v1.11.5
	github.com/aws/aws-sdk-go-v2/service/kafka v1.19.0
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.17.1
	github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2 v1.16.0
	github.com/aws/aws-sdk-go-v2/service/kinesisvideo v1.15.1
	github.com/aws/aws-sdk-go-v2/service/kms v1.20.1
	github.com/aws/aws-sdk-go-v2/service/lambda v1.29.0
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.25.0
	github.com/aws/aws-sdk-go-v2/service/macie2 v1.26.0
	github.com/aws/aws-sdk-go-v2/service/mediastore v1.13.0
	github.com/aws/aws-sdk-go-v2/service/mgn v1.17.0
	github.com/aws/aws-sdk-go-v2/service/neptune v1.19.1
	github.com/aws/aws-sdk-go-v2/service/networkfirewall v1.24.0
	github.com/aws/aws-sdk-go-v2/service/oam v1.1.1
	github.com/aws/aws-sdk-go-v2/service/opensearch v1.13.1
	github.com/aws/aws-sdk-go-v2/service/organizations v1.18.0
	github.com/aws/aws-sdk-go-v2/service/pinpoint v1.18.0
	github.com/aws/aws-sdk-go-v2/service/pipes v1.1.0
	github.com/aws/aws-sdk-go-v2/service/pricing v1.18.0
	github.com/aws/aws-sdk-go-v2/service/ram v1.17.0
	github.com/aws/aws-sdk-go-v2/service/rds v1.40.1
	github.com/aws/aws-sdk-go-v2/service/redshift v1.27.1
	github.com/aws/aws-sdk-go-v2/service/redshiftserverless v1.4.1
	github.com/aws/aws-sdk-go-v2/service/resourceexplorer2 v1.2.1
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.14.1
	github.com/aws/aws-sdk-go-v2/service/route53 v1.27.0
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.14.0
	github.com/aws/aws-sdk-go-v2/service/route53resolver v1.16.1
	github.com/aws/aws-sdk-go-v2/service/s3 v1.30.1
	github.com/aws/aws-sdk-go-v2/service/s3control v1.29.1
	github.com/aws/aws-sdk-go-v2/service/sagemaker v1.66.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.18.2
	github.com/aws/aws-sdk-go-v2/service/securityhub v1.27.1
	github.com/aws/aws-sdk-go-v2/service/securitylake v1.2.0
	github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository v1.12.0
	github.com/aws/aws-sdk-go-v2/service/servicecatalog v1.25.5
	github.com/aws/aws-sdk-go-v2/service/servicediscovery v1.21.4
	github.com/aws/aws-sdk-go-v2/service/servicequotas v1.14.0
	github.com/aws/aws-sdk-go-v2/service/ses v1.15.0
	github.com/aws/aws-sdk-go-v2/service/sfn v1.17.1
	github.com/aws/aws-sdk-go-v2/service/simspaceweaver v1.1.0
	github.com/aws/aws-sdk-go-v2/service/sns v1.19.1
	github.com/aws/aws-sdk-go-v2/service/sqs v1.20.1
	github.com/aws/aws-sdk-go-v2/service/ssm v1.35.1
	github.com/aws/aws-sdk-go-v2/service/ssoadmin v1.16.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.18.2
	github.com/aws/aws-sdk-go-v2/service/transfer v1.33.7
	github.com/aws/aws-sdk-go-v2/service/waf v1.12.0
	github.com/aws/aws-sdk-go-v2/service/wafregional v1.13.1
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.24.2
	github.com/aws/aws-sdk-go-v2/service/wellarchitected v1.20.1
	github.com/aws/aws-sdk-go-v2/service/workspaces v1.28.0
	github.com/aws/smithy-go v1.19.0
	github.com/gocarina/gocsv v0.0.0-20201208093247-67c824bc04d4
	github.com/golang/protobuf v1.5.3
	github.com/turbot/go-kit v0.9.0-rc.1
	github.com/turbot/steampipe-plugin-sdk/v5 v5.8.0
	golang.org/x/text v0.13.0
)

require (
	cloud.google.com/go v0.110.7 // indirect
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.1 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/apparentlymart/go-textseg/v15 v15.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.2.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.5.9 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.28 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.18 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.21 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/eko/gocache/lib/v4 v4.1.5 // indirect
	github.com/eko/gocache/store/bigcache/v4 v4.2.1 // indirect
	github.com/eko/gocache/store/ristretto/v4 v4.2.1 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.3.0 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/glog v1.1.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/s2a-go v0.1.4 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.2.3 // indirect
	github.com/googleapis/gax-go/v2 v2.11.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.16.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.7.3 // indirect
	github.com/hashicorp/go-hclog v1.6.1 // indirect
	github.com/hashicorp/go-plugin v1.6.0 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.19.1 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
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
	github.com/zclconf/go-cty v1.14.1 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.opentelemetry.io/otel v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.44.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.21.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.21.0 // indirect
	go.opentelemetry.io/otel/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk v1.21.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.21.0 // indirect
	go.opentelemetry.io/otel/trace v1.21.0 // indirect
	go.opentelemetry.io/proto/otlp v1.0.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/exp v0.0.0-20230522175609-2e198f4a06a1 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.11.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/time v0.5.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.126.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/grpc v1.59.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
