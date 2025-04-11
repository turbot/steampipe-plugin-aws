---
title: "Steampipe Table: aws_dax_cluster - Query AWS DAX Clusters using SQL"
description: "Allows users to query AWS DAX Clusters to fetch details about their configurations, status, nodes, and other associated metadata."
folder: "DAX"
---

# Table: aws_dax_cluster - Query AWS DAX Clusters using SQL

The AWS DAX Cluster is a fully managed, highly available, in-memory cache for DynamoDB that delivers up to a 10x read performance improvement. It operates fully in-memory and is compatible with existing DynamoDB API calls. DAX does all the heavy lifting to deliver accelerated read performance and can be used without application changes.

## Table Usage Guide

The `aws_dax_cluster` table in Steampipe provides you with information about AWS DAX Clusters. This table allows you, as a DevOps engineer, to query cluster-specific details, including cluster names, node types, status, and associated metadata. You can utilize this table to gather insights on clusters, such as cluster configurations, status, nodes, and more. The schema outlines the various attributes of the DAX cluster for you, including the cluster name, ARN, status, node type, and associated tags.

## Examples

### Basic info
Determine the status and region of active nodes in your AWS DAX clusters to understand their configuration and performance. This helps in managing resources and planning for scalability.

```sql+postgres
select
  cluster_name,
  description,
  active_nodes,
  iam_role_arn,
  status,
  region
from
  aws_dax_cluster;
```

```sql+sqlite
select
  cluster_name,
  description,
  active_nodes,
  iam_role_arn,
  status,
  region
from
  aws_dax_cluster;
```


### List clusters that does not enforce server-side encryption (SSE)
Determine the areas in your AWS DAX clusters where server-side encryption is not enforced. This is beneficial for identifying potential security vulnerabilities within your system.

```sql+postgres
select
  cluster_name,
  description,
  sse_description ->> 'Status' as sse_status
from
  aws_dax_cluster
where
  sse_description ->> 'Status' = 'DISABLED';
```

```sql+sqlite
select
  cluster_name,
  description,
  json_extract(sse_description, '$.Status') as sse_status
from
  aws_dax_cluster
where
  json_extract(sse_description, '$.Status') = 'DISABLED';
```


### List clusters provisioned with undesired (for example, cache.m5.large and cache.m4.4xlarge are desired) node types
Determine the areas in which clusters are provisioned with non-preferred node types to optimize resource allocation and cost efficiency.

```sql+postgres
select
  cluster_name,
  node_type,
  count(*) as count
from
  aws_dax_cluster
where
  node_type not in ('cache.m5.large', 'cache.m4.4xlarge')
group by
  cluster_name, node_type;
```

```sql+sqlite
select
  cluster_name,
  node_type,
  count(*) as count
from
  aws_dax_cluster
where
  node_type not in ('cache.m5.large', 'cache.m4.4xlarge')
group by
  cluster_name, node_type;
```


### Get the network details for each cluster
Discover the segments that provide detailed network information for each cluster, including security group identifiers and availability zones. This can be useful for understanding the network configuration of your clusters and ensuring they are set up correctly.

```sql+postgres
select
  cluster_name,
  subnet_group,
  sg ->> 'SecurityGroupIdentifier' as sg_id,
  n ->> 'AvailabilityZone' as az_name,
  cluster_discovery_endpoint ->> 'Address' as cluster_discovery_endpoint_address,
  cluster_discovery_endpoint ->> 'Port' as cluster_discovery_endpoint_port
from
  aws_dax_cluster,
  jsonb_array_elements(security_groups) as sg,
  jsonb_array_elements(nodes) as n;
```

```sql+sqlite
select
  cluster_name,
  subnet_group,
  json_extract(sg.value, '$.SecurityGroupIdentifier') as sg_id,
  json_extract(n.value, '$.AvailabilityZone') as az_name,
  json_extract(cluster_discovery_endpoint, '$.Address') as cluster_discovery_endpoint_address,
  json_extract(cluster_discovery_endpoint, '$.Port') as cluster_discovery_endpoint_port
from
  aws_dax_cluster,
  json_each(security_groups) as sg,
  json_each(nodes) as n;
```