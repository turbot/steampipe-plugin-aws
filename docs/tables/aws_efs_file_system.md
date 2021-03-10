# Table: aws_efs_file_system

Amazon EFS provides a durable, high throughput file system for content management systems and web serving applications that store and serve information for a range of applications like websites, online publications, and archives.

## Examples

### List of unencrypted elastic file systems

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


### Size of the data stored in each file system

```sql
select
  file_system_id,
  size_in_bytes ->> 'Value' as data_size,
  size_in_bytes ->> 'Timestamp' as data_size_timestamp,
  size_in_bytes ->> 'ValueInIA' as data_size_infrequent_access_storage,
  size_in_bytes ->> 'ValueInIA' as data_size_standard_storage
from
  aws_efs_file_system;
```


### List of file systems which has root access

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


### List of file systems that DO NOT enforce encryption in transit

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