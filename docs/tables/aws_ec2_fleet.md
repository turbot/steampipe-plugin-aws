---
title: "Steampipe Table: aws_ec2_fleet - Query AWS EC2 Fleets using SQL"
description: "Allows users to query AWS EC2 Fleets to retrieve information about EC2 Fleet configurations, capacity, and instances."
folder: "EC2"
---

# Table: aws_ec2_fleet - Query AWS EC2 Fleets using SQL

An AWS EC2 Fleet contains the configuration information to launch a fleet—or group—of instances. In a single API call, a fleet can launch multiple instance types across multiple Availability Zones, using the On-Demand Instance, Reserved Instance, and Spot Instance purchasing options together. Using EC2 Fleet, you can define separate On-Demand and Spot capacity targets, specify the instance types that work best for your applications, and specify how Amazon EC2 should distribute your fleet capacity within each purchasing option.

## Table Usage Guide

The `aws_ec2_fleet` table in Steampipe provides you with information about EC2 Fleets within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query fleet-specific details, including fleet state, capacity specifications, launch template configurations, and associated metadata. You can utilize this table to gather insights on fleets, such as fleet capacity utilization, instance distribution, spot vs on-demand allocation, and more. The schema outlines the various attributes of the EC2 Fleet for you, including the fleet ID, creation time, target capacity, and associated tags.

**Important:** EC2 Fleets of type `instant` are not returned by the list operation. To query an instant fleet, you must specify the `fleet_id` in the WHERE clause.

## Examples

### Basic info
Explore the status and type of your EC2 fleets to understand their current operational state and configuration. This can help in managing resources and planning capacity.

```sql+postgres
select
  fleet_id,
  arn,
  fleet_state,
  type,
  create_time
from
  aws_ec2_fleet;
```

```sql+sqlite
select
  fleet_id,
  arn,
  fleet_state,
  type,
  create_time
from
  aws_ec2_fleet;
```

### List active EC2 fleets
Identify active EC2 fleets to monitor ongoing fleet operations and ensure they are functioning as expected.

```sql+postgres
select
  fleet_id,
  arn,
  fleet_state,
  activity_status,
  type
from
  aws_ec2_fleet
where
  fleet_state = 'active';
```

```sql+sqlite
select
  fleet_id,
  arn,
  fleet_state,
  activity_status,
  type
from
  aws_ec2_fleet
where
  fleet_state = 'active';
```

### List fleets with unfulfilled capacity
Discover fleets that have not yet reached their target capacity, which may indicate issues with instance availability or configuration.

```sql+postgres
select
  fleet_id,
  fleet_state,
  fulfilled_capacity,
  target_capacity_specification ->> 'TotalTargetCapacity' as total_target_capacity
from
  aws_ec2_fleet
where
  fulfilled_capacity < (target_capacity_specification ->> 'TotalTargetCapacity')::float;
```

```sql+sqlite
select
  fleet_id,
  fleet_state,
  fulfilled_capacity,
  json_extract(target_capacity_specification, '$.TotalTargetCapacity') as total_target_capacity
from
  aws_ec2_fleet
where
  fulfilled_capacity < cast(json_extract(target_capacity_specification, '$.TotalTargetCapacity') as real);
```

### Get target capacity details for each fleet
Analyze the target capacity configuration of your EC2 fleets to understand the distribution between On-Demand and Spot instances.

```sql+postgres
select
  fleet_id,
  fleet_state,
  target_capacity_specification ->> 'TotalTargetCapacity' as total_target_capacity,
  target_capacity_specification ->> 'OnDemandTargetCapacity' as on_demand_target_capacity,
  target_capacity_specification ->> 'SpotTargetCapacity' as spot_target_capacity,
  target_capacity_specification ->> 'DefaultTargetCapacityType' as default_target_capacity_type
from
  aws_ec2_fleet;
```

