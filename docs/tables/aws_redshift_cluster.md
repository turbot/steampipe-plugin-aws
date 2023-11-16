---
title: "Table: aws_redshift_cluster - Query AWS Redshift Clusters using SQL"
description: "Allows users to query AWS Redshift Clusters and retrieve comprehensive information about each cluster, including its configuration, status, performance, and security settings."
---

# Table: aws_redshift_cluster - Query AWS Redshift Clusters using SQL

The `aws_redshift_cluster` table in Steampipe provides information about Redshift clusters within Amazon Web Services. This table allows DevOps engineers to query cluster-specific details, such as cluster status, node type, number of nodes, and associated metadata. Users can utilize this table to gather insights on clusters, including their availability, performance, and security settings. The schema outlines the various attributes of the Redshift cluster, including the cluster identifier, creation time, database name, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_redshift_cluster` table, you can use the `.inspect aws_redshift_cluster` command in Steampipe.

### Key columns:

- `cluster_identifier`: This is the unique identifier of the cluster. It can be used to join this table with other tables that also contain Redshift cluster identifiers.
- `db_name`: This is the name of the initial database that was created when the cluster was created. It can be used to join with tables that contain information about specific databases in Redshift.
- `vpc_id`: This is the identifier of the virtual private cloud (VPC) that the cluster is in. It can be used to join this table with other tables that contain VPC-specific information.

## Examples

### Basic info

```sql
select
  cluster_identifier,
  arn,
  node_type,
  region
from
  aws_redshift_cluster;
```

### List clusters that are publicly accessible

```sql
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  publicly_accessible
from
  aws_redshift_cluster
where
  publicly_accessible;
```

### List clusters that are not in a VPC

```sql
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  vpc_id
from
  aws_redshift_cluster
where
  vpc_id is null;
```

### List clusters whose storage is not encrypted

```sql
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  encrypted
from
  aws_redshift_cluster
where
  not encrypted;
```

### Get logging status for each cluster

```sql
select
  cluster_identifier,
  logging_status -> 'LoggingEnabled' as LoggingEnabled
from
  aws_redshift_cluster
```
