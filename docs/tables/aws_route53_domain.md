---
title: "Table: aws_route53_domain - Query AWS Route 53 Domains using SQL"
description: "Allows users to query AWS Route 53 Domains for detailed information about domain names, including their status, expiration date, and associated tags."
---

# Table: aws_route53_domain - Query AWS Route 53 Domains using SQL

The `aws_route53_domain` table in Steampipe allows users to query detailed information about domain names registered with Route 53, a scalable Domain Name System (DNS) web service in AWS. This table provides DevOps engineers with domain-specific details such as domain name, status, expiration date, and associated tags. Users can utilize this table to retrieve details about domain names, including their registration, renewal, and transfer status, and also to verify the associated tags. The schema outlines various attributes of the domain, including the domain name, auto renew, transfer lock, expiry date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_domain` table, you can use the `.inspect aws_route53_domain` command in Steampipe.

**Key columns**:

- `name`: The name of the domain. This is the primary key of the table and can be used to join this table with other tables that contain domain-specific information.
- `expiry`: The expiration date of the domain. This is useful for tracking when domains are due for renewal.
- `tags`: The metadata tags associated with the domain. These can be used to filter and categorize domains based on user-defined criteria.

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


### List domains with auto-renewal enabled

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
