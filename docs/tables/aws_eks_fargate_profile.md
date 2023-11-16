---
title: "Table: aws_eks_fargate_profile - Query AWS EKS Fargate Profiles using SQL"
description: "Allows users to query AWS EKS Fargate Profiles and retrieve data such as the Fargate profile name, ARN, status, and more."
---

# Table: aws_eks_fargate_profile - Query AWS EKS Fargate Profiles using SQL

The `aws_eks_fargate_profile` table in Steampipe provides information about Fargate Profiles within Amazon Elastic Kubernetes Service (EKS). This table allows DevOps engineers to query profile-specific details, including the profile name, ARN, status, and the EKS cluster to which it belongs. Users can utilize this table to gather insights on Fargate profiles, such as profiles associated with a specific EKS cluster, the status of the profiles, and more. The [schema](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eks_fargate_profile) outlines the various attributes of the EKS Fargate profile, including the profile name, ARN, status, EKS cluster, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_eks_fargate_profile` table, you can use the `.inspect aws_eks_fargate_profile` command in Steampipe.

### Key columns:

* `name`: The name of the Fargate profile. This is a unique identifier and can be used to join with other tables that reference the Fargate profile by name.
* `arn`: The Amazon Resource Number (ARN) of the Fargate profile. This unique identifier is useful for joining with other tables that reference the Fargate profile by its ARN.
* `cluster_name`: The name of the EKS cluster to which the Fargate profile belongs. This column is useful for joining with other tables that hold information about EKS clusters.

## Examples

### Basic info

```sql
select
  fargate_profile_name,
  fargate_profile_arn,
  cluster_name,
  created_at,
  status,
  tags
from
  aws_eks_fargate_profile;
```

### List fargate profiles which are inactive
```sql
select
  fargate_profile_name,
  fargate_profile_arn,
  cluster_name,
  created_at,
  status
from
  aws_eks_fargate_profile
where
  status <> 'ACTIVE';
```

### Get the subnet configuration for each fargate profile

```sql
select
  f.fargate_profile_name,
  f.cluster_name,
  f.status as fargate_profile_status,
  s.availability_zone,
  s.available_ip_address_count,
  s.cidr_block,
  s.vpc_id
from
  aws_eks_fargate_profile as f,
  aws_vpc_subnet as s,
  jsonb_array_elements(f.subnets) as subnet_id
where
  s.subnet_id = subnet_id;
```

### List fargate profiles for clusters not running Kubernetes version greater than 1.19

```sql
select
  c.name as cluster_name,
  c.arn as cluster_arn,
  c.version as cluster_version,
  f.fargate_profile_name as fargate_profile_name,
  f.fargate_profile_arn as fargate_profile_arn,
  f.created_at as created_at,
  f.pod_execution_role_arn as pod_execution_role_arn,
  f.status as fargate_profile_status
from
  aws_eks_fargate_profile as f,
  aws_eks_cluster as c
where
  c.version::float > 1.19 and f.cluster_name = c.name;
```
