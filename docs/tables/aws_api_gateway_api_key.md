---
title: "Table: aws_api_gateway_api_key - Query AWS API Gateway API Keys using SQL"
description: "Allows users to query API Keys in AWS API Gateway. The `aws_api_gateway_api_key` table in Steampipe provides information about API Keys within AWS API Gateway. This table allows DevOps engineers to query API Key-specific details, including its ID, value, enabled status, and associated metadata. Users can utilize this table to gather insights on API Keys, such as keys that are enabled, keys associated with specific stages, and more. The schema outlines the various attributes of the API Key, including the key ID, creation date, enabled status, and associated tags."
---

# Table: aws_api_gateway_api_key - Query AWS API Gateway API Keys using SQL

The `aws_api_gateway_api_key` table in Steampipe provides information about API Keys within AWS API Gateway. This table allows DevOps engineers to query API Key-specific details, including its ID, value, enabled status, and associated metadata. Users can utilize this table to gather insights on API Keys, such as keys that are enabled, keys associated with specific stages, and more. The schema outlines the various attributes of the API Key, including the key ID, creation date, enabled status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gateway_api_key` table, you can use the `.inspect aws_api_gateway_api_key` command in Steampipe.

### Key columns:

- `id`: The identifier of the API Key. This can be used to join this table with other tables to get more detailed information about the API Key.
- `value`: The value of the API Key. This is important as it is the unique value used to authenticate requests to your APIs.
- `enabled`: The status of the API Key. This is useful to know whether the API Key is currently active or not.

## Examples

### API gateway API key basic info

```sql
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

```sql
select
  name,
  id,
  customer_id
from
  aws_api_gateway_api_key
where
  not enabled;
```
