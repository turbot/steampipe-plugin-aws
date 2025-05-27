---
title: "Steampipe Table: aws_drs_recovery_instance - Query AWS Disaster Recovery Service Recovery Instances using SQL"
description: "Allows users to query AWS Disaster Recovery Service Recovery Instances to retrieve information about recovery instances, including instance type, recovery instance ARN, and associated tags."
folder: "Elastic Disaster Recovery (DRS)"
---

# Table: aws_drs_recovery_instance - Query AWS Disaster Recovery Service Recovery Instances using SQL

The AWS Disaster Recovery Service Recovery Instance is a component of the AWS Disaster Recovery Service, which aids in the recovery of applications and data in the event of a disaster. It allows for the rapid recovery of your IT infrastructure and data by utilizing AWS's robust, scalable, and secure global infrastructure. This service supports recovery scenarios ranging from small customer workload data loss to a complete site outage.

## Table Usage Guide

The `aws_drs_recovery_instance` table in Steampipe provides you with information about recovery instances within AWS Disaster Recovery Service (DRS). This table allows you, as a DevOps engineer, to query recovery instance-specific details, including instance type, recovery instance ARN, and associated tags. You can utilize this table to gather insights on recovery instances, such as instance type, recovery instance ARN, and associated tags. The schema outlines the various attributes of the recovery instance for you, including the instance type, recovery instance ARN, and associated tags.

## Examples

### Basic Info
Uncover the details of AWS Disaster Recovery Service's recovery instances, such as their current state and associated EC2 instances. This can be useful for maintaining an overview of your disaster recovery setup and ensuring everything is functioning as expected.

```sql+postgres
select
  recovery_instance_id,
  arn,
  source_server_id,
  ec2_instance_id,
  ec2_instance_state
from
  aws_drs_recovery_instance;
```

```sql+sqlite
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
Explore the characteristics of each recovery instance, such as CPU usage, disk activity, identification hints, update time, network interfaces, operating system, and RAM usage. This can help in assessing the performance and resource usage of each instance, aiding in efficient resource management and troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  recovery_instance_id,
  arn,
  json_extract(recovery_instance_properties, '$.Cpus') as recovery_instance_cpus,
  json_extract(recovery_instance_properties, '$.Disks') as recovery_instance_disks,
  json_extract(recovery_instance_properties, '$.IdentificationHints') as recovery_instance_identification_hints,
  json_extract(recovery_instance_properties, '$.LastUpdatedDateTime') as recovery_instance_last_updated_date_time,
  json_extract(recovery_instance_properties, '$.NetworkInterfaces') as recovery_instance_network_interfaces,
  json_extract(recovery_instance_properties, '$.Os') as recovery_instance_os,
  json_extract(recovery_instance_properties, '$.RamBytes') as recovery_instance_ram_bytes
from
  aws_drs_recovery_instance;
```

### Get failback details of each recovery instance
Determine the status and details of each recovery instance's failback process in your AWS Disaster Recovery Service. This allows you to understand the progress and potential issues in your data recovery efforts.

```sql+postgres
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

```sql+sqlite
select
  recovery_instance_id,
  arn,
  source_server_id,
  ec2_instance_id,
  json_extract(failback, '$.AgentLastSeenByServiceDateTime') as agent_last_seen_by_service_date_time,
  json_extract(failback, '$.ElapsedReplicationDuration') as elapsed_replication_duration,
  json_extract(failback, '$.FailbackClientID') as failback_client_id,
  json_extract(failback, '$.FailbackClientLastSeenByServiceDateTime') as failback_client_last_seen_by_service_date_time,
  json_extract(failback, '$.FailbackInitiationTime') as failback_initiation_time,
  json_extract(failback, '$.FailbackJobID') as failback_job_id,
  json_extract(failback, '$.FailbackLaunchType') as failback_launch_type,
  json_extract(failback, '$.FailbackToOriginalServer') as failback_to_original_server,
  json_extract(failback, '$.FirstByteDateTime') as failback_first_byte_date_time,
  json_extract(failback, '$.State') as failback_state
from
  aws_drs_recovery_instance;
```

### Get data replication info of each recovery instance
Determine the areas in which data replication is occurring within each recovery instance. This can help assess the status and health of your data recovery operations, and identify any potential issues or delays in the replication process.

```sql+postgres
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

```sql+sqlite
select
  recovery_instance_id,
  arn,
  json_extract(data_replication_info, '$.DataReplicationInitiation.StartDateTime') as data_replication_start_date_time,
  json_extract(data_replication_info, '$.DataReplicationInitiation.NextAttemptDateTime') as data_replication_next_attempt_date_time,
  json_extract(data_replication_info, '$.DataReplicationError') as data_replication_error,
  json_extract(data_replication_info, '$.DataReplicationState') as data_replication_state,
  json_extract(data_replication_info, '$.ReplicatedDisks') as data_replication_replicated_disks
from
  aws_drs_recovery_instance;
```

### List recovery instances that are created for an actual recovery event
Determine the instances created for actual recovery events, allowing you to focus on real-time disaster recovery efforts rather than drills or tests. This is beneficial in tracking and managing genuine recovery instances for efficient resource allocation and response.

```sql+postgres
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

```sql+sqlite
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
  is_drill = 0;
```