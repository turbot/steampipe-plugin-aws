---
title: "Steampipe Table: aws_ec2_spot_fleet_request - Query AWS EC2 Spot Fleet Requests using SQL"
description: "Allows users to query AWS EC2 Spot Fleet Requests for comprehensive data on each spot fleet request, including configuration, state, capacity, and more."
folder: "EC2"
---

# Table: aws_ec2_spot_fleet_request - Query AWS EC2 Spot Fleet Requests using SQL

The AWS EC2 Spot Fleet Request is a component of Amazon's Elastic Compute Cloud (EC2) that allows you to request and manage a fleet of Spot Instances. Spot Fleet is a collection of Spot Instances and optionally On-Demand Instances that automatically maintains the target capacity of your workload. It provides a way to launch and maintain the desired number of instances by automatically replenishing any interrupted Spot Instances.

## Table Usage Guide

The `aws_ec2_spot_fleet_request` table in Steampipe provides you with information about Spot Fleet requests within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query spot fleet request-specific details, including request state, configuration, capacity settings, and associated metadata. You can utilize this table to gather insights on spot fleet requests, such as their current state, allocation strategies, pricing configurations, and more. The schema outlines the various attributes of the EC2 spot fleet request for you, including the request ID, state, target capacity, and associated tags.

## Examples

### Basic info
Explore which spot fleet requests exist in your AWS EC2 service and their current states. This query can be particularly useful in understanding the distribution of spot fleet resources and their operational status.

```sql+postgres
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  activity_status,
  type,
  target_capacity,
  allocation_strategy,
  create_time
from
  aws_ec2_spot_fleet_request;
```

```sql+sqlite
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  activity_status,
  type,
  target_capacity,
  allocation_strategy,
  create_time
from
  aws_ec2_spot_fleet_request;
```

### List active spot fleet requests
Identify spot fleet requests that are currently active and running. This can be useful for monitoring active workloads and understanding current resource utilization.

```sql+postgres
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  activity_status,
  target_capacity,
  fulfilled_capacity,
  on_demand_target_capacity,
  on_demand_fulfilled_capacity
from
  aws_ec2_spot_fleet_request
where
  spot_fleet_request_state = 'active';
```

```sql+sqlite
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  activity_status,
  target_capacity,
  fulfilled_capacity,
  on_demand_target_capacity,
  on_demand_fulfilled_capacity
from
  aws_ec2_spot_fleet_request
where
  spot_fleet_request_state = 'active';
```

### Spot fleet requests with pricing information
Analyze spot fleet requests to understand their pricing configuration and cost management settings. This helps in optimizing costs and understanding spending patterns.

```sql+postgres
select
  spot_fleet_request_id,
  spot_price,
  spot_max_total_price,
  on_demand_max_total_price,
  allocation_strategy,
  target_capacity
from
  aws_ec2_spot_fleet_request
where
  spot_price is not null
  or spot_max_total_price is not null
  or on_demand_max_total_price is not null;
```

```sql+sqlite
select
  spot_fleet_request_id,
  spot_price,
  spot_max_total_price,
  on_demand_max_total_price,
  allocation_strategy,
  target_capacity
from
  aws_ec2_spot_fleet_request
where
  spot_price is not null
  or spot_max_total_price is not null
  or on_demand_max_total_price is not null;
```

### Spot fleet requests by allocation strategy
Group spot fleet requests by their allocation strategy to understand how different strategies are being used across your infrastructure. This can help in optimizing resource allocation and cost management.

```sql+postgres
select
  allocation_strategy,
  count(*) as request_count,
  sum(target_capacity) as total_target_capacity,
  avg(target_capacity) as avg_target_capacity
from
  aws_ec2_spot_fleet_request
group by
  allocation_strategy
order by
  request_count desc;
```

```sql+sqlite
select
  allocation_strategy,
  count(*) as request_count,
  sum(target_capacity) as total_target_capacity,
  avg(target_capacity) as avg_target_capacity
from
  aws_ec2_spot_fleet_request
group by
  allocation_strategy
order by
  request_count desc;
```

### Spot fleet requests with maintenance strategies
Identify spot fleet requests that have maintenance strategies configured. This is important for understanding how your infrastructure handles instance interruptions and maintenance events.

```sql+postgres
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  spot_maintenance_strategies,
  instance_interruption_behavior,
  replace_unhealthy_instances
from
  aws_ec2_spot_fleet_request
where
  spot_maintenance_strategies is not null;
```

```sql+sqlite
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  spot_maintenance_strategies,
  instance_interruption_behavior,
  replace_unhealthy_instances
from
  aws_ec2_spot_fleet_request
where
  spot_maintenance_strategies is not null;
```

### Spot fleet requests with load balancer configuration
Find spot fleet requests that are configured with load balancers. This helps in understanding how your spot fleet workloads are distributed and load balanced.

```sql+postgres
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  load_balancers_config,
  target_capacity
from
  aws_ec2_spot_fleet_request
where
  load_balancers_config is not null;
```

```sql+sqlite
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  load_balancers_config,
  target_capacity
from
  aws_ec2_spot_fleet_request
where
  load_balancers_config is not null;
```

### Spot fleet requests with expiration dates
Identify spot fleet requests that have expiration dates set. This is important for understanding which requests will automatically terminate and when.

```sql+postgres
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  valid_from,
  valid_until,
  terminate_instances_with_expiration
from
  aws_ec2_spot_fleet_request
where
  valid_until is not null
order by
  valid_until;
```

```sql+sqlite
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  valid_from,
  valid_until,
  terminate_instances_with_expiration
from
  aws_ec2_spot_fleet_request
where
  valid_until is not null
order by
  valid_until;
```

### Spot fleet requests with launch template configurations
Analyze spot fleet requests that use launch templates for instance configuration. This provides insights into how instances are being launched and configured.

```sql+postgres
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  launch_template_configs,
  launch_specifications
from
  aws_ec2_spot_fleet_request
where
  launch_template_configs is not null
  or launch_specifications is not null;
```

```sql+sqlite
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  launch_template_configs,
  launch_specifications
from
  aws_ec2_spot_fleet_request
where
  launch_template_configs is not null
  or launch_specifications is not null;
```

### Spot fleet requests by creation date
Analyze spot fleet requests by their creation date to understand trends in spot fleet usage over time. This can help in capacity planning and resource management.

```sql+postgres
select
  date_trunc('day', create_time) as creation_date,
  count(*) as request_count,
  sum(target_capacity) as total_target_capacity
from
  aws_ec2_spot_fleet_request
group by
  date_trunc('day', create_time)
order by
  creation_date desc;
```

```sql+sqlite
select
  date(create_time) as creation_date,
  count(*) as request_count,
  sum(target_capacity) as total_target_capacity
from
  aws_ec2_spot_fleet_request
group by
  date(create_time)
order by
  creation_date desc;
```

### Spot fleet requests with excess capacity termination policy
Identify spot fleet requests that have specific policies for handling excess capacity. This is important for understanding how your infrastructure scales down when demand decreases.

```sql+postgres
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  excess_capacity_termination_policy,
  target_capacity,
  fulfilled_capacity
from
  aws_ec2_spot_fleet_request
where
  excess_capacity_termination_policy is not null;
```

```sql+sqlite
select
  spot_fleet_request_id,
  spot_fleet_request_state,
  excess_capacity_termination_policy,
  target_capacity,
  fulfilled_capacity
from
  aws_ec2_spot_fleet_request
where
  excess_capacity_termination_policy is not null;
```
