---
title: "Steampipe Table: aws_eks_cluster - Query AWS Elastic Kubernetes Service Cluster using SQL"
description: "Allows users to query AWS Elastic Kubernetes Service Cluster data, including cluster configurations, statuses, and associated metadata."
folder: "EKS"
---

# Table: aws_eks_cluster - Query AWS Elastic Kubernetes Service Cluster using SQL

The AWS Elastic Kubernetes Service (EKS) Cluster is a managed service that simplifies the deployment, management, and scaling of containerized applications using Kubernetes, an open-source system. EKS runs Kubernetes control plane instances across multiple AWS availability zones to ensure high availability, automatically detects and replaces unhealthy control plane instances, and provides on-demand, zero downtime upgrades and patching. It integrates with AWS services to provide scalability and security for your applications, including Elastic Load Balancing for load distribution, IAM for authentication, and Amazon VPC for isolation.

## Table Usage Guide

The `aws_eks_cluster` table in Steampipe provides you with information about EKS clusters within AWS Elastic Kubernetes Service (EKS). This table enables you, as a DevOps engineer, to query cluster-specific details, including cluster name, status, endpoint, and associated metadata. You can utilize this table to gather insights on clusters, such as their current status, role ARN, VPC configurations, and more. The schema outlines the various attributes of the EKS cluster, including the cluster ARN, creation date, attached security groups, and associated tags for you.

## Examples

### Basic info
Determine the status and identity of your Amazon EKS clusters to assess their operational condition and identify any potential issues. This can help maintain optimal performance and security within your AWS environment.

```sql+postgres
select
  name,
  arn,
  endpoint,
  identity,
  status
from
  aws_eks_cluster;
```

```sql+sqlite
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
This query helps to assess the configuration of each cluster's Virtual Private Cloud (VPC) in an AWS EKS setup. It can be used to gain insights into the cluster's security group ID, endpoint access details, CIDR blocks for public access, associated security group IDs, subnet IDs, and the VPC ID, which can be crucial for managing network accessibility and security.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(resources_vpc_config, '$.ClusterSecurityGroupId') as cluster_security_group_id,
  json_extract(resources_vpc_config, '$.EndpointPrivateAccess') as endpoint_private_access,
  json_extract(resources_vpc_config, '$.EndpointPublicAccess') as endpoint_public_access,
  json_extract(resources_vpc_config, '$.PublicAccessCidrs') as public_access_cidrs,
  json_extract(resources_vpc_config, '$.SecurityGroupIds') as security_group_ids,
  json_extract(resources_vpc_config, '$.SubnetIds') as subnet_ids,
  json_extract(resources_vpc_config, '$.VpcId') as vpc_id
from
  aws_eks_cluster;
```


### List disabled log types for each cluster
Determine the areas in which log types are disabled for each cluster in AWS EKS service. This is useful for identifying potential gaps in your logging strategy, ensuring comprehensive coverage for effective monitoring and debugging.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(i.value, '$.Enabled') as enabled,
  json_extract(i.value, '$.Types') as types
from
  aws_eks_cluster,
  json_each(logging, 'ClusterLogging') as i
where
  json_extract(i.value, '$.Enabled') = 'false';
```


### List clusters not running Kubernetes version 1.19
Identify those clusters within your AWS EKS environment that are not operating on Kubernetes version 1.19. This can be useful to ensure compliance with specific version requirements or to plan for necessary upgrades.

```sql+postgres
select
  name,
  arn,
  version
from
  aws_eks_cluster
where
  version <> '1.19';
```

```sql+sqlite
select
  name,
  arn,
  version
from
  aws_eks_cluster
where
  version != '1.19';
```