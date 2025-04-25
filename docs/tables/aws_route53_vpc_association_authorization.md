---
title: "Steampipe Table: aws_route53_vpc_association_authorization - Query AWS Route53 for other-account VPCs."
description: "Gets a list of the VPCs that were created by other accounts and that can be associated with a specified hosted zone because you've submitted one or more `CreateVPCAssociationAuthorization` requests."
folder: "Route 53"
---

# Table: aws_route53_vpc_association_authorization - Query AWS Route53 for other-account VPCs using SQL

Amazon Route 53 is a highly available and scalable Domain Name System (DNS) web service.

## Table Usage Guide

The `aws_route53_vpc_association_authorization` table in Steampipe provides you with information VPCs in other AWS accounts that are authorized to be associated with a specified `hosted_zone_id`.

## Examples

### Basic info
Check which other-account VPCs are authorized

```sql+postgres
select
  hosted_zone_id,
  vpc_id,
  vpc_region
from
  aws_route53_vpc_association_authorization
where
  hosted_zone_id = 'Z3M3LMPEXAMPLE';
```

```sql+sqlite
select
  hosted_zone_id,
  vpc_id,
  vpc_region
from
  aws_route53_vpc_association_authorization
where
  hosted_zone_id = 'Z3M3LMPEXAMPLE';
```

### Sort VPCs descending by region name

```sql+postgres
select
  hosted_zone_id,
  vpc_id,
  vpc_region
from
  aws_route53_vpc_association_authorization
where
  hosted_zone_id = 'Z3M3LMPEXAMPLE'
order by
  vpc_region desc;
```

```sql+sqlite
select
  hosted_zone_id,
  vpc_id,
  vpc_region
from
  aws_route53_vpc_association_authorization
where
  hosted_zone_id = 'Z3M3LMPEXAMPLE'
order by
  vpc_region desc;
```

### Retrieve VPC Association Authorizations for available Hosted Zones

You can combine multiple tables to query or get fields such as the zone domain name (something the AWS API does not provide by default).

```sql+postgres
select
  auth.hosted_zone_id,
  z.name,
  auth.vpc_id,
  auth.vpc_region
from
  aws_route53_vpc_association_authorization auth
inner join
  aws_route53_zone z on auth.hosted_zone_id = z.id
where z.name = 'mycooldomain.xyz';
```

```sql+sqlite
select
  auth.hosted_zone_id,
  z.name,
  auth.vpc_id,
  auth.vpc_region
from
  aws_route53_vpc_association_authorization auth
inner join
  aws_route53_zone z on auth.hosted_zone_id = z.id
where z.name = 'mycooldomain.xyz';
```
