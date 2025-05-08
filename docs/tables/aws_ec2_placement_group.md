---
title: "Steampipe Table: aws_ec2_placement_group - Query AWS EC2 Placement Groups using SQL"
description: "Allows users to query AWS EC2 Placement Groups, providing information about placement strategies, state, and configuration."
folder: "EC2"
---

# Table: aws_ec2_placement_group - Query AWS EC2 Placement Groups using SQL

The AWS EC2 Placement Group is a logical grouping of instances within a single Availability Zone or across partitions or racks, depending on the strategy. Placement groups help you influence the placement of a group of interdependent instances to meet the needs of your workload.

## Table Usage Guide

The `aws_ec2_placement_group` table in Steampipe provides you with information about Placement Groups within AWS EC2. This table allows you, as a DevOps engineer, cloud architect, or system administrator, to query placement group-specific details, including strategy, partition count, spread level, state, and associated tags. You can utilize this table to gather insights on placement groups, such as their configuration, current state, and more. The schema outlines the various attributes of the EC2 placement group for you, including the group name, strategy, and tags.

## Examples

### List all placement groups and their strategies
Discover all EC2 placement groups in your account, including their placement strategy and current state. This helps you understand how your workloads are distributed and managed.

```sql+postgres
select
  group_name,
  strategy,
  state,
  region
from
  aws_ec2_placement_group;
```

```sql+sqlite
select
  group_name,
  strategy,
  state,
  region
from
  aws_ec2_placement_group;
```

### List all available partition placement groups
Find all partition placement groups that are currently available. This is useful for identifying groups that can be used for high-availability workloads.

```sql+postgres
select
  group_name,
  partition_count,
  state
from
  aws_ec2_placement_group
where
  strategy = 'partition'
  and state = 'available';
```

```sql+sqlite
select
  group_name,
  partition_count,
  state
from
  aws_ec2_placement_group
where
  strategy = 'partition'
  and state = 'available';
```

### Count of placement groups by strategy
See how many placement groups use each strategy (cluster, partition, or spread). This helps you assess your architecture's distribution and redundancy.

```sql+postgres
select
  strategy,
  count(*) as count
from
  aws_ec2_placement_group
group by
  strategy;
```

```sql+sqlite
select
  strategy,
  count(*) as count
from
  aws_ec2_placement_group
group by
  strategy;
```

### List placement groups with fewer than 3 partitions
Identify partition placement groups that have fewer than 3 partitions. This can help you find groups that may not meet your high-availability requirements.

```sql+postgres
select
  group_name,
  partition_count,
  strategy
from
  aws_ec2_placement_group
where
  partition_count < 3;
```

```sql+sqlite
select
  group_name,
  partition_count,
  strategy
from
  aws_ec2_placement_group
where
  partition_count < 3;
```