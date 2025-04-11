---
title: "Steampipe Table: aws_emr_instance_fleet - Query AWS EMR Instance Fleets using SQL"
description: "Allows users to query AWS EMR Instance Fleets to obtain detailed information about each instance fleet, including its configuration, instance type specifications, target capacities, and associated metadata."
folder: "EMR"
---

# Table: aws_emr_instance_fleet - Query AWS EMR Instance Fleets using SQL

An AWS EMR Instance Fleet is a component of Amazon EMR that you use to specify the EC2 instances and instance types that Amazon EMR provisions to create a cluster. Instance fleets provide a way to specify a diverse set of instances to accommodate workloads that benefit from a variety of EC2 instance types. They allow you to specify target capacities for On-Demand and Spot instances in terms of instances, vCPUs, or ECUs.

## Table Usage Guide

The `aws_emr_instance_fleet` table in Steampipe provides you with information about instance fleets within AWS Elastic MapReduce (EMR). This table allows you, as a DevOps engineer, to query instance fleet-specific details, including instance type specifications, target capacities, and associated metadata. You can utilize this table to gather insights on instance fleets, such as the current status of instance fleets, the target and provisioned capacities of on-demand and spot instances, and the instance type configurations. The schema outlines the various attributes of the EMR instance fleet for you, including the fleet ID, cluster ID, name, state, instance type specifications, target capacities, and associated tags.

## Examples

### Basic info
Explore the status and type of your EMR instance fleets in AWS to understand their current operational state and configuration.

```sql+postgres
select
  id,
  arn,
  cluster_id,
  instance_fleet_type,
  state
from
  aws_emr_instance_fleet;
```

```sql+sqlite
select
  id,
  arn,
  cluster_id,
  instance_fleet_type,
  state
from
  aws_emr_instance_fleet;
```

### Get the cluster details of the instance fleets
Discover the segments that provide a comprehensive view of the state and name of your clusters within your instance fleets. This can be useful to manage and monitor the health and status of your fleets in AWS Elastic MapReduce (EMR) service.

```sql+postgres
select
  cluster_id,
  c.name as cluster_name,
  c.state as cluster_state
from
  aws_emr_instance_fleet as f,
  aws_emr_cluster as c
where
  f.cluster_id = c.id;
```

```sql+sqlite
select
  cluster_id,
  c.name as cluster_name,
  c.state as cluster_state
from
  aws_emr_instance_fleet as f
join
  aws_emr_cluster as c
on
  f.cluster_id = c.id;
```

### Get the provisioned & target on demand capacity of the instance fleets
Determine the current and target capacity of your instance fleets to effectively manage resources and plan for future capacity needs. This allows you to optimize your resource utilization and avoid potential over-provisioning or under-provisioning.

```sql+postgres
select
  cluster_id,
  provisioned_on_demand_capacity,
  target_on_demand_capacity
from
  aws_emr_instance_fleet;
```

```sql+sqlite
select
  cluster_id,
  provisioned_on_demand_capacity,
  target_on_demand_capacity
from
  aws_emr_instance_fleet;
```