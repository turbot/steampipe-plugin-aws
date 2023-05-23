# Table: aws_rds_db_instance_automated_backup

Amazon RDS creates and saves automated backups of your DB instance or Multi-AZ DB cluster during the backup window of your database. RDS creates a storage volume snapshot of your database, backing up the entire database and not just individual databases.

## Examples

### Basic info

```sql
select
  db_instance_identifier,
  arn,
  status,
  allocated_storage,
  encrypted,
  engine
from
  aws_rds_db_instance_automated_backup;
```

### List DB instance automated backups that are not encrypted

```sql
select
  db_instance_identifier,
  arn,
  status,
  backup_target,
  instance_create_time,
  encrypted,
  engine
from
  aws_rds_db_instance_automated_backup
where
  not encrypted;
```

### List DB instance automated backups that are not authenticated through IAM users and roles

```sql
select
  db_instance_identifier,
  iam_database_authentication_enabled,
  status,
  availability_zone,
  dbi_resource_id
from
  aws_rds_db_instance_automated_backup
where
  not iam_database_authentication_enabled;
```

### Get VPC and subnet info for each DB instance automated backup

```sql
select
  b.arn,
  b.vpc_id,
  v.cidr_block,
  v.is_default,
  v.instance_tenancy
from
  aws_rds_db_instance_automated_backup as b,
  aws_vpc as v
where
  v.vpc_id = b.vpc_id;
```

### List DB instance automated backups of deleted instances

```sql
select
  db_instance_identifier,
  arn,
  engine,
  engine_version,
  availability_zone,
  backup_retention_period,
  status
from
  aws_rds_db_instance_automated_backup
where
  status = 'retained';
```

### Get KMS key details of each DB instance automated backup

```sql
select
  b.db_instance_identifier,
  b.arn as automated_backup_arn,
  b.engine,
  b.kms_key_id,
  k.creation_date as kms_key_creation_date,
  k.key_state,
  k.key_rotation_enabled
from
  aws_rds_db_instance_automated_backup as b,
  aws_kms_key as k
where
  k.id = b.kms_key_id;
```
