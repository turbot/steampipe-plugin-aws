module github.com/turbot/steampipe-plugin-aws

go 1.19

require (
	github.com/aws/aws-sdk-go v1.44.153
	github.com/aws/aws-sdk-go-v2 v1.17.3
	github.com/aws/aws-sdk-go-v2/config v1.17.8
	github.com/aws/aws-sdk-go-v2/credentials v1.12.21
	github.com/aws/aws-sdk-go-v2/service/accessanalyzer v1.16.0
	github.com/aws/aws-sdk-go-v2/service/account v1.7.8
	github.com/aws/aws-sdk-go-v2/service/acm v1.14.8
	github.com/aws/aws-sdk-go-v2/service/amplify v1.11.18
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.15.10
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.12.8
	github.com/aws/aws-sdk-go-v2/service/appconfig v1.13.7
	github.com/aws/aws-sdk-go-v2/service/applicationautoscaling v1.15.18
	github.com/aws/aws-sdk-go-v2/service/auditmanager v1.20.4
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.23.10
	github.com/aws/aws-sdk-go-v2/service/backup v1.18.0
	github.com/aws/aws-sdk-go-v2/service/cloudcontrol v1.10.13
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.22.10
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.20.0
	github.com/aws/aws-sdk-go-v2/service/cloudsearch v1.13.19
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.21.1
	github.com/aws/aws-sdk-go-v2/service/cloudwatch v1.21.6
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.15.14
	github.com/aws/aws-sdk-go-v2/service/codeartifact v1.13.6
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.19.13
	github.com/aws/aws-sdk-go-v2/service/codecommit v1.13.17
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.14.16
	github.com/aws/aws-sdk-go-v2/service/codepipeline v1.13.15
	github.com/aws/aws-sdk-go-v2/service/configservice v1.28.0
	github.com/aws/aws-sdk-go-v2/service/costexplorer v1.19.2
	github.com/aws/aws-sdk-go-v2/service/databasemigrationservice v1.21.10
	github.com/aws/aws-sdk-go-v2/service/dax v1.11.15
	github.com/aws/aws-sdk-go-v2/service/directoryservice v1.14.11
	github.com/aws/aws-sdk-go-v2/service/dlm v1.12.4
	github.com/aws/aws-sdk-go-v2/service/docdb v1.19.11
	github.com/aws/aws-sdk-go-v2/service/drs v1.9.0
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.15.9
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.74.1
	github.com/aws/aws-sdk-go-v2/service/ecr v1.17.16
	github.com/aws/aws-sdk-go-v2/service/ecrpublic v1.13.15
	github.com/aws/aws-sdk-go-v2/service/ecs v1.21.0
	github.com/aws/aws-sdk-go-v2/service/efs v1.17.15
	github.com/aws/aws-sdk-go-v2/service/eks v1.26.0
	github.com/aws/aws-sdk-go-v2/service/elasticache v1.22.10
	github.com/aws/aws-sdk-go-v2/service/elasticbeanstalk v1.14.18
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.14.12
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.18.12
	github.com/aws/aws-sdk-go-v2/service/elasticsearchservice v1.16.10
	github.com/aws/aws-sdk-go-v2/service/emr v1.20.11
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.16.15
	github.com/aws/aws-sdk-go-v2/service/firehose v1.14.19
	github.com/aws/aws-sdk-go-v2/service/fsx v1.24.14
	github.com/aws/aws-sdk-go-v2/service/glacier v1.13.17
	github.com/aws/aws-sdk-go-v2/service/globalaccelerator v1.15.2
	github.com/aws/aws-sdk-go-v2/service/glue v1.38.1
	github.com/aws/aws-sdk-go-v2/service/guardduty v1.15.9
	github.com/aws/aws-sdk-go-v2/service/health v1.15.22
	github.com/aws/aws-sdk-go-v2/service/iam v1.18.9
	github.com/aws/aws-sdk-go-v2/service/identitystore v1.15.5
	github.com/aws/aws-sdk-go-v2/service/inspector v1.12.15
	github.com/aws/aws-sdk-go-v2/service/kafka v1.17.15
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.15.19
	github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2 v1.14.18
	github.com/aws/aws-sdk-go-v2/service/kinesisvideo v1.12.14
	github.com/aws/aws-sdk-go-v2/service/kms v1.18.11
	github.com/aws/aws-sdk-go-v2/service/lambda v1.26.0
	github.com/aws/aws-sdk-go-v2/service/lightsail v1.23.0
	github.com/aws/aws-sdk-go-v2/service/macie2 v1.23.4
	github.com/aws/aws-sdk-go-v2/service/mediastore v1.12.17
	github.com/aws/aws-sdk-go-v2/service/mgn v1.16.2
	github.com/aws/aws-sdk-go-v2/service/neptune v1.17.12
	github.com/aws/aws-sdk-go-v2/service/networkfirewall v1.20.0
	github.com/aws/aws-sdk-go-v2/service/oam v1.0.2
	github.com/aws/aws-sdk-go-v2/service/opensearch v1.10.10
	github.com/aws/aws-sdk-go-v2/service/organizations v1.16.8
	github.com/aws/aws-sdk-go-v2/service/pinpoint v1.17.10
	github.com/aws/aws-sdk-go-v2/service/pipes v1.0.1
	github.com/aws/aws-sdk-go-v2/service/pricing v1.16.8
	github.com/aws/aws-sdk-go-v2/service/ram v1.16.18
	github.com/aws/aws-sdk-go-v2/service/rds v1.26.1
	github.com/aws/aws-sdk-go-v2/service/redshift v1.26.10
	github.com/aws/aws-sdk-go-v2/service/redshiftserverless v1.2.9
	github.com/aws/aws-sdk-go-v2/service/resourceexplorer2 v1.0.0
	github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi v1.13.19
	github.com/aws/aws-sdk-go-v2/service/route53 v1.24.0
	github.com/aws/aws-sdk-go-v2/service/route53domains v1.12.17
	github.com/aws/aws-sdk-go-v2/service/route53resolver v1.15.19
	github.com/aws/aws-sdk-go-v2/service/s3 v1.27.1
	github.com/aws/aws-sdk-go-v2/service/s3control v1.21.9
	github.com/aws/aws-sdk-go-v2/service/sagemaker v1.48.0
	github.com/aws/aws-sdk-go-v2/service/secretsmanager v1.16.2
	github.com/aws/aws-sdk-go-v2/service/securityhub v1.26.0
	github.com/aws/aws-sdk-go-v2/service/securitylake v1.0.0
	github.com/aws/aws-sdk-go-v2/service/serverlessapplicationrepository v1.11.17
	github.com/aws/aws-sdk-go-v2/service/servicequotas v1.13.18
	github.com/aws/aws-sdk-go-v2/service/ses v1.14.18
	github.com/aws/aws-sdk-go-v2/service/sfn v1.14.1
	github.com/aws/aws-sdk-go-v2/service/simspaceweaver v1.0.2
	github.com/aws/aws-sdk-go-v2/service/sns v1.17.9
	github.com/aws/aws-sdk-go-v2/service/sqs v1.19.10
	github.com/aws/aws-sdk-go-v2/service/ssm v1.30.0
	github.com/aws/aws-sdk-go-v2/service/ssoadmin v1.15.11
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.19
	github.com/aws/aws-sdk-go-v2/service/waf v1.11.17
	github.com/aws/aws-sdk-go-v2/service/wafregional v1.12.18
	github.com/aws/aws-sdk-go-v2/service/wafv2 v1.22.9
	github.com/aws/aws-sdk-go-v2/service/wellarchitected v1.16.11
	github.com/aws/aws-sdk-go-v2/service/workspaces v1.23.0
	github.com/aws/smithy-go v1.13.5
	github.com/gocarina/gocsv v0.0.0-20201208093247-67c824bc04d4
	github.com/golang/protobuf v1.5.2
	github.com/turbot/go-kit v0.5.0
	github.com/turbot/steampipe-plugin-sdk/v5 v5.2.0-rc.2
	golang.org/x/text v0.4.0
)

