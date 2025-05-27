---
title: "Steampipe Table: aws_api_gatewayv2_route - Query AWS API Gateway V2 Routes using SQL"
description: "Allows users to query AWS API Gateway V2 Routes and obtain detailed information about each route, including the route key, route response selection expression, and target."
folder: "API Gateway"
---

# Table: aws_api_gatewayv2_route - Query AWS API Gateway V2 Routes using SQL

The AWS API Gateway V2 Routes is a feature within the Amazon API Gateway service. It allows you to define the paths that a client application can take to access your API. This feature is integral to the process of creating, deploying, and managing your APIs in a secure and scalable manner.

## Table Usage Guide

The `aws_api_gatewayv2_route` table in Steampipe provides you with information about routes within AWS API Gateway V2. This table allows you, as a DevOps engineer, to query route-specific details, including the route key, route response selection expression, and target. You can utilize this table to gather insights on routes, such as route configurations, route response behaviors, and more. The schema outlines the various attributes of the route for you, including the API identifier, route ID, route key, and associated metadata.

## Examples

### Basic info
Determine the areas in which your AWS API Gateway is managed and if an API key is required. This can help in identifying potential security risks and ensuring appropriate access controls are in place.

```sql+postgres
select
  route_key,
  api_id,
  route_id,
  api_gateway_managed,
  api_key_required
from
  aws_api_gatewayv2_route;
```

```sql+sqlite
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
Explore which routes are associated with a specific API to better manage and optimize your API Gateway. This can be particularly useful for troubleshooting or for identifying opportunities for API performance enhancement.

```sql+postgres
select
  route_key,
  api_id,
  route_id
from
  aws_api_gatewayv2_route
where
  api_id = 'w5n71b2m85';
```

```sql+sqlite
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
Identify the instances where the default endpoint is enabled in APIs, allowing you to understand and manage the routes that are directly accessible.

```sql+postgres
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

```sql+sqlite
select
  r.route_id,
  a.name,
  a.api_id,
  a.api_endpoint
from
  aws_api_gatewayv2_route as r,
  aws_api_gatewayv2_api as a
where
  a.disable_execute_api_endpoint != 1;
```