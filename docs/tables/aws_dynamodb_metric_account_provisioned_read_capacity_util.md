---
title: "Steampipe Table: aws_dynamodb_metric_account_provisioned_read_capacity_util - Query AWS DynamoDB Metrics using SQL"
description: "Allows users to query DynamoDB Metrics on account provisioned read capacity utilization."
folder: "DynamoDB"
---

# Table: aws_dynamodb_metric_account_provisioned_read_capacity_util - Query AWS DynamoDB Metrics using SQL

The AWS DynamoDB Metrics service provides detailed performance metrics for your DynamoDB tables. One such metric is the Account Provisioned Read Capacity Utilization, which measures the percentage of provisioned read capacity units that your application consumes. This allows you to monitor your application's read activity and optimize your provisioned read capacity units for cost-effectiveness and performance.

## Table Usage Guide

The `aws_dynamodb_metric_account_provisioned_read_capacity_util` table in Steampipe provides you with information about account provisioned read capacity utilization metrics within AWS DynamoDB. This table allows you, as a DevOps engineer, to query metric-specific details, including the average, maximum, and minimum read capacity utilization. You can utilize this table to gather insights on DynamoDB performance, such as understanding the read capacity utilization of your DynamoDB tables, identifying potential performance bottlenecks, and planning capacity accordingly. The schema outlines the various attributes of the DynamoDB metric, including the region, account_id, and timestamp for you.

## Examples

### Basic info
Determine the areas in which your AWS DynamoDB account's provisioned read capacity is being utilized. This can help in monitoring resource usage over time and planning for future capacity needs.

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
  aws_dynamodb_metric_account_provisioned_read_capacity_util
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
  aws_dynamodb_metric_account_provisioned_read_capacity_util
order by
  timestamp;
```

### Intervals where throughput exceeds 80 percent
Analyze the instances where the provisioned read capacity utilization of your AWS DynamoDB account exceeds 80 percent. This can help in identifying periods of high demand and assist in capacity planning to ensure optimal performance.

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
  aws_dynamodb_metric_account_provisioned_read_capacity_util
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
  aws_dynamodb_metric_account_provisioned_read_capacity_util
where
  maximum > 80
order by
  timestamp;
```