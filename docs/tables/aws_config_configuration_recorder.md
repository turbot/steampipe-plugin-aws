# Table: aws_config_configuration_recorder

AWS Config uses the configuration recorder to detect changes in your resource configurations and capture these changes as configuration items. You must create a configuration recorder before AWS Config can track your resource configurations.

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
