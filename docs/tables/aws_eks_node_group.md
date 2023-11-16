---
title: "Table: aws_eks_node_group - Query AWS EKS Node Group using SQL"
description: "Allows users to query AWS EKS Node Group data, providing information about each node group within an AWS Elastic Kubernetes Service (EKS) cluster."
---

# Table: aws_eks_node_group - Query AWS EKS Node Group using SQL

The `aws_eks_node_group` table in Steampipe provides information about each node group within an AWS Elastic Kubernetes Service (EKS) cluster. This table allows DevOps engineers, system administrators, and other technical professionals to query node-group-specific details, including the node group ARN, creation timestamp, health status, and associated metadata. Users can utilize this table to gather insights on node groups, such as the status of each node, the instance types used, and more. The schema outlines the various attributes of the EKS node group, including the node role, subnets, scaling configuration, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_eks_node_group` table, you can use the `.inspect aws_eks_node_group` command in Steampipe.

**Key columns**:

- `nodegroup_name`: The name of the node group. This can be used to join with other tables that contain node group-specific information.
- `cluster_name`: The name of the cluster that the node group belongs to. This can be used to join with other tables that contain cluster-specific information.
- `nodegroup_arn`: The Amazon Resource Name (ARN) of the node group. This can be used to join with other tables that contain ARN-specific information, such as access policies.

## Examples

### Basic info

```sql
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

```sql
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

### Get health status of the node groups

```sql
select
  nodegroup_name,
  cluster_name,
  jsonb_pretty(health) as health
from
  aws_eks_node_group;
```

### Get launch template details of the node groups

```sql
select
  nodegroup_name,
  cluster_name,
  jsonb_pretty(launch_template) as launch_template
from
  aws_eks_node_group;
```
