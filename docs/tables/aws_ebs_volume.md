---
title: "Steampipe Table: aws_ebs_volume - Query AWS Elastic Block Store (EBS) using SQL"
description: "Allows users to query AWS Elastic Block Store (EBS) volumes for detailed information about their configuration, status, and associated tags."
folder: "EBS"
---

# Table: aws_ebs_volume - Query AWS Elastic Block Store (EBS) using SQL

The AWS Elastic Block Store (EBS) is a high-performance block storage service designed for use with Amazon Elastic Compute Cloud (EC2) for both throughput and transaction intensive workloads at any scale. It provides persistent block-level storage volumes for use with Amazon EC2 instances. EBS volumes are highly available and reliable storage volumes that can be attached to any running instance and used like a physical hard drive.

## Table Usage Guide

The `aws_ebs_volume` table in Steampipe provides you with information about volumes within AWS Elastic Block Store (EBS). This table allows you, as a DevOps engineer, to query volume-specific details, including size, state, type, and associated metadata. You can utilize this table to gather insights on volumes, such as their encryption status, IOPS performance, and snapshot details. The schema outlines the various attributes of the EBS volume for you, including the volume ID, creation time, attached instances, and associated tags.

## Examples

### List of unencrypted EBS volumes
Identify instances where EBS volumes in your AWS environment are not encrypted. This is crucial for security audits and ensuring compliance with data protection policies.

```sql+postgres
select
  volume_id,
  encrypted
from
  aws_ebs_volume
where
  not encrypted;
```

```sql+sqlite
select
  volume_id,
  encrypted
from
  aws_ebs_volume
where
  encrypted = 0;
```

### List of unattached EBS volumes
Identify instances where EBS volumes in AWS are not attached to any instances. This could help in optimizing resource usage and managing costs by removing unnecessary volumes.

```sql+postgres
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  jsonb_array_length(attachments) = 0;
```

```sql+sqlite
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  json_array_length(attachments) = 0;
```

### List of Provisioned IOPS SSD (io1) volumes
Determine the areas in which Provisioned IOPS SSD (io1) volumes are being used in your AWS infrastructure. This information can help optimize storage performance and costs by identifying potential areas for volume type adjustment.

```sql+postgres
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  volume_type = 'io1';
```

```sql+sqlite
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  volume_type = 'io1';
```

### List of EBS volumes with size more than 100GiB
Identify instances where AWS EBS volumes exceed 100GiB in size. This is useful to manage storage resources and prevent excessive usage.

```sql+postgres
select
  volume_id,
  size
from
  aws_ebs_volume
where
  size > '100';
```

```sql+sqlite
select
  volume_id,
  size
from
  aws_ebs_volume
where
  size > 100;
```

### Count the number of EBS volumes by volume type
Identify the distribution of different types of EBS volumes in your AWS environment. This helps in understanding the usage patterns and planning for cost optimization.

```sql+postgres
select
  volume_type,
  count(volume_type) as count
from
  aws_ebs_volume
group by
  volume_type;
```

```sql+sqlite
select
  volume_type,
  count(volume_type) as count
from
  aws_ebs_volume
group by
  volume_type;
```

### Find EBS Volumes Attached To Stopped EC2 Instances
Discover the segments that include EBS volumes attached to EC2 instances that are currently in a stopped state. This information can be beneficial to optimize resource allocation and reduce unnecessary costs.

```sql+postgres
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

```sql+sqlite
select
  volume_id,
  size,
  json_extract(att.value, '$.InstanceId') as instance_id
from
  aws_ebs_volume
  join json_each(attachments) as att
  join aws_ec2_instance as i on i.instance_id = json_extract(att.value, '$.InstanceId')
where
  instance_state = 'stopped';
```

### List of Provisioned IOPS SSD (io1) volumes
Identify instances where the SSD volumes with provisioned IOPS (IO1) are being used. This could be beneficial for performance optimization and cost management.

```sql+postgres
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  volume_type = 'io1';
```

```sql+sqlite
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  volume_type = 'io1';
```