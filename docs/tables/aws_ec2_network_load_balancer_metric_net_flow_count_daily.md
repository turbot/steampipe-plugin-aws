---
title: "Table: aws_ec2_network_load_balancer_metric_net_flow_count_daily - Query AWS EC2 Network Load Balancer Metrics using SQL"
description: "Allows users to query Network Load Balancer Metrics in EC2, specifically the daily net flow count, providing insights into network traffic patterns and potential anomalies."
---

# Table: aws_ec2_network_load_balancer_metric_net_flow_count_daily - Query AWS EC2 Network Load Balancer Metrics using SQL

The `aws_ec2_network_load_balancer_metric_net_flow_count_daily` table in Steampipe provides information about network load balancer metrics within AWS Elastic Compute Cloud (EC2). This table allows DevOps engineers to query daily net flow count details, including timestamp, average, minimum, maximum, and sum. Users can utilize this table to gather insights on network traffic patterns, detect potential network anomalies, and optimize load balancing strategies. The schema outlines the various attributes of the network load balancer metric, including the load balancer name, namespace, region, and metric unit.

The `aws_ec2_network_load_balancer_metric_net_flow_count_daily` table provides metric statistics at 24 hour intervals for the most recent 1 year.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_network_load_balancer_metric_net_flow_count_daily` table, you can use the `.inspect aws_ec2_network_load_balancer_metric_net_flow_count_daily` command in Steampipe.

### Key columns:

- `load_balancer_name`: This is the name of the load balancer. It is crucial for identifying the specific load balancer for which the metrics are being queried.
- `timestamp`: This column records the time at which the metrics were collected. It is useful for tracking network traffic patterns over time.
- `average`: This column shows the average net flow count for the day. It is important for understanding the typical load on the network.

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
  aws_ec2_network_load_balancer_metric_net_flow_count_daily
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
  aws_ec2_network_load_balancer_metric_net_flow_count_daily
where
  average < 100
order by
  name,
  timestamp;
```
