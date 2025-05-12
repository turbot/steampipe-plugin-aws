---
title: "Steampipe Table: aws_dynamodb_metric_account_provisioned_write_capacity_util - Query AWS DynamoDB Metrics using SQL"
description: "Allows users to query AWS DynamoDB Metrics for account provisioned write capacity utilization."
folder: "DynamoDB"
---

# Table: aws_dynamodb_metric_account_provisioned_write_capacity_util - Query AWS DynamoDB Metrics using SQL

The AWS DynamoDB Metrics service allows you to monitor the performance characteristics of DynamoDB tables. One such metric is the Account Provisioned Write Capacity Utilization, which provides information about the write capacity units consumed by all tables in your AWS account. This helps you manage your resources effectively and optimize your database's performance.

## Table Usage Guide

The `aws_dynamodb_metric_account_provisioned_write_capacity_util` table in Steampipe provides you with information about the provisioned write capacity utilization metrics at the account level within Amazon DynamoDB. This table allows you, as a DevOps engineer, to query details related to the provisioned write capacity, such as the average, maximum, and minimum write capacity units consumed by all tables in your account. You can utilize this table to monitor the utilization of provisioned write capacity, ensuring optimal performance and identifying potential bottlenecks or over-provisioning. The schema outlines the various attributes of the metric, including your account id, region, timestamp, and the provisioned write capacity units.

## Examples

### Basic info
Determine the areas in which your AWS DynamoDB is being utilized by understanding the provisioned write capacity over time. This query can help you manage resources more efficiently by identifying peak usage times and patterns.

```sql+postgres
select
  account_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_dynamodb_metric_account_provisioned_write_capacity_util
order by
  timestamp;
```

```sql+sqlite
select
  account_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_dynamodb_metric_account_provisioned_write_capacity_util
order by
  timestamp;
```

### Intervals where throughput exceeds 80 percent
Determine the instances where the provisioned write capacity of your AWS DynamoDB exceeds 80 percent. This can be useful to identify periods of high demand and potentially optimize your resource allocation for improved performance.

```sql+postgres
select
  account_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_dynamodb_metric_account_provisioned_write_capacity_util
where
  maximum > 80
order by
  timestamp;
```

```sql+sqlite
select
  account_id,
  timestamp,
  minimum,
  maximum,
  average,
  sum,
  sample_count
from
  aws_dynamodb_metric_account_provisioned_write_capacity_util
where
  maximum > 80
order by
  timestamp;
```