# Table: aws_efs_file_system

Amazon Elastic File System is a cloud storage service provided by Amazon Web Services designed to provide scalable, elastic, concurrent with some restrictions, and encrypted file storage for use with both AWS cloud services and on-premises resources.

## Examples

### Basic info

```sql
select
  name,
  file_system_id,
  owner_id,
  automatic_backups,
  creation_token,
  creation_time,
  life_cycle_state,
  number_of_mount_targets,
  performance_mode,
  throughput_mode
from
  aws_efs_file_system;
```


### List file systems which are not encrypted at rest

```sql
select
  file_system_id,
  encrypted,
  kms_key_id,
  region
from
  aws_efs_file_system
where
  not encrypted;
```


### Get the size of the data stored in each file system

```sql
select
  file_system_id,
  size_in_bytes ->> 'Value' as data_size,
  size_in_bytes ->> 'Timestamp' as data_size_timestamp,
  size_in_bytes ->> 'ValueInIA' as data_size_infrequent_access_storage,
  size_in_bytes ->> 'ValueInStandard' as data_size_standard_storage
from
  aws_efs_file_system;
```


### List file systems which have root access

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_efs_file_system,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  a in ('elasticfilesystem:clientrootaccess');
```


### List file systems that do not enforce encryption in transit

```sql
select
  title
from
  aws_efs_file_system
where
  title not in (
    select
      title
    from
      aws_efs_file_system,
      jsonb_array_elements(policy_std -> 'Statement') as s,
      jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
      jsonb_array_elements_text(s -> 'Action') as a,
      jsonb_array_elements_text(
        s -> 'Condition' -> 'Bool' -> 'aws:securetransport'
      ) as ssl
    where
      p = '*'
      and s ->> 'Effect' = 'Deny'
      and ssl :: bool = false
  );
```


### List file systems with automatic backups enabled

```sql
select
  name,
  automatic_backups,
  arn,
  file_system_id
from
  aws_efs_file_system
where
  automatic_backups = 'enabled';
```
