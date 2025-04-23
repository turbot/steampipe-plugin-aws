---
title: "Steampipe Table: aws_api_gatewayv2_integration - Query AWS API Gateway Integrations using SQL"
description: "Allows users to query AWS API Gateway Integrations to retrieve detailed information about each integration within the API Gateway."
folder: "API Gateway"
---

# Table: aws_api_gatewayv2_integration - Query AWS API Gateway Integrations using SQL

The AWS API Gateway Integrations is a feature within the Amazon API Gateway service that allows you to integrate backend operations such as Lambda functions, HTTP endpoints, and other AWS services into your API. These integrations enable your API to interact with these services, processing incoming requests and returning responses to the client. This functionality aids in creating efficient, scalable, and secure APIs.

## Table Usage Guide

The `aws_api_gatewayv2_integration` table in Steampipe provides you with information about each integration within AWS API Gateway. This table allows you as a DevOps engineer to query integration-specific details, including the integration type, API Gateway ID, integration method, and more. You can utilize this table to gather insights on integrations, such as integration protocols, request templates, and connection type. The schema outlines the various attributes of the integration for you, including the integration ID, integration response selection expression, integration subtype, and associated tags.

## Examples

### Basic info
Determine the areas in which specific integrations are being used within AWS API Gateway. This can help in understanding the scope and purpose of these integrations, aiding in efficient system management and optimization.Explore the different types of integrations within your AWS API Gateway and understand their respective roles and characteristics. This can help you manage and optimize your API configurations effectively.

```sql+postgres
select
  integration_id,
  api_id,
  integration_type,
  integration_uri,
  description
from
  aws_api_gatewayv2_integration;
```

```sql+sqlite
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
Explore which APIs have the most integrations to identify potential areas of complexity or high usage. This can help in managing resources and planning future developments.Analyze the distribution of integrations across different APIs to understand the extent of their utilization. This can help identify heavily used APIs, potentially indicating areas for performance optimization or resource allocation.


```sql+postgres
select 
  api_id,
  count(integration_id) as integration_count
from 
  aws_api_gatewayv2_integration
group by
  api_id;
```

```sql+sqlite
select 
  api_id,
  count(integration_id) as integration_count
from 
  aws_api_gatewayv2_integration
group by
  api_id;
```