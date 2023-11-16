---
title: "Table: aws_opensearch_domain - Query AWS OpenSearch Service Domains using SQL"
description: "Allows users to query AWS OpenSearch Service Domains for detailed information on their configuration, status, and associated resources."
---

# Table: aws_opensearch_domain - Query AWS OpenSearch Service Domains using SQL

The `aws_opensearch_domain` table in Steampipe provides information about domains within the AWS OpenSearch Service. This table allows DevOps engineers to query domain-specific details, including configurations, access policies, and associated metadata. Users can utilize this table to gather insights on domains, such as their encryption status, node-to-node encryption options, automated snapshot settings, and more. The schema outlines the various attributes of the OpenSearch domain, including the domain ARN, domain ID, created date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_opensearch_domain` table, you can use the `.inspect aws_opensearch_domain` command in Steampipe.

### Key columns:

- `domain_name`: This is the name of the OpenSearch domain. It is a key column because it is a unique identifier for each domain and can be used to join this table with other tables.
- `arn`: This is the Amazon Resource Name (ARN) of the OpenSearch domain. It is a key column because it provides a consistent way to identify AWS resources across all AWS services.
- `domain_id`: This is the unique identifier for the OpenSearch domain. It is a key column because it can be used to join this table with other tables that contain information about OpenSearch domains.

## Examples

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
  aws_opensearch_domain
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