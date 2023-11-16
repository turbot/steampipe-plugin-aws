---
title: "Table: aws_ec2_application_load_balancer_metric_request_count - Query AWS EC2 Application Load Balancer Metrics using SQL"
description: "Allows users to query AWS EC2 Application Load Balancer Metrics, specifically the request count."
---

# Table: aws_ec2_application_load_balancer_metric_request_count - Query AWS EC2 Application Load Balancer Metrics using SQL

The `aws_ec2_application_load_balancer_metric_request_count` table in Steampipe provides information about the request count metrics of Application Load Balancers within Amazon Elastic Compute Cloud (EC2). This table allows DevOps engineers, system administrators, and other technical professionals to query specific details about the number of requests processed by their Application Load Balancers. Users can utilize this table to gather insights on load balancing performance and to monitor the traffic their applications are receiving. The schema outlines the various attributes of the request count metric, including the load balancer name, namespace, metric name, and dimensions.

The `aws_ec2_application_load_balancer_metric_request_count` table provides metric statistics at 5 min intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_application_load_balancer_metric_request_count` table, you can use the `.inspect aws_ec2_application_load_balancer_metric_request_count` command in Steampipe.

Key columns:

- `load_balancer_name`: The name of the load balancer. This column is useful for joining with other tables that contain load balancer details.
- `namespace`: The namespace for the AWS metric data. This column can be used to filter results by the metric namespace.
- `dimensions`: The dimensions for the metric data. This column can be used to filter or aggregate results by specific metric dimensions.

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
  aws_ec2_application_load_balancer_metric_request_count
order by
  name,
  timestamp;
```

### Intervals averaging less than 100 net flow count

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
  aws_ec2_application_load_balancer_metric_request_count
where
  average < 100
order by
  name,
  timestamp;
```
