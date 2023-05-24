# Table: aws_ssm_inventory_entry

AWS SSM (Systems Manager) Inventory Entries refer to the collected metadata information about the managed instances in your AWS environment. It includes details such as installed software, configuration settings, and other attributes of the instances. The inventory entries provide a comprehensive view of the software and configuration across your fleet of managed instances, allowing you to effectively manage and track the resources in your AWS infrastructure.

## Examples

### Basic info

```sql
select
  instance_id,
  type_name,
  capture_time,
  schema_version,
  entries
from
  aws_ssm_inventory_entry;
```

### List inventory entries in the last 30 days

```sql
select
  instance_id,
  type_name,
  capture_time,
  schema_version,
  entries
from
  aws_ssm_inventory_entry
where
  capture_time >= time() - interval '30 day';
```

### Get inventory details of entries

```sql
select
  e.instance_id,
  e.type_name,
  i.schema_version,
  i.schema
from
  aws_ssm_inventory_entry as e,
  aws_ssm_inventory as i
where
  i.id = e.instance_id;
```

### Get managed instance details of inventory entries

```sql
select
  e.instance_id,
  e.type_name,
  i.resource_type,
  i.association_status,
  i.computer_name,
  i.ip_address,
  i.is_latest_version
from
  aws_ssm_inventory_entry as e,
  aws_ssm_managed_instance as i
where
  i.instance_id = e.instance_id;
```