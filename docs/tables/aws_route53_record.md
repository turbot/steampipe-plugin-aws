---
title: "Steampipe Table: aws_route53_record - Query AWS Route 53 Record using SQL"
description: "Allows users to query Route 53 DNS records within Amazon Web Services. The `aws_route53_record` table in Steampipe provides information about DNS records within AWS Route 53. This table allows DevOps engineers to query record-specific details, including type, name, TTL, and associated metadata. Users can utilize this table to gather insights on DNS records, such as record types, verification of TTL values, and more."
folder: "Route 53"
---

# Table: aws_route53_record - Query AWS Route 53 Record using SQL

The AWS Route 53 Record is a component of Amazon's scalable and highly available Domain Name System (DNS) web service. It is designed to provide reliable and cost-effective domain registration, DNS routing, and health checking of resources within your environment. It translates domain names into the numeric IP addresses that computers use to connect to each other, thus facilitating the routing of internet traffic to your applications and services.

## Table Usage Guide

The `aws_route53_record` table in Steampipe provides you with information about DNS records within AWS Route 53. This table allows you, as a DevOps engineer, to query record-specific details, including type, name, TTL (Time to Live), and associated metadata. You can utilize this table to gather insights on DNS records, such as record types, verification of TTL values, and more. The schema outlines the various attributes of the DNS record for you, including the record name, type, set identifier, TTL, and associated resource records.

**Important Notes**
- We recommend specifying the `name` and `type` columns when querying zones with a large number of records to reduce the query time.

## Examples

### Basic info
Explore which type of records are associated with your Route 53 DNS entries to gain insights into your AWS environment's DNS configuration. This can help in identifying potential misconfigurations or understanding the distribution of different record types.

```sql+postgres
select
  name,
  type,
  records,
  alias_target
from
  aws_route53_record;
```

```sql+sqlite
select
  name,
  type,
  records,
  alias_target
from
  aws_route53_record;
```

### List all test.com records in a zone
Determine the areas in which specific 'test.com' records exist within a zone, enabling better management and organization of your domain records.

```sql+postgres
select
  r.name,
  r.type,
  record
from
  aws_route53_record as r,
  jsonb_array_elements_text(records) as record
where
  name = 'test.com.';
```

```sql+sqlite
select
  r.name,
  r.type,
  json_extract(record.value, '$') as record
from
  aws_route53_record as r,
  json_each(records) as record
where
  name = 'test.com.';
```

### List all NS records in a zone
Identify instances where you need to analyze all the NS records within a specific zone. This can be particularly useful when managing DNS configurations and ensuring accurate routing of internet traffic.

```sql+postgres
select
  r.name,
  r.type,
  record
from
  aws_route53_record as r,
  jsonb_array_elements_text(records) as record
where
  r.type = 'NS';
```

```sql+sqlite
select
  r.name,
  r.type,
  json_extract(record.value, '$') as record
from
  aws_route53_record as r,
  json_each(r.records) as record
where
  r.type = 'NS';
```

### Get test.com NS record in a zone
Determine the specific Name Server (NS) record associated with 'test.com' in a DNS zone. This is useful for verifying correct DNS configuration or troubleshooting DNS issues.

```sql+postgres
select
  r.name,
  r.type,
  record
from
  aws_route53_record as r,
  jsonb_array_elements_text(records) as record
where
  r.name = 'test.com.'
  and r.type = 'NS';
```

```sql+sqlite
select
  r.name,
  r.type,
  json_extract(record.value, '$') as record
from
  aws_route53_record as r,
  json_each(r.records) as record
where
  r.name = 'test.com.'
  and r.type = 'NS';
```

### Count records by type
Analyze the distribution of different record types in your AWS Route53 configuration to understand which types are most commonly used. This information can be useful for optimizing DNS setup and identifying potential areas for improvement.

```sql+postgres
select
  type,
  count(*)
from
  aws_route53_record
group by
  type
order by
  count desc;
```

```sql+sqlite
select
  type,
  count(*)
from
  aws_route53_record
group by
  type
order by
  count(*) desc;
```

### List geo-location routing information
Explore geo-location routing information to gain insights into the distribution of your web traffic. This can help optimize your network strategy by identifying which continents and countries are accessing your resources the most.

```sql+postgres
select
  name,
  type,
  records,
  alias_target,
  geo_location ->> 'ContinentCode' as continent,
  geo_location ->> 'CountryCode' as country,
  geo_location ->> 'SubdivisionCode' as subdivision
from
  aws_route53_record
where
  geo_location is not null
order by
  name;
```

```sql+sqlite
select
  name,
  type,
  records,
  alias_target,
  json_extract(geo_location, '$.ContinentCode') as continent,
  json_extract(geo_location, '$.CountryCode') as country,
  json_extract(geo_location, '$.SubdivisionCode') as subdivision
from
  aws_route53_record
where
  geo_location is not null
order by
  name;
```

### Count of records by name and type
Determine the frequency of different record types within your AWS Route53 service. This can help in understanding the distribution and usage patterns of various record types, aiding in effective DNS management.

```sql+postgres
select
  name,
  type,
  count(*)
from
  aws_route53_record
  left join jsonb_array_elements_text(records) as record on true
group by
  name,
  type;
```

```sql+sqlite
select
  name,
  type,
  count(*)
from
  aws_route53_record
  left join (
    select
      value as record
    from
      aws_route53_record,
      json_each(records)
  )
group by
  name,
  type;
```