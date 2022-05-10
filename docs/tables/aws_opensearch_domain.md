# Table: aws_opensearch_domain

Amazon OpenSearch Service is a managed service that makes it easy to deploy, operate, and scale OpenSearch clusters in the AWS Clouds. Domains are clusters with the settings, instance types, instance counts, and storage resources that you specify.

## Example

### Basic info

```sql
select
  domain_name,
  domain_id,
  arn,
  engine_version,
  created
from
  aws_opensearch_domain;
```


### List domains that are not encrypted at rest

```sql
select
  domain_name,
  domain_id,
  encryption_at_rest_options ->> 'Enabled' as enabled,
  encryption_at_rest_options ->> 'KmsKeyId' as kms_key_id
from
  aws_opensearch_domain
where
  encryption_at_rest_options ->> 'Enabled' = 'false';
```


### Get storage details for domains that are using EBS storage type

```sql
select
  domain_name,
  domain_id,
  ebs_options ->> 'VolumeSize' as volume_size,
  ebs_options ->> 'VolumeType' as volume_type,
  ebs_options ->> 'EBSEnabled' as ebs_enabled
from
  aws_opensearch_domain
where
  ebs_options ->> 'EBSEnabled' = 'true';
```


### Get network details for each domain

```sql
select
  domain_name,
  vpc_options ->> 'AvailabilityZones' as availability_zones,
  vpc_options ->> 'SecurityGroupIds' as security_group_ids,
  vpc_options ->> 'SubnetIds' as subnet_ids,
  vpc_options ->> 'VPCId' as vpc_id
from
  aws_opensearch_domain
where
  vpc_options ->> 'AvailabilityZones' is not null;
```


### Get the instance details for each domain

```sql
select
  domain_name,
  domain_id,
  cluster_config ->> 'InstanceType' as instance_type,
  cluster_config ->> 'InstanceCount' as instance_count
from
  aws_opensearch_domain;
```


### List domains that are publicly accessible

```sql
select
  domain_name,
  domain_id,
  arn,
  engine_version,
  created
from
  aws_opensearch_domain,
where
  vpc_options is null;
```


### List domain log publishing options

```sql
select
  domain_name,
  domain_id,
  log_publishing_options
from
  aws_opensearch_domain;
```


### List domain Search slow logs details

```sql
select
  domain_name,
  domain_id,
  log_publishing_options -> 'SEARCH_SLOW_LOGS' -> 'Enabled' as enabled,
  log_publishing_options -> 'SEARCH_SLOW_LOGS' -> 'CloudWatchLogsLogGroupArn' as cloud_watch_logs_log_group_arn
from
  aws_opensearch_domain;
```
