---
title: "Steampipe Table: aws_route53_domain - Query AWS Route 53 Domains using SQL"
description: "Allows users to query AWS Route 53 Domains for detailed information about domain names, including their status, expiration date, and associated tags."
folder: "Route 53"
---

# Table: aws_route53_domain - Query AWS Route 53 Domains using SQL

The AWS Route 53 Domain service is a scalable and highly available domain name system (DNS) web service. It is designed to give developers and businesses an extremely reliable and cost-effective way to route end users to Internet applications by translating human readable names like www.example.com into the numeric IP addresses like 192.0.2.1 that computers use to connect to each other. Route 53 effectively connects user requests to infrastructure running in AWS, such as Amazon EC2 instances, ELB load balancers, or Amazon S3 buckets, and can also be used to route users to infrastructure outside of AWS.

## Table Usage Guide

The `aws_route53_domain` table in Steampipe allows you to query detailed information about domain names registered with Route 53, a scalable Domain Name System (DNS) web service in AWS. This table provides you, as a DevOps engineer, with domain-specific details such as domain name, status, expiration date, and associated tags. You can utilize this table to retrieve details about domain names, including their registration, renewal, and transfer status, and also to verify the associated tags. The schema outlines various attributes of the domain, including the domain name, auto renew, transfer lock, expiry date, and associated tags.

## Examples

### Basic info
Determine the status of domain renewals in your AWS Route53 service to anticipate and manage upcoming expiration dates. This query is useful for maintaining domain continuity and avoiding unexpected service disruptions.

```sql+postgres
select
  domain_name,
  auto_renew,
  expiration_date
from
  aws_route53_domain;
```

```sql+sqlite
select
  domain_name,
  auto_renew,
  expiration_date
from
  aws_route53_domain;
```


### List domains with auto-renewal enabled
Discover the domains that have the auto-renewal feature enabled to understand the areas that may need attention for timely renewals, thus preventing any potential service disruptions. This can be particularly useful in managing domain registrations and ensuring continuity of web services.

```sql+postgres
select
  domain_name,
  auto_renew,
  expiration_date
from
  aws_route53_domain
where
  auto_renew;
```

```sql+sqlite
select
  domain_name,
  auto_renew,
  expiration_date
from
  aws_route53_domain
where
  auto_renew = 1;
```


### List domains with transfer lock enabled
Explore which domains have the transfer lock feature enabled, which is useful for maintaining domain security by preventing unauthorized transfers. This can be particularly beneficial for organizations looking to safeguard their domains from potential cyber threats or unauthorized changes.

```sql+postgres
select
  domain_name,
  expiration_date,
  transfer_lock
from
  aws_route53_domain
where
  transfer_lock;
```

```sql+sqlite
select
  domain_name,
  expiration_date,
  transfer_lock
from
  aws_route53_domain
where
  transfer_lock = 1;
```