# Table: aws_dms_endpoint

AWS Database Migration Service (DMS) Endpoint refers to a specific component within the AWS DMS that defines the connection information for a source or a target database in a database migration activity. Endpoints in AWS DMS are crucial for specifying where your data is coming from (source) and where it is going to (target).

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  endpoint_identifier,
  arn,
  engine_name,
  instance_create_time,
  my_sql_settings
from
  aws_dms_endpoint
where
  engine_name = 'mysql';
```

### List endpoints that has SSL enabled

```sql
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

```sql
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
