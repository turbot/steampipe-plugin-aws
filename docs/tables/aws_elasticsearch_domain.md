---
title: "Table: aws_elasticsearch_domain - Query AWS Elasticsearch Service Domain using SQL"
description: "Allows users to query AWS Elasticsearch Service Domains for detailed information related to the configuration, status, and access policies of the Elasticsearch domains."
---

# Table: aws_elasticsearch_domain - Query AWS Elasticsearch Service Domain using SQL

The `aws_elasticsearch_domain` table in Steampipe provides information about Elasticsearch domains within AWS Elasticsearch Service. This table enables DevOps engineers to query domain-specific details, including configuration settings, access policies, and associated metadata. Users can utilize this table to gather insights on domains, such as the domain's configuration, access and security settings, and more. The schema outlines the various attributes of the Elasticsearch domain, including the domain name, domain ID, ARN, created and deleted status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticsearch_domain` table, you can use the `.inspect aws_elasticsearch_domain` command in Steampipe.

### Key columns:

- `domain_name`: The name of the Elasticsearch domain. It is a unique key that can be used to join this table with other tables.
- `domain_id`: The unique identifier for the domain. This is important for tracking the domain across different AWS services.
- `arn`: The Amazon Resource Name (ARN) of the Elasticsearch domain. This is a globally unique identifier that is important for managing access and security.

## Example

### Basic info

```sql
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

```sql
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


### Get storage details for domains that are using EBS storage type

```sql
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


### Get network details for each domain

```sql
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


### Get the instance details for each domain

```sql
select
  domain_name,
  domain_id,
  elasticsearch_cluster_config ->> 'InstanceType' as instance_type,
  elasticsearch_cluster_config ->> 'InstanceCount' as instance_count
from
  aws_elasticsearch_domain;
```


### List domains that grant anonymous access

```sql
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


### List domain log publishing options

```sql
select
  domain_name,
  domain_id,
  log_publishing_options
from
  aws_elasticsearch_domain;
```


### List domain Search slow logs details

```sql
select
  domain_name,
  domain_id,
  log_publishing_options -> 'SEARCH_SLOW_LOGS' -> 'Enabled' as enabled,
  log_publishing_options -> 'SEARCH_SLOW_LOGS' -> 'CloudWatchLogsLogGroupArn' as cloud_watch_logs_log_group_arn
from
  aws_elasticsearch_domain;
```