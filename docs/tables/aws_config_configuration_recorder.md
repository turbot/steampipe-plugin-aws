---
title: "Table: aws_config_configuration_recorder - Query AWS Config Configuration Recorder using SQL"
description: "Allows users to query AWS Config Configuration Recorder"
---

# Table: aws_config_configuration_recorder - Query AWS Config Configuration Recorder using SQL

The `aws_config_configuration_recorder` table in Steampipe provides information about Configuration Recorders within AWS Config. This table allows DevOps engineers, security analysts, and cloud administrators to query configuration recorder-specific details, including its current status, associated role ARN, and whether it is recording all resource types. Users can utilize this table to gather insights on configuration recorders, such as which resources are being recorded, the recording status, and more. The schema outlines the various attributes of the Configuration Recorder, including the name, role ARN, resource types, and recording group.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_config_configuration_recorder` table, you can use the `.inspect aws_config_configuration_recorder` command in Steampipe.

### Key columns:

- `name`: The name of the configuration recorder. It can be used to join with other tables that contain information about the configuration recorder.
- `role_arn`: The Amazon Resource Name (ARN) of the role used to make read or write requests to the delivery channel. This is useful for joining with tables that contain IAM role information.
- `recording_group_all_supported`: Indicates whether the configuration recorder is recording all resources. It can be used to filter or join with other tables based on the recording status.

## Examples

### Basic info

```sql
select
  name,
  role_arn,
  status,
  recording_group,
  status_recording,
  akas,
  title
from
  aws_config_configuration_recorder;
```


### List configuration recorders that are not recording

```sql
select
  name,
  role_arn,
  status_recording,
  title
from
  aws_config_configuration_recorder
where
  not status_recording;
```


### List configuration recorders with failed deliveries

```sql
select
  name,
  status ->> 'LastStatus' as last_status,
  status ->> 'LastStatusChangeTime' as last_status_change_time,
  status ->> 'LastErrorCode' as last_error_code,
  status ->> 'LastErrorMessage' as last_error_message
from
  aws_config_configuration_recorder
where
  status ->> 'LastStatus' = 'FAILURE';
```
