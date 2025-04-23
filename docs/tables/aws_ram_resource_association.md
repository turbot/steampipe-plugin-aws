---
title: "Steampipe Table: aws_ram_resource_association - Query AWS RAM Resource Associations using SQL"
description: "Allows users to query AWS RAM Resource Associations to retrieve information about the associations between resources and resource shares."
folder: "Resource Access Manager"
---

# Table: aws_ram_resource_association - Query AWS RAM Resource Associations using SQL

The AWS RAM (Resource Access Manager) Resource Associations allow you to share your resources with any AWS account or within your AWS Organization. It simplifies the sharing process of AWS Transit Gateways, Subnets, and AWS License Manager configurations. This allows you to share resources across accounts to reduce operational overhead.

## Table Usage Guide

The `aws_ram_resource_association` table in Steampipe provides you with information about resource associations within AWS Resource Access Manager (RAM). This table lets you, as a DevOps engineer, query association-specific details, including the associated resource ARN, resource share ARN, association type, and status. You can utilize this table to gather insights on resource associations, such as resources associated with a specific resource share, status of the association, and more. The schema outlines the various attributes of the resource association for you, including the resource ARN, resource share ARN, association type, status, and creation time.

## Examples

### Basic info
Analyze the settings to understand the status and associations of your shared AWS resources. This can help in managing access and resource allocation, ensuring optimal resource utilization.

```sql+postgres
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_resource_association;
```

```sql+sqlite
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_resource_association;
```

### List permissions attached with each shared resource associated
Determine the areas in which shared resources are associated with specific permissions. This is useful for managing access control and ensuring proper resource allocation within your AWS environment.

```sql+postgres
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  p ->> 'Arn' as resource_share_permission_arn,
  p ->> 'Status' as resource_share_permission_status
from
  aws_ram_resource_association,
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
  aws_ram_resource_association,
  json_each(resource_share_permission) as p;
```

### Get resources that failed association
Identify instances where resource sharing has failed within your AWS environment. This can be useful for troubleshooting and maintaining efficient resource allocation.

```sql+postgres
select
  resource_share_name,
  resource_share_arn,
  associated_entity,
  status
from
  aws_ram_resource_association
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
  aws_ram_resource_association
where
  status = 'FAILED';
```