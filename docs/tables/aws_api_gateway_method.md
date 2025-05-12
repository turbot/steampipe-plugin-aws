---
title: "Steampipe Table: aws_api_gateway_method - Query AWS API Gateway Methods using SQL"
description: "Represents a client-facing interface by which the client calls the API to access back-end resources. A Method resource is integrated with an Integration resource. Both consist of a request and one or more responses. The method request takes the client input that is passed to the back end through the integration request. A method response returns the output from the back end to the client through an integration response. A method request is embodied in a Method resource, whereas an integration request is embodied in an Integration resource. On the other hand, a method response is represented by a MethodResponse resource, whereas an integration response is represented by an IntegrationResponse resource."
folder: "API Gateway"
---

# Table: aws_api_gateway_method - Query AWS API Gateway Methods using SQL

Represents a client-facing interface by which the client calls the API to access back-end resources. A Method resource is integrated with an Integration resource. Both consist of a request and one or more responses. The method request takes the client input that is passed to the back end through the integration request. A method response returns the output from the back end to the client through an integration response. A method request is embodied in a Method resource, whereas an integration request is embodied in an Integration resource. On the other hand, a method response is represented by a MethodResponse resource, whereas an integration response is represented by an IntegrationResponse resource.

## Table Usage Guide

The `aws_api_gateway_method` table in Steampipe allows users to query information about AWS API Gateway Methods. These methods represent client-facing interfaces for accessing back-end resources. Users can retrieve details such as the REST API ID, resource ID, HTTP method, path, and whether API key authorization is required. Additionally, users can query methods with specific criteria, such as HTTP method type or authorization type.

## Examples

### Basic info
Retrieve basic information about AWS API Gateway Methods, including the REST API ID, resource ID, HTTP method, path, and whether API key authorization is required. This query provides an overview of the methods in your AWS API Gateway.

```sql+postgres
select
  rest_api_id,
  resource_id,
  http_method,
  path,
  api_key_required
from
  aws_api_gateway_method;
```

```sql+sqlite
select
  rest_api_id,
  resource_id,
  http_method,
  path,
  api_key_required
from
  aws_api_gateway_method;
```

### List API Gateway GET methods
Identify AWS API Gateway Methods that use the HTTP GET method. This query helps you filter and view specific types of methods in your API Gateway.

```sql+postgres
select
  rest_api_id,
  resource_id,
  http_method,
  operation_name
from
  aws_api_gateway_method
where
  http_method = 'GET';
```

```sql+sqlite
select
  rest_api_id,
  resource_id,
  http_method,
  operation_name
from
  aws_api_gateway_method
where
  http_method = 'GET';
```

### List methods with open access
Retrieve AWS API Gateway Methods that do not require any authorization. This query helps you identify methods with open access settings.

```sql+postgres
select
  rest_api_id,
  resource_id,
  http_method,
  path,
  authorization_type,
  authorizer_id
from
  aws_api_gateway_method
where
  authorization_type = 'none';
```

```sql+sqlite
select
  rest_api_id,
  resource_id,
  http_method,
  path,
  authorization_type,
  authorizer_id
from
  aws_api_gateway_method
where
  authorization_type = 'none';
```

### Get integration details of methods
Retrieve detailed integration configuration information for AWS API Gateway Methods. This query includes information such as cache key parameters, cache namespace, connection ID, connection type, content handling, credentials, HTTP method, passthrough behavior, request parameters, request templates, timeout in milliseconds, TLS configuration, integration type, URI, and integration responses.

```sql+postgres
select
  rest_api_id,
  resource_id,
  http_method,
  method_integration -> 'CacheKeyParameters' as cache_key_parameters,
  method_integration ->> 'CacheNamespace' as cache_namespace,
  method_integration ->> 'ConnectionId' as connection_id,
  method_integration ->> 'ConnectionType' as connection_type,
  method_integration ->> 'ContentHandling' as content_handling,
  method_integration ->> 'Credentials' as credentials,
  method_integration ->> 'HttpMethod' as http_method,
  method_integration ->> 'PassthroughBehavior' as passthrough_behavior,
  method_integration ->> 'RequestParameters' as request_parameters,
  method_integration -> 'RequestTemplates' as request_templates,
  method_integration ->> 'TimeoutInMillis' as timeout_in_millis,
  method_integration ->> 'tls_config' as tls_config,
  method_integration ->> 'Type' as type,
  method_integration ->> 'Uri' as uri,
  method_integration -> 'IntegrationResponses' as integration_responses
from
  aws_api_gateway_method;
```

```sql+sqlite
select
  rest_api_id,
  resource_id,
  http_method,
  json_extract(method_integration, '$.CacheKeyParameters') as cache_key_parameters,
  json_extract(method_integration, '$.CacheNamespace') as cache_namespace,
  json_extract(method_integration, '$.ConnectionId') as connection_id,
  json_extract(method_integration, '$.ConnectionType') as connection_type,
  json_extract(method_integration, '$.ContentHandling') as content_handling,
  json_extract(method_integration, '$.Credentials') as credentials,
  json_extract(method_integration, '$.HttpMethod') as http_method,
  json_extract(method_integration, '$.PassthroughBehavior') as passthrough_behavior,
  json_extract(method_integration, '$.RequestParameters') as request_parameters,
  json_extract(method_integration, '$.RequestTemplates') as request_templates,
  json_extract(method_integration, '$.TimeoutInMillis') as timeout_in_millis,
  json_extract(method_integration, '$.tls_config') as tls_config,
  json_extract(method_integration, '$.Type') as type,
  json_extract(method_integration, '$.Uri') as uri,
  json_extract(method_integration, '$.IntegrationResponses') as integration_responses
from
  aws_api_gateway_method;
```
