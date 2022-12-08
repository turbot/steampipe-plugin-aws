# Table: aws_drs_source_server

AWS Elastic Disaster Recovery (AWS DRS) minimizes downtime and data loss with fast, reliable recovery of on-premises and cloud-based applications using affordable storage, minimal compute, and point-in-time recovery.

## Examples

### Basic Info

```sql
select
  arn,
  last_launch_result,
  source_server_id,
  title
from
  aws_drs_source_server;
```

### Get source cloud properties of all source servers

```sql
select
  arn,
  title,
  source_cloud_properties ->> 'OriginAccountID' as source_cloud_origin_account_id,
  source_cloud_properties ->> 'OriginAvailabilityZone' as source_cloud_origin_availability_zone,
  source_cloud_properties ->> 'OriginRegion' as source_cloud_origin_region
from
  aws_drs_source_server;
```

### Get source properties of all source servers

```sql
select
  arn,
  title,
  source_properties ->> 'Cpus' as source_cpus,
  source_properties ->> 'Disks' as source_disks,
  source_properties -> 'IdentificationHints' ->> 'Hostname' as source_hostname,
  source_properties ->> 'NetworkInterfaces' as source_network_interfaces,
  source_properties -> 'Os' ->> 'FullString' as source_os,
  source_properties -> 'RamBytes' as source_ram_bytes,
  source_properties -> 'RecommendedInstanceType' as source_recommended_instance_type,
  source_properties -> 'LastUpdatedDateTime' as source_last_updated_date_time
from
  aws_drs_source_server;
```

### Get data replication info of all source servers

```sql
select
  arn,
  title,
  data_replication_info -> 'DataReplicationInitiation' ->> 'StartDateTime' as data_replication_start_date_time,
  data_replication_info -> 'DataReplicationInitiation' ->> 'NextAttemptDateTime' as data_replication_next_attempt_date_time,
  data_replication_info ->> 'DataReplicationError' as data_replication_error,
  data_replication_info ->> 'DataReplicationState' as data_replication_state,
  data_replication_info ->> 'ReplicatedDisks' as data_replication_replicated_disks
from
  aws_drs_source_server;
```

### List source servers that failed last recovery launch

```sql
select
  title,
  arn,
  last_launch_result,
  source_server_id
from
  aws_drs_source_server
where
  last_launch_result = 'FAILED';
```

### List disconnected source servers

```sql
select
  title,
  arn,
  data_replication_info ->> 'DataReplicationState' as data_replication_state,
  data_replication_info ->> 'DataReplicationError' as data_replication_error,
  data_replication_info -> 'DataReplicationInitiation' ->> 'StartDateTime' as data_replication_start_date_time,
  data_replication_info -> 'DataReplicationInitiation' ->> 'NextAttemptDateTime' as data_replication_next_attempt_date_time
from
  aws_drs_source_server
where
  data_replication_info ->> 'DataReplicationState' = 'DISCONNECTED';
```
