---
title: "Steampipe Table: aws_eks_addon - Query AWS EKS Add-Ons using SQL"
description: "Allows users to query AWS EKS Add-Ons to retrieve information about add-ons associated with each Amazon EKS cluster."
folder: "EKS"
---

# Table: aws_eks_addon - Query AWS EKS Add-Ons using SQL

The AWS EKS Add-Ons are additional software components that enhance the functionality of your Amazon Elastic Kubernetes Service (EKS) clusters. They provide a way to deploy and manage Kubernetes applications, improve cluster security, and simplify cluster management. Using AWS EKS Add-Ons, you can automate time-consuming tasks such as patching, updating, and scaling.

## Table Usage Guide

The `aws_eks_addon` table in Steampipe provides you with information about add-ons associated with each Amazon EKS cluster. This table allows you, as a DevOps engineer, to query add-on-specific details, including add-on versions, status, and associated metadata. You can utilize this table to gather insights on add-ons, such as the current version of each add-on, the health of add-ons, and more. The schema outlines the various attributes of the EKS add-on for you, including the add-on name, add-on version, service account role ARN, and associated tags.

## Examples

### Basic info
Explore the status of various add-ons within your AWS EKS clusters to understand their versions and associated roles. This can be beneficial for assessing the current configuration and ensuring the optimal functionality of your clusters.

```sql+postgres
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

```sql+sqlite
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
Identify instances where certain add-ons in the AWS EKS service are not active. This can help in monitoring and managing resources effectively by pinpointing inactive add-ons that may need attention or removal.

```sql+postgres
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

```sql+sqlite
select
  addon_name,
  arn,
  cluster_name,
  status
from
  aws_eks_addon
where
  status != 'ACTIVE';
```

### Get count of add-ons by cluster
Determine the total number of add-ons per cluster within your AWS EKS environment to better manage resources and understand utilization.

```sql+postgres
select
  cluster_name,
  count(addon_name) as addon_count
from
  aws_eks_addon
group by
  cluster_name;
```

```sql+sqlite
select
  cluster_name,
  count(addon_name) as addon_count
from
  aws_eks_addon
group by
  cluster_name;
```