---
title: "Steampipe Table: aws_eks_fargate_profile - Query AWS EKS Fargate Profiles using SQL"
description: "Allows users to query AWS EKS Fargate Profiles and retrieve data such as the Fargate profile name, ARN, status, and more."
folder: "EKS"
---

# Table: aws_eks_fargate_profile - Query AWS EKS Fargate Profiles using SQL

The AWS EKS Fargate Profile is a component of Amazon Elastic Kubernetes Service (EKS) that allows you to run Kubernetes pods on AWS Fargate. With Fargate, you can focus on designing and building your applications instead of managing the infrastructure that runs them. It eliminates the need for you to choose server types, decide when to scale your node groups, or optimize cluster packing.

## Table Usage Guide

The `aws_eks_fargate_profile` table in Steampipe provides you with information about Fargate Profiles within Amazon Elastic Kubernetes Service (EKS). This table allows you as a DevOps engineer to query profile-specific details, including the profile name, ARN, status, and the EKS cluster to which it belongs. You can utilize this table to gather insights on Fargate profiles, such as profiles associated with a specific EKS cluster, the status of the profiles, and more. The [schema](https://hub.steampipe.io/plugins/turbot/aws/tables/aws_eks_fargate_profile) outlines the various attributes of the EKS Fargate profile for you, including the profile name, ARN, status, EKS cluster, and associated tags.

## Examples

### Basic info
Determine the areas in which AWS EKS Fargate profiles are being utilized. This query can help you assess the status and creation date of these profiles, offering insights for resource management and optimization.

```sql+postgres
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

```sql+sqlite
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
Identify instances where AWS EKS Fargate profiles are not currently active. This can be useful for troubleshooting, maintenance, or resource optimization purposes.

```sql+postgres
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

```sql+sqlite
select
  fargate_profile_name,
  fargate_profile_arn,
  cluster_name,
  created_at,
  status
from
  aws_eks_fargate_profile
where
  status != 'ACTIVE';
```

### Get the subnet configuration for each fargate profile
Explore the configurations of various Fargate profiles within your EKS clusters to understand the availability and IP address count for each associated subnet. This can be beneficial to manage resources efficiently and ensure optimal performance of your applications.

```sql+postgres
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

```sql+sqlite
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
  aws_vpc_subnet as s
where
  s.subnet_id IN (select value from json_each(f.subnets));
```

### List fargate profiles for clusters not running Kubernetes version greater than 1.19
Explore which Fargate profiles are associated with clusters not running Kubernetes version greater than 1.19. This can be beneficial in identifying outdated clusters, facilitating necessary upgrades to improve system performance and security.

```sql+postgres
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

```sql+sqlite
Error: The corresponding SQLite query is unavailable.
```