---
title: "Steampipe Table: aws_shield_protection - Query AWS Shield Advanced Protections using SQL"
description: "Allows users to query AWS Shield Advanced Protections and retrieve detailed information about each protection's settings."
folder: "Shield"
---

# Table: aws_shield_protection - Query AWS Shield Advanced Protections using SQL

AWS Shield Advanced Protections are safeguards provided by AWS to protect AWS resources against Distributed Denial of Service (DDoS) attacks.

## Table Usage Guide

The `aws_shield_protection` table in Steampipe allows you to query AWS Shield Advanced Protections and retrieve detailed information about each protection's settings. This table provides you with insights into the protections that are currently active in your AWS environment, including the ARN of the resource that is protected and the automatic application layer DDoS mitigation setting. You can use this table to monitor the status of your AWS Shield Advanced Protections and ensure that your resources are protected against DDoS attacks. For more information about the individual fields, please refer to the [AWS Shield Advanced API documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_DescribeProtection.html#API_DescribeProtection_ResponseSyntax).

**Note:** The column `resource_type` only has a value when it was part of the where clause. For a list of valid values for filtering by `resource_type`, please refer to the [AWS documentation](https://docs.aws.amazon.com/waf/latest/DDOSAPIReference/API_InclusionProtectionFilters.html#AWSShield-Type-InclusionProtectionFilters-ResourceTypes).

## Examples

### Basic info

Discover the protections that are currently active in your account and what kind of resources they are protecting.

```sql+postgres
select
  name,
  resource_arn
from
  aws_shield_protection;
```

```sql+sqlite
select
  name,
  resource_arn
from
  aws_shield_protection;
```

### Identify Protections without Route 53 Health Checks

Identify the protections that are missing Route 53 Health Checks. This information can be useful to see which protections still need Health Checks in order to be covered by the proactive engagement of the Shield Response Team.

```sql+postgres
select
  name,
  resource_arn
from
  aws_shield_protection
where
  health_check_ids is null;
```

```sql+sqlite
select
  name,
  resource_arn
from
  aws_shield_protection
where
  health_check_ids is null;
```

### List Protections for Route 53 Hosted Zones

```sql+postgres
select
  name,
  resource_arn
from
  aws_shield_protection
where
  resource_type = 'ROUTE_53_HOSTED_ZONE';
```

```sql+sqlite
select
  name,
  resource_arn
from
  aws_shield_protection
where
  resource_type = 'ROUTE_53_HOSTED_ZONE';
```

### Identify Protections with automatic Application Layer DDoS Mitigation enabled

```sql+postgres
select
  title,
  resource_arn
from
  aws_shield_protection
where
  application_layer_automatic_response_configuration ->> 'Status' = 'ENABLED'
  and application_layer_automatic_response_configuration -> 'Action' -> 'Block' is not null;
```

```sql+sqlite
select
  title,
  resource_arn
from
  aws_shield_protection
where
  application_layer_automatic_response_configuration ->> 'Status' = 'ENABLED'
  and application_layer_automatic_response_configuration -> 'Action' -> 'Block' is not null;
```

### Check if all Shield protected CloudFront distributions are protected by Shield's automatic Application-Layer-DDoS-Mitigation

```sql+postgres
select
  protection.name as protection_name,
  distribution.arn,
  distribution.aliases ->> 'Items' as aliases,
  web_acl_id is not null as has_web_acl,
  protection.application_layer_automatic_response_configuration ->> 'Status' = 'ENABLED' as auto_mitigation_enabled
from
  aws_shield_protection as protection
join
  aws_cloudfront_distribution as distribution
on
  protection.resource_arn = distribution.arn;
```

```sql+sqlite
select
  protection.name as protection_name,
  distribution.arn,
  distribution.aliases ->> 'Items' as aliases,
  web_acl_id is not null as has_web_acl,
  protection.application_layer_automatic_response_configuration ->> 'Status' = 'ENABLED' as auto_mitigation_enabled
from
  aws_shield_protection as protection
join
  aws_cloudfront_distribution as distribution
on
  protection.resource_arn = distribution.arn;
```
