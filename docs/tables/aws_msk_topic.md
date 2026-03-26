---
title: "Steampipe Table: aws_msk_topic - Query AWS MSK Topics using SQL"
description: "Allows users to query Kafka topics on AWS Managed Streaming for Apache Kafka (MSK) clusters."
folder: "MSK"
---

# Table: aws_msk_topic - Query AWS MSK Topics using SQL

An AWS MSK topic is a category or feed name to which records are published in an Amazon Managed Streaming for Apache Kafka cluster. Topics are partitioned and replicated across the brokers in the cluster.

## Table Usage Guide

The `aws_msk_topic` table in Steampipe provides you with information about Kafka topics within AWS MSK clusters. This table allows you to query topic-specific details, including partition count, replication factor, and synchronization status. You can utilize this table to monitor topic health, identify under-replicated topics, and review topic configurations. Note that `cluster_arn` is a required key column for all queries.

## Examples

### Basic info
Explore the topics configured on a specific MSK cluster.

```sql+postgres
select
  topic_name,
  topic_arn,
  partition_count,
  replication_factor,
  cluster_arn
from
  aws_msk_topic
where
  cluster_arn = 'arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abcd1234';
```

```sql+sqlite
select
  topic_name,
  topic_arn,
  partition_count,
  replication_factor,
  cluster_arn
from
  aws_msk_topic
where
  cluster_arn = 'arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abcd1234';
```

### List topics with out-of-sync replicas
Identify topics that have out-of-sync replicas, which may indicate potential data durability issues.

```sql+postgres
select
  topic_name,
  partition_count,
  replication_factor,
  out_of_sync_replica_count,
  cluster_arn
from
  aws_msk_topic
where
  cluster_arn = 'arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abcd1234'
  and out_of_sync_replica_count > 0;
```

```sql+sqlite
select
  topic_name,
  partition_count,
  replication_factor,
  out_of_sync_replica_count,
  cluster_arn
from
  aws_msk_topic
where
  cluster_arn = 'arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abcd1234'
  and out_of_sync_replica_count > 0;
```

### Get topic details including status
Retrieve detailed information about a specific topic including its status and configuration.

```sql+postgres
select
  topic_name,
  topic_arn,
  partition_count,
  replication_factor,
  status,
  configs,
  cluster_arn
from
  aws_msk_topic
where
  cluster_arn = 'arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abcd1234'
  and topic_name = 'my-topic';
```

```sql+sqlite
select
  topic_name,
  topic_arn,
  partition_count,
  replication_factor,
  status,
  configs,
  cluster_arn
from
  aws_msk_topic
where
  cluster_arn = 'arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abcd1234'
  and topic_name = 'my-topic';
```

### List topics with low replication factor
Find topics with a replication factor less than 3, which may not meet high-availability requirements.

```sql+postgres
select
  topic_name,
  partition_count,
  replication_factor,
  cluster_arn
from
  aws_msk_topic
where
  cluster_arn = 'arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abcd1234'
  and replication_factor < 3;
```

```sql+sqlite
select
  topic_name,
  partition_count,
  replication_factor,
  cluster_arn
from
  aws_msk_topic
where
  cluster_arn = 'arn:aws:kafka:us-east-1:123456789012:cluster/my-cluster/abcd1234'
  and replication_factor < 3;
```
