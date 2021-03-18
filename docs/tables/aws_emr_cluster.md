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


### List of EMR clusters whose auto terminate is disabled

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


### List of EMR clusters which got terminated due to some error

```sql
select
  id,
  name,
  status ->> 'State' as state,
  status -> 'StateChangeReason' ->> 'Message'  as  state_change_reason
from
  aws_emr_cluster
where
  status ->> 'State' = 'TERMINATED_WITH_ERRORS';
```


### List of application names and its versions installed on EMR clusters

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