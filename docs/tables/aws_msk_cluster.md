---
title: "Table: aws_msk_cluster - Query AWS MSK Clusters using SQL"
description: "Allows users to query AWS Managed Streaming for Apache Kafka (MSK) clusters."
---

# Table: aws_msk_cluster - Query AWS MSK Clusters using SQL

The `aws_msk_cluster` table in Steampipe provides information about Managed Streaming for Apache Kafka (MSK) clusters within AWS. This table allows DevOps engineers to query cluster-specific details, including the cluster ARN, creation time, and associated metadata. Users can utilize this table to gather insights on clusters, such as the number of broker nodes, the version of Apache Kafka, the state of the cluster, and more. The schema outlines the various attributes of the MSK cluster, including the broker node group info, encryption info, open monitoring status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_msk_cluster` table, you can use the `.inspect aws_msk_cluster` command in Steampipe.

**Key columns**:

- `arn`: The Amazon Resource Number (ARN) of the MSK cluster. It can be used to join with other tables that contain ARN references.
- `cluster_name`: The name of the MSK cluster. This can be used to join with other tables that reference the cluster by name.
- `state`: The current state of the MSK cluster. This is useful for filtering clusters based on their current state.

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
  aws_msk_cluster;
```

### List inactive clusters

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  state <> 'ACTIVE';
```

### List clusters that allow public access

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  provisioned -> 'BrokerNodeGroupInfo' -> 'ConnectivityInfo' -> 'PublicAccess' ->> 'Type' <> 'DISABLED';
```

### List clusters with encryption at rest disabled

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  provisioned -> 'EncryptionInfo' -> 'EncryptionAtRest' is null;
```

### List clusters with encryption in transit disabled

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  provisioned -> 'EncryptionInfo' -> 'EncryptionInTransit' is null;
```

### List clusters with logging disabled

```sql
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  provisioned -> 'LoggingInfo' is null;
```

### Get total storage used by all the clusters

```sql
select
  sum((provisioned -> 'BrokerNodeGroupInfo' -> 'StorageInfo' -> 'EbsStorageInfo' ->> 'VolumeSize')::int) as total_storage
from
  aws_msk_cluster;
```
