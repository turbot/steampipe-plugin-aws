---
title: "Steampipe Table: aws_dms_endpoint - Query AWS DMS Endpoints using SQL"
description: "Query AWS DMS Endpoints to retrieve connection information for source or target databases in database migration activities."
folder: "DMS"
---

# Table: aws_dms_endpoint - Query AWS DMS Endpoints using SQL

AWS Database Migration Service (DMS) Endpoints are a pivotal component within AWS DMS, delineating the connection details for source or target databases involved in migration tasks. These endpoints are essential for defining the data's origin (source) and destination (target).

## Table Usage Guide

The `aws_dms_endpoint` table in Steampipe allows you to query connection-specific information, such as the endpoint identifier, ARN, database name, endpoint type, and the database engine details. This table is invaluable for DevOps engineers and database administrators overseeing database migrations, as it facilitates the monitoring and management of endpoint configurations and ensures the smooth execution of migration tasks.

## Examples

### Basic info
Retrieve basic information about AWS DMS Endpoint, including their identifiers, ARNs, certificate, database, endpoint type, engine name, and regions.

```sql+postgres
select
  endpoint_identifier,
  arn,
  certificate_arn,
  database_name,
  endpoint_type,
  engine_display_name,
  engine_name
from
  aws_dms_endpoint;
```

```sql+sqlite
select
  endpoint_identifier,
  arn,
  certificate_arn,
  database_name,
  endpoint_type,
  engine_display_name,
  engine_name
from
  aws_dms_endpoint;
```

### List source endpoints
Identify all source endpoints in AWS DMS, showcasing their identifiers, ARNs, display names, types, and engine names.

```sql+postgres
select
  endpoint_identifier,
  arn,
  engine_display_name,
  endpoint_type,
  engine_name
from
  aws_dms_endpoint
where
  endpoint_type = 'SOURCE';
```

```sql+sqlite
select
  endpoint_identifier,
  arn,
  engine_display_name,
  endpoint_type,
  engine_name
from
  aws_dms_endpoint
where
  endpoint_type = 'SOURCE';
```

### List MySQL endpoints
Retrieve a comprehensive list of AWS DMS endpoints configured for MySQL databases, including their identifiers, ARNs, engine names, creation times, and MySQL-specific settings."

```sql+postgres
select
  endpoint_identifier,
  arn,
  engine_name,
  my_sql_settings
from
  aws_dms_endpoint
where
  engine_name = 'mysql';
```

```sql+sqlite
select
  endpoint_identifier,
  arn,
  engine_name,
  my_sql_settings
from
  aws_dms_endpoint
where
  engine_name = 'mysql';
```

### List endpoints that have SSL enabled
Display all AWS DMS endpoints with SSL encryption enabled, detailing their identifiers, KMS key IDs, server names, service access role ARNs, and SSL modes."

```sql+postgres
select
  endpoint_identifier,
  kms_key_id,
  server_name,
  service_access_role_arn,
  ssl_mode
from
  aws_dms_endpoint
where
  ssl_mode <> 'none';
```

```sql+sqlite
select
  endpoint_identifier,
  kms_key_id,
  server_name,
  service_access_role_arn,
  ssl_mode
from
  aws_dms_endpoint
where
  ssl_mode <> 'none';
```

### Get MySQL setting details for MySQL endpoints
Extract detailed MySQL settings for AWS DMS endpoints configured for MySQL, including connection scripts, metadata settings, database names, and other MySQL-specific configurations.

```sql+postgres
select
  endpoint_identifier,
  arn,
  my_sql_settings ->> 'AfterConnectScript' as after_connect_script,
  (my_sql_settings ->> 'CleanSourceMetadataOnMismatch')::boolean as clean_source_metadata_on_mismatch,
  my_sql_settings ->> 'DatabaseName' as database_name,
  (my_sql_settings ->> 'EventsPollInterval')::integer as events_poll_interval,
  (my_sql_settings ->> 'ExecuteTimeout')::integer as execute_timeout,
  (my_sql_settings ->> 'MaxFileSize')::integer as max_file_size,
  (my_sql_settings ->> 'ParallelLoadThreads')::integer as parallel_load_threads,
  my_sql_settings ->> 'Password' as password,
  (my_sql_settings ->> 'Port')::integer as port,
  my_sql_settings ->> 'SecretsManagerAccessRoleArn' as secrets_manager_access_role_arn,
  my_sql_settings ->> 'SecretsManagerSecretId' as secrets_manager_secret_id,
  my_sql_settings ->> 'ServerName' as server_name,
  my_sql_settings ->> 'ServerTimezone' as server_timezone,
  my_sql_settings ->> 'TargetDbType' as target_db_type,
  my_sql_settings ->> 'Username' as username
from
  aws_dms_endpoint
where
  engine_name = 'mysql';
```

```sql+sqlite
select
  endpoint_identifier,
  arn,
  my_sql_settings ->> 'AfterConnectScript' as after_connect_script,
  cast(json_extract(my_sql_settings, '$.CleanSourceMetadataOnMismatch') as boolean) as clean_source_metadata_on_mismatch,
  my_sql_settings ->> 'DatabaseName' as database_name,
  cast(json_extract(my_sql_settings, '$.EventsPollInterval') as integer) as events_poll_interval,
  cast(json_extract(my_sql_settings, '$.ExecuteTimeout') as integer) as execute_timeout,
  cast(json_extract(my_sql_settings, '$.MaxFileSize') as integer) as max_file_size,
  cast(json_extract(my_sql_settings, '$.ParallelLoadThreads') as integer) as parallel_load_threads,
  my_sql_settings ->> 'Password' as password,
  cast(json_extract(my_sql_settings, '$.Port') as integer) as port,
  my_sql_settings ->> 'SecretsManagerAccessRoleArn' as secrets_manager_access_role_arn,
  my_sql_settings ->> 'SecretsManagerSecretId' as secrets_manager_secret_id,
  my_sql_settings ->> 'ServerName' as server_name,
  my_sql_settings ->> 'ServerTimezone' as server_timezone,
  my_sql_settings ->> 'TargetDbType' as target_db_type,
  my_sql_settings ->> 'Username' as username
from
  aws_dms_endpoint
where
  engine_name = 'mysql';
```
