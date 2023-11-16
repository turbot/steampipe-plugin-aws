---
title: "Table: aws_ec2_network_load_balancer_metric_net_flow_count - Query AWS EC2 Network Load Balancer Metrics using SQL"
description: "Allows users to query AWS EC2 Network Load Balancer Metrics for net flow count data. This includes information such as the number of new or terminated flows per minute from a network load balancer."
---

# Table: aws_ec2_network_load_balancer_metric_net_flow_count - Query AWS EC2 Network Load Balancer Metrics using SQL

The `aws_ec2_network_load_balancer_metric_net_flow_count` table in Steampipe provides information about the net flow count metrics of AWS EC2 Network Load Balancers. This table allows DevOps engineers to query net flow count-specific details, including the number of new or terminated flows per minute. Users can utilize this table to gather insights on network load balancing, such as monitoring the amount of traffic processed by their load balancer, identifying trends in network traffic, and more. The schema outlines the various attributes of the net flow count metric, including the load balancer name, namespace, metric name, and dimensions.

The `aws_ec2_network_load_balancer_metric_net_flow_count` table provides metric statistics at 5 min intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_network_load_balancer_metric_net_flow_count` table, you can use the `.inspect aws_ec2_network_load_balancer_metric_net_flow_count` command in Steampipe.

**Key columns**:

- `name`: The name of the load balancer. This can be used to join with other tables that contain load balancer information.
- `namespace`: The namespace for the AWS service that the metric data is associated with. This can be used to filter metrics from a specific AWS service.
- `dimensions`: The dimensions for the metric. This can be used to provide specific characteristics for a metric.

## Examples

### Basic info

```sql
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  sample_count,
  timestamp
from
  aws_ec2_network_load_balancer_metric_net_flow_count
order by
  name,
  timestamp;
```

### Intervals where net flow count < 100

```sql
select
  name,
  metric_name,
  namespace,
  maximum,
  minimum,
  average,
  sample_count,
  timestamp
from
  aws_ec2_network_load_balancer_metric_net_flow_count
where
  average < 100
order by
  name,
  timestamp;
```
