---
title: "Steampipe Table: aws_ec2_network_load_balancer_metric_net_flow_count_daily - Query AWS EC2 Network Load Balancer Metrics using SQL"
description: "Allows users to query Network Load Balancer Metrics in EC2, specifically the daily net flow count, providing insights into network traffic patterns and potential anomalies."
folder: "ELB"
---

# Table: aws_ec2_network_load_balancer_metric_net_flow_count_daily - Query AWS EC2 Network Load Balancer Metrics using SQL

The AWS EC2 Network Load Balancer is a fully managed service that automatically distributes incoming traffic across multiple targets, such as Amazon EC2 instances, containers, IP addresses, and Lambda functions. It can handle the varying load of your applications in a single Availability Zone or across multiple Availability Zones. The 'NetFlowCount' metric provides the total number of new TCP/UDP flows established from clients to targets in a specified time period.

## Table Usage Guide

The `aws_ec2_network_load_balancer_metric_net_flow_count_daily` table in Steampipe provides you with information about network load balancer metrics within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query daily net flow count details, including timestamp, average, minimum, maximum, and sum. You can utilize this table to gather insights on network traffic patterns, detect potential network anomalies, and optimize load balancing strategies. The schema outlines the various attributes of the network load balancer metric for you, including the load balancer name, namespace, region, and metric unit.

The `aws_ec2_network_load_balancer_metric_net_flow_count_daily` table provides you with metric statistics at 24-hour intervals for the most recent 1 year.

## Examples

### Basic info
Explore which network load balancers in your AWS EC2 environment have the highest and lowest daily net flow counts. This allows you to gain insights into the performance and usage patterns of your load balancers over time.

```sql+postgres
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

```sql+sqlite
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
Determine the intervals where the average daily network flow count is less than 100 for AWS EC2 Network Load Balancer. This can be useful in identifying periods of low traffic, which could indicate underutilization or potential opportunities for cost-saving.

```sql+postgres
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

```sql+sqlite
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