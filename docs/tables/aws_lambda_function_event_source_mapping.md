# Table: aws_lambda_function_event_source_mapping

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
  aws_lambda_function_event_source_mapping;
```

### List of lambda functions which disabled

```sql
select
    function_name, 
    state, 
    last_modified, 
    state_transition_reason
from
    aws_lambda_function_event_source_mapping
where
    state = 'Disabled';
```

### Count of lambda function event source mappings by runtime engines of the function

```sql
select
  runtime,
  count(*)
from
    aws_lambda_function_event_source_mapping, 
    aws_lambda_function
where
    arn = function_arn
group by
  runtime;
```

### List Kafka Bootstrap Servers from a self-managed Kafka cluster feeding events to a Lambda function

```sql
select
    function_name, jsonb_array_elements_text(jsonb_extract_path(self_managed_event_source, 'Endpoints', 'KAFKA_BOOTSTRAP_SERVERS'))
from
    aws_lambda_function_event_source_mapping
where
    function_name = 'MyFunctionName'
```
