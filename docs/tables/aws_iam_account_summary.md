# Table: aws_iam_account_summary

Information about IAM entity usage and IAM quotas in the AWS account.

The number and size of IAM resources in an AWS account are limited. For more information, see [IAM and STS Quotas](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_iam-quotas.html) in the IAM User Guide.

## Examples

### List the IAM summary for the account 
```sql
select * from aws_iam_account_summary;
```

### CIS v1 > 1 Identity and Access Management > 1.13 Ensure MFA is enabled for the "root" account (Scored)
```sql
select
  account_mfa_enabled as cis_v1_1_13
from
  aws_iam_account_summary;
```