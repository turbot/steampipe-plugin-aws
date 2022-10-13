module github.com/turbot/steampipe-plugin-aws

go 1.19

require (
	github.com/aws/aws-sdk-go v1.44.49
	github.com/aws/aws-sdk-go-v2 v1.16.16
	github.com/aws/aws-sdk-go-v2/config v1.17.1
	github.com/aws/aws-sdk-go-v2/credentials v1.12.14
	github.com/aws/aws-sdk-go-v2/service/accessanalyzer v1.15.12
	github.com/aws/aws-sdk-go-v2/service/account v1.7.8
	github.com/aws/aws-sdk-go-v2/service/acm v1.14.8
	github.com/aws/aws-sdk-go-v2/service/apigateway v1.15.10
	github.com/aws/aws-sdk-go-v2/service/apigatewayv2 v1.12.8
	github.com/aws/aws-sdk-go-v2/service/appconfig v1.13.7
	github.com/aws/aws-sdk-go-v2/service/auditmanager v1.20.4
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.23.10
	github.com/aws/aws-sdk-go-v2/service/backup v1.17.5
	github.com/aws/aws-sdk-go-v2/service/cloudfront v1.20.0
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.16.8
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.15.14
	github.com/aws/aws-sdk-go-v2/service/codeartifact v1.13.6
	github.com/aws/aws-sdk-go-v2/service/codebuild v1.19.13
	github.com/aws/aws-sdk-go-v2/service/codedeploy v1.14.16
	github.com/aws/aws-sdk-go-v2/service/costexplorer v1.19.2
	github.com/aws/aws-sdk-go-v2/service/dax v1.11.15
	github.com/aws/aws-sdk-go-v2/service/docdb v1.19.8
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.15.9
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.52.1
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.14.12
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.18.12
	github.com/aws/aws-sdk-go-v2/service/firehose v1.14.19
	github.com/aws/aws-sdk-go-v2/service/glacier v1.13.17
	github.com/aws/aws-sdk-go-v2/service/iam v1.18.9
	github.com/aws/aws-sdk-go-v2/service/kafka v1.17.15
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.15.19
	github.com/aws/aws-sdk-go-v2/service/kinesisanalyticsv2 v1.14.18
	github.com/aws/aws-sdk-go-v2/service/kinesisvideo v1.12.14
	github.com/aws/aws-sdk-go-v2/service/opensearch v1.10.10
	github.com/aws/aws-sdk-go-v2/service/organizations v1.16.8
	github.com/aws/aws-sdk-go-v2/service/redshiftserverless v1.2.9
	github.com/aws/aws-sdk-go-v2/service/s3 v1.27.1
	github.com/aws/aws-sdk-go-v2/service/s3control v1.21.9
	github.com/aws/aws-sdk-go-v2/service/sns v1.17.9
	github.com/aws/smithy-go v1.13.3
	github.com/gocarina/gocsv v0.0.0-20201208093247-67c824bc04d4
	github.com/golang/protobuf v1.5.2
	github.com/turbot/go-kit v0.4.0
	github.com/turbot/steampipe-plugin-sdk/v4 v4.1.7
	golang.org/x/text v0.3.7
)

require (
	github.com/XiaoMi/pegasus-go-client v0.0.0-20210427083443-f3b6b08bc4c2 // indirect
	github.com/acarl005/stripansi v0.0.0-20180116102854-5a71ef0e047d // indirect
	github.com/agext/levenshtein v1.2.2 // indirect
	github.com/allegro/bigcache/v3 v3.0.2 // indirect
	github.com/apparentlymart/go-textseg/v13 v13.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.23 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.17 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.19 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.9 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.7.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.14 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.13 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bradfitz/gomemcache v0.0.0-20220106215444-fb4bf637b56d // indirect
	github.com/btubbs/datetime v0.1.1 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgraph-io/ristretto v0.1.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/eko/gocache/v3 v3.1.1 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/gertd/go-pluralize v0.2.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/hashicorp/go-hclog v1.2.2 // indirect
	github.com/hashicorp/go-plugin v1.4.4 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/hashicorp/hcl/v2 v2.13.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20181012175058-2f1d1f20f75d // indirect
	github.com/iancoleman/strcase v0.2.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/pegasus-kv/thrift v0.13.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.33.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sethvargo/go-retry v0.1.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/stevenle/topsort v0.0.0-20130922064739-8130c1d7596b // indirect
	github.com/tkrajina/go-reflector v0.5.4 // indirect
	github.com/zclconf/go-cty v1.10.0 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric v0.30.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc v0.30.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.7.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	go.opentelemetry.io/proto/otlp v0.16.0 // indirect
	golang.org/x/exp v0.0.0-20220518171630-0b5c67f07fdf // indirect
	golang.org/x/net v0.0.0-20220412020605-290c469a71a5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	google.golang.org/genproto v0.0.0-20220407144326-9054f6ed7bac // indirect
	google.golang.org/grpc v1.48.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/tomb.v2 v2.0.0-20161208151619-d5d1b5820637 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/apimachinery v0.23.5 // indirect
)
