# Table: aws_lambda_alias

A Lambda alias is like a pointer to a specific function version.

## Examples

### Lambda alias basic info

```sql
select
  name,
  function_name,
  function_version,
  authorizer_credentials
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
