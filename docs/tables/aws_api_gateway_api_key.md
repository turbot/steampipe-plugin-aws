---
title: "Steampipe Table: aws_api_gateway_api_key - Query AWS API Gateway API Keys using SQL"
description: "Allows users to query API Keys in AWS API Gateway. The `aws_api_gateway_api_key` table in Steampipe provides information about API Keys within AWS API Gateway. This table allows DevOps engineers to query API Key-specific details, including its ID, value, enabled status, and associated metadata. Users can utilize this table to gather insights on API Keys, such as keys that are enabled, keys associated with specific stages, and more. The schema outlines the various attributes of the API Key, including the key ID, creation date, enabled status, and associated tags."
folder: "API Gateway"
---

# Table: aws_api_gateway_api_key - Query AWS API Gateway API Keys using SQL

AWS API Gateway API Keys are used to control and track API usage in Amazon API Gateway. They are associated with API stages to manage access and can be used in conjunction with usage plans to authorize access to specific APIs. API keys are not meant for client-side security, but rather for tracking and controlling how your customers use your API.

## Table Usage Guide

The `aws_api_gateway_api_key` table in Steampipe provides you with information about API Keys within AWS API Gateway. This table allows you, as a DevOps engineer, to query API Key-specific details, including its ID, value, enabled status, and associated metadata. You can utilize this table to gather insights on API Keys, such as keys that are enabled, keys associated with specific stages, and more. The schema outlines the various attributes of the API Key for you, including the key ID, creation date, enabled status, and associated tags.

## Examples

### API gateway API key basic info
Discover the segments that utilize the API gateway key within the AWS infrastructure. This query can provide insights into the status and usage of API keys, which can be beneficial for monitoring security and optimizing resource utilization.

```sql+postgres
select
  name,
  id,
  enabled,
  created_date,
  last_updated_date,
  customer_id,
  stage_keys
from
  aws_api_gateway_api_key;
```

```sql+sqlite
select
  name,
  id,
  enabled,
  created_date,
  last_updated_date,
  customer_id,
  stage_keys
from
  aws_api_gateway_api_key;
```


### List of API keys which are not enabled
Determine the areas in which API keys are not activated to assess potential security risks or unused resources within your AWS API Gateway.

```sql+postgres
select
  name,
  id,
  customer_id
from
  aws_api_gateway_api_key
where
  not enabled;
```

```sql+sqlite
select
  name,
  id,
  customer_id
from
  aws_api_gateway_api_key
where
  enabled = 0;
```