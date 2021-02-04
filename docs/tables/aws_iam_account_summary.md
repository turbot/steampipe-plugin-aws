# Table: aws_iam_account_summary

Information about IAM entity usage and IAM quotas in the AWS account.

The number and size of IAM resources in an AWS account are limited. For more information, see [IAM and STS Quotas](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_iam-quotas.html) in the IAM User Guide.

## Examples

### List the IAM summary for the account 
```sql
select
  *
from
  aws_iam_account_summary;
```

### Ensure MFA is enabled for the "root" account (CIS v1.1.13)
```sql
select
  account_mfa_enabled
from
  aws_iam_account_summary;
```




### Summary report - Total number of IAM resources in the account by type
```sql
select
  users,
  groups,
  roles,
  policies
from
  aws_iam_account_summary;
```