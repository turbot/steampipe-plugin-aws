---
title: "Table: aws_emr_instance_fleet - Query AWS EMR Instance Fleets using SQL"
description: "Allows users to query AWS EMR Instance Fleets to obtain detailed information about each instance fleet, including its configuration, instance type specifications, target capacities, and associated metadata."
---

# Table: aws_emr_instance_fleet - Query AWS EMR Instance Fleets using SQL

The `aws_emr_instance_fleet` table in Steampipe provides information about instance fleets within AWS Elastic MapReduce (EMR). This table allows DevOps engineers to query instance fleet-specific details, including instance type specifications, target capacities, and associated metadata. Users can utilize this table to gather insights on instance fleets, such as the current status of instance fleets, the target and provisioned capacities of on-demand and spot instances, and the instance type configurations. The schema outlines the various attributes of the EMR instance fleet, including the fleet ID, cluster ID, name, state, instance type specifications, target capacities, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_emr_instance_fleet` table, you can use the `.inspect aws_emr_instance_fleet` command in Steampipe.

**Key columns**:

- `fleet_id`: The unique identifier of the instance fleet. This column is important as it can be used to join this table with other tables to fetch more specific data related to each instance fleet.
- `cluster_id`: The unique identifier of the cluster. This column is useful to join with other tables to obtain detailed information about the cluster associated with each instance fleet.
- `instance_fleet_type`: The type of instance fleet. This column is useful to filter the data based on the type of instance fleet (TASK, CORE, or MASTER).

## Examples

### Basic info

```sql
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

```sql
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

### Get the provisioned & target on demand capacity of the instance fleets

```sql
select
  cluster_id,
  provisioned_on_demand_capacity,
  target_on_demand_capacity
from
  aws_emr_instance_fleet;
```
