---
title: "Steampipe Table: aws_api_gateway_authorizer - Query AWS API Gateway Authorizer using SQL"
description: "Allows users to query AWS API Gateway Authorizer and access data about API Gateway Authorizers in an AWS account. This data includes the authorizer's ID, name, type, provider ARNs, and other configuration details."
folder: "API Gateway"
---

# Table: aws_api_gateway_authorizer - Query AWS API Gateway Authorizer using SQL

The AWS API Gateway Authorizer is a crucial component in Amazon API Gateway that validates incoming requests before they reach the backend systems. It verifies the caller's identity and checks if the caller has permission to execute the requested operation. This feature enhances the security of your APIs by preventing unauthorized access to your resources.

## Table Usage Guide

The `aws_api_gateway_api_authorizer` table in Steampipe provides you with information about API Gateway Authorizers within AWS API Gateway. This table allows you, as a DevOps engineer, to query authorizer-specific details, including the authorizer's ID, name, type, provider ARNs, and other configuration details. You can utilize this table to gather insights on authorizers, such as the authorizer's type, the ARN of the authorizer's provider, and more. The schema outlines the various attributes of the API Gateway Authorizer for you, including the authorizer's ID, name, type, provider ARNs, and associated metadata.

## Examples

### API gateway API authorizer basic info
Explore the core details of an API gateway's authorizer configuration, such as its ID, name, and authorization type. This can help you understand the security measures in place for your API gateway and can be useful for auditing purposes.

```sql+postgres
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

```sql+sqlite
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
Explore which API authorizers are utilizing Cognito user pools for API call authorization. This can help in assessing the security configuration of your APIs and identify any potential areas for improvement.

```sql+postgres
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

```sql+sqlite
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