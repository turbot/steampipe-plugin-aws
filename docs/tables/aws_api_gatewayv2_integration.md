---
title: "Table: aws_api_gatewayv2_integration - Query AWS API Gateway Integrations using SQL"
description: "Allows users to query AWS API Gateway Integrations to retrieve detailed information about each integration within the API Gateway."
---

# Table: aws_api_gatewayv2_integration - Query AWS API Gateway Integrations using SQL

The `aws_api_gatewayv2_integration` table in Steampipe provides information about each integration within AWS API Gateway. This table allows DevOps engineers to query integration-specific details, including the integration type, API Gateway ID, integration method, and more. Users can utilize this table to gather insights on integrations, such as integration protocols, request templates, and connection type. The schema outlines the various attributes of the integration, including the integration ID, integration response selection expression, integration subtype, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gatewayv2_integration` table, you can use the `.inspect aws_api_gatewayv2_integration` command in Steampipe.

**Key columns**:

- `api_id`: The API identifier. This can be used to join with other tables that contain API-specific information.
- `integration_id`: The identifier of the integration. This can be used to join with other tables that contain integration-specific information.
- `integration_type`: The type of the integration. This can be useful for filtering or grouping integrations based on their type.

## Examples

### Basic info

```sql
select
  integration_id,
  api_id,
  integration_type,
  integration_uri,
  description
from
  aws_api_gatewayv2_integration;
```

### Count of integrations per API

```sql
select 
  api_id,
  count(integration_id) as integration_count
from 
  aws_api_gatewayv2_integration
group by
  api_id;
```
