---
title: "Steampipe Table: aws_emr_cluster - Query AWS Elastic MapReduce Cluster using SQL"
description: "Allows users to query AWS Elastic MapReduce Cluster data."
folder: "EMR"
---

# Table: aws_emr_cluster - Query AWS Elastic MapReduce Cluster using SQL

The AWS Elastic MapReduce (EMR) Cluster is a web service that makes it easy to process large amounts of data efficiently. EMR uses Hadoop processing combined with several AWS products to do tasks such as web indexing, data mining, log file analysis, machine learning, scientific simulation, and data warehousing. Users can interactively analyze their data to achieve faster time-to-insights.

## Table Usage Guide

The `aws_emr_cluster` table in Steampipe provides you with information about clusters within AWS Elastic MapReduce (EMR). This table allows you as a data engineer to query cluster-specific details, including cluster status, hardware and software configurations, VPC settings, and associated metadata. You can utilize this table to gather insights on EMR clusters, such as cluster states, hardware and software configurations, and verification of VPC settings. The schema outlines the various attributes of the EMR cluster for you, including the cluster ID, name, status, normalized instance hours, and associated tags.

## Examples

### Basic info
Explore the status and termination settings of your AWS EMR clusters to manage resources effectively. This helps in identifying clusters that are in use and those that can be terminated to save costs.

```sql+postgres
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

```sql+sqlite
select
  id,
  cluster_arn,
  name,
  auto_terminate,
  json_extract(status, '$.State') as state,
  tags
from
  aws_emr_cluster;
```


### List clusters with auto-termination disabled
Determine the areas in which clusters are operating with auto-termination disabled, which could potentially lead to unnecessary resource usage and increased costs.

```sql+postgres
select
  name,
  cluster_arn,
  auto_terminate
from
  aws_emr_cluster
where
  not auto_terminate;
```

```sql+sqlite
select
  name,
  cluster_arn,
  auto_terminate
from
  aws_emr_cluster
where
  auto_terminate = 0;
```


### List clusters which have terminated with errors
Identify instances where clusters have ended with errors. This allows you to pinpoint specific locations where issues have occurred, enabling efficient troubleshooting and problem resolution.

```sql+postgres
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

```sql+sqlite
select
  id,
  name,
  json_extract(status, '$.State') as state,
  json_extract(json_extract(status, '$.StateChangeReason'), '$.Message') as state_change_reason
from
  aws_emr_cluster
where
  json_extract(status, '$.State') = 'TERMINATED_WITH_ERRORS';
```


### Get application names and versions installed for each cluster
Determine the applications and their respective versions installed across different clusters. This is useful for tracking software versions and ensuring consistency across your cluster environment.

```sql+postgres
select
  name,
  cluster_arn,
  a ->> 'Name' as application_name,
  a ->> 'Version' as application_version
from
  aws_emr_cluster,
  jsonb_array_elements(applications) as a;
```

```sql+sqlite
select
  name,
  cluster_arn,
  json_extract(a.value, '$.Name') as application_name,
  json_extract(a.value, '$.Version') as application_version
from
  aws_emr_cluster,
  json_each(applications) as a;
```


### List clusters with logging disabled
Determine the areas in which logging is disabled in your clusters. This is useful for identifying potential gaps in your data tracking and ensuring comprehensive monitoring across all clusters.

```sql+postgres
select
  name,
  cluster_arn,
  log_uri
from
  aws_emr_cluster
where
  log_uri is null
```

```sql+sqlite
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
Explore clusters where logging is activated but without the added security layer of log encryption. This can help identify potential vulnerabilities in your data security practices.

```sql+postgres
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

```sql+sqlite
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