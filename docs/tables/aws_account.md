# Table: aws_account

An AWS account is a container for your AWS resources. You create and manage your AWS resources in an AWS account, and the AWS account provides administrative capabilities for access and billing.

## Examples

### Basic AWS account info

```sql
select
  alias,
  organization_id,
  organization_master_account_email,
  organization_master_account_id
from
  aws_account
  cross join jsonb_array_elements(account_aliases) as alias;
```


### Organization policy of aws account

```sql
select
  organization_id,
  policy ->> 'Type' as policy_type,
  policy ->> 'Status' as policy_status
from
  aws_account
  cross join jsonb_array_elements(organization_available_policy_types) as policy;
```
