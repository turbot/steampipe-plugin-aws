# Table: aws_cloudcontrol_resource

The Cloud Control resource table allows you to list and read a wide range of AWS and third-party resources. A full list of supported AWS resource types can be found at [Resource types that support Cloud Control API](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/supported-resources.html).

In order to list resources, the `type_name` column must be specified. Some resources also require additional information, which is specified in the `resource_model` column. For more information on these resource types, please see [Resources that require additional information](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/resource-operations-list.html#resource-operations-list-containers).

In order to read a resource, the `type_name` and `identifier` columns must be specified. The identifier for each resource type is different, for more information on identifiers please see [Identifying resources](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/resource-identifier.html).

_We recommend you use native Steampipe tables when available, but this table is helpful to query uncommon resources not yet supported._

## Known limitations

* `AWS::S3::Bucket` will only include detailed information if an identifier is provided. There is no way to determine the region of a bucket from the list result, so full information cannot be automatically hydrated.
* Global resources like `AWS::IAM::Role` will return duplicate results per region. Specify `region = 'us-east-1'` (or similar) in the where clause to avoid.

For more information on other Cloud Control limitations and caveats, please see [A deep dive into AWS Cloud Control for asset inventory](https://steampipe.io/blog/aws-cloud-control).

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
  aws_cloudcontrol_resource
where
  type_name = 'AWS::Lambda::Function';
```

### List ELBv2 listeners for a load balancer

Listeners are a sub-resource, so can only be listed if passed the `LoadBalancerArn` data.

Warning: This does not work with multi-account in Steampipe. The query will be run
against all accounts and Cloud Control returns a GeneralServiceException (rather than
NotFound), making it difficult to handle.

Warning: If using multi-region in Steampipe then you MUST specify the region in
the query. Otherwise, the request will be tried against each region. This would
be slow anyway, but because Cloud Control returns a GeneralServiceException (rather
than NotFound), we cannot handle it automatically.

```sql
select
  identifier,
  properties ->> 'AlpnPolicy' as alpn_policy,
  properties ->> 'Certificates' as certificates,
  properties ->> 'Port' as port,
  properties ->> 'Protocol' as protocol,
  region,
  account_id
from
  aws_cloudcontrol_resource
where
  type_name = 'AWS::ElasticLoadBalancingV2::Listener'
  and resource_model = '{"LoadBalancerArn": "arn:aws:elasticloadbalancing:us-east-1:123456789012:loadbalancer/app/test-lb/4e695b8755d7003c"}'
  and region = 'us-east-1';
```

### Get details for a CloudTrail trail

Get a single specific resource by setting the identifier.

```sql
select
  identifier,
  properties ->> 'IncludeGlobalServiceEvents' as include_global_service_events,
  properties ->> 'IsLogging' as is_logging,
  properties ->> 'IsMultiRegionTrail' as is_multi_region_trail,
  region
from
  aws_cloudcontrol_resource
where
  type_name = 'AWS::CloudTrail::Trail'
  and identifier = 'my-trail';
```

### List global resources using a single region

Global resources (e.g. `AWS::IAM::Role`) are returned by each region endpoint.
When working with a multi-region configuration in Steampipe this creates
duplicate rows. To avoid the duplicates, you can specify a region qualifier.

```sql
select
  properties ->> 'RoleName' as name
from
  aws_cloudcontrol_resource
where
  type_name = 'AWS::IAM::Role'
  and region = 'us-east-1'
order by
  name;
```
