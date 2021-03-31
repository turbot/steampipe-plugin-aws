# Table: aws_api_gateway_stage

A stage is a named reference to a deployment, which is a snapshot of the API.

## Examples

### Count of api gateway stages per rest APIs

```sql
select
  rest_api_id,
  count(name) stage_count
from
  aws_api_gateway_stage
group by
  rest_api_id;
```


### List of API gateway stages where API caching is enabled

```sql
select
  name,
  rest_api_id,
  cache_cluster_enabled,
  cache_cluster_size
from
  aws_api_gateway_stage
where
  cache_cluster_enabled;
```


### List of web acls associated with the gateway stages

```sql
select
  name,
  split_part(web_acl_arn, '/', 3) as web_acl_name
from
  aws_api_gateway_stage;
```


### List stages where CloudWatch logging is disabled

```sql
select
  deployment_id,
  name,
  tracing_enabled,
  method_settings -> '*/*' ->> 'LoggingLevel' as cloudwatch_log_level
from
  aws_api_gateway_stage
where
  method_settings -> '*/*' ->> 'LoggingLevel' = 'OFF';
```
