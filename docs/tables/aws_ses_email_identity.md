# Table: aws_ses_email_identity

Amazon SES is an email platform that provides an easy, cost-effective way for you to send and receive email using your own email addresses.

## Examples

### Basic info

```sql
select
  identity,
  arn,
  region,
  akas
from
  aws_ses_email_identity;
```

### List email identities which failed verification

```sql
select
  identity,
  region,
  verification_status
from
  aws_ses_email_identity
where
  verification_status = 'Failed';
```
