# Table: aws_dms_replication_instance

AWS Database Migration Service (DMS) replication instance is used to connect to your source data store, read the source data, and format the data for consumption by the target data store along with loading the data into the target data store.

## Examples

### Basic info

```sql
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  kms_key_id,
  publicly_accessible,
  region
from
  aws_dms_replication_instance;
```


### List replication instances with auto minor version upgrades disabled

```sql
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  auto_minor_version_upgrade,
  region
from
  aws_dms_replication_instance
where
  not auto_minor_version_upgrade;
```


### List replication instances provisioned with undesired (for example, dms.r5.16xlarge and dms.r5.24xlarge are not desired) instance classes

```sql
select
  replication_instance_identifier,
  arn,
  engine_version,
  instance_create_time,
  replication_instance_class,
  region
from
  aws_dms_replication_instance
where
  replication_instance_class not in ('dms.r5.16xlarge', 'dms.r5.24xlarge');
```


### List publicly accessible replication instances

```sql
select
  replication_instance_identifier,
  arn,
  publicly_accessible,
  region
from
  aws_dms_replication_instance
where
  publicly_accessible;
```


### List replication instances not using multi-AZ deployment configurations

```sql
select
  replication_instance_identifier,
  arn,
  publicly_accessible,
  multi_az,
  region
from
  aws_dms_replication_instance
where
  not multi_az;
```
