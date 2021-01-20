# Table: aws_api_gateway_api_key

API keys are alphanumeric string values that are shared to give access to APIs.

## Examples

### API gateway API key basic info

```sql
select
  name,
  id,
  enabled,
  created_date,
  last_updated_date,
  customer_id,
  stage_keys
from
  aws_api_gateway_api_key;
```


### List of API keys which are not enabled

```sql
select
  name,
  id,
  customer_id
from
  aws_api_gateway_api_key
where
  not enabled;
```
