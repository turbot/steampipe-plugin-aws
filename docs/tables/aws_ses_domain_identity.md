# Table: aws_ses_domain_identity

Amazon SES is an email platform that provides an easy, cost-effective way for you to send and receive email using your own domain.

## Examples

### Basic info

```sql
select
  name,
  arn,
  region,
  akas
from
  aws_ses_domain_identity;
```

### List domain identities which failed verification

```sql
select
  name,
  region,
  verification_status
from
  aws_ses_domain_identity
where
  verification_status = 'Failed';
```
