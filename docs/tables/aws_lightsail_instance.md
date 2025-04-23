---
title: "Steampipe Table: aws_lightsail_instance - Query AWS Lightsail Instances using SQL"
description: "Allows users to query AWS Lightsail Instances and retrieve detailed information about each instance such as instance state, type, associated bundles, and more."
folder: "Lightsail"
---

# Table: aws_lightsail_instance - Query AWS Lightsail Instances using SQL

The AWS Lightsail Instance is a part of Amazon Lightsail service, providing a simple virtual private server solution. It offers easy-to-use instances with a variety of applications and stacks, including WordPress, Joomla, and more. These instances are ideal for simpler workloads, quick deployments, and developers seeking a smooth transition to the cloud.

## Table Usage Guide

The `aws_lightsail_instance` table in Steampipe provides you with information about instances within AWS Lightsail. This table allows you, as a DevOps engineer, to query instance-specific details, including instance state, attached static IP, associated bundles, and associated tags. You can utilize this table to gather insights on instances, such as instances with specific tags, instances in a certain state, instances associated with a specific bundle, and more. The schema outlines the various attributes of the Lightsail instance for you, including the instance name, creation timestamp, location, blueprint ID, and more.

## Examples

### Instance count in each availability zone
Determine the distribution of instances across different availability zones to effectively manage resources and optimize performance.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which AWS Lightsail instances have been inactive for more than 30 days. This query is useful for identifying potential resource inefficiencies and cost-saving opportunities.

```sql+postgres
select
  name,
  state_name
from
  aws_lightsail_instance
where
  state_name = 'stopped'
  and created_at <= (current_date - interval '30' day);
```

```sql+sqlite
select
  name,
  state_name
from
  aws_lightsail_instance
where
  state_name = 'stopped'
  and created_at <= date('now','-30 day');
```

### List public instances
Identify instances where your AWS Lightsail instances are publicly accessible. This is useful for reviewing your network security and ensuring your instances are not exposed to unnecessary risks.

```sql+postgres
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

```sql+sqlite
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
Analyze your AWS Lightsail instances to identify those that do not have an 'application' tag assigned. This can help streamline your resource management and ensure consistent tagging practices across your cloud environment.

```sql+postgres
select
  name,
  tags
from
  aws_lightsail_instance
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  name,
  tags
from
  aws_lightsail_instance
where
  not json_valid(tags) or json_extract(tags, '$.application') is null;
```

### Hardware specifications of the instances
Explore the hardware specifications of your instances to assess their computing power and memory capacity. This is particularly useful in optimizing resource allocation and performance in your AWS Lightsail instances.

```sql+postgres
select
  name,
  hardware ->> 'CpuCount' as "CPU Count",
  hardware ->> 'RamSizeInGb' as "RAM Size (in GB)"
from
  aws_lightsail_instance;
```

```sql+sqlite
select
  name,
  json_extract(hardware, '$.CpuCount') as "CPU Count",
  json_extract(hardware, '$.RamSizeInGb') as "RAM Size (in GB)"
from
  aws_lightsail_instance;
```