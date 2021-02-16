# Table: aws_lambda_function

AWS Lambda is a compute service that lets you run code without provisioning or managing servers.

## Examples

### List of lambda functions which are not encrypted with CMK

```sql
select
  name,
  kms_key_arn
from
  aws_lambda_function
where
  kms_key_arn is null;
```


### Count of lambda functions by runtime engines

```sql
select
  runtime,
  count(*)
from
  aws_lambda_function
group by
  runtime;
```


### List of lambda function whose retention period is less than 30 days

```sql
select
  fn.name,
  lg.name,
  lg.retention_in_days
from
  aws_lambda_function as fn
  inner join aws_cloudwatch_log_group as lg on (
    (lg.name = '/aws/lambda/')
    or (lg.name = fn.name)
  )
where
  lg.retention_in_days < 30;
```


### Availability zone count for each VPC lambda function

```sql
select
  fn.name,
  count (availability_zone) as zone_count
from
  aws_lambda_function as fn
  cross join jsonb_array_elements_text(vpc_subnet_ids) as vpc_subnet
  join aws_vpc_subnet as sub on sub.subnet_id = vpc_subnet
group by
  fn.name,
  sub.availability_zone;
```


### List all the actions allowed by managed policies for a Lambda execution role
```sql
select
  f.name,
  f.role,
  a.action,
  a.access_level,
  a.description
from 
  aws_lambda_function as f,
  aws_iam_role as r,
  jsonb_array_elements_text(r.attached_policy_arns) as pol_arn,
  aws_iam_policy as p,
  jsonb_array_elements(p.policy_std -> 'Statement') as stmt,
  jsonb_array_elements_text(stmt -> 'Action') as action_glob,
  glob(action_glob) as action_regex
  join aws_iam_action a ON a.action LIKE action_regex
where
  f.role = r.arn
  and pol_arn = p.arn 
  and stmt ->> 'Effect' = 'Allow'
  and f.name = 'hellopython';
```