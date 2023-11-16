---
title: "Table: aws_ssm_inventory - Query AWS Systems Manager Inventory using SQL"
description: "Allows users to query AWS Systems Manager Inventory, providing information about managed instances in AWS Systems Manager."
---

# Table: aws_ssm_inventory - Query AWS Systems Manager Inventory using SQL

The `aws_ssm_inventory` table in Steampipe provides information about managed instances within AWS Systems Manager. This table allows DevOps engineers to query instance-specific details, including instance name, type, platform type, and associated metadata. Users can utilize this table to gather insights on instances, such as instances' status, their associated tags, and more. The schema outlines the various attributes of the managed instance, including the instance ID, instance type, platform type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_inventory` table, you can use the `.inspect aws_ssm_inventory` command in Steampipe.

**Key columns**:
- `instance_id`: This is the unique identifier for the managed instance. It can be used to join with other tables that contain instance-specific information.
- `instance_type`: This column provides information about the type of the managed instance, which can be used to join with tables that contain instance type-specific information.
- `platform_type`: This column provides information about the platform of the managed instance, which can be used to join with tables that contain platform-specific information.

## Examples

### Basic info

```sql
select
  id,
  type_name,
  capture_time,
  schema_version,
  content,
  region
from
  aws_ssm_inventory;
```

### Get content details of a managed instance 

```sql
select
  id,
  c ->> 'AgentType' as agent_type,
  c ->> 'IpAddress' as ip_address,
  c ->> 'AgentVersion' as agent_version,
  c ->> 'ComputerName' as computer_name,
  c ->> 'PlatformName' as platform_name,
  c ->> 'PlatformType' as platform_type,
  c ->> 'ResourceType' as resource_type,
  c ->> 'InstanceStatus' as instance_status,
  c ->> 'PlatformVersion' as platform_version
from
  aws_ssm_inventory,
  jsonb_array_elements(content) as c
where
  id = 'i-0665a65b1a1c2b47g';
```

### List schema definitions of inventories

```sql
select
  id,
  s ->> 'Version' as schema_version,
  s ->> 'TypeName' as type_name,
  s ->> 'DisplayName' as display_name,
  jsonb_pretty(s -> 'Attributes') as attributes
from
  aws_ssm_inventory,
  jsonb_array_elements(schema) as s
order by 
  id;
```

### Get inventory details from the last 10 days

```sql
select
  id,
  type_name,
  capture_time,
  schema_version,
  content
from
  aws_ssm_inventory
where
  capture_time >= now() - interval '10' day;
```

### Get inventory content of all running instances

```sql
select
  v.id,
  i.instance_state,
  i.instance_type,
  c ->> 'AgentType' as agent_type,
  c ->> 'IpAddress' as ip_address,
  c ->> 'AgentVersion' as agent_version,
  c ->> 'ComputerName' as computer_name,
  c ->> 'PlatformName' as platform_name,
  c ->> 'PlatformType' as platform_type,
  c ->> 'ResourceType' as resource_type,
  c ->> 'InstanceStatus' as instance_status,
  c ->> 'PlatformVersion' as platform_version
from
  aws_ssm_inventory as v,
  aws_ec2_instance as i,
  jsonb_array_elements(content) as c
where
  v.id = i.instance_id
and
  i.instance_state = 'running';
```
