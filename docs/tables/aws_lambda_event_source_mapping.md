---
title: "Table: aws_lambda_event_source_mapping - Query AWS Lambda Event Source Mappings using SQL"
description: "Allows users to query AWS Lambda Event Source Mappings, providing a detailed view of each mapping configuration."
---

# Table: aws_lambda_event_source_mapping - Query AWS Lambda Event Source Mappings using SQL

The `aws_lambda_event_source_mapping` table in Steampipe provides information about event source mappings within AWS Lambda. This table allows DevOps engineers to query mapping-specific details, including the source ARN, function ARN, batch size, and last processing result. Users can utilize this table to gather insights on mappings, such as those with errors, the state of the mapping, the maximum record age, and more. The schema outlines the various attributes of the event source mapping, including the UUID, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_lambda_event_source_mapping` table, you can use the `.inspect aws_lambda_event_source_mapping` command in Steampipe.

Key columns:

- `uuid`: This is the unique identifier of the event source mapping. It can be used to join this table with other tables that contain event source mapping information.
- `function_arn`: This is the ARN of the Lambda function. It can be used to join this table with other tables that contain Lambda function information.
- `event_source_arn`: This is the ARN of the event source. It can be used to join this table with other tables that contain event source information.

## Examples

### Basic Info

```sql
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

### Retrieve bootstrap server endpoints from a self-managed Kafka cluster integrated with AWS Lambda

-- Returns the list of bootstrap servers for your Kafka brokers

--  function_name | jsonb_array_elements_text
-- ---------------+---------------------------
--  myFunction    | abc.xyz.com:xxxx
--  myFunction    | abc2.xyz.com:xxxx
```sql
select
  function_name,
  jsonb_array_elements_text(jsonb_extract_path(self_managed_event_source, 'Endpoints', 'KAFKA_BOOTSTRAP_SERVERS'))
from
  aws_lambda_event_source_mapping;
```

### Get source access configuration of event source mappings

```sql
select
  uuid,
  arn,
  a ->> 'Type' as source_access_type,
  a ->> 'URL' as source_access_url
from
  aws_lambda_event_source_mapping,
  jsonb_array_elements(source_access_configurations) as a;
```

### Get scaling configuration details of event source mappings

```sql
select
  uuid,
  arn,
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

### List AWS Lambda event source mappings with specific filter criteria patterns

```sql
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

### Get lambda function details of each event source mapping

```sql
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