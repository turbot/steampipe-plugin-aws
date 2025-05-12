---
title: "Steampipe Table: aws_route53_resolver_query_log_config - Query AWS Route 53 Resolver Query Log Config using SQL"
description: "Allows users to query AWS Route 53 Resolver Query Log Configurations."
folder: "Config"
---

# Table: aws_route53_resolver_query_log_config - Query AWS Route 53 Resolver Query Log Config using SQL

The AWS Route 53 Resolver Query Log Config enables DNS query logging in your Virtual Private Cloud (VPC). It logs the DNS queries that originate in your VPC and forwards them to CloudWatch Logs or S3 for safekeeping and analysis. This service aids in troubleshooting connectivity issues and understanding DNS querying behavior for security analysis.

## Table Usage Guide

The `aws_route53_resolver_query_log_config` table in Steampipe provides you with information about the query logging configurations within AWS Route 53 Resolver. This table allows you, as a DevOps engineer, security professional, or developer, to query configuration-specific details, including the destination, ownership, and status of the log configurations. You can utilize this table to gather insights on configurations, such as the AWS resource that logs are sent to, the number of VPCs that are associated with the configuration, and the ARN of the configuration. The schema outlines the various attributes of the Resolver Query Log Configuration, including the ID, creation time, destination, owner ID, and status.

**Important Notes**
- You must have `route53resolver:ListResolverQueryLogConfigs` permission to query the table.

## Examples

### Basic info
Explore the status and share status of your AWS Route53 Resolver query log configurations, along with their creation time. This can help in understanding their current state and managing them effectively.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which logging configurations are shared with your account but not owned by you. This can help you understand potential dependencies or collaborations within your AWS environment.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that have failed to properly log configurations on AWS Route53 Resolver. This is beneficial for identifying and rectifying any issues that may be causing the log configurations to fail.

```sql+postgres
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

```sql+sqlite
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
Explore which log configurations are shared with another account or organization. This can be useful to manage access and monitor the activity of shared logs.

```sql+postgres
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

```sql+sqlite
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
Determine recent additions to your log configurations by identifying those created within the past month. This allows you to stay updated on any new changes or additions made to your logging setup.

```sql+postgres
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

```sql+sqlite
select
  name,
  id,
  creation_time,
  destination_arn,
  status
from
  aws_route53_resolver_query_log_config
where
  creation_time >= datetime('now', '-30 day');
```