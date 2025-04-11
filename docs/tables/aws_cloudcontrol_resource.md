---
title: "Steampipe Table: aws_cloudcontrol_resource - Query AWS Cloud Control API Resource using SQL"
description: "Allows users to query AWS Cloud Control API Resource data, providing detailed insights into resource properties, types, and statuses."
folder: "CloudControl"
---

# Table: aws_cloudcontrol_resource - Query AWS Cloud Control API Resource using SQL

The AWS Cloud Control API Resource is a service that allows you to manage your cloud resources in a programmatic way. It provides a unified, consistent set of application programming interfaces (APIs) and extends the capabilities of AWS CloudFormation to support all AWS resource types. This service allows you to create, read, update, delete, and list resources across multiple AWS services from a single API endpoint.

## Table Usage Guide

The `aws_cloudcontrol_resource` table in Steampipe provides you with information about resources within the AWS Cloud Control API. This table allows you, as a DevOps engineer, to query resource-specific details, including resource properties, types, and statuses. You can utilize this table to gather insights on resources, such as the resource's specific properties, the type of the resource, and the current status of the resource. The schema outlines for you the various attributes of the AWS Cloud Control API resource, including the resource name, resource type, role ARN, and associated metadata.

**Important Notes**
- In order to list resources, the `type_name` column must be specified. Some resources also require additional information, which is specified in the `resource_model` column. For more information on these resource types, please see [Resources that require additional information](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/resource-operations-list.html#resource-operations-list-containers).

- In order to read a resource, the `type_name` and `identifier` columns must be specified. The identifier for each resource type is different, for more information on identifiers please see [Identifying resources](https://docs.aws.amazon.com/cloudcontrolapi/latest/userguide/resource-identifier.html).

_We recommend you use native Steampipe tables when available, but this table is helpful to query uncommon resources not yet supported._

**Known limitations**

* `AWS::S3::Bucket` will only include detailed information if an identifier is provided. There is no way to determine the region of a bucket from the list result, so full information cannot be automatically hydrated.
* Global resources like `AWS::IAM::Role` will return duplicate results per region. Specify `region = 'us-east-1'` (or similar) in the where clause to avoid.

For more information on other Cloud Control limitations and caveats, please see [A deep dive into AWS Cloud Control for asset inventory](https://steampipe.io/blog/aws-cloud-control).

## Examples

### List Lambda functions
Explore the Lambda functions within your AWS environment, focusing on aspects like their associated identifiers, regions, and runtime settings. This analysis can help in understanding the setup and distribution of your Lambda functions, which is crucial for optimizing resource allocation and troubleshooting.

```sql+postgres
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

```sql+sqlite
select
  identifier,
  json_extract(properties, '$.Arn') as arn,
  json_extract(properties, '$.MemorySize') as memory_size,
  json_extract(properties, '$.Runtime') as runtime,
  region
from
  aws_cloudcontrol_resource
where
  type_name = 'AWS::Lambda::Function';
```

### List ELBv2 listeners for a load balancer
Explore the settings of specific listeners within a load balancer to understand their protocols, ports, and certificates, particularly useful for auditing and optimizing network traffic management.
Listeners are a sub-resource, so can only be listed if passed the `LoadBalancerArn` data.

Warning: This does not work with multi-account in Steampipe. The query will be run
against all accounts and Cloud Control returns a GeneralServiceException (rather than
NotFound), making it difficult to handle.

Warning: If using multi-region in Steampipe then you MUST specify the region in
the query. Otherwise, the request will be tried against each region. This would
be slow anyway, but because Cloud Control returns a GeneralServiceException (rather
than NotFound), we cannot handle it automatically.


```sql+postgres
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

```sql+sqlite
select
  identifier,
  json_extract(properties, '$.AlpnPolicy') as alpn_policy,
  json_extract(properties, '$.Certificates') as certificates,
  json_extract(properties, '$.Port') as port,
  json_extract(properties, '$.Protocol') as protocol,
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
Determine the status and settings of a specific CloudTrail trail in your AWS environment. This can be essential in auditing and understanding your cloud resource configuration for security and compliance purposes.
Get a single specific resource by setting the identifier.


```sql+postgres
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

```sql+sqlite
select
  identifier,
  json_extract(properties, '$.IncludeGlobalServiceEvents') as include_global_service_events,
  json_extract(properties, '$.IsLogging') as is_logging,
  json_extract(properties, '$.IsMultiRegionTrail') as is_multi_region_trail,
  region
from
  aws_cloudcontrol_resource
where
  type_name = 'AWS::CloudTrail::Trail'
  and identifier = 'my-trail';
```

### List global resources using a single region
Determine the areas in which global resources are utilized through a single region. This is useful for managing and optimizing resource usage within a specified region in the AWS cloud environment.
Global resources (e.g. `AWS::IAM::Role`) are returned by each region endpoint.
When working with a multi-region configuration in Steampipe this creates
duplicate rows. To avoid the duplicates, you can specify a region qualifier.


```sql+postgres
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

```sql+sqlite
select
  json_extract(properties, '$.RoleName') as name
from
  aws_cloudcontrol_resource
where
  type_name = 'AWS::IAM::Role'
  and region = 'us-east-1'
order by
  name;
```