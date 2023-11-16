---
title: "Table: aws_emr_cluster - Query AWS Elastic MapReduce Cluster using SQL"
description: "Allows users to query AWS Elastic MapReduce Cluster data."
---

# Table: aws_emr_cluster - Query AWS Elastic MapReduce Cluster using SQL

The `aws_emr_cluster` table in Steampipe provides information about clusters within AWS Elastic MapReduce (EMR). This table allows data engineers to query cluster-specific details, including cluster status, hardware and software configurations, VPC settings, and associated metadata. Users can utilize this table to gather insights on EMR clusters, such as cluster states, hardware and software configurations, and verification of VPC settings. The schema outlines the various attributes of the EMR cluster, including the cluster ID, name, status, normalized instance hours, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_emr_cluster` table, you can use the `.inspect aws_emr_cluster` command in Steampipe.

### Key columns:

- `id`: The unique identifier of the cluster. This can be used to join this table with other tables that contain cluster information.
- `name`: The name of the cluster. This can be used to join with tables that contain cluster names for more human-readable queries.
- `status`: The current status of the cluster. This can be useful for joining with tables that contain status information to filter or sort clusters based on their current status.

## Examples

### Basic info

```sql
select
  id,
  cluster_arn,
  name,
  auto_terminate,
  status ->> 'State' as state,
  tags
from
  aws_emr_cluster;
```


### List clusters with auto-termination disabled

```sql
select
  name,
  cluster_arn,
  auto_terminate
from
  aws_emr_cluster
where
  not auto_terminate;
```


### List clusters which have terminated with errors

```sql
select
  id,
  name,
  status ->> 'State' as state,
  status -> 'StateChangeReason' ->> 'Message' as state_change_reason
from
  aws_emr_cluster
where
  status ->> 'State' = 'TERMINATED_WITH_ERRORS';
```


### Get application names and versions installed for each cluster

```sql
select
  name,
  cluster_arn,
  a ->> 'Name' as application_name,
  a ->> 'Version' as application_version
from
  aws_emr_cluster,
  jsonb_array_elements(applications) as a;
```


### List clusters with logging disabled

```sql
select
  name,
  cluster_arn,
  log_uri
from
  aws_emr_cluster
where
  log_uri is null
```


### List clusters with logging enabled but log encryption is disabled

```sql
select
  name,
  cluster_arn,
  log_uri,
  log_encryption_kms_key_id
from
  aws_emr_cluster
where
  log_uri is not null and log_encryption_kms_key_id is null;
```
