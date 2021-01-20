# Table: aws_api_gateway_usage_plan

A usage plan specifies who can access one or more deployed API stages and methods and also how much and how fast they can access them. The plan uses API keys to identify API clients and meters access to the associated API stages for each key.

## Examples

### Basic info

```sql
select
  name,
  id,
  product_code,
  description,
  api_stages
from
  aws_api_gateway_usage_plan;
```


### List the API gateway usage plans where quota ( i.e the number of api call a user can make within a time period) is disabled

```sql
select
  name,
  id,
  quota
from
  aws_api_gateway_usage_plan
where
  quota is null;
```


### List the API gateway usage plan where throttle ( i.e the rate at which user can make request ) is disabled

```sql
select
  name,
  id,
  throttle
from
  aws_api_gateway_usage_plan
where
  throttle is null;
```
