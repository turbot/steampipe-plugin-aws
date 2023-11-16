---
title: "Table: aws_emr_instance - Query AWS EMR Instances using SQL"
description: "Allows users to query AWS EMR Instances for detailed information about the status, configuration, and other metadata of each instance."
---

# Table: aws_emr_instance - Query AWS EMR Instances using SQL

The `aws_emr_instance` table in Steampipe provides information about instances within AWS Elastic MapReduce (EMR). This table allows DevOps engineers to query instance-specific details, including instance status, instance group ID, and associated metadata. Users can utilize this table to gather insights on instances, such as instance health status, instance configuration details, and more. The schema outlines the various attributes of the EMR instance, including the instance ID, EBS volumes, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_emr_instance` table, you can use the `.inspect aws_emr_instance` command in Steampipe.

**Key columns**:

- `id`: The unique identifier for the instance. This can be used to join with other tables that reference the instance ID.
- `instance_group_id`: The ID of the instance group to which this instance belongs. This can be useful for querying information about the instance group.
- `status_state`: The current state of the instance. This is important for monitoring and managing the lifecycle of the instance.

## Examples

### Basic info

```sql
select
  id,
  cluster_id,
  ec2_instance_id,
  instance_type,
  private_dns_name,
  private_ip_address
from
  aws_emr_instance;
```

### List instances by type

```sql
select
  id,
  ec2_instance_id,
  instance_type
from
  aws_emr_instance
where
  instance_type = 'm2.4xlarge';
```

### List instances for a cluster

```sql
select
  id,
  ec2_instance_id,
  instance_type
from
  aws_emr_instance
where
  cluster_id = 'j-21HIX5R2NZMXJ';
```

### Get volume details for an instance

```sql
select
  id,
  ec2_instance_id,
  instance_type,
  v -> 'Device' as device,
  v -> 'VolumeId' as volume_id
from
  aws_emr_instance,
  jsonb_array_elements(ebs_volumes) as v
where
  id = 'ci-ULCFS2ZN0FK7';
```