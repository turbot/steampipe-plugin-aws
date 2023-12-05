# Table: aws_api_gateway_method

Represents a client-facing interface by which the client calls the API to access back-end resources. A Method resource is integrated with an Integration resource. Both consist of a request and one or more responses. The method request takes the client input that is passed to the back end through the integration request. A method response returns the output from the back end to the client through an integration response. A method request is embodied in a Method resource, whereas an integration request is embodied in an Integration resource. On the other hand, a method response is represented by a MethodResponse resource, whereas an integration response is represented by an IntegrationResponse resource.

## Examples

### Basic info

```sql
select
  rest_api_id,
  resource_id,
  http_method,
  path,
  api_key_required
from
  aws_api_gateway_method;
```

### List API Gateway get methods

```sql
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

### List methods that are open access

```sql
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

```sql
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
