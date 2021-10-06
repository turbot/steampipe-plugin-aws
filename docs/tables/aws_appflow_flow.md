# Table: aws_appflow_flow

A flow transfers data between a source and a destination, e.g., from Salesforce to an Amazon S3 bucket.

## Examples

### Basic info

```sql
select
  flow_name,
  source_flow_config,
  destination_flow_config_list
from
  aws_appflow_flow;
```

### List all flows with Redshift or S3 sources

```sql
select
  flow_name,
  source_flow_config,
  destination_flow_config_list
from
  aws_appflow_flow
where
  source_flow_config ->> 'ConnectorType' in
  (
    'Redshift',
    'S3'
  );
```

### Get task details for each flow

```sql
select
  flow_name,
  task ->> 'ConnectorOperator' as connector_operator,
  task ->> 'SourceFields' as source_fields,
  task ->> 'TaskType' as task_type,
  task ->> 'TaskProperties' as task_properties,
  task ->> 'DestinationField' as destination_field
from
  aws_appflow_flow,
  jsonb_array_elements(tasks) as task;
```
