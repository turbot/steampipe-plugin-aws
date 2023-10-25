# Table: aws_lambda_event_source_mapping

AWS Lambda is a compute service that lets you run code without provisioning or managing servers.
Event Source Mappings are triggers that route events from services like Kafka, Kinesis, SQS or SNS to the configured AWS Lambda.

## Examples

### Basic Info

```sql
select
  event_source_arn,
  function_arn,
  function_name,
  last_processing_result,
  parallelization_factor,
  state,
  destination_config
from
  aws_lambda_event_source_mapping;
```

### List of lambda event source mappings that are disabled

```sql
select
  function_name,
  state,
  last_modified,
  state_transition_reason
from
  aws_lambda_event_source_mapping
where
  state = 'Disabled';
```

### List Kafka Bootstrap Servers from a self-managed Kafka cluster feeding events to Lambda functions

```sql
select
  function_name, jsonb_array_elements_text(jsonb_extract_path(self_managed_event_source, 'Endpoints', 'KAFKA_BOOTSTRAP_SERVERS'))
from
  aws_lambda_event_source_mapping
where
  function_name = 'MyFunctionName'
```

### Get source access configuration of event source mappings

```sql
select
  uuid,
  event_source_arn,
  a ->> 'Type' as source_access_type,
  a ->> 'URL' as source_access_url
from
  aws_lambda_event_source_mapping,
  jsonb_array_elements(source_access_configurations) as a;
```

### Get scaling configuration details of event sourcce mappings

```sql
select
  uuid,
  event_source_arn,
  scaling_config ->> 'MaximumConcurrency' as maximum_concurrency
from
  aws_lambda_event_source_mapping;
```

### Get destionation configuration of event source mappings

```sql
select
  uuid,
  function_name,
  destination_config ->> 'OnFailure' as on_failure,
  destination_config ->> 'OnSuccess' as on_success
from
  aws_lambda_event_source_mapping;
```

### List event source maps by filter criteria

```sql
select
  uuid,
  event_source_arn,
  function_arn,
  state,
  filter ->> 'Pattern' as filter_criteria_pattern
from
  aws_lambda_event_source_mapping,
  jsonb_array_elements(filter_criteria -> 'Filters') as filter
where
  filter ->> 'Pattern' like '{ \"Metadata\" : [ 1, 2 ]}';
```