```sql+sqlite
select
  fleet_id,
  fleet_state,
  json_extract(target_capacity_specification, '$.TotalTargetCapacity') as total_target_capacity,
  json_extract(target_capacity_specification, '$.OnDemandTargetCapacity') as on_demand_target_capacity,
  json_extract(target_capacity_specification, '$.SpotTargetCapacity') as spot_target_capacity,
  json_extract(target_capacity_specification, '$.DefaultTargetCapacityType') as default_target_capacity_type
from
  aws_ec2_fleet;
```

### List fleets that will terminate instances on expiration
Identify fleets configured to automatically terminate instances when the fleet expires, which is important for cost management.

```sql+postgres
select
  fleet_id,
  fleet_state,
  type,
  terminate_instances_with_expiration,
  valid_until
from
  aws_ec2_fleet
where
  terminate_instances_with_expiration = true;
```

```sql+sqlite
select
  fleet_id,
  fleet_state,
  type,
  terminate_instances_with_expiration,
  valid_until
from
  aws_ec2_fleet
where
  terminate_instances_with_expiration = 1;
```

### Get spot options for each fleet
Examine the Spot Instance configuration for your EC2 fleets to understand allocation strategies and pricing settings.

```sql+postgres
select
  fleet_id,
  fleet_state,
  spot_options ->> 'AllocationStrategy' as spot_allocation_strategy,
  spot_options ->> 'InstanceInterruptionBehavior' as instance_interruption_behavior,
  spot_options ->> 'MaxTotalPrice' as spot_max_total_price
from
  aws_ec2_fleet
where
  spot_options is not null;
```

```sql+sqlite
select
  fleet_id,
  fleet_state,
  json_extract(spot_options, '$.AllocationStrategy') as spot_allocation_strategy,
  json_extract(spot_options, '$.InstanceInterruptionBehavior') as instance_interruption_behavior,
  json_extract(spot_options, '$.MaxTotalPrice') as spot_max_total_price
from
  aws_ec2_fleet
where
  spot_options is not null;
```

### List fleets with errors
Identify fleets that encountered errors during instance launch, which can help troubleshoot capacity or configuration issues.

```sql+postgres
select
  fleet_id,
  fleet_state,
  e ->> 'ErrorCode' as error_code,
  e ->> 'ErrorMessage' as error_message
from
  aws_ec2_fleet,
  jsonb_array_elements(errors) as e
where
  errors is not null
  and jsonb_array_length(errors) > 0;
```

```sql+sqlite
select
  fleet_id,
  fleet_state,
  json_extract(e.value, '$.ErrorCode') as error_code,
  json_extract(e.value, '$.ErrorMessage') as error_message
from
  aws_ec2_fleet,
  json_each(errors) as e
where
  errors is not null
  and json_array_length(errors) > 0;
```

### Get EC2 fleet by ID
Retrieve details for a specific EC2 fleet using its fleet ID.

```sql+postgres
select
  fleet_id,
  arn,
  fleet_state,
  type,
  create_time,
  target_capacity_specification
from
  aws_ec2_fleet
where
  fleet_id = 'fleet-12345678901234567';
```

```sql+sqlite
select
  fleet_id,
  arn,
  fleet_state,
  type,
  create_time,
  target_capacity_specification
from
  aws_ec2_fleet
where
  fleet_id = 'fleet-12345678901234567';
```

### Get instant fleet by ID
Retrieve details for an instant EC2 fleet. Note that instant fleets are not returned by the list operation and must be queried by their fleet ID.

```sql+postgres
select
  fleet_id,
  arn,
  fleet_state,
  type,
  create_time,
  instances,
  target_capacity_specification
from
  aws_ec2_fleet
where
  fleet_id = 'fleet-219fbc9f-31a4-6e8d-aeb0-2b0854a6eea7';
```

```sql+sqlite
select
  fleet_id,
  arn,
  fleet_state,
  type,
  create_time,
  instances,
  target_capacity_specification
from
  aws_ec2_fleet
where
  fleet_id = 'fleet-219fbc9f-31a4-6e8d-aeb0-2b0854a6eea7';
```
