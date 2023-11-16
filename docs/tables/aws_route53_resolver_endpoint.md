---
title: "Table: aws_route53_resolver_endpoint - Query AWS Route 53 Resolver Endpoints using SQL"
description: "Allows users to query AWS Route 53 Resolver Endpoints, providing detailed information about each endpoint, including its ID, direction, IP addresses, and status, among other details."
---

# Table: aws_route53_resolver_endpoint - Query AWS Route 53 Resolver Endpoints using SQL

The `aws_route53_resolver_endpoint` table in Steampipe provides information about Resolver Endpoints within AWS Route 53. This table allows DevOps engineers to query endpoint-specific details, including the endpoint's direction (INBOUND or OUTBOUND), IP addresses, status, and associated metadata. Users can utilize this table to gather insights on endpoints, such as their security status, IP address associations, and more. The schema outlines the various attributes of the Resolver Endpoint, including the endpoint ID, ARN, direction, IP address count, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_route53_resolver_endpoint` table, you can use the `.inspect aws_route53_resolver_endpoint` command in Steampipe.

**Key columns**:

- `id`: The ID of the resolver endpoint. This can be used to join this table with others that reference resolver endpoints by their ID.
- `arn`: The Amazon Resource Name (ARN) of the resolver endpoint. This can be used to join this table with others that reference resolver endpoints by their ARN.
- `direction`: The direction of the resolver endpoint (INBOUND or OUTBOUND). This can be used to filter or group resolver endpoints based on their direction.

## Examples

### List all endpoints

```sql
select
  name,
  id,
  direction,
  ip_address_count
  status
from
  aws_route53_resolver_endpoint;
```

### Get a specific endpoint

```sql
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

### List unhealthy endpoints

```sql
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

```sql
select
  name,
  p ->> 'Ip' as ip,
  p ->> 'Status' as status,
  p ->> 'SubnetId' as subnet_id
from
  aws_route53_resolver_endpoint,
  jsonb_array_elements(ip_addresses) as p;
```
