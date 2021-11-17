# Table: aws_media_store_container

AWS Elemental MediaStore is a video origination and storage service that offers the high performance, predictable low latency, and immediate consistency required for live origination. You use containers in MediaStore to store your folders and objects. Related objects can be grouped in containers in the same way that you use a directory to group files in a file system. 

## Examples

### Basic info

```sql
select
  name,
  arn,
  status,
  access_logging_enabled,
  creation_time,
  endpoint
from
  aws_media_store_container;
```

### List containers which are in 'CREATING' state

```sql
select
  name,
  arn,
  status,
  access_logging_enabled,
  creation_time,
  endpoint
from
  aws_media_store_container
where
  status = 'CREATING';
```

### List policy details for the containers

```sql
select
  name,
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_media_store_container;
```

### List containers with access logging enabled

```sql
select
  name,
  arn,
  access_logging_enabled
from
  aws_media_store_container
where
  access_logging_enabled;
```
