# Table: aws_route53_resolver_query_log_config

AWS Route53 Resolver Log Config is a configuration that enables logging for DNS queries and responses in Route53 Resolver. Resolver logs provide detailed information about DNS traffic flowing through the Resolver service, including the source IP, query domain, query type, response code, and more. The log data can be collected and analyzed for monitoring, troubleshooting, and security purposes. Resolver Log Config in AWS Route53 allows you to define the log destination and log format, and associate it with a VPC or a Resolver rule to start capturing the DNS query and response logs.

## Examples

### Basic info

```sql
select
  name,
  id,
  arn,
  creation_time,
  share_status,
  status
from
  aws_route53_resolver_query_log_config;
```

### List shared log configs

```sql
select
  name,
  id,
  arn,
  creation_time,
  share_status,
  status,
  destination_arn
from
  aws_route53_resolver_query_log_config
where
  owner_id <> account_id;
```

### List failed log configurations

```sql
select
  name,
  id,
  creator_request_id,
  destination_arn
from
  aws_route53_resolver_query_log_config
where
  status = 'FAILED';
```

### List log configurations shared with another account

```sql
select
  name,
  id,
  share_status,
  association_count
from
  aws_route53_resolver_query_log_config
where
  share_status = 'SHARED';
```

### List log configs created in last 30 days

```sql
select
  name,
  id,
  creation_time,
  destination_arn,
  status
from
  aws_route53_resolver_query_log_config
where
  creation_time >= now() - interval '30' day;
```
