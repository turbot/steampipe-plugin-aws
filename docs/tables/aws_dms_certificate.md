# Table: aws_dms_certificate

AWS DMS (Database Migration Service) Certificate refers to an SSL/TLS certificate used in AWS DMS for encrypting data during the process of migrating databases. This certificate plays a crucial role in ensuring the security and integrity of the data as it is transferred between the source and target databases in a migration task.

## Examples

### Basic info

```sql
select
  certificate_identifier,
  arn,
  certificate_creation_date,
  signing_algorithm,
  valid_to_date,
  region
from
  aws_dms_certificate;
```

### List certificates that are expires in next 10 days

```sql
select
  certificate_identifier,
  arn,
  key_length,
  signing_algorithm,
  valid_to_date
from
  aws_dms_certificate
where
  valid_to_date <= now() - interval '10' day;
```

### List certificates with SHA256 signin algorithm

```sql
select
  certificate_identifier,
  arn,
  signing_algorithm,
  key_length,
  certificate_owner
from
  aws_dms_certificate
where
  signing_algorithm = 'SHA256withRSA';
```

### List List certificates that are not owned by the current account

```sql
select
  certificate_identifier,
  arn,
  certificate_owner,
  account_id
from
  aws_dms_certificate
where
  certificate_owner <> account_id;
```

### Get the number of days left for expire the certificates

```sql
select
  certificate_identifier,
  arn,
  certificate_owner,
  (valid_to_date - current_date) as days_left,
  region
from
  aws_dms_certificate;
```
