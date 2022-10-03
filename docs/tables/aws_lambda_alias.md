# Table: aws_lambda_alias

A Lambda alias is like a pointer to a specific function version.

## Examples

### Lambda alias basic info

```sql
select
  name,
  function_name,
  function_version
from
  aws_lambda_alias;
```

### Count of lambda alias per Lambda function

```sql
select
  function_name,
  count(function_name) count
from
  aws_lambda_alias
group by
  function_name;
```

### List policy details

```sql
select
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_lambda_alias;
```

### List URL configuration details for each alias

```sql
select
  name,
  function_name,
  jsonb_pretty(url_config) as url_config
from
  aws_lambda_alias;
```