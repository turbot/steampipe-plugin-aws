---
title: "Steampipe Table: aws_neptune_db_cluster - Query Amazon Neptune DB Clusters using SQL"
description: "Allows users to query Amazon Neptune DB clusters for comprehensive information about their configuration, status, and other relevant details."
folder: "Neptune"
---

# Table: aws_neptune_db_cluster - Query Amazon Neptune DB Clusters using SQL

Amazon Neptune DB Clusters are a scalable and reliable graph database service that is optimized for storing billions of relationships and querying the graph with milliseconds latency. Neptune supports popular graph models Property Graph and W3C's RDF, and their respective query languages Apache TinkerPop Gremlin and SPARQL. It is designed to offer fast and reliable access to your graph applications, allowing you to build and run applications that work with highly connected datasets.

## Table Usage Guide

The `aws_neptune_db_cluster` table in Steampipe provides you with information about DB clusters within Amazon Neptune. This table allows you, as a DevOps engineer, to query DB cluster-specific details, including configuration, status, and associated metadata. You can utilize this table to gather insights on DB clusters, such as their availability, security settings, backup policies, and more. The schema outlines the various attributes of the DB cluster for you, including the cluster identifier, creation time, enabled cloudwatch logs exports, and associated tags.

**Important Notes**
- This table only returns Neptune DB clusters for you, not RDS or DocumentDB DB clusters.

## Examples

### List of DB clusters which are not encrypted
Discover the segments of your database clusters that are not encrypted. This query is useful to identify potential security risks and ensure compliance with data protection standards.

```sql+postgres
select
  db_cluster_identifier,
  allocated_storage,
  kms_key_id
from
  aws_neptune_db_cluster
where
  kms_key_id is null;
```

```sql+sqlite
select
  db_cluster_identifier,
  allocated_storage,
  kms_key_id
from
  aws_neptune_db_cluster
where
  kms_key_id is null;
```

### List of DB clusters where backup retention period is greater than 7 days
Identify instances where database clusters have a backup retention period exceeding a week. This could be useful for ensuring data safety and compliance with company policies or regulations.

```sql+postgres
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_neptune_db_cluster
where
  backup_retention_period > 7;
```

```sql+sqlite
select
  db_cluster_identifier,
  backup_retention_period
from
  aws_neptune_db_cluster
where
  backup_retention_period > 7;
```

### Avalability zone count for each db instance
Explore which database clusters have varying availability zone counts, allowing you to assess the distribution and redundancy of your databases for better resource management and disaster recovery planning.

```sql+postgres
select
  db_cluster_identifier,
  jsonb_array_length(availability_zones) availability_zones_count
from
  aws_neptune_db_cluster;
```

```sql+sqlite
select
  db_cluster_identifier,
  json_array_length(availability_zones) as availability_zones_count
from
  aws_neptune_db_cluster;
```

### DB cluster Members info
Explore the configuration and status of various components within your AWS Neptune database clusters. This query will help you identify which instances are cluster writers and their respective promotion tiers, providing valuable insights for managing and optimizing your database operations.

```sql+postgres
select
  db_cluster_identifier,
  member ->> 'DBClusterParameterGroupStatus' as db_cluster_parameter_group_status,
  member ->> 'DBInstanceIdentifier' as db_instance_identifier,
  member ->> 'IsClusterWriter' as is_cluster_writer,
  member ->> 'PromotionTier' as promotion_tier
from
  aws_neptune_db_cluster
  cross join jsonb_array_elements(db_cluster_members) as member;
```

```sql+sqlite
select
  db_cluster_identifier,
  json_extract(member.value, '$.DBClusterParameterGroupStatus') as db_cluster_parameter_group_status,
  json_extract(member.value, '$.DBInstanceIdentifier') as db_instance_identifier,
  json_extract(member.value, '$.IsClusterWriter') as is_cluster_writer,
  json_extract(member.value, '$.PromotionTier') as promotion_tier
from
  aws_neptune_db_cluster,
  json_each(db_cluster_members) as member;
```