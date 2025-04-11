---
title: "Steampipe Table: aws_msk_cluster - Query AWS MSK Clusters using SQL"
description: "Allows users to query AWS Managed Streaming for Apache Kafka (MSK) clusters."
folder: "MSK"
---

# Table: aws_msk_cluster - Query AWS MSK Clusters using SQL

The AWS Managed Streaming for Apache Kafka (MSK) is a fully managed service that makes it easy to build and run applications that use Apache Kafka to process streaming data. AWS MSK provides the control-plane operations, such as those for creating, updating, and deleting clusters. It also takes care of the maintenance and operations of the underlying infrastructure, so you can focus on building and running your applications.

## Table Usage Guide

The `aws_msk_cluster` table in Steampipe provides you with information about Managed Streaming for Apache Kafka (MSK) clusters within AWS. This table allows you, as a DevOps engineer, to query cluster-specific details, including the cluster ARN, creation time, and associated metadata. You can utilize this table to gather insights on clusters, such as the number of broker nodes, the version of Apache Kafka, the state of the cluster, and more. The schema outlines the various attributes of the MSK cluster for you, including the broker node group info, encryption info, open monitoring status, and associated tags.

## Examples

### Basic Info
Explore the status and details of your AWS MSK clusters to understand their configuration and operational state. This can be useful for auditing purposes, or to identify potential issues with cluster setup or version compatibility.

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
  aws_msk_cluster;
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
  aws_msk_cluster;
```

### List inactive clusters
Identify instances where certain clusters within AWS MSK service are not in an active state. This could be useful for system administrators who need to manage resources or troubleshoot issues.

```sql+postgres
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

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  state != 'ACTIVE';
```

### List clusters that allow public access
Determine the areas in which public access is allowed for clusters. This is useful for identifying potential security risks and ensuring that access to sensitive data is appropriately restricted.

```sql+postgres
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

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  json_extract(provisioned, '$.BrokerNodeGroupInfo.ConnectivityInfo.PublicAccess.Type') <> 'DISABLED';
```

### List clusters with encryption at rest disabled
Determine the areas in which encryption at rest is disabled for clusters, allowing you to address potential security vulnerabilities by identifying clusters without this added layer of data protection.

```sql+postgres
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

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  json_extract(provisioned, '$.EncryptionInfo.EncryptionAtRest') is null;
```

### List clusters with encryption in transit disabled
Determine the areas in which encryption in transit is disabled for clusters. This can be useful for identifying potential security vulnerabilities and ensuring data safety.

```sql+postgres
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

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  json_extract(provisioned, '$.EncryptionInfo.EncryptionInTransit') is null;
```

### List clusters with logging disabled
Discover the segments that consist of clusters with logging disabled, allowing you to identify potential areas for enhancing security measures and ensuring compliance with logging policies.

```sql+postgres
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

```sql+sqlite
select
  arn,
  cluster_name,
  state,
  creation_time
from
  aws_msk_cluster
where
  json_extract(provisioned, '$.LoggingInfo') is null;
```

### Get total storage used by all the clusters
Determine the total storage utilized by all clusters to manage resources efficiently and optimize cost. This can help in understanding the overall storage usage and planning for future scaling needs.

```sql+postgres
select
  sum((provisioned -> 'BrokerNodeGroupInfo' -> 'StorageInfo' -> 'EbsStorageInfo' ->> 'VolumeSize')::int) as total_storage
from
  aws_msk_cluster;
```

```sql+sqlite
select
  sum(json_extract(provisioned, '$.BrokerNodeGroupInfo.StorageInfo.EbsStorageInfo.VolumeSize')) as total_storage
from
  aws_msk_cluster;
```