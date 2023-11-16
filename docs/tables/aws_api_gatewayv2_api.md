---
title: "Table: aws_api_gatewayv2_api - Query AWS API Gateway using SQL"
description: "Allows users to query API Gateway APIs and retrieve detailed information about each API, including its ID, name, protocol type, and more."
---

# Table: aws_api_gatewayv2_api - Query AWS API Gateway using SQL

The `aws_api_gatewayv2_api` table in Steampipe provides information about APIs within AWS API Gateway. This table allows DevOps engineers to query API-specific details, including the API ID, name, protocol type, route selection expression, and associated tags. Users can utilize this table to gather insights on APIs, such as their configuration details, associated resources, and more. The schema outlines the various attributes of the API, including the API key selection expression, CORS configuration, created date, and description.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gatewayv2_api` table, you can use the `.inspect aws_api_gatewayv2_api` command in Steampipe.

**Key columns**:

- `api_id`: The API identifier. This is a unique key that can be used to join this table with other tables to gather more detailed information about the API.
- `name`: The name of the API. This column can be useful for identifying specific APIs when joining with other tables.
- `protocol_type`: The protocol type of the API. This column can provide insights into the type of protocol (HTTP, WEBSOCKET) used by the API.


## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  name,
  api_id,
  api_endpoint
from
  aws_api_gatewayv2_api
where
  not disable_execute_api_endpoint;
```
