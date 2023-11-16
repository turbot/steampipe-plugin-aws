---
title: "Table: aws_route53_resolver_query_log_config - Query AWS Route 53 Resolver Query Log Config using SQL"
description: "Allows users to query AWS Route 53 Resolver Query Log Configurations."
---

# Table: aws_route53_resolver_query_log_config - Query AWS Route 53 Resolver Query Log Config using SQL

The `aws_route53_resolver_query_log_config` table in Steampipe provides information about the query logging configurations within AWS Route 53 Resolver. This table allows DevOps engineers, security professionals, and developers to query configuration-specific details, including the destination, ownership, and status of the log configurations. Users can utilize this table to gather insights on configurations, such as the AWS resource that logs are sent to, the number of VPCs that are associated with the configuration, and the ARN of the configuration. The schema outlines the various attributes of the Resolver Query Log Configuration, including the ID, creation time, destination, owner ID, and status.

**Note:** User must have `route53resolver:ListResolverQueryLogConfigs` permission for quering the table.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_resolver_query_log_config` table, you can use the `.inspect aws_route53_resolver_query_log_config` command in Steampipe.

**Key columns**:

- `id`: The ID of the configuration. This is a unique identifier and can be used to join this table with other tables that reference Resolver Query Log Configurations.
- `owner_id`: The AWS account ID of the account that created the configuration. This can be used to filter configurations by the account that owns them.
- `destination_arn`: The ARN of the resource that DNS queries are logged to, such as an Amazon S3 bucket, CloudWatch Logs log group, or Kinesis Data Firehose delivery stream. This can be used to join this table with tables that reference these resources.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  creation_time,
  share_status,
  status
from
  aws_route53_resolver_query_log_config;
```

### List log configs shared with my account

```sql
select
  name,
  id,
  arn,
  creation_time,
  share_status,
  status,
  destination_arn
from
  aws_route53_resolver_query_log_config
where
  owner_id <> account_id;
```

### List failed log configurations

```sql
select
  name,
  id,
  creator_request_id,
  destination_arn
from
  aws_route53_resolver_query_log_config
where
  status = 'FAILED';
```

### List log configurations shared with another account or organization

```sql
select
  name,
  id,
  share_status,
  association_count
from
  aws_route53_resolver_query_log_config
where
  share_status = 'SHARED';
```

### List log configs created in the last 30 days

```sql
select
  name,
  id,
  creation_time,
  destination_arn,
  status
from
  aws_route53_resolver_query_log_config
where
  creation_time >= now() - interval '30' day;
```
