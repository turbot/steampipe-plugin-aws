---
title: "Steampipe Table: aws_elasticsearch_domain - Query AWS Elasticsearch Service Domain using SQL"
description: "Allows users to query AWS Elasticsearch Service Domains for detailed information related to the configuration, status, and access policies of the Elasticsearch domains."
folder: "Elasticsearch"
---

# Table: aws_elasticsearch_domain - Query AWS Elasticsearch Service Domain using SQL

The AWS Elasticsearch Service Domain is a fully managed service that makes it easy for you to deploy, secure, operate, and scale Elasticsearch to search, analyze, and visualize data in real-time. With the service, you get direct access to the Elasticsearch APIs and can seamlessly scale your workloads to hundreds of thousands of events per second. It offers built-in integrations with Kibana, Logstash, AWS services including Amazon Kinesis Data Firehose, AWS Lambda, and Amazon CloudWatch, so you can go from raw data to actionable insights quickly.

## Table Usage Guide

The `aws_elasticsearch_domain` table in Steampipe provides you with information about Elasticsearch domains within AWS Elasticsearch Service. This table enables you, as a DevOps engineer, to query domain-specific details, including configuration settings, access policies, and associated metadata. You can utilize this table to gather insights on domains, such as the domain's configuration, access and security settings, and more. The schema outlines the various attributes of the Elasticsearch domain for you, including the domain name, domain ID, ARN, created and deleted status, and associated tags.

## Examples

### Basic info

```sql+postgres
select
  domain_name,
  domain_id,
  arn,
  elasticsearch_version,
  created
from
  aws_elasticsearch_domain;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  arn,
  elasticsearch_version,
  created
from
  aws_elasticsearch_domain;
```


### List domains that are not encrypted at rest

```sql+postgres
select
  domain_name,
  domain_id,
  encryption_at_rest_options ->> 'Enabled' as enabled,
  encryption_at_rest_options ->> 'KmsKeyId' as kms_key_id
from
  aws_elasticsearch_domain
where
  encryption_at_rest_options ->> 'Enabled' = 'false';
```

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(encryption_at_rest_options, '$.Enabled') as enabled,
  json_extract(encryption_at_rest_options, '$.KmsKeyId') as kms_key_id
from
  aws_elasticsearch_domain
where
  json_extract(encryption_at_rest_options, '$.Enabled') = 'false';
```


### Get storage details for domains that are using EBS storage type

```sql+postgres
select
  domain_name,
  domain_id,
  ebs_options ->> 'VolumeSize' as volume_size,
  ebs_options ->> 'VolumeType' as volume_type,
  ebs_options ->> 'EBSEnabled' as ebs_enabled
from
  aws_elasticsearch_domain
where
  ebs_options ->> 'EBSEnabled' = 'true';
```

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(ebs_options, '$.VolumeSize') as volume_size,
  json_extract(ebs_options, '$.VolumeType') as volume_type,
  json_extract(ebs_options, '$.EBSEnabled') as ebs_enabled
from
  aws_elasticsearch_domain
where
  json_extract(ebs_options, '$.EBSEnabled') = 'true';
```


### Get network details for each domain

```sql+postgres
select
  domain_name,
  vpc_options ->> 'AvailabilityZones' as availability_zones,
  vpc_options ->> 'SecurityGroupIds' as security_group_ids,
  vpc_options ->> 'SubnetIds' as subnet_ids,
  vpc_options ->> 'VPCId' as vpc_id
from
  aws_elasticsearch_domain
where
  vpc_options ->> 'AvailabilityZones' is not null;
```

```sql+sqlite
select
  domain_name,
  json_extract(vpc_options.value, '$.AvailabilityZones') as availability_zones,
  json_extract(vpc_options.value, '$.SecurityGroupIds') as security_group_ids,
  json_extract(vpc_options.value, '$.SubnetIds') as subnet_ids,
  json_extract(vpc_options.value, '$.VPCId') as vpc_id
from
  aws_elasticsearch_domain,
  json_each(vpc_options) as vpc_options
where
  json_extract(vpc_options.value, '$.AvailabilityZones') is not null;
```


### Get the instance details for each domain

```sql+postgres
select
  domain_name,
  domain_id,
  elasticsearch_cluster_config ->> 'InstanceType' as instance_type,
  elasticsearch_cluster_config ->> 'InstanceCount' as instance_count
from
  aws_elasticsearch_domain;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(elasticsearch_cluster_config, '$.InstanceType') as instance_type,
  json_extract(elasticsearch_cluster_config, '$.InstanceCount') as instance_count
from
  aws_elasticsearch_domain;
```


### List domains that grant anonymous access

```sql+postgres
select
  domain_name,
  p as principal,
  a as action,
  s ->> 'Effect' as effect
from
  aws_elasticsearch_domain,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  p = '*'
  and s ->> 'Effect' = 'Allow';
```

```sql+sqlite
select
  domain_name,
  p.value as principal,
  a.value as action,
  json_extract(s.value, '$.Effect') as effect
from
  aws_elasticsearch_domain,
  json_each(policy_std, '$.Statement') as s,
  json_each(s.value, '$.Principal.AWS') as p,
  json_each(s.value, '$.Action') as a
where
  p.value = '*'
  and json_extract(s.value, '$.Effect') = 'Allow';
```


### List domain log publishing options

```sql+postgres
select
  domain_name,
  domain_id,
  log_publishing_options
from
  aws_elasticsearch_domain;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  log_publishing_options
from
  aws_elasticsearch_domain;
```


### List domain Search slow logs details

```sql+postgres
select
  domain_name,
  domain_id,
  log_publishing_options -> 'SEARCH_SLOW_LOGS' -> 'Enabled' as enabled,
  log_publishing_options -> 'SEARCH_SLOW_LOGS' -> 'CloudWatchLogsLogGroupArn' as cloud_watch_logs_log_group_arn
from
  aws_elasticsearch_domain;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(json_extract(log_publishing_options, '$.SEARCH_SLOW_LOGS'), '$.Enabled') as enabled,
  json_extract(json_extract(log_publishing_options, '$.SEARCH_SLOW_LOGS'), '$.CloudWatchLogsLogGroupArn') as cloud_watch_logs_log_group_arn
from
  aws_elasticsearch_domain;
```