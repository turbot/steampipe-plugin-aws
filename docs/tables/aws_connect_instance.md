---
title: "Steampipe Table: aws_connect_instance - Query AWS Connect Instances using SQL"
description: "Allows users to query AWS Connect instances to retrieve detailed information about each contact center instance configuration."
folder: "Connect"
---

# Table: aws_connect_instance - Query AWS Connect Instances using SQL

AWS Connect is a cloud contact center service that makes it easy for any business to deliver better customer service at lower cost. The `aws_connect_instance` table provides information about Connect instances, including their configuration, status, and metadata.

## Table Usage Guide

The `aws_connect_instance` table in Steampipe provides you with information about Connect instances within AWS. This table allows you, as a DevOps engineer, to query instance-specific details, including instance ID, alias, status, identity management type, and associated metadata. You can utilize this table to gather insights on Connect instances, such as their current status, call configuration, access URLs, and associated tags.

## Examples

### Basic info
Explore the features and settings of your AWS Connect instances to better understand their configuration, such as status, identity management type, and regional distribution. This can help in assessing instance performance and operational efficiency.

```sql+postgres
select
  instance_alias,
  id,
  arn,
  instance_status,
  identity_management_type,
  region
from
  aws_connect_instance;
```

```sql+sqlite
select
  instance_alias,
  id,
  arn,
  instance_status,
  identity_management_type,
  region
from
  aws_connect_instance;
```

### List active Connect instances
Find Connect instances that are currently active to understand which contact centers are operational.

```sql+postgres
select
  instance_alias,
  id,
  instance_status,
  created_time
from
  aws_connect_instance
where
  instance_status = 'ACTIVE';
```

```sql+sqlite
select
  instance_alias,
  id,
  instance_status,
  created_time
from
  aws_connect_instance
where
  instance_status = 'ACTIVE';
```

### Find instances with both inbound and outbound calls enabled
Identify Connect instances that have both inbound and outbound calling capabilities enabled for comprehensive contact center functionality.

```sql+postgres
select
  instance_alias,
  id,
  inbound_calls_enabled,
  outbound_calls_enabled,
  instance_access_url
from
  aws_connect_instance
where
  inbound_calls_enabled = true
  and outbound_calls_enabled = true;
```

```sql+sqlite
select
  instance_alias,
  id,
  inbound_calls_enabled,
  outbound_calls_enabled,
  instance_access_url
from
  aws_connect_instance
where
  inbound_calls_enabled = 1
  and outbound_calls_enabled = 1;
```

### Get instance details with tags
Retrieve detailed information about Connect instances including their associated tags for resource management and cost allocation.

```sql+postgres
select
  instance_alias,
  id,
  service_role,
  tags
from
  aws_connect_instance
where
  tags is not null;
```

```sql+sqlite
select
  instance_alias,
  id,
  service_role,
  tags
from
  aws_connect_instance
where
  tags is not null;
```

### List recently created instances
Find Connect instances that were created recently to track new deployments and changes in your contact center environment.

```sql+postgres
select
  instance_alias,
  id,
  created_time,
  instance_status
from
  aws_connect_instance
where
  created_time >= now() - interval '30 days'
order by
  created_time desc;
```

```sql+sqlite
select
  instance_alias,
  id,
  created_time,
  instance_status
from
  aws_connect_instance
where
  created_time >= datetime('now', '-30 days')
order by
  created_time desc;
```

### Find instances with status issues
Identify Connect instances that may have status issues or failed creation attempts for troubleshooting.

```sql+postgres
select
  instance_alias,
  id,
  instance_status,
  status_reason
from
  aws_connect_instance
where
  status_reason is not null;
```

```sql+sqlite
select
  instance_alias,
  id,
  instance_status,
  status_reason
from
  aws_connect_instance
where
  status_reason is not null;
```

### Find instances with specific call configurations
Identify Connect instances with specific call configurations, such as inbound and outbound calling capabilities.

```sql+postgres
select
  instance_alias,
  id,
  inbound_calls_enabled,
  outbound_calls_enabled,
  instance_access_url
from
  aws_connect_instance
where
  inbound_calls_enabled = true
  and outbound_calls_enabled = true;
```

```sql+sqlite
select
  instance_alias,
  id,
  inbound_calls_enabled,
  outbound_calls_enabled,
  instance_access_url
from
  aws_connect_instance
where
  inbound_calls_enabled = 1
  and outbound_calls_enabled = 1;
```

### Get instance configuration summary
Retrieve a summary of Connect instance configurations including status, creation time, and service role.

```sql+postgres
select
  instance_alias,
  id,
  instance_status,
  created_time,
  service_role,
  identity_management_type
from
  aws_connect_instance
order by
  created_time desc;
```

```sql+sqlite
select
  instance_alias,
  id,
  instance_status,
  created_time,
  service_role,
  identity_management_type
from
  aws_connect_instance
order by
  created_time desc;
```

### Find instances by identity management type
Identify Connect instances by their identity management type to understand authentication configurations.

```sql+postgres
select
  instance_alias,
  id,
  identity_management_type,
  instance_status
from
  aws_connect_instance
where
  identity_management_type = 'SAML';
```

```sql+sqlite
select
  instance_alias,
  id,
  identity_management_type,
  instance_status
from
  aws_connect_instance
where
  identity_management_type = 'SAML';
```

### Check instance status distribution
Analyze the distribution of Connect instance statuses across your AWS environment.

```sql+postgres
select
  instance_status,
  count(*) as instance_count
from
  aws_connect_instance
group by
  instance_status
order by
  instance_count desc;
```

```sql+sqlite
select
  instance_status,
  count(*) as instance_count
from
  aws_connect_instance
group by
  instance_status
order by
  instance_count desc;
```

### Find instances with status issues
Identify Connect instances that may have status issues or failed creation attempts for troubleshooting.

```sql+postgres
select
  instance_alias,
  id,
  instance_status,
  status_reason
from
  aws_connect_instance
where
  status_reason is not null;
```

```sql+sqlite
select
  instance_alias,
  id,
  instance_status,
  status_reason
from
  aws_connect_instance
where
  status_reason is not null;
```

### Get attribute information for all instances
Combine instance data with attribute information to get a comprehensive view of Connect instances and their configurations.

```sql+postgres
select
  i.instance_alias,
  i.id,
  i.instance_status,
  a.attribute_type,
  a.value
from
  aws_connect_instance as i
  right join aws_connect_instance_attribute as a on a.instance_id = i.id;
```

```sql+sqlite
select
  i.instance_alias,
  i.id,
  i.instance_status,
  a.attribute_type,
  a.value
from
  aws_connect_instance as i
  right aws_connect_instance_attribute as a on i.id = a.instance_id;
```
