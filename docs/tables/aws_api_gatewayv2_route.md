---
title: "Table: aws_api_gatewayv2_route - Query AWS API Gateway V2 Routes using SQL"
description: "Allows users to query AWS API Gateway V2 Routes and obtain detailed information about each route, including the route key, route response selection expression, and target."
---

# Table: aws_api_gatewayv2_route - Query AWS API Gateway V2 Routes using SQL

The `aws_api_gatewayv2_route` table in Steampipe provides information about routes within AWS API Gateway V2. This table allows DevOps engineers to query route-specific details, including the route key, route response selection expression, and target. Users can utilize this table to gather insights on routes, such as route configurations, route response behaviors, and more. The schema outlines the various attributes of the route, including the API identifier, route ID, route key, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gatewayv2_route` table, you can use the `.inspect aws_api_gatewayv2_route` command in Steampipe.

### Key columns:

- `api_id`: The API identifier. This column is important as it uniquely identifies the API and can be used to join with other tables related to the API.
- `route_id`: The route ID. This column is useful because it uniquely identifies the route within the API and can be used to join with other tables related to the route.
- `route_key`: The route key. This column is important because it defines the method and route path for the route and can be used to join with other tables related to the route.

## Examples

### Basic info

```sql
select
  route_key,
  api_id,
  route_id,
  api_gateway_managed,
  api_key_required
from
  aws_api_gatewayv2_route;
```

### List routes by API

```sql
select
  route_key,
  api_id,
  route_id
from
  aws_api_gatewayv2_route
where
  api_id = 'w5n71b2m85';
```

### List routes with default endpoint enabled APIs

```sql
select
  r.route_id,
  a.name,
  a.api_id,
  a.api_endpoint
from
  aws_api_gatewayv2_route as r,
  aws_api_gatewayv2_api as a
where
  not a.disable_execute_api_endpoint;
```
