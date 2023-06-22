# Table: aws_route53_query_log

AWS Route 53 query logging is a feature that allows you to capture detailed information about DNS queries that are made to your Route 53 hosted zones. It provides visibility into the DNS queries and responses that occur for your domain names, enabling you to monitor and analyze DNS traffic.

## Examples

### Basic info

```sql
select
  id,
  hosted_zone_id,
  cloud_watch_logs_log_group_arn,
  title,
  akas
from
  aws_route53_query_log;
```

### Get hosted zone details of each query log

```sql
select
  l.id,
  l.hosted_zone_id,
  z.private_zone,
  z.resource_record_set_count
from
  aws_route53_query_log as l,
  aws_route53_zone as z
where
  z.id = l.hosted_zone_id;
```

### Count the number of query logs by hosted zone

```sql
select
  hosted_zone_id,
  count(id)
from
  aws_route53_query_log
group by
  hosted_zone_id;
```