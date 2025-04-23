---
title: "Steampipe Table: aws_dms_replication_task - Query AWS DMS Replication Tasks using SQL"
description: "Enables users to query AWS DMS Replication Tasks to retrieve detailed information on data migration activities between source and target databases."
folder: "DMS"
---

# Table: aws_dms_replication_task - Query AWS DMS Replication Tasks using SQL

AWS Database Migration Service (DMS) Replication Tasks play a critical role in managing data migrations between source and target databases. These tasks facilitate the entire migration process, supporting various migration types, including full load migrations, ongoing replication to synchronize source and target databases and change data capture (CDC) for applying data modifications.

The `aws_dms_replication_task` table in Steampipe allows for in-depth analysis of replication tasks, providing details such as task identifiers, status, migration types, settings, and endpoint ARNs. This table proves essential for database administrators and DevOps engineers overseeing database migrations, offering comprehensive insights into each task's configuration, progress, and performance.

## Examples

### Basic Info
Query to fetch basic details about DMS replication tasks.

```sql+postgresql
select
  replication_task_identifier,
  arn,
  migration_type,
  status,
  replication_task_creation_date
from
  aws_dms_replication_task;
```

```sql+sqlite
select
  replication_task_identifier,
  arn,
  migration_type,
  status,
  replication_task_creation_date
from
  aws_dms_replication_task;
```

### Tasks with specific migration types
List replication tasks by a specific migration type, such as 'full-load'.

```sql+postgresql
select
  replication_task_identifier,
  migration_type,
  status
from
  aws_dms_replication_task
where
  migration_type = 'full-load';
```

```sql+sqlite
select
  replication_task_identifier,
  migration_type,
  status
from
  aws_dms_replication_task
where
  migration_type = 'full-load';
```

### Replication tasks with failures
Identify replication tasks that have failed, focusing on the last failure message.

```sql+postgresql
select
  replication_task_identifier,
  status,
  last_failure_message
from
  aws_dms_replication_task
where
  status = 'failed';
```

```sql+sqlite
select
  replication_task_identifier,
  status,
  last_failure_message
from
  aws_dms_replication_task
where
  status = 'failed';
```

### Task performance statistics
Examine detailed performance statistics of replication tasks.

```sql+postgresql
select
  replication_task_identifier,
  status,
  replication_task_stats -> 'ElapsedTimeMillis' as elapsed_time_millis,
  replication_task_stats -> 'FreshStartDate' as fresh_start_date,
  replication_task_stats -> 'FullLoadFinishDate' as full_load_finish_date,
  replication_task_stats -> 'FullLoadProgressPercent' as full_load_progress_percent,
  replication_task_stats -> 'FullLoadStartDate' as full_load_start_date,
  replication_task_stats -> 'StartDate' as start_date,
  replication_task_stats -> 'StopDate' as stop_date,
  replication_task_stats -> 'TablesErrored' as tables_errored,
  replication_task_stats -> 'TablesLoaded' as tables_loaded,
  replication_task_stats -> 'TablesLoading' as tables_loading,
  replication_task_stats -> 'TablesQueued' as tables_queued
from
  aws_dms_replication_task;
```

```sql+sqlite
select
  replication_task_identifier,
  status,
  json_extract(replication_task_stats, '$.ElapsedTimeMillis') as elapsed_time_millis,
  json_extract(replication_task_stats, '$.FreshStartDate') as fresh_start_date,
  json_extract(replication_task_stats, '$.FullLoadFinishDate') as full_load_finish_date,
  json_extract(replication_task_stats, '$.FullLoadProgressPercent') as full_load_progress_percent,
  json_extract(replication_task_stats, '$.FullLoadStartDate') as full_load_start_date,
  json_extract(replication_task_stats, '$.StartDate') as start_date,
  json_extract(replication_task_stats, '$.StopDate') as stop_date,
  json_extract(replication_task_stats, '$.TablesErrored') as tables_errored,
  json_extract(replication_task_stats, '$.TablesLoaded') as tables_loaded,
  json_extract(replication_task_stats, '$.TablesLoading') as tables_loading,
  json_extract(replication_task_stats, '$.TablesQueued') as tables_queued
from
  aws_dms_replication_task;
```

### Get replication instance details
Retrieve replication instance details for the tasks.

```sql+postgresql
select
  t.replication_task_identifier,
  t.arn as task_arn,
  i.replication_instance_class,
  i.engine_version,
  i.publicly_accessible,
  i.dns_name_servers
from
  aws_dms_replication_task t
join aws_dms_replication_instance i on t.replication_instance_arn = i.arn;
```

```sql+sqlite
select
  t.replication_task_identifier,
  t.arn as task_arn,
  i.replication_instance_class,
  i.engine_version,
  i.publicly_accessible,
  i.dns_name_servers
from
  aws_dms_replication_task as t
join
  aws_dms_replication_instance as i on t.replication_instance_arn = i.arn;
```