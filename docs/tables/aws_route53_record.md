# Table: aws_route53_record

A Route 53 record contains authoritative DNS information for a specified DNS name. DNS records are most commonly used to map a name to an IP Address

We recommend specifying the `name` and `type` columns when querying zones with a large number of records to reduce the query time.

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
