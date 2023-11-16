---
title: "Table: aws_dynamodb_metric_account_provisioned_read_capacity_util - Query AWS DynamoDB Metrics using SQL"
description: "Allows users to query DynamoDB Metrics on account provisioned read capacity utilization."
---

# Table: aws_dynamodb_metric_account_provisioned_read_capacity_util - Query AWS DynamoDB Metrics using SQL

The `aws_dynamodb_metric_account_provisioned_read_capacity_util` table in Steampipe provides information about account provisioned read capacity utilization metrics within AWS DynamoDB. This table allows DevOps engineers to query metric-specific details, including the average, maximum, and minimum read capacity utilization. Users can utilize this table to gather insights on DynamoDB performance, such as understanding the read capacity utilization of their DynamoDB tables, identifying potential performance bottlenecks, and planning capacity accordingly. The schema outlines the various attributes of the DynamoDB metric, including the region, account_id, and timestamp.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dynamodb_metric_account_provisioned_read_capacity_util` table, you can use the `.inspect aws_dynamodb_metric_account_provisioned_read_capacity_util` command in Steampipe.

**Key columns**:

- `region`: This is the AWS region in which the DynamoDB table resides. It is vital for understanding the geographical distribution of your DynamoDB tables.
- `account_id`: This column contains the AWS account ID. It's useful for querying metrics across multiple AWS accounts.
- `timestamp`: This column provides the timestamp for the data point. It is important for tracking the metrics over a period of time.

## Examples

### Basic info

```sql
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

```sql
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
