# Table: aws_lambda_version

Versions are used to manage the deployment of lambda functions.

## Examples

### Runtime info of each lambda version

```sql
select
  function_name,
  version,
  runtime,
  handler
from
  aws_lambda_version;
```

### List of lambda versions where code run timout is more than 2 mins

```sql
select
  function_name,
  version,
  timeout
from
  aws_lambda_version
where
  timeout :: int > 120;
```

### VPC info of each lambda version

```sql
select
  function_name,
  version,
  vpc_id,
  vpc_security_group_ids,
  vpc_subnet_ids
from
  aws_lambda_version;
```

### List policy details

```sql
select
  jsonb_pretty(policy) as policy,
  jsonb_pretty(policy_std) as policy_std
from
  aws_lambda_version;
```
