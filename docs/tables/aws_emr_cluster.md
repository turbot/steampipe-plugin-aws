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
  status,
  applications,
  configurations,
  log_uri,
  security_configuration,
  ebs_root_volume_size
from
  aws_new.aws_emr_cluster;
```


### List of emr clusters which are in terminating or terminated state

```sql
select
  id,
  name,
  cluster_arn,
  status ->>  'State' as state
from
  aws_new.aws_emr_cluster
where
  status ->>  'State' IN ('TERMINATED','TERMINATED_WITH_ERRORS','TERMINATING');
```


### List of Emr clusters whose auto terminate is enabled

```sql
select
  id,
  name,
  cluster_arn,
  auto_terminate
from
  aws_new.aws_emr_cluster
where
  auto_terminate = 'true';
```


### List of applications and its versions installed on emr clusters

```sql
select
  name,
  id,
  a ->> 'Name' as name,
  a ->> 'Version' as version
from
  aws_new.aws_emr_cluster,
  jsonb_array_elements(applications) as a
```

