---
title: "Steampipe Table: aws_config_configuration_recorder - Query AWS Config Configuration Recorder using SQL"
description: "Allows users to query AWS Config Configuration Recorder"
folder: "Config"
---

# Table: aws_config_configuration_recorder - Query AWS Config Configuration Recorder using SQL

The AWS Config Configuration Recorder is a feature that enables you to record the resource configurations in your AWS account. It captures and tracks changes to the configuration of your AWS resources, allowing you to assess, audit, and evaluate the configurations of your AWS resources. This helps ensure that your resource configurations are in compliance with your organization's policies and best practices.

## Table Usage Guide

The `aws_config_configuration_recorder` table in Steampipe provides you with information about Configuration Recorders within AWS Config. This table allows you, as a DevOps engineer, security analyst, or cloud administrator, to query configuration recorder-specific details, including its current status, associated role ARN, and whether it is recording all resource types. You can utilize this table to gather insights on configuration recorders, such as which resources are being recorded, the recording status, and more. The schema outlines the various attributes of the Configuration Recorder for you, including the name, role ARN, resource types, and recording group.

## Examples

### Basic info
Explore which AWS configuration recorders are active and recording, to better understand and manage your AWS resources and their configurations. This can be particularly useful for auditing, compliance, and operational troubleshooting purposes.

```sql+postgres
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

```sql+sqlite
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
Discover segments of configuration recorders that are currently inactive. This is beneficial in identifying potential gaps in your AWS Config setup, ensuring all necessary configuration changes are being tracked.

```sql+postgres
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

```sql+sqlite
select
  name,
  role_arn,
  status_recording,
  title
from
  aws_config_configuration_recorder
where
  status_recording != 1;
```


### List configuration recorders with failed deliveries
Discover the segments that have experienced delivery failures in AWS Configuration Recorder. This is beneficial for identifying and resolving issues in the system to ensure smooth operations.

```sql+postgres
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

```sql+sqlite
select
  name,
  json_extract(status, '$.LastStatus') as last_status,
  json_extract(status, '$.LastStatusChangeTime') as last_status_change_time,
  json_extract(status, '$.LastErrorCode') as last_error_code,
  json_extract(status, '$.LastErrorMessage') as last_error_message
from
  aws_config_configuration_recorder
where
  json_extract(status, '$.LastStatus') = 'FAILURE';
```