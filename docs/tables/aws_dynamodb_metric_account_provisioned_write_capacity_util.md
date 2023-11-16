---
title: "Table: aws_dynamodb_metric_account_provisioned_write_capacity_util - Query AWS DynamoDB Metrics using SQL"
description: "Allows users to query AWS DynamoDB Metrics for account provisioned write capacity utilization."
---

# Table: aws_dynamodb_metric_account_provisioned_write_capacity_util - Query AWS DynamoDB Metrics using SQL

The `aws_dynamodb_metric_account_provisioned_write_capacity_util` table in Steampipe provides information about the provisioned write capacity utilization metrics at the account level within Amazon DynamoDB. This table allows DevOps engineers to query details related to the provisioned write capacity, such as the average, maximum, and minimum write capacity units consumed by all tables in the account. Users can utilize this table to monitor the utilization of provisioned write capacity, ensuring optimal performance and identifying potential bottlenecks or over-provisioning. The schema outlines the various attributes of the metric, including the account id, region, timestamp, and the provisioned write capacity units.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dynamodb_metric_account_provisioned_write_capacity_util` table, you can use the `.inspect aws_dynamodb_metric_account_provisioned_write_capacity_util` command in Steampipe.

Key columns:

- `account_id`: The AWS account id associated with the metric. This column can be used to join with other tables that contain account-specific information.
- `region`: The AWS region in which the metric was recorded. This column can be used to join with other tables that contain region-specific information.
- `timestamp`: The timestamp when the metric data was recorded. This column can be used to join with other tables that contain time-specific information, allowing for a chronological analysis of data.

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
  aws_dynamodb_metric_account_provisioned_write_capacity_util
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
  aws_dynamodb_metric_account_provisioned_write_capacity_util
where
  maximum > 80
order by
  timestamp;
```
