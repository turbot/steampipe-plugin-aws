---
title: "Table: aws_ebs_volume - Query AWS Elastic Block Store (EBS) using SQL"
description: "Allows users to query AWS Elastic Block Store (EBS) volumes for detailed information about their configuration, status, and associated tags."
---

# Table: aws_ebs_volume - Query AWS Elastic Block Store (EBS) using SQL

The `aws_ebs_volume` table in Steampipe provides information about volumes within AWS Elastic Block Store (EBS). This table allows DevOps engineers to query volume-specific details, including size, state, type, and associated metadata. Users can utilize this table to gather insights on volumes, such as their encryption status, IOPS performance, and snapshot details. The schema outlines the various attributes of the EBS volume, including the volume ID, creation time, attached instances, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ebs_volume` table, you can use the `.inspect aws_ebs_volume` command in Steampipe.

**Key columns**:

- `volume_id`: This is the unique identifier for the EBS volume. It is a key column for joining with other tables that reference EBS volumes.
- `state`: This column provides the current state of the EBS volume (e.g., creating, available, in-use, deleting, deleted, error). This information can be useful for monitoring and managing EBS volumes.
- `tags`: This column contains metadata that you can create and assign to your AWS resources. It can be used to manage and filter EBS volumes based on custom criteria.

## Examples

### List of unencrypted EBS volumes

```sql
select
  volume_id,
  encrypted
from
  aws_ebs_volume
where
  not encrypted;
```

### List of unattached EBS volumes

```sql
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  jsonb_array_length(attachments) = 0;
```

### List of Provisioned IOPS SSD (io1) volumes

```sql
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  volume_type = 'io1';
```

### List of EBS volumes with size more than 100GiB

```sql
select
  volume_id,
  size
from
  aws_ebs_volume
where
  size > '100';
```

### Count the number of EBS volumes by volume type

```sql
select
  volume_type,
  count(volume_type) as count
from
  aws_ebs_volume
group by
  volume_type;
```

### Find EBS Volumes Attached To Stopped EC2 Instances

```sql
select
  volume_id,
  size,
  att ->> 'InstanceId' as instance_id
from
  aws_ebs_volume
  cross join jsonb_array_elements(attachments) as att
  join aws_ec2_instance as i on i.instance_id = att ->> 'InstanceId'
where
  instance_state = 'stopped';
```
