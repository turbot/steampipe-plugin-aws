---
title: "Steampipe Table: aws_datasync_task - Query AWS DataSync Tasks using SQL"
description: "Allows users to query AWS DataSync Tasks to retrieve detailed information about each task configuration."
folder: "DataSync"
---

# Table: aws_datasync_task - Query AWS DataSync Tasks using SQL

AWS DataSync is a data transfer service that makes it easy and fast to move data between on-premises storage systems and AWS Storage services, as well as between AWS Storage services. DataSync Tasks define the data transfer operations, including source and destination locations, transfer options, and scheduling.

## Table Usage Guide

The `aws_datasync_task` table in Steampipe provides you with information about DataSync tasks within AWS. This table allows you, as a DevOps engineer, to query task-specific details, including task ARN, status, source and destination locations, transfer options, and associated metadata. You can utilize this table to gather insights on tasks, such as their current status, error details, transfer configurations, and associated tags.

## Examples

### Basic info
Explore the features and settings of your AWS DataSync tasks to better understand their configuration, such as status, source and destination locations, and regional distribution. This can help in assessing task performance and operational efficiency.

```sql+postgres
select
  name,
  arn,
  status,
  source_location_arn,
  destination_location_arn,
  creation_time,
  region
from
  aws_datasync_task;
```

```sql+sqlite
select
  name,
  arn,
  status,
  source_location_arn,
  destination_location_arn,
  creation_time,
  region
from
  aws_datasync_task;
```

### Get task options configuration details
Determine the transfer options for each task to understand how data is being transferred, including bandwidth limits, verification settings, and other configuration options.

```sql+postgres
select
  name,
  status,
  options::json ->> 'VerifyMode' as verify_mode,
  options::json ->> 'Atime' as atime,
  options::json ->> 'Mtime' as mtime,
  options::json ->> 'Uid' as uid,
  options::json ->> 'Gid' as gid,
  options::json ->> 'PreserveDeletedFiles' as preserve_deleted_files,
  options::json ->> 'PreserveDevices' as preserve_devices,
  options::json ->> 'PosixPermissions' as posix_permissions,
  options::json ->> 'BytesPerSecond' as bytes_per_second,
  options::json ->> 'TaskQueueing' as task_queueing,
  options::json ->> 'LogLevel' as log_level,
  options::json ->> 'TransferMode' as transfer_mode
from
  aws_datasync_task;
```

```sql+sqlite
select
  name,
  status,
  json_extract(options, '$.VerifyMode') as verify_mode,
  json_extract(options, '$.Atime') as atime,
  json_extract(options, '$.Mtime') as mtime,
  json_extract(options, '$.Uid') as uid,
  json_extract(options, '$.Gid') as gid,
  json_extract(options, '$.PreserveDeletedFiles') as preserve_deleted_files,
  json_extract(options, '$.PreserveDevices') as preserve_devices,
  json_extract(options, '$.PosixPermissions') as posix_permissions,
  json_extract(options, '$.BytesPerSecond') as bytes_per_second,
  json_extract(options, '$.TaskQueueing') as task_queueing,
  json_extract(options, '$.LogLevel') as log_level,
  json_extract(options, '$.TransferMode') as transfer_mode
from
  aws_datasync_task;
```

### List tasks with errors
Identify DataSync tasks that have encountered errors during execution. This is useful for troubleshooting and monitoring task health.

```sql+postgres
select
  name,
  arn,
  status,
  error_code,
  error_detail,
  creation_time
from
  aws_datasync_task
where
  error_code is not null;
```

```sql+sqlite
select
  name,
  arn,
  status,
  error_code,
  error_detail,
  creation_time
from
  aws_datasync_task
where
  error_code is not null;
```

### Find tasks with specific transfer modes
Find DataSync tasks with specific transfer modes to understand how data is being transferred and optimize transfer strategies.

```sql+postgres
select
  name,
  status,
  options::json ->> 'TransferMode' as transfer_mode,
  source_location_arn,
  destination_location_arn
from
  aws_datasync_task
where
  options::json ->> 'TransferMode' = 'CHANGED';
```

```sql+sqlite
select
  name,
  status,
  json_extract(options, '$.TransferMode') as transfer_mode,
  source_location_arn,
  destination_location_arn
from
  aws_datasync_task
where
  json_extract(options, '$.TransferMode') = 'CHANGED';
```

### List tasks with bandwidth limits
Identify DataSync tasks that have bandwidth limits configured to understand network usage and optimize transfer performance.

```sql+postgres
select
  name,
  status,
  options::json ->> 'BytesPerSecond' as bytes_per_second,
  creation_time
from
  aws_datasync_task
where
  options::json ->> 'BytesPerSecond' is not null;
```

```sql+sqlite
select
  name,
  status,
  json_extract(options, '$.BytesPerSecond') as bytes_per_second,
  creation_time
from
  aws_datasync_task
where
  json_extract(options, '$.BytesPerSecond') is not null;
```
