---
title: "Steampipe Table: aws_route53_query_log - Query AWS Route 53 Query Log using SQL"
description: "Allows users to query AWS Route 53 Query Log data, providing insights into DNS queries made to Route 53 hosted zones."
folder: "Route 53"
---

# Table: aws_route53_query_log - Query AWS Route 53 Query Log using SQL

The AWS Route 53 Query Log is a feature of Amazon Route 53 that lets you log the DNS queries that Route 53 receives. It provides detailed records of the DNS queries that Amazon Route 53 receives and includes information like the domain or subdomain that was requested, the date and time of the request, and the DNS record type. This service is useful for troubleshooting and auditing purposes.

## Table Usage Guide

The `aws_route53_query_log` table in Steampipe provides you with information about DNS queries made to Route 53 hosted zones within AWS Route 53. This table allows you, as a network administrator or DevOps engineer, to query DNS query-specific details, including the hosted zone, query name, query type, and response code. You can utilize this table to gather insights on DNS query patterns, troubleshoot DNS issues, and analyze DNS traffic. The schema outlines the various attributes of the Route 53 query log for you, including the query timestamp, query name, query type, query class, and response code.

## Examples

### Basic info
Determine the areas in which AWS Route53 query logs are being used. This can provide insights into the distribution and usage of specific hosted zones and associated log groups, assisting in resource management and security monitoring.

```sql+postgres
select
  id,
  hosted_zone_id,
  cloud_watch_logs_log_group_arn,
  title,
  akas
from
  aws_route53_query_log;
```

```sql+sqlite
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
Gain insights into the characteristics of each query log's hosted zone, such as whether it's private and its record set count. This is useful for understanding the properties and usage of different hosted zones in your AWS Route53 service.

```sql+postgres
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

```sql+sqlite
select
  l.id,
  l.hosted_zone_id,
  z.private_zone,
  z.resource_record_set_count
from
  aws_route53_query_log as l
join
  aws_route53_zone as z
on
  z.id = l.hosted_zone_id;
```

### Count the number of query logs by hosted zone
Explore the distribution of query logs across different hosted zones to understand the areas of high activity and potential issues. This can aid in identifying zones that may require additional resources or troubleshooting.

```sql+postgres
select
  hosted_zone_id,
  count(id)
from
  aws_route53_query_log
group by
  hosted_zone_id;
```

```sql+sqlite
select
  hosted_zone_id,
  count(id)
from
  aws_route53_query_log
group by
  hosted_zone_id;
```