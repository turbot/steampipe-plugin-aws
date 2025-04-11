---
title: "Steampipe Table: aws_msk_serverless_cluster - Query AWS Managed Streaming for Kafka (MSK) Serverless Clusters using SQL"
description: "Allows users to query AWS MSK Serverless Clusters to retrieve detailed information about each cluster."
folder: "MSK"
---

# Table: aws_msk_serverless_cluster - Query AWS Managed Streaming for Kafka (MSK) Serverless Clusters using SQL

The AWS Managed Streaming for Kafka (MSK) Serverless Cluster is a fully managed service that makes it easy to build and run applications that use Apache Kafka to process streaming data. It takes care of the underlying Kafka infrastructure, so you can focus on application development. With MSK, you can use native Apache Kafka APIs to populate data lakes, stream changes to and from databases, and power machine learning and analytics applications.

## Table Usage Guide

The `aws_msk_serverless_cluster` table in Steampipe provides you with information about serverless clusters within AWS Managed Streaming for Kafka (MSK). This table allows you, as a DevOps engineer, to query cluster-specific details, including the cluster ARN, creation time, and associated metadata. You can utilize this table to gather insights on clusters, such as their current state, the number of brokers, and more. The schema outlines for you the various attributes of the MSK serverless cluster, including the cluster name, tags, and the version of Apache Kafka.

## Examples

### Basic Info
Determine the areas in which AWS MSK Serverless clusters are located and their current state to assess their availability and performance. This helps in understanding the distribution and health of your clusters across different regions.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that contain inactive clusters, allowing you to understand the areas in your AWS MSK serverless cluster that are not currently active. This can be useful for resource management and optimizing your cloud infrastructure.

```sql+postgres
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

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_serverless_cluster
where
  state != 'ACTIVE';
```

### List clusters created within the last 90 days
Identify recently created clusters to monitor their performance and manage resources efficiently. This query is useful for tracking system growth and planning future capacity needs.

```sql+postgres
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

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_serverless_cluster
where
  creation_time >= date('now','-90 day')
order by
  creation_time;
```

### Get VPC details of each cluster
Analyze the settings to understand the virtual private cloud (VPC) configurations for each of your AWS serverless clusters. This can help you ensure the security and network performance of your clusters.

```sql+postgres
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

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  json_extract(vpc.value, '$.SubnetIds') as subnet_ids,
  json_extract(vpc.value, '$.SecurityGroupIds') as security_group_ids
from
  aws_msk_serverless_cluster,
  json_each(json_extract(serverless, '$.VpcConfigs')) as vpc
```

### List clusters with IAM authentication disabled
Explore which clusters have their IAM authentication disabled. This can be used to identify potential security risks and ensure that all clusters are appropriately secured.

```sql+postgres
select
  arn,
  cluster_name,
  state,
  serverless -> 'ClientAuthentication' as client_authentication
from
  aws_msk_serverless_cluster
where
  (serverless -> 'ClientAuthentication' -> 'Sasl' -> 'Iam' ->> 'Enabled')::boolean = false;
```

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  json_extract(serverless, '$.ClientAuthentication') as client_authentication
from
  aws_msk_serverless_cluster
where
  json_extract(json_extract(json_extract(json_extract(serverless, '$.ClientAuthentication'), '$.Sasl'), '$.Iam'), '$.Enabled') = 'false';
```