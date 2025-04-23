---
title: "Steampipe Table: aws_route53_resolver_endpoint - Query AWS Route 53 Resolver Endpoints using SQL"
description: "Allows users to query AWS Route 53 Resolver Endpoints, providing detailed information about each endpoint, including its ID, direction, IP addresses, and status, among other details."
folder: "Route 53"
---

# Table: aws_route53_resolver_endpoint - Query AWS Route 53 Resolver Endpoints using SQL

The AWS Route 53 Resolver Endpoint is a component of Amazon's Route 53 service, which provides highly scalable and reliable domain name system (DNS) web services. The Resolver Endpoint specifically enables recursive DNS for your Amazon VPCs and your on-premises networks over a Direct Connect or VPN connection. It offers DNS resolution between virtual networks, improved response times, and manageability of DNS data.

## Table Usage Guide

The `aws_route53_resolver_endpoint` table in Steampipe provides you with information about Resolver Endpoints within AWS Route 53. This table allows you, as a DevOps engineer, to query endpoint-specific details, including the endpoint's direction (INBOUND or OUTBOUND), IP addresses, status, and associated metadata. You can utilize this table to gather insights on endpoints, such as their security status, IP address associations, and more. The schema outlines the various attributes of the Resolver Endpoint for you, including the endpoint ID, ARN, direction, IP address count, and associated tags.

## Examples

### List all endpoints
Explore the various endpoints in your AWS Route53 Resolver to assess their status and direction, which aids in managing your network traffic effectively.

```sql+postgres
select
  name,
  id,
  direction,
  ip_address_count
  status
from
  aws_route53_resolver_endpoint;
```

```sql+sqlite
select
  name,
  id,
  direction,
  ip_address_count,
  status
from
  aws_route53_resolver_endpoint;
```

### Get a specific endpoint
Determine the details of a specific network traffic flow direction, the number of IP addresses, and the current status within your Amazon Route 53 Resolver. This can be particularly useful to troubleshoot or optimize your DNS resolution strategy.

```sql+postgres
select
  name,
  id,
  direction,
  ip_address_count
  status
from
  aws_route53_resolver_endpoint
where
  id = 'rslvr-out-ebb7db0b7498463eb';
```

```sql+sqlite
select
  name,
  id,
  direction,
  ip_address_count,
  status
from
  aws_route53_resolver_endpoint
where
  id = 'rslvr-out-ebb7db0b7498463eb';
```

### List unhealthy endpoints
Determine the areas in which your AWS Route53 Resolver endpoints require action. This query helps in identifying endpoints that are experiencing issues, enabling you to address them promptly for a smoother network operation.

```sql+postgres
select
  name,
  id,
  direction,
  status,
  status_message
from
  aws_route53_resolver_endpoint
where
  status = 'ACTION_NEEDED';
```

```sql+sqlite
select
  name,
  id,
  direction,
  status,
  status_message
from
  aws_route53_resolver_endpoint
where
  status = 'ACTION_NEEDED';
```

### Get IP address details for each endpoint
Determine the status and location of each endpoint in your network by analyzing their IP addresses. This can help in network troubleshooting and optimization.

```sql+postgres
select
  name,
  p ->> 'Ip' as ip,
  p ->> 'Status' as status,
  p ->> 'SubnetId' as subnet_id
from
  aws_route53_resolver_endpoint,
  jsonb_array_elements(ip_addresses) as p;
```

```sql+sqlite
select
  name,
  json_extract(p.value, '$.Ip') as ip,
  json_extract(p.value, '$.Status') as status,
  json_extract(p.value, '$.SubnetId') as subnet_id
from
  aws_route53_resolver_endpoint,
  json_each(ip_addresses) as p;
```