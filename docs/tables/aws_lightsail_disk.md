---
title: "Steampipe Table: aws_lightsail_disk - Query AWS Lightsail Disks using SQL"
description: "Allows users to query AWS Lightsail Disks for detailed information about block storage volumes, including their size, state, and associated metadata."
folder: "Lightsail"
---

# Table: aws_lightsail_disk - Query AWS Lightsail Disks using SQL

The AWS Lightsail Disk is a feature of Amazon Lightsail that provides block storage volumes for your Lightsail instances. These disks can be attached to Lightsail instances to provide additional storage capacity. AWS Lightsail offers both SSD and HDD disk options with various size configurations.

## Table Usage Guide

The `aws_lightsail_disk` table in Steampipe provides you with information about AWS Lightsail disks. This table allows you as a DevOps engineer to query disk-specific details, including the disk name, size, state, and associated metadata. You can utilize this table to gather insights on disks, such as their attachment status, IOPS performance, and any associated tags. The schema outlines the various attributes of the AWS Lightsail disk for you, including the disk ARN, location information, and associated tags.

## Examples

### Basic info
Explore the basic information about your Lightsail disks, including their names, sizes, and states. This can help you understand the current state of your disks and identify any that might need attention.

```sql+postgres
select
  name,
  size_in_gb,
  state,
  is_attached,
  iops
from
  aws_lightsail_disk;
```

```sql+sqlite
select
  name,
  size_in_gb,
  state,
  is_attached,
  iops
from
  aws_lightsail_disk;
```

### List unattached disks
Identify disks that are not currently attached to any instance to help manage your storage resources and potentially reduce costs.

```sql+postgres
select
  name,
  size_in_gb,
  state,
  created_at
from
  aws_lightsail_disk
where
  not is_attached;
```

```sql+sqlite
select
  name,
  size_in_gb,
  state,
  created_at
from
  aws_lightsail_disk
where
  not is_attached;
```

### List disks by size
Analyze the distribution of disks by their size to understand your storage allocation patterns.

```sql+postgres
select
  size_in_gb,
  count(*) as disk_count
from
  aws_lightsail_disk
group by
  size_in_gb
order by
  size_in_gb;
```

```sql+sqlite
select
  size_in_gb,
  count(*) as disk_count
from
  aws_lightsail_disk
group by
  size_in_gb
order by
  size_in_gb;
```

### List disks with specific tags
Find disks that have specific tags associated with them to help organize and manage your storage resources based on custom criteria.

```sql+postgres
select
  name,
  size_in_gb,
  state,
  tags
from
  aws_lightsail_disk
where
  tags ->> 'Environment' = 'Production';
```

```sql+sqlite
select
  name,
  size_in_gb,
  state,
  tags
from
  aws_lightsail_disk
where
  json_extract(tags, '$.Environment') = 'Production';
```

### List disks with high IOPS
Identify disks with high IOPS values to understand your high-performance storage requirements.

```sql+postgres
select
  name,
  size_in_gb,
  iops,
  state
from
  aws_lightsail_disk
where
  iops > 1000
order by
  iops desc;
```

```sql+sqlite
select
  name,
  size_in_gb,
  iops,
  state
from
  aws_lightsail_disk
where
  iops > 1000
order by
  iops desc;
```

### List disks by location
View the distribution of disks across different AWS regions and availability zones.

```sql+postgres
select
  name,
  size_in_gb,
  state,
  location
from
  aws_lightsail_disk
order by
  location;
```

```sql+sqlite
select
  name,
  size_in_gb,
  state,
  location
from
  aws_lightsail_disk
order by
  location;
``` 