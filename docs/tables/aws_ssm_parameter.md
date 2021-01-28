# Table: aws_ssm_parameter

AWS Systems Manager Parameter Store (SSM) provides a secure way to store config variables for applications. SSM can store plaintext parameters or KMS encrypted secure strings.

## Examples

### SSM parameter basic info

```sql
select
  name,
  type,
  data_type,
  tier,
  region
from
  aws_ssm_parameter;
```


### Policy details of advanced tier ssm parameter

```sql
select
  name,
  tier,
  p ->> 'PolicyType' as policy_type,
  p ->> 'PolicyStatus' as Policy_status,
  p ->> 'PolicyText' as policy_text
from
  aws_ssm_parameter,
  jsonb_array_elements(policies) as p;
```


### List of SSM parameters which do not have owner or app_id tag key

```sql
select
  name
from
  aws_ssm_parameter
where
  tags -> 'owner' is null
  or tags -> 'app_id' is null;
```