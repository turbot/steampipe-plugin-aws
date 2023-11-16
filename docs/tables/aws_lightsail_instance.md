---
title: "Table: aws_lightsail_instance - Query AWS Lightsail Instances using SQL"
description: "Allows users to query AWS Lightsail Instances and retrieve detailed information about each instance such as instance state, type, associated bundles, and more."
---

# Table: aws_lightsail_instance - Query AWS Lightsail Instances using SQL

The `aws_lightsail_instance` table in Steampipe provides information about instances within AWS Lightsail. This table allows DevOps engineers to query instance-specific details, including instance state, attached static IP, associated bundles, and associated tags. Users can utilize this table to gather insights on instances, such as instances with specific tags, instances in a certain state, instances associated with a specific bundle, and more. The schema outlines the various attributes of the Lightsail instance, including the instance name, creation timestamp, location, blueprint ID, and more.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lightsail_instance` table, you can use the `.inspect aws_lightsail_instance` command in Steampipe.

### Key columns:

- `name`: The name of the Lightsail instance. This column is useful as it provides a unique identifier for each instance, enabling you to easily join this table with others.
- `state_name`: The state of the Lightsail instance (running, stopped, etc.). This column is important for monitoring and managing instance states.
- `blueprint_id`: The ID of the blueprint that was used to create the instance. This column is useful for identifying the type of software that the instance is running.

## Examples

### Instance count in each availability zone

```sql
select
  availability_zone as az,
  bundle_id,
  count(*)
from
  aws_lightsail_instance
group by
  availability_zone,
  bundle_id;
```

### List stopped instances created for more than 30 days

```sql
select
  name,
  state_name
from
  aws_lightsail_instance
where
  state_name = 'stopped'
  and created_at <= (current_date - interval '30' day);
```

### List public instances

```sql
select
  name,
  state_name,
  bundle_id,
  region
from
  aws_lightsail_instance
where
  public_ip_address is not null;
```

### List of instances without application tag key

```sql
select
  name,
  tags
from
  aws_lightsail_instance
where
  not tags :: JSONB ? 'application';
```

### Hardware specifications of the instances

```sql
select
  name,
  hardware ->> 'CpuCount' as "CPU Count",
  hardware ->> 'RamSizeInGb' as "RAM Size (in GB)"
from
  aws_lightsail_instance;
```