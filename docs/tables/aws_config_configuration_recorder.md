# Table: aws_config_configuration_recorder

AWS Config uses the configuration recorder to detect changes in your resource configurations and capture these changes as configuration items. You must create a configuration recorder before AWS Config can track your resource configurations.

## Examples

### Basic AWS config configuration recorder info

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
	aws.aws_config_configuration_recorder;
```


### List of configuration recorders with recording is disabled

```sql
select
  name,
  role_arn,
  status_recording,
  title
from
  aws.aws_config_configuration_recorder
where
  status_recording = 'false';
```


### List of configuration recorders with status in pending or failure state

```sql
select
  name,
  status ->> 'LastStatus' as last_status
from
  aws.aws_config_configuration_recorder
where
  status ->> 'LastStatus' IN ('PENDING', 'FAILURE');
```


### List of resource types for which recording is enabled

```sql
select
  name,
  recording_group ->> 'ResourceTypes' as resource_types
from
  aws.aws_config_configuration_recorder;
```