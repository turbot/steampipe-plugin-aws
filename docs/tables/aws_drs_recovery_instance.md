# Table: aws_drs_recovery_instance

Once your source servers have been added to Elastic Disaster Recovery (DRS), you can monitor and interact with them from the Source Servers page. The Source Servers page is the default view in the Elastic Disaster Recovery Console, and will be the page that you interact with the most. On the Source Servers page, you can view all of your source servers, monitor their recovery readiness and data replication state, see the last recovery result, see any pending actions, and sort your servers by a variety of categories. You can also perform a variety of commands from the Source Servers page through the command menus. These menus allow you to fully control your servers by launching Drill and Recovery instances and performing a variety of actions, such as adding servers, editing settings, disconnecting, and deleting servers.

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

### Get recovery instance properties for each recovery instance

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

### Get failback details for each recovery instance

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

### Get data replication info for each recovery instance

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
where not is_drill;
```