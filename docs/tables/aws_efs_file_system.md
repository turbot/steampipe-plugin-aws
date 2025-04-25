---
title: "Steampipe Table: aws_efs_file_system - Query AWS Elastic File System using SQL"
description: "Allows users to query AWS Elastic File System (EFS) file systems, providing detailed information about each file system such as its ID, ARN, creation token, performance mode, and lifecycle state."
folder: "EFS"
---

# Table: aws_efs_file_system - Query AWS Elastic File System using SQL

The AWS Elastic File System (EFS) is a scalable file storage for use with Amazon EC2 instances. It's easy to use and offers a simple interface that allows you to create and configure file systems quickly and easily. With EFS, you have the flexibility to store and retrieve data across different AWS regions and availability zones.

## Table Usage Guide

The `aws_efs_file_system` table in Steampipe provides you with information about file systems within AWS Elastic File System (EFS). This table allows you, as a DevOps engineer, to query file system-specific details, including its ID, ARN, creation token, performance mode, lifecycle state, and associated metadata. You can utilize this table to gather insights on file systems, such as their performance mode, lifecycle state, and more. The schema outlines the various attributes of the EFS file system for you, including the file system ID, creation token, tags, and associated mount targets.

## Examples

### Basic info
Discover the segments that have automatic backups enabled in your AWS Elastic File System (EFS). This helps in assessing the elements within your system that are safeguarded and those that might need additional data protection measures.

```sql+postgres
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

```sql+sqlite
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
Discover the segments of your AWS Elastic File System that are not encrypted, allowing you to identify potential security risks and take necessary action to ensure data protection.

```sql+postgres
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

```sql+sqlite
select
  file_system_id,
  encrypted,
  kms_key_id,
  region
from
  aws_efs_file_system
where
  encrypted = 0;
```


### Get the size of the data stored in each file system
Assess the elements within your file system to understand the distribution of data storage. This is useful for managing storage resources effectively and identifying opportunities for cost optimization.

```sql+postgres
select
  file_system_id,
  size_in_bytes ->> 'Value' as data_size,
  size_in_bytes ->> 'Timestamp' as data_size_timestamp,
  size_in_bytes ->> 'ValueInIA' as data_size_infrequent_access_storage,
  size_in_bytes ->> 'ValueInStandard' as data_size_standard_storage
from
  aws_efs_file_system;
```

```sql+sqlite
select
  file_system_id,
  json_extract(size_in_bytes, '$.Value') as data_size,
  json_extract(size_in_bytes, '$.Timestamp') as data_size_timestamp,
  json_extract(size_in_bytes, '$.ValueInIA') as data_size_infrequent_access_storage,
  json_extract(size_in_bytes, '$.ValueInStandard') as data_size_standard_storage
from
  aws_efs_file_system;
```


### List file systems which have root access
Identify instances where file systems have root access, which can be critical in understanding the security posture of your AWS Elastic File System, and ensuring that only authorized users have such elevated privileges.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(principal.value, '$') as principal,
  json_extract(action.value, '$') as action,
  json_extract(statement.value, '$.Effect') as effect,
  json_extract(statement.value, '$.Condition') as conditions
from
  aws_efs_file_system,
  json_each(policy_std, '$.Statement') as statement,
  json_each(json_extract(statement.value, '$.Principal.AWS')) as principal,
  json_each(json_extract(statement.value, '$.Action')) as action
where
  json_extract(action.value, '$') = 'elasticfilesystem:clientrootaccess';
```

### List file systems that do not enforce encryption in transit
Discover the segments of your AWS Elastic File System that are not enforcing encryption in transit. This can help improve your system's security by identifying potential vulnerabilities.

```sql+postgres
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

```sql+sqlite
select
  title
from
  aws_efs_file_system
where
  title not in (
    select
      title
    from
      aws_efs_file_system
    where
      json_extract(policy_std, '$.Statement[*].Principal.AWS') = '*'
      and json_extract(policy_std, '$.Statement[*].Effect') = 'Deny'
      and json_extract(policy_std, '$.Statement[*].Condition.Bool.aws:securetransport') = 'false'
  );
```


### List file systems with automatic backups enabled
Gain insights into the file systems that have automatic backups enabled. This is useful for ensuring that your data is being regularly backed up for recovery purposes.

```sql+postgres
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

```sql+sqlite
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