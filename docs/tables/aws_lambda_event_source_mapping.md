---
title: "Steampipe Table: aws_lambda_event_source_mapping - Query AWS Lambda Event Source Mappings using SQL"
description: "Allows users to query AWS Lambda Event Source Mappings, providing a detailed view of each mapping configuration."
folder: "Lambda"
---

# Table: aws_lambda_event_source_mapping - Query AWS Lambda Event Source Mappings using SQL

The AWS Lambda Event Source Mapping is a service that enables the connection of a Lambda function to specific AWS services as an event source, triggering the function when records are available. This service allows you to read batches of records from the event source and invoke your function synchronously with an event that contains stream records. This mapping service is integral for applications that need to respond to data modifications in Amazon DynamoDB tables, Amazon Kinesis streams, and Amazon Simple Queue Service queues.

## Table Usage Guide

The `aws_lambda_event_source_mapping` table in Steampipe provides you with information about event source mappings within AWS Lambda. This table allows you, as a DevOps engineer, to query mapping-specific details, including the source ARN, function ARN, batch size, and last processing result. You can utilize this table to gather insights on mappings, such as those with errors, the state of the mapping, the maximum record age, and more. The schema outlines for you the various attributes of the event source mapping, including the UUID, creation date, and associated tags.

## Examples

### Basic Info
Explore the status and configuration of AWS Lambda event source mappings to understand the efficiency of your serverless architecture. This can help identify any potential issues or bottlenecks in your application's event-driven processing.

```sql+postgres
select
  arn,
  function_arn,
  function_name,
  last_processing_result,
  parallelization_factor,
  state,
  destination_config
from
  aws_lambda_event_source_mapping;
```

```sql+sqlite
select
  arn,
  function_arn,
  function_name,
  last_processing_result,
  parallelization_factor,
  state,
  destination_config
from
  aws_lambda_event_source_mapping;
```

### List lambda event source mappings that are disabled
Identify instances where Lambda event source mappings are disabled to understand the areas in your AWS infrastructure that might not be processing events as expected.

```sql+postgres
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

```sql+sqlite
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

### Retrieve bootstrap server endpoints from a self-managed Kafka cluster integrated with AWS Lambda
Identify the bootstrap server endpoints within a self-managed Kafka cluster that's integrated with AWS Lambda. This can be useful for gaining insights into the event source mapping of your AWS Lambda functions.
-- Returns the list of bootstrap servers for your Kafka brokers

--  function_name | jsonb_array_elements_text
-- ---------------+---------------------------
--  myFunction    | abc.xyz.com:xxxx
--  myFunction    | abc2.xyz.com:xxxx

```sql+postgres
select
  function_name,
  jsonb_array_elements_text(jsonb_extract_path(self_managed_event_source, 'Endpoints', 'KAFKA_BOOTSTRAP_SERVERS'))
from
  aws_lambda_event_source_mapping;
```

```sql+sqlite
select
  function_name,
  json_extract(json_extract(self_managed_event_source, '$.Endpoints'), '$.KAFKA_BOOTSTRAP_SERVERS')
from
  aws_lambda_event_source_mapping;
```

### Get source access configuration of event source mappings
Discover the configurations of event source access in your AWS Lambda setup. This query is useful to understand the types and URLs of source access, which can help in managing and troubleshooting your Lambda functions.

```sql+postgres
select
  uuid,
  arn,
  a ->> 'Type' as source_access_type,
  a ->> 'URL' as source_access_url
from
  aws_lambda_event_source_mapping,
  jsonb_array_elements(source_access_configurations) as a;
```

```sql+sqlite
select
  uuid,
  arn,
  json_extract(a.value, '$.Type') as source_access_type,
  json_extract(a.value, '$.URL') as source_access_url
from
  aws_lambda_event_source_mapping,
  json_each(source_access_configurations) as a;
```

### Get scaling configuration details of event source mappings
Analyze the scaling configuration of event source mappings to understand the maximum concurrency level. This is useful in managing and optimizing the performance of AWS Lambda functions.

```sql+postgres
select
  uuid,
  arn,
  scaling_config ->> 'MaximumConcurrency' as maximum_concurrency
from
  aws_lambda_event_source_mapping;
```

```sql+sqlite
select
  uuid,
  arn,
  json_extract(scaling_config, '$.MaximumConcurrency') as maximum_concurrency
from
  aws_lambda_event_source_mapping;
```

### Get destionation configuration of event source mappings
Explore the success and failure configurations of event source mappings in AWS Lambda. This could be useful to understand how your system behaves under different outcomes of the event source mapping.

```sql+postgres
select
  uuid,
  function_name,
  destination_config ->> 'OnFailure' as on_failure,
  destination_config ->> 'OnSuccess' as on_success
from
  aws_lambda_event_source_mapping;
```

```sql+sqlite
select
  uuid,
  function_name,
  json_extract(destination_config, '$.OnFailure') as on_failure,
  json_extract(destination_config, '$.OnSuccess') as on_success
from
  aws_lambda_event_source_mapping;
```

### List AWS Lambda event source mappings with specific filter criteria patterns
Determine the areas in which specific filter criteria patterns are applied in AWS Lambda event source mappings. This can be useful to pinpoint specific locations where certain metadata patterns are used, allowing for more targeted management and optimization of serverless computing services.

```sql+postgres
select
  uuid,
  arn,
  function_arn,
  state,
  filter ->> 'Pattern' as filter_criteria_pattern
from
  aws_lambda_event_source_mapping,
  jsonb_array_elements(filter_criteria -> 'Filters') as filter
where
  filter ->> 'Pattern' like '{ \"Metadata\" : [ 1, 2 ]}';
```

```sql+sqlite
select
  uuid,
  arn,
  function_arn,
  state,
  json_extract(filter.value, '$.Pattern') as filter_criteria_pattern
from
  aws_lambda_event_source_mapping,
  json_each(filter_criteria, '$.Filters') as filter
where
  json_extract(filter.value, '$.Pattern') like '{ "Metadata" : [ 1, 2 ]}';
```

### Get lambda function details of each event source mapping
Explore the relationship between Lambda functions and their respective event source mappings in AWS to understand how they interact. This is particularly useful for auditing and optimizing your serverless architecture.

```sql+postgres
select
  m.arn,
  m.function_arn,
  f.runtime,
  f.handler,
  f.architectures
from
  aws_lambda_event_source_mapping as m,
  aws_lambda_function as f
where
  f.name = m.function_name;
```

```sql+sqlite
select
  m.arn,
  m.function_arn,
  f.runtime,
  f.handler,
  f.architectures
from
  aws_lambda_event_source_mapping as m
join aws_lambda_function as f
on f.name = m.function_name;
```