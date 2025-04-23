---
title: "Steampipe Table: aws_ssm_inventory_entry - Query AWS Systems Manager Inventory Entry using SQL"
description: "Allows users to query AWS Systems Manager Inventory Entry to fetch information about the inventory entries of managed instances. The table provides details such as instance ID, type name, schema version, capture time, and inventory data."
folder: "Systems Manager (SSM)"
---

# Table: aws_ssm_inventory_entry - Query AWS Systems Manager Inventory Entry using SQL

The AWS Systems Manager Inventory provides visibility into your Amazon EC2 and on-premises computing environment. It collects metadata from your managed instances, such as installed applications, system configuration, network configurations, and Windows updates. This allows you to manage your systems at scale, and helps you to quickly diagnose and troubleshoot operational issues.

## Table Usage Guide

The `aws_ssm_inventory_entry` table in Steampipe provides you with information about the inventory entries of managed instances within AWS Systems Manager. This table allows you, as a DevOps engineer, to query inventory-specific details, including the instance ID, type name, schema version, capture time, and the actual inventory data. You can utilize this table to gather insights on inventory entries, such as the software installed on instances, network configurations, Windows updates status, and more. The schema outlines the various attributes of the inventory entry for you, including the instance ID, type name, schema version, capture time, and inventory data.

## Examples

### Basic info
Explore which AWS Simple Systems Manager (SSM) instances are being inventoried, by determining the type, capture time, and schema version of each entry. This can be particularly useful for managing and auditing your AWS resources effectively.

```sql+postgres
select
  instance_id,
  type_name,
  capture_time,
  schema_version,
  entries
from
  aws_ssm_inventory_entry;
```

```sql+sqlite
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
Explore recent changes in your system by identifying inventory entries made within the last 30 days. This is particularly useful for tracking system modifications and maintaining up-to-date records.

```sql+postgres
select
  instance_id,
  type_name,
  capture_time,
  schema_version,
  entries
from
  aws_ssm_inventory_entry
where
  capture_time >= now() - interval '30 day';
```

```sql+sqlite
select
  instance_id,
  type_name,
  capture_time,
  schema_version,
  entries
from
  aws_ssm_inventory_entry
where
  capture_time >= datetime('now', '-30 day');
```

### Get inventory details of entries
This example helps to identify the type and schema version of specific inventory entries within the AWS SSM service. It can be useful for auditing purposes or to ensure compliance with specific schema versions.

```sql+postgres
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

```sql+sqlite
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
Explore the details of managed instances within your inventory to understand their association status, resource type, and whether they are running the latest version. This can help in maintaining an up-to-date and efficient inventory system.

```sql+postgres
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

```sql+sqlite
select
  e.instance_id,
  e.type_name,
  i.resource_type,
  i.association_status,
  i.computer_name,
  i.ip_address,
  i.is_latest_version
from
  aws_ssm_inventory_entry as e
join
  aws_ssm_managed_instance as i
on
  i.instance_id = e.instance_id;
```

### List custom inventory entries of an instance
Determine the areas in which custom inventory entries of a specific instance are listed. This is useful to understand and analyze the custom configurations made to an instance for better management and optimization of resources.

```sql+postgres
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

```sql+sqlite
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