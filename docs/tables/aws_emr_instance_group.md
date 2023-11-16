---
title: "Table: aws_emr_instance_group - Query AWS EMR Instance Groups using SQL"
description: "Allows users to query AWS EMR Instance Groups to fetch details about each instance group within an EMR cluster."
---

# Table: aws_emr_instance_group - Query AWS EMR Instance Groups using SQL

The `aws_emr_instance_group` table in Steampipe provides information about instance groups within AWS Elastic MapReduce (EMR). This table allows DevOps engineers to query instance group-specific details, including instance group ID, instance type, instance count, and associated metadata. Users can utilize this table to gather insights on instance groups, such as their current status, market type, and more. The schema outlines the various attributes of the EMR instance group, including the cluster ID, instance group type, EBS volumes, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_emr_instance_group` table, you can use the `.inspect aws_emr_instance_group` command in Steampipe.

**Key columns**:

- `cluster_id`: The unique identifier for the cluster. This column can be used to join with other tables that contain information about EMR clusters.
- `instance_group_id`: The unique identifier for the instance group. This is critical when joining with tables that contain additional details on specific instance groups within an EMR cluster.
- `instance_type`: The type of instance group. This column is useful to understand the resources allocated to the instance group, which can help in capacity planning and cost analysis.

## Examples

### Basic info

```sql
select
  id,
  arn,
  cluster_id,
  instance_group_type,
  state
from
  aws_emr_instance_group;
```

### Get the master instance type used for a cluster

```sql
select
  ig.id as instance_group_id,
  ig.cluster_id,
  c.name as cluster_name,
  ig.instance_type
from
  aws_emr_instance_group as ig,
  aws_emr_cluster as c
where
  ig.cluster_id = c.id
  and ig.instance_group_type = 'MASTER';
```

### Get the count of running instances (core and master) per cluster

```sql
select
  cluster_id,
  sum(running_instance_count) as running_instance_count
from
  aws_emr_instance_group
where
  state = 'RUNNING'
group by cluster_id;
```
