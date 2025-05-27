---
title: "Steampipe Table: aws_ecs_cluster - Query AWS ECS Clusters using SQL"
description: "Allows users to query AWS ECS Clusters to retrieve detailed information about each cluster's configuration, status, and associated resources."
folder: "ECS"
---

# Table: aws_ecs_cluster - Query AWS ECS Clusters using SQL

The AWS ECS Cluster is a regional, logical grouping of services in Amazon Elastic Container Service (ECS). It allows you to manage and scale a group of tasks or services, and determine their placement across a set of Amazon EC2 instances. ECS Clusters help in running applications and services on a managed cluster of EC2 instances, eliminating the need to install, operate, and scale your own cluster management infrastructure.

## Table Usage Guide

The `aws_ecs_cluster` table in Steampipe provides you with information about clusters within AWS Elastic Container Service (ECS). This table allows you, as a DevOps engineer, to query cluster-specific details, including its configuration, status, and associated resources. You can utilize this table to gather insights on clusters, such as cluster capacity providers, default capacity provider strategy, and more. The schema outlines for you the various attributes of the ECS cluster, including the cluster ARN, cluster name, status, and associated tags.

## Examples

### Basic info
Analyze the settings to understand the overall status and active services of your AWS ECS clusters. This is useful for maintaining optimal cluster performance and identifying any potential issues.

```sql+postgres
select
  cluster_arn,
  cluster_name,
  active_services_count,
  attachments,
  attachments_status,
  status
from
  aws_ecs_cluster;
```

```sql+sqlite
select
  cluster_arn,
  cluster_name,
  active_services_count,
  attachments,
  attachments_status,
  status
from
  aws_ecs_cluster;
```


### List clusters that have failed to provision resources
Identify instances where resource provisioning has failed in certain clusters. This can be useful in troubleshooting and understanding the reasons for failure in resource allocation.

```sql+postgres
select
  cluster_arn,
  status
from
  aws_ecs_cluster
where
  status = 'FAILED';
```

```sql+sqlite
select
  cluster_arn,
  status
from
  aws_ecs_cluster
where
  status = 'FAILED';
```


### Get details of resources attached to each cluster
Explore the status and type of resources linked to each cluster in your AWS ECS setup. This helps you monitor the health and functionality of various components within your clusters.

```sql+postgres
select
  cluster_arn,
  attachment ->> 'id' as attachment_id,
  attachment ->> 'status' as attachment_status,
  attachment ->> 'type' as attachment_type
from
  aws_ecs_cluster,
  jsonb_array_elements(attachments) as attachment;
```

```sql+sqlite
select
  cluster_arn,
  json_extract(attachment.value, '$.id') as attachment_id,
  json_extract(attachment.value, '$.status') as attachment_status,
  json_extract(attachment.value, '$.type') as attachment_type
from
  aws_ecs_cluster,
  json_each(attachments) as attachment;
```


### List clusters with CloudWatch Container Insights disabled
Determine the areas in your AWS ECS clusters where CloudWatch Container Insights is disabled. This is beneficial in understanding and managing the monitoring capabilities of your clusters.

```sql+postgres
select
  cluster_arn,
  setting ->> 'Name' as name,
  setting ->> 'Value' as value
from
  aws_ecs_cluster,
  jsonb_array_elements(settings) as setting
where
  setting ->> 'Value' = 'disabled';
```

```sql+sqlite
select
  cluster_arn,
  json_extract(setting.value, '$.Name') as name,
  json_extract(setting.value, '$.Value') as value
from
  aws_ecs_cluster,
  json_each(settings) as setting
where
  json_extract(setting, '$.Value') = 'disabled';
```