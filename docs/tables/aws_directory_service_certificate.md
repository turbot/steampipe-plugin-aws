# Table: aws_directory_service_certificate

AWS Directory Service is a managed service provided by Amazon Web Services (AWS) that allows you to connect and integrate your AWS resources with an existing on-premises Microsoft Active Directory (AD) or establish a new, standalone directory in the AWS Cloud.

When setting up AWS Directory Service, you have the option to use Secure Sockets Layer (SSL) certificates to secure the communication between your on-premises directory and the AWS Cloud. This is especially important if you have a hybrid environment where you need to establish a secure connection between your on-premises AD and the AWS Cloud.

The AWS Directory Service certificate refers to the SSL certificate that is used to secure the communication between your on-premises AD and AWS. This certificate is typically issued by a trusted certificate authority (CA) and ensures that the communication is encrypted and secure.

## Examples

### Basic Info

```sql
select
  directory_id,
  certificate_id,
  common_name,
  type,
  state,
  expiry_date_time
from
  aws_directory_service_certificate;
```

### List 'MicrosoftAD' type directories

```sql
select
  c.certificate_id,
  c.common_name,
  c.directory_id,
  c.type as certificate_type,
  d.name as directory_name,
  d.type as directory_type
from
  aws_directory_service_certificate c,
  aws_directory_service_directory d
where
  d.type = 'MicrosoftAD';
```

### List deregistered certificates

```sql
select
  common_name,
  directory_id,
  type,
  state
from
  aws_directory_service_certificate
where
  state = 'Deregistered';
```

### List certificates that will expire in the coming 7 days

```sql
select
  directory_id,
  certificate_id,
  common_name,
  type,
  state,
  expiry_date_time
from
  aws_directory_service_certificate
where
  expiry_date_time >= now() + interval '7' day;
```

### Get client certificate auth settings of each certificate

```sql
select
  directory_id,
  certificate_id,
  common_name,
  client_cert_auth_settings -> 'OCSPUrl' as ocsp_url
from
  aws_directory_service_certificate;
```

### Retrieve the number of certificates registered in each directory

```sql
select
  directory_id,
  count(*) as certificate_count
from
  aws_directory_service_certificate
group by
  directory_id;
```

### List all certificates that were registered more than a year ago and have not been deregistered

```sql
select
  common_name,
  directory_id,
  type,
  state
from
  aws_directory_service_certificate
where
  registered_date_time <= now() - interval '1 year'
  and state not like 'Deregister%';
```

### Find the certificate with the latest registration date in each AWS partition

```sql
select
  distinct partition,
  registered_date_time
from
  aws_directory_service_certificate
order by
  partition,
  registered_date_time desc;
```
