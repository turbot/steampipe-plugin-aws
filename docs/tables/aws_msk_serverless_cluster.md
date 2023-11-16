---
title: "Table: aws_msk_serverless_cluster - Query AWS Managed Streaming for Kafka (MSK) Serverless Clusters using SQL"
description: "Allows users to query AWS MSK Serverless Clusters to retrieve detailed information about each cluster."
---

# Table: aws_msk_serverless_cluster - Query AWS Managed Streaming for Kafka (MSK) Serverless Clusters using SQL

The `aws_msk_serverless_cluster` table in Steampipe provides information about serverless clusters within AWS Managed Streaming for Kafka (MSK). This table allows DevOps engineers to query cluster-specific details, including the cluster ARN, creation time, and associated metadata. Users can utilize this table to gather insights on clusters, such as their current state, the number of brokers, and more. The schema outlines the various attributes of the MSK serverless cluster, including the cluster name, tags, and the version of Apache Kafka.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_msk_serverless_cluster` table, you can use the `.inspect aws_msk_serverless_cluster` command in Steampipe.

Key columns:

- `arn`: The Amazon Resource Name (ARN) of the cluster. This unique identifier can be used to join this table with other AWS resources.
- `cluster_name`: The name of the cluster. This can be used to filter or join data based on the cluster name.
- `state`: The current state of the cluster. This can be useful for identifying clusters that are not in a desired state.

## Examples

### Basic Info

```sql
select
  arn,
  cluster_name,
  state,
  cluster_type,
  creation_time,
  current_version,
  region,
  tags
from
  aws_msk_serverless_cluster;
```

### List inactive clusters

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_serverless_cluster
where
  state <> 'ACTIVE';
```

### List clusters created within the last 90 days

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_serverless_cluster
where
  creation_time >= (current_date - interval '90' day)
order by
  creation_time;
```

### Get VPC details of each cluster

```sql
select
  arn,
  cluster_name,
  state,
  vpc ->> 'SubnetIds' as subnet_ids,
  vpc ->> 'SecurityGroupIds' as security_group_ids
from
  aws_msk_serverless_cluster,
  jsonb_array_elements(serverless -> 'VpcConfigs') as vpc
```

### List clusters with IAM authentication disabled

```sql
select
  arn,
  cluster_name,
  state,
  serverless -> 'ClientAuthentication' as client_authentication
from
  aws_msk_serverless_cluster
where
  (serverless -> 'ClientAuthentication' -> 'Sasl' -> 'Iam' ->> 'Enabled')::boolean = false
```
