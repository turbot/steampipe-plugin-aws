# Table: aws_sts_caller_identity

The details about the AWS IAM user or role whose credentials are used to call the operation.

## Examples

### Basic info

```sql
select
  arn,
  user_id,
  title,
  account_id,
  akas
from
  aws_sts_caller_identity;
```