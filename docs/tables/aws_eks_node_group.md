---
title: "Steampipe Table: aws_eks_node_group - Query AWS EKS Node Group using SQL"
description: "Allows users to query AWS EKS Node Group data, providing information about each node group within an AWS Elastic Kubernetes Service (EKS) cluster."
folder: "EKS"
---

# Table: aws_eks_node_group - Query AWS EKS Node Group using SQL

The AWS EKS Node Group is a resource within Amazon Elastic Kubernetes Service (EKS). It represents a group of nodes within a cluster that all share the same configuration, making it easier to manage and scale your applications. Node groups are associated with a specific Amazon EKS cluster and can be customized according to your workload requirements.

## Table Usage Guide

The `aws_eks_node_group` table in Steampipe provides you with information about each node group within an AWS Elastic Kubernetes Service (EKS) cluster. This table allows you, as a DevOps engineer, system administrator, or other technical professional, to query node-group-specific details, including the node group ARN, creation timestamp, health status, and associated metadata. You can utilize this table to gather insights on node groups, such as the status of each node, the instance types used, and more. The schema outlines the various attributes of the EKS node group for you, including the node role, subnets, scaling configuration, and associated tags.

## Examples

### Basic info
Explore the status and creation details of node groups within your Amazon EKS clusters. This allows you to track the health and longevity of your Kubernetes resources, aiding in efficient resource management.

```sql+postgres
select
  nodegroup_name,
  arn,
  created_at,
  cluster_name,
  status
from
  aws_eks_node_group;
```

```sql+sqlite
select
  nodegroup_name,
  arn,
  created_at,
  cluster_name,
  status
from
  aws_eks_node_group;
```

### List node groups that are not active
Identify instances where certain node groups within your AWS EKS service are not active. This can help in managing resources effectively by pinpointing potential areas of concern or underutilization.

```sql+postgres
select
  nodegroup_name,
  arn,
  created_at,
  cluster_name,
  status
from
  aws_eks_node_group
where
  status <> 'ACTIVE';
```

```sql+sqlite
select
  nodegroup_name,
  arn,
  created_at,
  cluster_name,
  status
from
  aws_eks_node_group
where
  status != 'ACTIVE';
```

### Get health status of the node groups
Assess the health status of various node groups within your AWS EKS clusters. This can help identify any potential issues or anomalies, ensuring optimal performance and stability of your Kubernetes workloads.

```sql+postgres
select
  nodegroup_name,
  cluster_name,
  jsonb_pretty(health) as health
from
  aws_eks_node_group;
```

```sql+sqlite
select
  nodegroup_name,
  cluster_name,
  health
from
  aws_eks_node_group;
```

### Get launch template details of the node groups
Determine the configuration details of node groups within a cluster to understand the settings and specifications of each node group.

```sql+postgres
select
  nodegroup_name,
  cluster_name,
  jsonb_pretty(launch_template) as launch_template
from
  aws_eks_node_group;
```

```sql+sqlite
select
  nodegroup_name,
  cluster_name,
  launch_template
from
  aws_eks_node_group;
```