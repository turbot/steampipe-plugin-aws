---
title: "Table: aws_api_gateway_usage_plan - Query AWS API Gateway Usage Plans using SQL"
description: "Allows users to query AWS API Gateway Usage Plans in order to retrieve information about the usage plans configured in the AWS API Gateway service."
---

# Table: aws_api_gateway_usage_plan - Query AWS API Gateway Usage Plans using SQL

The `aws_api_gateway_usage_plan` table in Steampipe provides information about usage plans within AWS API Gateway. This table allows DevOps engineers to query usage plan specific details, including associated API stages, throttle and quota limits, and associated metadata. Users can utilize this table to gather insights on usage plans, such as plans with specific rate limits, the number of requests clients can make per a given period, and more. The schema outlines the various attributes of the usage plan, including the plan ID, name, description, associated API keys, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gateway_usage_plan` table, you can use the `.inspect aws_api_gateway_usage_plan` command in Steampipe.

Key columns:

- `id`: The identifier of the usage plan. This can be used to join this table with other tables that contain usage plan information.
- `name`: The name of the usage plan. This can be useful for joining with tables that contain usage plan names as identifiers.
- `product_code`: The AWS product code associated with the usage plan. This can be used to join with other AWS product-related tables.

## Examples

### Basic info

```sql
select
  name,
  id,
  product_code,
  description,
  api_stages
from
  aws_api_gateway_usage_plan;
```


### List the API gateway usage plans where quota ( i.e the number of api call a user can make within a time period) is disabled

```sql
select
  name,
  id,
  quota
from
  aws_api_gateway_usage_plan
where
  quota is null;
```


### List the API gateway usage plan where throttle ( i.e the rate at which user can make request ) is disabled

```sql
select
  name,
  id,
  throttle
from
  aws_api_gateway_usage_plan
where
  throttle is null;
```
