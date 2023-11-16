---
title: "Table: aws_ecs_cluster - Query AWS ECS Clusters using SQL"
description: "Allows users to query AWS ECS Clusters to retrieve detailed information about each cluster's configuration, status, and associated resources."
---

# Table: aws_ecs_cluster - Query AWS ECS Clusters using SQL

The `aws_ecs_cluster` table in Steampipe provides information about clusters within AWS Elastic Container Service (ECS). This table allows DevOps engineers to query cluster-specific details, including its configuration, status, and associated resources. Users can utilize this table to gather insights on clusters, such as cluster capacity providers, default capacity provider strategy, and more. The schema outlines the various attributes of the ECS cluster, including the cluster ARN, cluster name, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ecs_cluster` table, you can use the `.inspect aws_ecs_cluster` command in Steampipe.

**Key columns**:

- `cluster_name`: The name of the cluster. It can be used to join with other tables that reference the cluster by its name.
- `cluster_arn`: The Amazon Resource Name (ARN) that identifies the cluster. This can be used to join with other tables that reference the cluster by its ARN.
- `status`: The status of the cluster (e.g., "ACTIVE"). This can be used to filter or join with other tables based on the cluster's status.

## Examples

### Basic info

```sql
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

```sql
select
  cluster_arn,
  status
from
  aws_ecs_cluster
where
  status = 'FAILED';
```


### Get details of resources attached to each cluster

```sql
select
  cluster_arn,
  attachment ->> 'id' as attachment_id,
  attachment ->> 'status' as attachment_status,
  attachment ->> 'type' as attachment_type
from
  aws_ecs_cluster,
  jsonb_array_elements(attachments) as attachment;
```


### List clusters with CloudWatch Container Insights disabled

```sql
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
