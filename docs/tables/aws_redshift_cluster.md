---
title: "Steampipe Table: aws_redshift_cluster - Query AWS Redshift Clusters using SQL"
description: "Allows users to query AWS Redshift Clusters and retrieve comprehensive information about each cluster, including its configuration, status, performance, and security settings."
folder: "Redshift"
---

# Table: aws_redshift_cluster - Query AWS Redshift Clusters using SQL

The AWS Redshift Cluster is a fully managed, petabyte-scale data warehouse service in the cloud. It enables users to analyze all their data using their existing business intelligence tools by leveraging SQL and existing business intelligence tools. It provides a cost-effective solution for analyzing large-scale data sets as it uses columnar storage technology and parallel query execution.

## Table Usage Guide

The `aws_redshift_cluster` table in Steampipe provides you with information about Redshift clusters within Amazon Web Services. This table allows you, as a DevOps engineer, to query cluster-specific details, such as cluster status, node type, number of nodes, and associated metadata. You can utilize this table to gather insights on clusters, including their availability, performance, and security settings. The schema outlines the various attributes of the Redshift cluster for you, including the cluster identifier, creation time, database name, and associated tags.

## Examples

### Basic info
Analyze the settings to understand the different types of nodes and their geographical distribution in your AWS Redshift clusters. This can help optimize resource allocation and improve data management across different regions.

```sql+postgres
select
  cluster_identifier,
  arn,
  node_type,
  region
from
  aws_redshift_cluster;
```

```sql+sqlite
select
  cluster_identifier,
  arn,
  node_type,
  region
from
  aws_redshift_cluster;
```

### List clusters that are publicly accessible
Determine the areas in which your clusters are publicly accessible, enabling you to identify potential security risks and take necessary precautions.

```sql+postgres
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

```sql+sqlite
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  publicly_accessible
from
  aws_redshift_cluster
where
  publicly_accessible = 1;
```

### List clusters that are not in a VPC
Discover the segments that are not part of any Virtual Private Cloud (VPC) in your AWS Redshift clusters. This can be useful in identifying potential security risks or compliance issues, as clusters not in a VPC may be more exposed to external threats.

```sql+postgres
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

```sql+sqlite
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
Identify the instances where your clusters' storage is not encrypted. This can help enhance your data security by pinpointing areas that need encryption.

```sql+postgres
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

```sql+sqlite
select
  cluster_identifier,
  node_type,
  number_of_nodes,
  encrypted
from
  aws_redshift_cluster
where
  encrypted = 0;
```

### Get logging status for each cluster
Determine the logging status for each cluster in your AWS Redshift environment to ensure proper logging practices are being followed for data security and compliance purposes.

```sql+postgres
select
  cluster_identifier,
  logging_status -> 'LoggingEnabled' as LoggingEnabled
from
  aws_redshift_cluster
```

```sql+sqlite
select
  cluster_identifier,
  json_extract(logging_status, '$.LoggingEnabled') as LoggingEnabled
from
  aws_redshift_cluster
```