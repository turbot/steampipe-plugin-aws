---
title: "Table: aws_api_gateway_api_authorizer - Query AWS API Gateway Authorizer using SQL"
description: "Allows users to query AWS API Gateway Authorizer and access data about API Gateway Authorizers in an AWS account. This data includes the authorizer's ID, name, type, provider ARNs, and other configuration details."
---

# Table: aws_api_gateway_api_authorizer - Query AWS API Gateway Authorizer using SQL

The `aws_api_gateway_api_authorizer` table in Steampipe provides information about API Gateway Authorizers within AWS API Gateway. This table allows DevOps engineers to query authorizer-specific details, including the authorizer's ID, name, type, provider ARNs, and other configuration details. Users can utilize this table to gather insights on authorizers, such as the authorizer's type, the ARN of the authorizer's provider, and more. The schema outlines the various attributes of the API Gateway Authorizer, including the authorizer's ID, name, type, provider ARNs, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gateway_api_authorizer` table, you can use the `.inspect aws_api_gateway_api_authorizer` command in Steampipe.

**Key columns**:

- `api_id`: The API identifier. This can be used to join with other tables that contain API Gateway information.
- `authorizer_id`: The identifier of the authorizer. This can be used to join with other tables that contain authorizer information.
- `name`: The name of the authorizer. This can be used to join with other tables that contain authorizer information.

## Examples

### API gateway API authorizer basic info

```sql
select
  id,
  name,
  rest_api_id,
  auth_type,
  authorizer_credentials,
  identity_validation_expression,
  identity_source
from
  aws_api_gateway_authorizer;
```


### List the API authorizers that uses cognito user pool to authorize API calls

```sql
select
  id,
  name,
  rest_api_id,
  auth_type
from
  aws_api_gateway_authorizer
where
  auth_type = 'cognito_user_pools';
```
