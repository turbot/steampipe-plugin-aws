---
title: "Table: aws_eks_addon_version - Query AWS EKS Add-On Versions using SQL"
description: "Allows users to query AWS EKS Add-On Versions."
---

# Table: aws_eks_addon_version - Query AWS EKS Add-On Versions using SQL

The `aws_eks_addon_version` table in Steampipe provides information about Add-On versions within Amazon Elastic Kubernetes Service (EKS). This table allows DevOps engineers to query add-on specific details, including addon name, addon version, architecture, and associated metadata. Users can utilize this table to gather insights on add-ons, such as the add-on version status, the specific architectures it supports, and more. The schema outlines the various attributes of the EKS add-on, including the add-on name, add-on version, and supported architectures.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_eks_addon_version` table, you can use the `.inspect aws_eks_addon_version` command in Steampipe.

### Key columns:

- `addon_name`: The name of the add-on. This can be used to join this table with other tables that contain information about EKS add-ons.
- `addon_version`: The version of the add-on. This is important when comparing different versions of the same add-on.
- `architecture`: The architectures that the add-on version supports. This is useful for ensuring compatibility with your EKS cluster's architecture.

## Examples

### Basic info

```sql
select
  addon_name,
  addon_version,
  type
from
  aws_eks_addon_version;
```

### Count the number of add-on versions by add-on

```sql
select
  addon_name,
  count(addon_version) as addon_version_count
from
  aws_eks_addon_version
group by
  addon_name;
```

### Get configuration details of each add-on version

```sql
select
  addon_name,
  addon_version,
  addon_configuration -> '$defs' -> 'extraVolumeTags' ->> 'description' as addon_configuration_def_description,
  addon_configuration -> '$defs' -> 'extraVolumeTags' -> 'propertyNames' as addon_configuration_def_property_names,
  addon_configuration -> '$defs' -> 'extraVolumeTags' -> 'patternProperties' as addon_configuration_def_pattern_properties,
  addon_configuration -> 'properties' as addon_configuration_properties
from
  aws_eks_addon_version limit 10;
```