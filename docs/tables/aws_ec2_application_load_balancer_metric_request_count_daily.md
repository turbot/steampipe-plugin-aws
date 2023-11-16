---
title: "Table: aws_ec2_application_load_balancer_metric_request_count_daily - Query AWS EC2 Application Load Balancer using SQL"
description: "Allows users to query daily request count metrics of the AWS EC2 Application Load Balancer."
---

# Table: aws_ec2_application_load_balancer_metric_request_count_daily - Query AWS EC2 Application Load Balancer using SQL

The `aws_ec2_application_load_balancer_metric_request_count_daily` table in Steampipe provides information about the daily request count metrics of the AWS EC2 Application Load Balancer. This table allows DevOps engineers to query and analyze the number of requests processed by the Application Load Balancer on a daily basis. Users can utilize this table to monitor the load on the balancer, identify potential spikes in traffic, and plan capacity accordingly. The schema outlines the various attributes of the metrics, including the load balancer name, namespace, metric name, and the timestamp of the metric.

The `aws_ec2_application_load_balancer_metric_request_count_daily` table provides metric statistics at 24 hour intervals for the most recent 1 year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_application_load_balancer_metric_request_count_daily` table, you can use the `.inspect aws_ec2_application_load_balancer_metric_request_count_daily` command in Steampipe.

Key columns:

- `title`: The title of the metric. This is useful for identifying the specific metric in question.
- `load_balancer_name`: The name of the load balancer. This can be used to join with other tables that contain load balancer information.
- `timestamp`: The timestamp of the metric. This can be used to track the metrics over time and identify patterns or anomalies.

## Examples

### Basic info

```sql
select
  name,
  metric_name,
  namespace,
  average,
  maximum,
  minimum,
  sample_count,
  timestamp
from
  aws_ec2_application_load_balancer_metric_request_count_daily
order by
  name,
  timestamp;
```

### Intervals averaging less than 100 request count

```sql
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  average
  sample_count,
  timestamp
from
  aws_ec2_application_load_balancer_metric_request_count_daily
where
  average < 100
order by
  name,
  timestamp;
```
