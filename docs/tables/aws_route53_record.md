# Table: aws_route53_record

A Route 53 record contains authoritative DNS information for a specified DNS name.  DNS records are most commonly used to map a name to an IP Address

Note that you ***must*** specify a single `zone_id` in a where clause in order to use this table.  Also, there is a [known issue](https://github.com/turbot/steampipe-postgres-fdw/issues/3) with nested select queries (select where in (select ...)) and joins on tables with required key columns (see below).

## Examples

### List records in a zone
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



## NOTE: Issue with nested select queries and joins on tables with required key columns
Currently, there is a [known issue](https://github.com/turbot/steampipe-postgres-fdw/issues/3) with nested select queries (select where in (select ...)) and joins on tables with required key columns. It seems that the qualifiers are not passed to the parent query because the nested query is executed in parallel. We are actively working to resolve this issue.

For example, this works as you would expect:

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


This SHOULD work but currently doesn't:
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
```
Error: pq: cannot iterate: there was an error executing scanIterator: rpc error: code = Internal desc = 'List' call requires an '=' qual for column: zone_id
```

This SHOULD also work but currently doesn't:
```sql
select 
  name,
  type,
  records,
  alias_target
from 
  aws_route53_record
where 
  zone_id in (select id from aws_route53_zone);
```
```
Error: pq: rpc error: code = Internal desc = 'List' call requires an '=' qual for column: zone_id
```