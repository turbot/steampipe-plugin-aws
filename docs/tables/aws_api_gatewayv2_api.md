---
title: "Steampipe Table: aws_api_gatewayv2_api - Query AWS API Gateway using SQL"
description: "Allows users to query API Gateway APIs and retrieve detailed information about each API, including its ID, name, protocol type, and more."
folder: "API Gateway"
---

# Table: aws_api_gatewayv2_api - Query AWS API Gateway using SQL

The AWS API Gateway is a fully managed service that makes it easy for developers to create, publish, maintain, monitor, and secure APIs at any scale. It handles all the tasks involved in accepting and processing up to hundreds of thousands of concurrent API calls, including traffic management, authorization and access control, monitoring, and API version management. With the use of SQL, you can query and manage your API Gateway effectively.

## Table Usage Guide

The `aws_api_gatewayv2_api` table in Steampipe provides you with information about APIs within AWS API Gateway. This table allows you, as a DevOps engineer, to query API-specific details, including the API ID, name, protocol type, route selection expression, and associated tags. You can utilize this table to gather insights on APIs, such as their configuration details, associated resources, and more. The schema outlines the various attributes of the API for you, including the API key selection expression, CORS configuration, created date, and description.

## Examples

### Basic info
Explore the configuration of your AWS API Gateway to gain insights into its protocol type and endpoint. This allows for a better understanding of how your API is set up and can assist in troubleshooting or optimizing API performance."Explore the essential details of your AWS API Gateway configurations to understand the protocols used and how routes and keys are selected. This information can aid in optimizing your API setup and troubleshooting issues."


```sql+postgres
select
  name,
  api_id,
  api_endpoint,
  protocol_type,
  api_key_selection_expression,
  route_selection_expression
from
  aws_api_gatewayv2_api;
```

```sql+sqlite
select
  name,
  api_id,
  api_endpoint,
  protocol_type,
  api_key_selection_expression,
  route_selection_expression
from
  aws_api_gatewayv2_api;
```

### List APIs with protocol type WEBSOCKET
Uncover the details of APIs that are using the WebSocket protocol. This can be useful for identifying which APIs may need specific handling or monitoring due to their protocol type."Identify instances where AWS APIs are using the WebSocket protocol. This allows you to understand which APIs are designed for real-time, two-way interactive communication."

```sql+postgres
select
  name,
  api_id,
  protocol_type
from
  aws_api_gatewayv2_api
where
  protocol_type = 'WEBSOCKET';
```

```sql+sqlite
select
  name,
  api_id,
  protocol_type
from
  aws_api_gatewayv2_api
where
  protocol_type = 'WEBSOCKET';
```

### List APIs with default endpoint enabled
Determine the areas in which APIs are operating with the default endpoint enabled. This can be particularly useful for identifying potential security risks and ensuring best practices in endpoint configuration."Identify all APIs in your AWS environment where the default endpoint is enabled. This can be useful to ensure that no unnecessary endpoints are open, potentially reducing the risk of security breaches."


```sql+postgres
select
  name,
  api_id,
  api_endpoint
from
  aws_api_gatewayv2_api
where
  not disable_execute_api_endpoint;
```

```sql+sqlite
select
  name,
  api_id,
  api_endpoint
from
  aws_api_gatewayv2_api
where
  disable_execute_api_endpoint = 0;
```