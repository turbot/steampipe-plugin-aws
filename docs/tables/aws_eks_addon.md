---
title: "Table: aws_eks_addon - Query AWS EKS Add-Ons using SQL"
description: "Allows users to query AWS EKS Add-Ons to retrieve information about add-ons associated with each Amazon EKS cluster."
---

# Table: aws_eks_addon - Query AWS EKS Add-Ons using SQL

The `aws_eks_addon` table in Steampipe provides information about add-ons associated with each Amazon EKS cluster. This table allows DevOps engineers to query add-on-specific details, including add-on versions, status, and associated metadata. Users can utilize this table to gather insights on add-ons, such as the current version of each add-on, the health of add-ons, and more. The schema outlines the various attributes of the EKS add-on, including the add-on name, add-on version, service account role ARN, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_eks_addon` table, you can use the `.inspect aws_eks_addon` command in Steampipe.

**Key columns**:

- `addon_name`: The name of the add-on. This column is important as it can be used to join the table with other tables that contain information about specific add-ons.
- `cluster_name`: The name of the cluster that the add-on is associated with. This column is useful for joining with other tables that contain cluster-specific information.
- `status`: The status of the add-on. This column is useful for understanding the current state of the add-on and can be used to filter for add-ons that are in certain states.

## Examples

### Basic info

```sql
select
  addon_name,
  arn,
  addon_version,
  cluster_name,
  status,
  service_account_role_arn
from
  aws_eks_addon;
```


### List add-ons that are not active

```sql
select
  addon_name,
  arn,
  cluster_name,
  status
from
  aws_eks_addon
where
  status <> 'ACTIVE';
```


### Get count of add-ons by cluster

```sql
select
  cluster_name,
  count(addon_name) as addon_count
from
  aws_eks_addon
group by
  cluster_name;
```
