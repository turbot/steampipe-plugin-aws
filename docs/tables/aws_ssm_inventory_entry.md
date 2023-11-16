---
title: "Table: aws_ssm_inventory_entry - Query AWS Systems Manager Inventory Entry using SQL"
description: "Allows users to query AWS Systems Manager Inventory Entry to fetch information about the inventory entries of managed instances. The table provides details such as instance ID, type name, schema version, capture time, and inventory data."
---

# Table: aws_ssm_inventory_entry - Query AWS Systems Manager Inventory Entry using SQL

The `aws_ssm_inventory_entry` table in Steampipe provides information about the inventory entries of managed instances within AWS Systems Manager. This table allows DevOps engineers to query inventory-specific details, including the instance ID, type name, schema version, capture time, and the actual inventory data. Users can utilize this table to gather insights on inventory entries, such as the software installed on instances, network configurations, Windows updates status, and more. The schema outlines the various attributes of the inventory entry, including the instance ID, type name, schema version, capture time, and inventory data.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ssm_inventory_entry` table, you can use the `.inspect aws_ssm_inventory_entry` command in Steampipe.

**Key columns**:

- `instance_id`: This is the ID of the managed instance. It is useful for linking inventory entries to their respective instances.
- `type_name`: This represents the type of inventory item. It helps in categorizing and filtering inventory entries.
- `capture_time`: This is the timestamp when the inventory information was collected. It is useful for tracking the inventory status over time.

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

### List custom inventory entries of an instance

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
  instance_id = 'i-1234567890abcwd4f'
and
  type_name like 'Custom%';
```
