---
title: "Steampipe Table: aws_ram_principal_association - Query AWS RAM Principal Associations using SQL"
description: "Allows users to query AWS RAM Principal Associations. The `aws_ram_principal_association` table in Steampipe provides information about principal associations within AWS Resource Access Manager (RAM). This table allows DevOps engineers to query principal-specific details, including resource share ARN, principal ARN, creation time, and associated tags. Users can utilize this table to gather insights on principal associations, such as their status, external status, and more. The schema outlines the various attributes of the principal association, including the resource share ARN, principal ARN, creation time, and associated tags."
folder: "Resource Access Manager"
---

# Table: aws_ram_principal_association - Query AWS RAM Principal Associations using SQL

The AWS RAM Principal Association is a component of AWS Resource Access Manager (RAM) that enables you to share your resources with any AWS account or within your AWS Organization. It allows you to centrally manage who can access your shared resources, thereby improving the efficiency and security of your cross-account resource sharing. This simplifies the process of sharing your resources while maintaining the existing resource permissions.

## Table Usage Guide

The `aws_ram_principal_association` table in Steampipe provides you with information about principal associations within AWS Resource Access Manager (RAM). This table allows you, as a DevOps engineer, to query principal-specific details, including resource share ARN, principal ARN, creation time, and associated tags. You can utilize this table to gather insights on principal associations, such as their status, external status, and more. The schema outlines the various attributes of the principal association for you, including the resource share ARN, principal ARN, creation time, and associated tags.

## Examples

### Basic info
Explore which AWS Resource Access Manager (RAM) principals are associated with your resources to determine their current status. This could be useful in managing resource permissions and identifying any potential issues.

```sql+postgres
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_principal_association;
```

```sql+sqlite
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_principal_association;
```

### List permissions attached with each principal associated
This query is used to gain insights into the permissions linked with each principal associated in AWS Resource Access Manager. It is useful for reviewing the configuration of access controls and ensuring appropriate permissions are in place.

```sql+postgres
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  p ->> 'Arn' as resource_share_permission_arn,
  p ->> 'Status' as resource_share_permission_status
from
  aws_ram_principal_association,
  jsonb_array_elements(resource_share_permission) p;
```

```sql+sqlite
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  json_extract(p.value, '$.Arn') as resource_share_permission_arn,
  json_extract(p.value, '$.Status') as resource_share_permission_status
from
  aws_ram_principal_association,
  json_each(resource_share_permission) as p;
```

### Get principals that failed association
Identify instances where the association of principals to resources within AWS Resource Access Manager (RAM) has failed. This can be useful in troubleshooting and resolving access issues within your AWS environment.

```sql+postgres
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_principal_association
where
  status = 'FAILED';
```

```sql+sqlite
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_principal_association
where
  status = 'FAILED';
```