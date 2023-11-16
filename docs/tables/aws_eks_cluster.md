---
title: "Table: aws_eks_cluster - Query AWS Elastic Kubernetes Service Cluster using SQL"
description: "Allows users to query AWS Elastic Kubernetes Service Cluster data, including cluster configurations, statuses, and associated metadata."
---

# Table: aws_eks_cluster - Query AWS Elastic Kubernetes Service Cluster using SQL

The `aws_eks_cluster` table in Steampipe provides information about EKS clusters within AWS Elastic Kubernetes Service (EKS). This table allows DevOps engineers to query cluster-specific details, including cluster name, status, endpoint, and associated metadata. Users can utilize this table to gather insights on clusters, such as their current status, role ARN, VPC configurations, and more. The schema outlines the various attributes of the EKS cluster, including the cluster ARN, creation date, attached security groups, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_eks_cluster` table, you can use the `.inspect aws_eks_cluster` command in Steampipe.

**Key columns**:

- `name`: The name of the EKS cluster. This column can be used to join this table with other tables that reference the EKS cluster by name.
- `arn`: The Amazon Resource Number (ARN) of the EKS cluster. This column is useful for joining with other tables that reference the EKS cluster by ARN.
- `status`: The status of the EKS cluster. This column is important for understanding the current state of the cluster and can be used to filter or join with other tables based on cluster status.

## Examples

### Basic info

```sql
select
  name,
  arn,
  endpoint,
  identity,
  status
from
  aws_eks_cluster;
```


### Get the VPC configuration for each cluster

```sql
select
  name,
  resources_vpc_config ->> 'ClusterSecurityGroupId' as cluster_security_group_id,
  resources_vpc_config ->> 'EndpointPrivateAccess' as endpoint_private_access,
  resources_vpc_config ->> 'EndpointPublicAccess' as endpoint_public_access,
  resources_vpc_config ->> 'PublicAccessCidrs' as public_access_cidrs,
  resources_vpc_config ->> 'SecurityGroupIds' as security_group_ids,
  resources_vpc_config -> 'SubnetIds' as subnet_ids,
  resources_vpc_config ->> 'VpcId' as vpc_id
from
  aws_eks_cluster;
```


### List disabled log types for each cluster

```sql
select
  name,
  i ->> 'Enabled' as enabled,
  i ->> 'Types' as types
from
  aws_eks_cluster,
  jsonb_array_elements(logging -> 'ClusterLogging') as i
where
  i ->> 'Enabled' = 'false';
```


### List clusters not running Kubernetes version 1.19

```sql
select
  name,
  arn,
  version
from
  aws_eks_cluster
where
  version <> '1.19';
```
