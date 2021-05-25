# Table: aws_route53_domain

Amazon Route 53 enables you to register and transfer domain names using your AWS account.

## Examples

### Basic info

```sql
select
  domain_name,
  auto_renew,
  expiration_date
from
  aws_route53_domain;
```


### List Route53 domains with auto-renewal enabled

```sql
select
  domain_name,
  auto_renew,
  expiration_date
from
  aws_route53_domain
where
  auto_renew;
```


### List domains with transfer lock enabled

```sql
select
  domain_name,
  expiration_date,
  transfer_lock
from
  aws_route53_domain
where
  transfer_lock;
```
