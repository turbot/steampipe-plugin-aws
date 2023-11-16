---
title: "Table: aws_route53_record - Query AWS Route 53 Record using SQL"
description: "Allows users to query Route 53 DNS records within Amazon Web Services. The `aws_route53_record` table in Steampipe provides information about DNS records within AWS Route 53. This table allows DevOps engineers to query record-specific details, including type, name, TTL, and associated metadata. Users can utilize this table to gather insights on DNS records, such as record types, verification of TTL values, and more."
---

# Table: aws_route53_record - Query AWS Route 53 Record using SQL

The `aws_route53_record` table in Steampipe provides information about DNS records within AWS Route 53. This table allows DevOps engineers to query record-specific details, including type, name, TTL (Time to Live), and associated metadata. Users can utilize this table to gather insights on DNS records, such as record types, verification of TTL values, and more. The schema outlines the various attributes of the DNS record, including the record name, type, set identifier, TTL, and associated resource records.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_record` table, you can use the `.inspect aws_route53_record` command in Steampipe.

**Key columns**:

- `name`: The name of the DNS record. This column is useful for joining this table with others that contain DNS record names.
- `type`: The type of the DNS record (e.g., A, AAAA, CNAME, MX, NS, PTR, SOA, SPF, SRV, TXT). This column is important as it allows users to filter or join with other tables based on the DNS record type.
- `ttl`: The Time to Live (TTL) of the DNS record. This column is important as it allows users to filter or join with other tables based on the TTL value.

## Examples

### Basic info

```sql
select
  name,
  type,
  records,
  alias_target
from
  aws_route53_record;
```

### List all test.com records in a zone

```sql
select
  name,
  type,
  record
from
  aws_route53_record,
  jsonb_array_elements_text(records) as record
where
  name = 'test.com.';
```

### List all NS records in a zone

```sql
select
  name,
  type,
  record
from
  aws_route53_record,
  jsonb_array_elements_text(records) as record
where
  type = 'NS';
```

### Get test.com NS record in a zone

```sql
select
  name,
  type,
  record
from
  aws_route53_record,
  jsonb_array_elements_text(records) as record
where
  name = 'test.com.'
  and type = 'NS';
```

### Count records by type

```sql
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

### List geo-location routing information

```sql
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

### Count of records by name and type

```sql
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
