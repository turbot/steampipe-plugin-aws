# Table: aws_emr_cluster

EMR automatically configures EC2 firewall settings, controlling network access to instances and launches clusters in an Amazon Virtual Private Cloud (VPC).

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


### List clusters which have auto terminate disabled

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


### List application names and its versions installed on each cluster

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


### List clusters for which logging is disabled

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


### List clusters for which logging is enabled but log encryption is disabled

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
