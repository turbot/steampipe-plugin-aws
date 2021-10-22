# Table: aws_fsx_file_system

Amazon FSx makes it easy and cost effective to launch, run, and scale feature-rich, high-performance file systems in the cloud. It supports a wide range of workloads with its reliability, security, scalability, and broad set of capabilities.

## Examples

### Basic info

```sql
select
  file_system_id,
  arn,
  dns_name,
  owner_id,
  creation_time,
  lifecycle,
  storage_capacity
from
  aws_fsx_file_system;
```

### List file systems which are encrypted

```sql
select
  file_system_id,
  kms_key_id,
  region
from
  aws_fsx_file_system
where
  kms_key_id is not null;
```