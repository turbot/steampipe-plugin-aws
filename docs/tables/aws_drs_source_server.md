---
title: "Steampipe Table: aws_drs_source_server - Query AWS Database Migration Service Source Server using SQL"
description: "Allows users to query AWS Database Migration Service Source Servers for detailed information about the replication servers used in database migrations."
folder: "Elastic Disaster Recovery (DRS)"
---

# Table: aws_drs_source_server - Query AWS Database Migration Service Source Server using SQL

The AWS Database Migration Service (DMS) Source Server is a component of AWS DMS that facilitates the migration of databases to AWS in a secure and efficient manner. It supports homogeneous migrations such as Oracle to Oracle, as well as heterogeneous migrations between different database platforms, such as Oracle to Amazon Aurora. The Source Server is the database instance from which the migration or replication tasks are initiated.

## Table Usage Guide

The `aws_drs_source_server` table in Steampipe provides you with information about source servers within AWS Database Migration Service (DMS). This table allows you, as a DevOps engineer, to query server-specific details, including server type, replication job status, associated replication tasks, and more. You can utilize this table to gather insights on source servers, such as server status, assigned replication tasks, and server configuration details. The schema outlines the various attributes of the source server for you, including the server ID, status, replication job details, and associated tags.

## Examples

### Basic Info
Explore the status of your last server launch and identify the source server details to understand the origin of your data. This can help in tracing data lineage or diagnosing issues related to specific server launches.

```sql+postgres
select
  arn,
  last_launch_result,
  source_server_id,
  title
from
  aws_drs_source_server;
```

```sql+sqlite
select
  arn,
  last_launch_result,
  source_server_id,
  title
from
  aws_drs_source_server;
```

### Get source cloud properties of all source servers
Explore the origin details of all source servers on your cloud platform. This information can be useful to understand the geographical distribution and account ownership of your servers, aiding in resource allocation and risk management strategies.

```sql+postgres
select
  arn,
  title,
  source_cloud_properties ->> 'OriginAccountID' as source_cloud_origin_account_id,
  source_cloud_properties ->> 'OriginAvailabilityZone' as source_cloud_origin_availability_zone,
  source_cloud_properties ->> 'OriginRegion' as source_cloud_origin_region
from
  aws_drs_source_server;
```

```sql+sqlite
select
  arn,
  title,
  json_extract(source_cloud_properties, '$.OriginAccountID') as source_cloud_origin_account_id,
  json_extract(source_cloud_properties, '$.OriginAvailabilityZone') as source_cloud_origin_availability_zone,
  json_extract(source_cloud_properties, '$.OriginRegion') as source_cloud_origin_region
from
  aws_drs_source_server;
```

### Get source properties of all source servers
This query helps you gain insights into the properties of all source servers, including CPU, disk details, network interfaces, RAM, and the recommended instance type. It's useful for understanding the capabilities of each server and for making informed decisions on server management and resource allocation.

```sql+postgres
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

```sql+sqlite
select
  arn,
  title,
  json_extract(source_properties, '$.Cpus') as source_cpus,
  json_extract(source_properties, '$.Disks') as source_disks,
  json_extract(source_properties, '$.IdentificationHints.Hostname') as source_hostname,
  json_extract(source_properties, '$.NetworkInterfaces') as source_network_interfaces,
  json_extract(source_properties, '$.Os.FullString') as source_os,
  json_extract(source_properties, '$.RamBytes') as source_ram_bytes,
  json_extract(source_properties, '$.RecommendedInstanceType') as source_recommended_instance_type,
  json_extract(source_properties, '$.LastUpdatedDateTime') as source_last_updated_date_time
from
  aws_drs_source_server;
```

### Get data replication info of all source servers
Explore the status of data replication across all source servers, identifying any errors and assessing when the next replication attempt will occur. This information can be crucial in ensuring data integrity and timely updates across your network.

```sql+postgres
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

```sql+sqlite
select
  arn,
  title,
  json_extract(data_replication_info, '$.DataReplicationInitiation.StartDateTime') as data_replication_start_date_time,
  json_extract(data_replication_info, '$.DataReplicationInitiation.NextAttemptDateTime') as data_replication_next_attempt_date_time,
  json_extract(data_replication_info, '$.DataReplicationError') as data_replication_error,
  json_extract(data_replication_info, '$.DataReplicationState') as data_replication_state,
  json_extract(data_replication_info, '$.ReplicatedDisks') as data_replication_replicated_disks
from
  aws_drs_source_server;
```

### Get launch configuration settings of all source servers
Explore the launch configuration settings of all source servers to understand their setup and configuration. This can help in assessing the current state of servers for auditing, troubleshooting, or planning purposes.

```sql+postgres
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

```sql+sqlite
select
  arn,
  title,
  json_extract(launch_configuration, '$.Name') as launch_configuration_name,
  json_extract(launch_configuration, '$.CopyPrivateIp') as launch_configuration_copy_private_ip,
  json_extract(launch_configuration, '$.CopyTags') as launch_configuration_copy_tags,
  json_extract(launch_configuration, '$.Ec2LaunchTemplateID') as launch_configuration_ec2_launch_template_id,
  json_extract(launch_configuration, '$.LaunchDisposition') as launch_configuration_disposition,
  json_extract(launch_configuration, '$.TargetInstanceTypeRightSizingMethod') as launch_configuration_target_instance_type_right_sizing_method,
  json_extract(launch_configuration, '$.Licensing') as launch_configuration_licensing,
  json_extract(launch_configuration, '$.ResultMetadata') as launch_configuration_result_metadata
from
  aws_drs_source_server;
```

### List source servers that failed last recovery launch
Identify instances where the last recovery launch of source servers was unsuccessful, which is crucial for troubleshooting and ensuring the robustness of the disaster recovery system.

```sql+postgres
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

```sql+sqlite
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
Identify instances where source servers have become disconnected. This is useful for troubleshooting and maintaining data integrity across your AWS infrastructure.

```sql+postgres
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

```sql+sqlite
select
  title,
  arn,
  json_extract(data_replication_info, '$.DataReplicationState') as data_replication_state,
  json_extract(data_replication_info, '$.DataReplicationError') as data_replication_error,
  json_extract(data_replication_info, '$.DataReplicationInitiation.StartDateTime') as data_replication_start_date_time,
  json_extract(data_replication_info, '$.DataReplicationInitiation.NextAttemptDateTime') as data_replication_next_attempt_date_time
from
  aws_drs_source_server
where
  json_extract(data_replication_info, '$.DataReplicationState') = 'DISCONNECTED';
```