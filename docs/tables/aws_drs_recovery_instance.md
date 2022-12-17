# Table: aws_drs_recovery_instance

AWS Elastic Disaster Recovery (AWS DRS) minimizes downtime and data loss with fast, reliable recovery of on-premises and cloud-based applications using affordable storage, minimal compute, and point-in-time recovery. If you need to recover applications, you can launch recovery instances on AWS within minutes, using the most up-to-date server state or a previous point in time.

## Examples

### Basic Info

```sql
select
  recovery_instance_id,
  arn,
  source_server_id,
  ec2_instance_id,
  ec2_instance_state
from
  aws_drs_recovery_instance;
```

### Get recovery instance properties of each recovery instance

```sql
select
  recovery_instance_id
  arn,
  recovery_instance_properties ->> 'Cpus' as recovery_instance_cpus,
  recovery_instance_properties ->> 'Disks' as recovery_instance_disks,
  recovery_instance_properties ->> 'IdentificationHints' as recovery_instance_identification_hints,
  recovery_instance_properties ->> 'LastUpdatedDateTime' as recovery_instance_last_updated_date_time,
  recovery_instance_properties ->> 'NetworkInterfaces' as recovery_instance_network_interfaces,
  recovery_instance_properties ->> 'Os' as recovery_instance_os,
  recovery_instance_properties ->> 'RamBytes' as recovery_instance_ram_bytes
from
  aws_drs_recovery_instance;
```

### Get failback details of each recovery instance

```sql
select
  recovery_instance_id,
  arn,
  source_server_id,
  ec2_instance_id,
  failback ->> 'AgentLastSeenByServiceDateTime' as agent_last_seen_by_service_date_time,
  failback ->> 'ElapsedReplicationDuration' as elapsed_replication_duration,
  failback ->> 'FailbackClientID' as failback_client_id,
  failback ->> 'FailbackClientLastSeenByServiceDateTime' as failback_client_last_seen_by_service_date_time,
  failback ->> 'FailbackInitiationTime' as failback_initiation_time,
  failback -> 'FailbackJobID' as failback_job_id,
  failback -> 'FailbackLaunchType' as failback_launch_type,
  failback -> 'FailbackToOriginalServer' as failback_to_original_server,
  failback -> 'FirstByteDateTime' as failback_first_byte_date_time,
  failback -> 'State' as failback_state
from
  aws_drs_recovery_instance;
```

### Get data replication info of each recovery instance

```sql
select
  recovery_instance_id,
  arn,
  data_replication_info -> 'DataReplicationInitiation' ->> 'StartDateTime' as data_replication_start_date_time,
  data_replication_info -> 'DataReplicationInitiation' ->> 'NextAttemptDateTime' as data_replication_next_attempt_date_time,
  data_replication_info ->> 'DataReplicationError' as data_replication_error,
  data_replication_info ->> 'DataReplicationState' as data_replication_state,
  data_replication_info ->> 'ReplicatedDisks' as data_replication_replicated_disks
from
  aws_drs_recovery_instance;
```

### List recovery instances which are created for an actual recovery event

```sql
select
  recovery_instance_id,
  arn,
  source_server_id,
  ec2_instance_id,
  ec2_instance_state,
  is_drill,
  job_id
from
  aws_drs_recovery_instance
where
  not is_drill;
```