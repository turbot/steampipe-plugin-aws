---
title: "Table: aws_drs_recovery_instance - Query AWS Disaster Recovery Service Recovery Instances using SQL"
description: "Allows users to query AWS Disaster Recovery Service Recovery Instances to retrieve information about recovery instances, including instance type, recovery instance ARN, and associated tags."
---

# Table: aws_drs_recovery_instance - Query AWS Disaster Recovery Service Recovery Instances using SQL

The `aws_drs_recovery_instance` table in Steampipe provides information about recovery instances within AWS Disaster Recovery Service (DRS). This table allows DevOps engineers to query recovery instance-specific details, including instance type, recovery instance ARN, and associated tags. Users can utilize this table to gather insights on recovery instances, such as instance type, recovery instance ARN, and associated tags. The schema outlines the various attributes of the recovery instance, including the instance type, recovery instance ARN, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_drs_recovery_instance` table, you can use the `.inspect aws_drs_recovery_instance` command in Steampipe.

### Key columns:

- `instance_type`: This is the type of the recovery instance. It is an important column as it helps in understanding the capacity and capability of the recovery instance.
- `recovery_instance_arn`: This is the Amazon Resource Name (ARN) of the recovery instance. It is a unique identifier for the recovery instance and can be used to join this table with other tables that contain recovery instance ARN.
- `tags`: These are the metadata tags assigned to the recovery instance. They can be used for filtering and organizing recovery instances within AWS DRS.

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

### List recovery instances that are created for an actual recovery event

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
