# Table: aws_ssm_inventory

AWS Systems Manager Inventory provides visibility into your AWS computing environment. You can use Inventory to collect metadata from your managed nodes.

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
