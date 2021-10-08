# Table: aws_cloudcontrol_resource

The Cloud Control resource table allows you to list and read a wide range of AWS and third-party resources. A full list of supported AWS resource types can be found at [Resource types that support Cloud Control API](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/supported-resources.html).

In order to list resources, the `type_name` column must be specified. Some resources also require additional information, which is specified in the `resource_model` column. For more information on these resource types, please see [Resources that require additional information](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/resource-operations-list.html#resource-operations-list-containers).

In order to read a resource, the `type_name` and `identifier` columns must be specified. The identifier for each resource type is different, for more information on identifiers please see [Identifying resources](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/resource-identifier.html).

## Examples

### List Lambda functions

```sql
select
  identifier,
  properties ->> 'Arn' as arn,
  properties ->> 'MemorySize' as memory_size,
  properties ->> 'Runtime' as runtime,
  region
from
  aws_cloudcontrolapi_resource
where
  type_name = 'AWS::Lambda::Function';
```

### List ELBv2 listeners for a load balancer

```sql
select
  identifier,
  properties ->> 'AlpnPolicy' as alpn_policy,
  properties ->> 'Certificates' as certificates,
  properties ->> 'Port' as port,
  properties ->> 'Protocol' as protocol,
  region
from
  aws_cloudcontrolapi_resource
where
  type_name = 'AWS::ElasticLoadBalancingV2::Listener'
  and resource_model = '{"LoadBalancerArn": "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/4e695b8755d7003c"}';
```

### Get details for a CloudTrail trail

```sql
select
  identifier,
  properties ->> 'IncludeGlobalServiceEvents' as include_global_service_events,
  properties ->> 'IsLogging' as is_logging,
  properties ->> 'IsMultiRegionTrail' as is_multi_region_trail,
  region
from
  aws_cloudcontrolapi_resource
where
  type_name = 'AWS::CloudTrail::Trail'
  and identifier = 'my-trail';
```
