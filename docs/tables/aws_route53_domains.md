# Table: aws_route53_domains

Amazon Route 53 enables you to register and transfer domain names using your AWS account.

## Examples

### Basic info

```sql
select
  domain_name,
  auto_renew,
  expiry
from
  aws_route53_domains;
```

### List Route53 domain names are renewed before their expiration.

```sql
select
  domain_name,
  auto_renew,
  expiry
from
  aws_route53_domains
where
  auto_renew = 'true';
```

### List Route53 domain names are renewed before Expiry 7 Days

```sql
select
  domain_name,
  auto_renew,
  expiry
from
  aws_route53_domains
where
  expiry <= (expiry - interval '7' day);
```

### Check your domain names have the Transfer Lock feature enabled in order to keep them secure.

```sql
select
  domain_name,
  expiry,
  transfer_lock
from
  aws_route53_domains
where
  transfer_lock = 'true';
```
