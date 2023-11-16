---
title: "Table: aws_dax_cluster - Query AWS DAX Clusters using SQL"
description: "Allows users to query AWS DAX Clusters to fetch details about their configurations, status, nodes, and other associated metadata."
---

# Table: aws_dax_cluster - Query AWS DAX Clusters using SQL

The `aws_dax_cluster` table in Steampipe provides information about AWS DAX Clusters. This table allows DevOps engineers to query cluster-specific details, including cluster names, node types, status, and associated metadata. Users can utilize this table to gather insights on clusters, such as cluster configurations, status, nodes, and more. The schema outlines the various attributes of the DAX cluster, including the cluster name, ARN, status, node type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_dax_cluster` table, you can use the `.inspect aws_dax_cluster` command in Steampipe.

**Key columns**:

- `cluster_name` - The name of the DAX cluster. It can be used to join with other tables where the cluster name is needed.
- `arn` - The Amazon Resource Name (ARN) of the DAX cluster. It can be used to join with other tables where the cluster ARN is required.
- `status` - The current status of the cluster. Useful for monitoring and managing the cluster's state.

## Examples

### Basic info

```sql
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

```sql
select
  cluster_name,
  description,
  sse_description ->> 'Status' as sse_status
from
  aws_dax_cluster
where
  sse_description ->> 'Status' = 'DISABLED';
```


### List clusters provisioned with undesired (for example, cache.m5.large and cache.m4.4xlarge are desired) node types

```sql
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

```sql
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
