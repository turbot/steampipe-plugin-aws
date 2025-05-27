---
title: "Steampipe Table: aws_ssm_inventory - Query AWS Systems Manager Inventory using SQL"
description: "Allows users to query AWS Systems Manager Inventory, providing information about managed instances in AWS Systems Manager."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_inventory - Query AWS Systems Manager Inventory using SQL

The AWS Systems Manager Inventory provides visibility into your Amazon EC2 and on-premises compute infrastructure. It collects metadata from your managed instances about applications, files, Windows updates, network configurations, and other details. This collected data assists in managing your systems, tracking software inventory, and applying patches.

## Table Usage Guide

The `aws_ssm_inventory` table in Steampipe provides you with information about managed instances within AWS Systems Manager. This table enables you, as a DevOps engineer, to query instance-specific details, including instance name, type, platform type, and associated metadata. You can utilize this table to gather insights on instances, such as their status, their associated tags, and more. The schema outlines for you the various attributes of the managed instance, including the instance ID, instance type, platform type, and associated tags.

## Examples

### Basic info
Explore which AWS Simple Systems Manager (SSM) inventory items have been captured at a specific time, allowing you to understand the historical state of your resources and their schema versions across different regions. This information can aid in resource management and tracking changes in your AWS environment.

```sql+postgres
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

```sql+sqlite
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
Explore the essential characteristics of a particular managed instance, such as its platform type, agent version, and status. This information can be useful for understanding the instance's current configuration and performance, as well as for troubleshooting potential issues.

```sql+postgres
select
  si.id,
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
  aws_ssm_inventory as si,
  jsonb_array_elements(content) as c
where
  id = 'i-0665a65b1a1c2b47g';
```

```sql+sqlite
select
  si.id,
  json_extract(c, '$.AgentType') as agent_type,
  json_extract(c, '$.IpAddress') as ip_address,
  json_extract(c, '$.AgentVersion') as agent_version,
  json_extract(c, '$.ComputerName') as computer_name,
  json_extract(c, '$.PlatformName') as platform_name,
  json_extract(c, '$.PlatformType') as platform_type,
  json_extract(c, '$.ResourceType') as resource_type,
  json_extract(c, '$.InstanceStatus') as instance_status,
  json_extract(c, '$.PlatformVersion') as platform_version
from
  aws_ssm_inventory as si,
  json_each(content) as c
where
  id = 'i-0665a65b1a1c2b47g';
```

### List schema definitions of inventories
This query helps you gain insights into the structure and organization of your AWS Systems Manager (SSM) inventories. It's useful for understanding the types of data stored in each inventory and how they are presented, which can aid in managing and utilizing your SSM resources effectively.

```sql+postgres
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

```sql+sqlite
select
  si.id,
  json_extract(s.value, '$.Version') as schema_version,
  json_extract(s.value, '$.TypeName') as type_name,
  json_extract(s.value, '$.DisplayName') as display_name,
  json_extract(s.value, '$.Attributes') as attributes
from
  aws_ssm_inventory as si,
  json_each(schema) as s
order by 
  si.id;
```

### Get inventory details from the last 10 days
Explore recent changes in your AWS inventory by identifying items that have been added or modified in the last 10 days. This is useful for keeping track of inventory updates and ensuring system integrity.

```sql+postgres
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

```sql+sqlite
select
  id,
  type_name,
  capture_time,
  schema_version,
  content
from
  aws_ssm_inventory
where
  capture_time >= datetime('now', '-10 day');
```

### Get inventory content of all running instances
Explore the specific attributes of all operational instances, including details such as their agent type, IP address, platform, and status. This is useful for effectively managing and monitoring your active instances in a cloud environment.

```sql+postgres
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

```sql+sqlite
select
  v.id,
  i.instance_state,
  i.instance_type,
  json_extract(c.value, '$.AgentType') as agent_type,
  json_extract(c.value, '$.IpAddress') as ip_address,
  json_extract(c.value, '$.AgentVersion') as agent_version,
  json_extract(c.value, '$.ComputerName') as computer_name,
  json_extract(c.value, '$.PlatformName') as platform_name,
  json_extract(c.value, '$.PlatformType') as platform_type,
  json_extract(c.value, '$.ResourceType') as resource_type,
  json_extract(c.value, '$.InstanceStatus') as instance_status,
  json_extract(c.value, '$.PlatformVersion') as platform_version
from
  aws_ssm_inventory as v,
  aws_ec2_instance as i,
  json_each(v.content) as c
where
  v.id = i.instance_id
and
  i.instance_state = 'running';
```