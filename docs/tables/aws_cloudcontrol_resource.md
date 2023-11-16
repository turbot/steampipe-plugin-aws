---
title: "Table: aws_cloudcontrol_resource - Query AWS Cloud Control API Resource using SQL"
description: "Allows users to query AWS Cloud Control API Resource data, providing detailed insights into resource properties, types, and statuses."
---

# Table: aws_cloudcontrol_resource - Query AWS Cloud Control API Resource using SQL

The `aws_cloudcontrol_resource` table in Steampipe provides information about resources within AWS Cloud Control API. This table allows DevOps engineers to query resource-specific details, including resource properties, types, and statuses. Users can utilize this table to gather insights on resources, such as the resource's specific properties, the type of the resource, and the current status of the resource. The schema outlines the various attributes of the AWS Cloud Control API resource, including the resource name, resource type, role ARN, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudcontrol_resource` table, you can use the `.inspect aws_cloudcontrol_resource` command in Steampipe.

**Key columns**:

- `resource_name`: This is the unique identifier of the resource. It is important as it can be used to join this table with other tables that also contain resource identifiers.
- `resource_type`: This column provides information about the type of the resource. It is useful for filtering or grouping resources based on their types.
- `status`: This column indicates the current status of the resource. It can be used to filter or group resources based on their statuses.

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
