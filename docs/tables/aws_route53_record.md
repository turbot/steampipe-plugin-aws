# Table: aws_route53_record

A Route 53 record contains authoritative DNS information for a specified DNS name. DNS records are most commonly used to map a name to an IP Address

You **_must_** specify a single `zone_id` in a where or join clause in order to use this table.

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
  aws_route53_record
where
  zone_id = 'Z09145482OD83AIAO253B';
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
  zone_id = 'Z09145482OD83AIAO253B'
  and name = 'test.com.';
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
  zone_id = 'Z09145482OD83AIAO253B'
  and type = 'NS';
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
  zone_id = 'Z09145482OD83AIAO253B'
  and name = 'test.com.'
  and type = 'NS';
```

### Count records by type

```sql
select
  type,
  count(*)
from
  aws_route53_record
where
  zone_id = 'Z09145482OD83AIAO253B'
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
  zone_id = 'Z09145482OD83AIAO253B'
  and geo_location is not null
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
where
  zone_id = 'Z09145482OD83AIAO253B'
group by
  name,
  type;
```

### List all records in all zones

```sql
select
  r.name,
  r.type,
  r.records,
  r.alias_target
from
  aws_route53_zone as z,
  aws_route53_record as r
where
  r.zone_id = z.id ;
```
