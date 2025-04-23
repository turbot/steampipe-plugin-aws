---
title: "Steampipe Table: aws_emr_instance_group - Query AWS EMR Instance Groups using SQL"
description: "Allows users to query AWS EMR Instance Groups to fetch details about each instance group within an EMR cluster."
folder: "EMR"
---

# Table: aws_emr_instance_group - Query AWS EMR Instance Groups using SQL

The AWS Elastic MapReduce (EMR) Instance Group is a component of Amazon EMR that organizes EC2 instances in a cluster. It is used to host big data frameworks like Apache Hadoop, Spark, HBase, and others for processing vast amounts of data. These groups can be resized manually or automatically, depending on the work requirements.

## Table Usage Guide

The `aws_emr_instance_group` table in Steampipe provides you with information about instance groups within AWS Elastic MapReduce (EMR). This table allows you, as a DevOps engineer, to query instance group-specific details, including instance group ID, instance type, instance count, and associated metadata. You can utilize this table to gather insights on instance groups, such as their current status, market type, and more. The schema outlines the various attributes of the EMR instance group, including the cluster ID, instance group type, EBS volumes, and associated tags for your convenience.

## Examples

### Basic info
Explore which Amazon EMR instance groups are currently active and their types. This can be useful for managing resources and understanding the state of your EMR clusters.

```sql+postgres
select
  id,
  arn,
  cluster_id,
  instance_group_type,
  state
from
  aws_emr_instance_group;
```

```sql+sqlite
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
Identify the type of master instances used in a cluster to better understand your resource usage and optimize your configurations.

```sql+postgres
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

```sql+sqlite
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
Explore the distribution of active instances across different clusters to effectively manage resources and ensure optimal performance. This can help in identifying clusters that might be overburdened or underutilized.

```sql+postgres
select
  cluster_id,
  sum(running_instance_count) as running_instance_count
from
  aws_emr_instance_group
where
  state = 'RUNNING'
group by cluster_id;
```

```sql+sqlite
select
  cluster_id,
  sum(running_instance_count) as running_instance_count
from
  aws_emr_instance_group
where
  state = 'RUNNING'
group by cluster_id;
```