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
	akas,
	title
from
	aws_new.aws_config_configuration_recorder;
```

### List of configuration recorder with recording is enabled

```sql
select
  name,
  role_arn,
  status_recording,
  title
from
  aws_new.aws_config_configuration_recorder
where
  status_recording = 'false';
```

### Last status information of all the configuration recorders

```sql
select
  name,
  status ->> 'LastStatus' as last_status
from
  aws_new.aws_config_configuration_recorder
where
  status ->> 'LastStatus' IN ('pending', 'Failure');
```

### List of resource types for which recording is on

```sql
select
  name,
  recording_group ->> 'ResourceTypes' as resource_types
from
  aws_new.aws_config_configuration_recorder
```

