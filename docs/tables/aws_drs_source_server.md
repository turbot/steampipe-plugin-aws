---
title: "Table: aws_drs_source_server - Query AWS Database Migration Service Source Server using SQL"
description: "Allows users to query AWS Database Migration Service Source Servers for detailed information about the replication servers used in database migrations."
---

# Table: aws_drs_source_server - Query AWS Database Migration Service Source Server using SQL

The `aws_drs_source_server` table in Steampipe provides information about source servers within AWS Database Migration Service (DMS). This table allows DevOps engineers to query server-specific details, including server type, replication job status, associated replication tasks, and more. Users can utilize this table to gather insights on source servers, such as server status, assigned replication tasks, and server configuration details. The schema outlines the various attributes of the source server, including the server ID, status, replication job details, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_drs_source_server` table, you can use the `.inspect aws_drs_source_server` command in Steampipe.

### Key columns:

- `server_id`: This is the unique identifier of the source server. This column can be used to join this table with other tables to get more specific details about a particular source server.
- `status`: This column provides the current status of the source server. It can be used to filter out servers based on their current status.
- `replication_job_id`: This column holds the ID of the replication job associated with the source server. It can be used to join this table with replication job tables to get detailed information about the replication tasks.

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

### Get launch configuration settings of all source servers

```sql
select
  arn,
  title,
  launch_configuration ->> 'Name' as launch_configuration_name,
  launch_configuration ->> 'CopyPrivateIp' as launch_configuration_copy_private_ip,
  launch_configuration ->> 'CopyTags' as launch_configuration_copy_tags,
  launch_configuration ->> 'Ec2LaunchTemplateID' as launch_configuration_ec2_launch_template_id,
  launch_configuration ->> 'LaunchDisposition' as launch_configuration_disposition,
  launch_configuration ->> 'TargetInstanceTypeRightSizingMethod' as launch_configuration_target_instance_type_right_sizing_method,
  launch_configuration -> 'Licensing' as launch_configuration_licensing,
  launch_configuration -> 'ResultMetadata' as launch_configuration_result_metadata
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
