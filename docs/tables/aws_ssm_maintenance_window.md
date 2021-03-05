# Table: aws_ssm_maintenance_window

AWS Systems Manager Maintenance Windows let you define a schedule for when to perform potentially disruptive actions on your instances such as patching an operating system, updating drivers, or installing software or patches.

## Examples

### SSM maintenance window basic info

```sql
select
	name,
	window_id,
	enabled,
	schedule,
	tags_src,
	region
from
	aws_ssm_maintenance_window;
```


### Target details of maintenance window

```sql
select
	name,
	p ->> 'WindowTargetId' as window_target_id,
	p ->> 'ResourceType' as resource_type,
	p ->> 'Name' as target_name
from
	aws_ssm_maintenance_window,
	jsonb_array_elements(targets) as p;
```


### Tasks details of maintenance window

```sql
select
  name,
  p ->> 'WindowTaskId' as window_task_id,
  p ->> 'ServiceRoleArn' as service_role_arn,
  p ->> 'Name' as task_name
from
  aws_ssm_maintenance_window,
  jsonb_array_elements(tasks) as p;
```
### Whether the maintenance window is enable.

```sql
select
  name,
  window_id
from
  aws_ssm_maintenance_window
where
  enabled