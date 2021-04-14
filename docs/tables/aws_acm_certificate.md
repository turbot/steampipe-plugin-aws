# Table: aws_acm_certificate

AWS Certificate Manager (ACM) handles the complexity of creating, storing, and renewing public and private SSL/TLS X.509 certificates and keys that protect the AWS websites and applications.

## Examples

### Basic info

```sql
select
  certificate_arn,
  domain_name,
  failure_reason,
  in_use_by,
  status,
  key_algorithm
from
  aws_acm_certificate;
```


### List of expired certificates

```sql
select
  certificate_arn,
  domain_name,
  status
from
  aws_acm_certificate
where
  status = 'expired';
```


### List certificates for which transparency logging is disabled

```sql
select
  certificate_arn,
  domain_name,
  status
from
  aws_acm_certificate
where
  certificate_transparency_logging_preference <> 'ENABLED';
```


### List certificates without application tag key

```sql
select
  certificate_arn,
  turbot_tags
from
  aws_acm_certificate
where
  not turbot_tags :: JSONB ? 'application';
```
