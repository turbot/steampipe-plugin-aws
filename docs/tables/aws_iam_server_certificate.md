# Table: aws_iam_server_certificate

SSL certificates are a set of small data which binds a cryptographic key to an organizations details. This enables a secure connection between a webserver and a browser.

## Examples

### Basic info

```sql
select
  name,
  arn,
  server_certificate_id,
  upload_date,
  expiration
from
  aws_iam_server_certificate;
```

### List expired certificates

```sql
select
  name,
  arn,
  expiration
from
  aws_iam_server_certificate
where
  expiration < now()::timestamp;
```
