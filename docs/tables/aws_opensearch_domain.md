---
title: "Steampipe Table: aws_opensearch_domain - Query AWS OpenSearch Service Domains using SQL"
description: "Allows users to query AWS OpenSearch Service Domains for detailed information on their configuration, status, and associated resources."
folder: "OpenSearch"
---

# Table: aws_opensearch_domain - Query AWS OpenSearch Service Domains using SQL

The AWS OpenSearch Service (successor to Amazon Elasticsearch Service) allows you to easily build, secure, and manage your own cost-effective search applications. With AWS OpenSearch Service, you can deploy and run OpenSearch (an open-source search and analytics suite) and its predecessor, Elasticsearch, on AWS without having to provision infrastructure or perform time-consuming setup and maintenance tasks. This service offers you direct access to the OpenSearch APIs and automatically takes care of the administrative operations.

## Table Usage Guide

The `aws_opensearch_domain` table in Steampipe provides you with information about domains within the AWS OpenSearch Service. This table allows you as a DevOps engineer to query domain-specific details, including configurations, access policies, and associated metadata. You can utilize this table to gather insights on domains, such as their encryption status, node-to-node encryption options, automated snapshot settings, and more. The schema outlines the various attributes of the OpenSearch domain for you, including the domain ARN, domain ID, created date, and associated tags.

## Examples

### Basic info
Explore which domains are currently active on your AWS OpenSearch service. This query is particularly useful for gaining insights into the engine versions being used and when they were created, allowing you to better manage and update your domains.

```sql+postgres
select
  domain_name,
  domain_id,
  arn,
  engine_version,
  created
from
  aws_opensearch_domain;
```

```sql+sqlite
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
Determine the areas in which domains are not encrypted at rest, allowing you to identify potential security vulnerabilities and take necessary actions to enhance data protection.

```sql+postgres
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

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(encryption_at_rest_options, '$.Enabled') as enabled,
  json_extract(encryption_at_rest_options, '$.KmsKeyId') as kms_key_id
from
  aws_opensearch_domain
where
  json_extract(encryption_at_rest_options, '$.Enabled') = 'false';
```

### Get storage details for domains that are using EBS storage type
Identify the domains utilizing EBS storage by assessing their storage details. This can help in management and optimization of storage resources within those domains.

```sql+postgres
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

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(ebs_options, '$.VolumeSize') as volume_size,
  json_extract(ebs_options, '$.VolumeType') as volume_type,
  json_extract(ebs_options, '$.EBSEnabled') as ebs_enabled
from
  aws_opensearch_domain
where
  json_extract(ebs_options, '$.EBSEnabled') = 'true';
```

### Get network details for each domain
Explore the network configuration of each domain to gain insights into their availability zones, security group IDs, subnet IDs, and VPC IDs. This could be useful for assessing the network structure and security measures implemented across different domains.

```sql+postgres
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

```sql+sqlite
select
  domain_name,
  json_extract(vpc_options, '$.AvailabilityZones') as availability_zones,
  json_extract(vpc_options, '$.SecurityGroupIds') as security_group_ids,
  json_extract(vpc_options, '$.SubnetIds') as subnet_ids,
  json_extract(vpc_options, '$.VPCId') as vpc_id
from
  aws_opensearch_domain
where
  json_extract(vpc_options, '$.AvailabilityZones') is not null;
```

### Get the instance details for each domain
Identify the configuration of each domain to understand its specific instance type and count. This can assist in managing resources and optimizing performance within the AWS OpenSearch service.

```sql+postgres
select
  domain_name,
  domain_id,
  cluster_config ->> 'InstanceType' as instance_type,
  cluster_config ->> 'InstanceCount' as instance_count
from
  aws_opensearch_domain;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(cluster_config, '$.InstanceType') as instance_type,
  json_extract(cluster_config, '$.InstanceCount') as instance_count
from
  aws_opensearch_domain;
```

### List domains that are publicly accessible
Discover the segments that are publicly accessible, allowing you to identify potential vulnerabilities and enhance security measures. This is useful for maintaining a secure environment by preventing unauthorized access.

```sql+postgres
select
  domain_name,
  domain_id,
  arn,
  engine_version,
  created
from
  aws_opensearch_domain
where
  vpc_options is null;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  arn,
  engine_version,
  created
from
  aws_opensearch_domain
where
  vpc_options is null;
```

### List domain log publishing options
Explore which AWS OpenSearch domains have specific log publishing options enabled. This can be useful in understanding the logging practices across your domains, helping ensure compliance with logging policies and troubleshoot any potential issues.

```sql+postgres
select
  domain_name,
  domain_id,
  log_publishing_options
from
  aws_opensearch_domain;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  log_publishing_options
from
  aws_opensearch_domain;
```

### List domain Search slow logs details
Explore which domains have slow search log publishing enabled and where these logs are stored. This is useful for identifying potential performance issues and ensuring logs are properly archived for future analysis.

```sql+postgres
select
  domain_name,
  domain_id,
  log_publishing_options -> 'SEARCH_SLOW_LOGS' -> 'Enabled' as enabled,
  log_publishing_options -> 'SEARCH_SLOW_LOGS' -> 'CloudWatchLogsLogGroupArn' as cloud_watch_logs_log_group_arn
from
  aws_opensearch_domain;
```

```sql+sqlite
select
  domain_name,
  domain_id,
  json_extract(json_extract(log_publishing_options, '$.SEARCH_SLOW_LOGS'), '$.Enabled') as enabled,
  json_extract(json_extract(log_publishing_options, '$.SEARCH_SLOW_LOGS'), '$.CloudWatchLogsLogGroupArn') as cloud_watch_logs_log_group_arn
from
  aws_opensearch_domain;
```