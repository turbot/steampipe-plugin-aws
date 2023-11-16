---
title: "Table: aws_route53_query_log - Query AWS Route 53 Query Log using SQL"
description: "Allows users to query AWS Route 53 Query Log data, providing insights into DNS queries made to Route 53 hosted zones."
---

# Table: aws_route53_query_log - Query AWS Route 53 Query Log using SQL

The `aws_route53_query_log` table in Steampipe provides information about DNS queries made to Route 53 hosted zones within AWS Route 53. This table allows network administrators and DevOps engineers to query DNS query-specific details, including the hosted zone, query name, query type, and response code. Users can utilize this table to gather insights on DNS query patterns, troubleshoot DNS issues, and analyze DNS traffic. The schema outlines the various attributes of the Route 53 query log, including the query timestamp, query name, query type, query class, and response code.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_query_log` table, you can use the `.inspect aws_route53_query_log` command in Steampipe.

**Key columns**:

- `hosted_zone_id`: The ID of the Amazon Route 53 hosted zone. This key column can be used to join this table with other tables that contain information about Route 53 hosted zones.
- `query_name`: The domain name that was specified in the DNS query. This key column can be used to filter the DNS queries made for a specific domain.
- `query_type`: The type of DNS query, such as A, AAAA, CNAME, MX, NS, PTR, SOA, SPF, SRV, TXT. This key column can be used to analyze the types of DNS queries made to the hosted zone.

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