require (
	cloud.google.com/go v0.65.0 // indirect
	cloud.google.com/go/storage v1.10.0 // indirect
	github.com/XiaoMi/pegasus-go-client v0.0.0-20210427083443-f3b6b08bc4c2 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/allegro/bigcache/v3 v3.1.0 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.27 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.21 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.24 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.20 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.23 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.13.6 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/go-netrc v0.0.0-20140422174119-9fd32a8b3d3d // indirect
	github.com/bradfitz/gomemcache v0.0.0-20221031212613-62deef7fc822 // indirect
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/danwakefield/fnmatch v0.0.0-20160403171240-cbb64ac3d964 // indirect
	github.com/dgraph-io/ristretto v0.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/eko/gocache/v3 v3.1.2 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/gertd/wild v0.0.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/googleapis/gax-go/v2 v2.0.5 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-getter v1.6.2 // indirect
	github.com/hashicorp/go-hclog v1.4.0 // indirect
	github.com/hashicorp/go-plugin v1.4.8 // indirect
	github.com/hashicorp/go-safetemp v1.0.0 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.15.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jstemmer/go-junit-report v0.9.1 // indirect
	github.com/klauspost/compress v1.11.2 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pegasus-kv/thrift v0.13.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sethvargo/go-retry v0.2.4 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/stevenle/topsort v0.0.0-20130922064739-8130c1d7596b // indirect
	github.com/tkrajina/go-reflector v0.5.6 // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	github.com/zclconf/go-cty v1.12.1 // indirect
	go.opencensus.io v0.22.4 // indirect
	go.opentelemetry.io/otel v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.30.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.30.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.10.0 // indirect
	go.opentelemetry.io/proto/otlp v0.16.0 // indirect
	golang.org/x/exp v0.0.0-20221110155412-d0897a79cd37 // indirect
	golang.org/x/lint v0.0.0-20200302205851-738671d3881b // indirect
	golang.org/x/mod v0.6.0 // indirect
	golang.org/x/net v0.2.0 // indirect
	golang.org/x/oauth2 v0.0.0-20220223155221-ee480838109b // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.2.0 // indirect
	golang.org/x/tools v0.2.0 // indirect
	google.golang.org/api v0.30.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac // indirect
	google.golang.org/grpc v1.51.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apimachinery v0.25.3 // indirect
)
