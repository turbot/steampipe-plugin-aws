---
title: "Steampipe Table: aws_wellarchitected_lens - Query AWS Well-Architected Lens using SQL"
description: "Allows users to query AWS Well-Architected Lens, providing details about each lens such as its name, description, and associated AWS Well-Architected Framework pillars."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_lens - Query AWS Well-Architected Lens using SQL

The AWS Well-Architected Lens is a component of the AWS Well-Architected Framework that provides a set of best practices and strategies to help you optimize your use of AWS resources and services. Each Lens provides a unique perspective, focusing on specific workloads or industry verticals. By using the Lens, you can ensure your architecture aligns with the best practices defined by AWS, improving your system's performance, reliability, and cost-effectiveness.

## Table Usage Guide

The `aws_wellarchitected_lens` table in Steampipe provides you with information about each Well-Architected Lens within AWS Well-Architected Tool. This table allows you, as a DevOps engineer, to query lens-specific details, including lens name, description, and associated AWS Well-Architected Framework pillars. You can utilize this table to gather insights on lenses, such as their associated pillars and descriptions. The schema outlines the various attributes of the Well-Architected Lens for you, including the lens name, lens version, lens status, and associated pillars.

## Examples

### Basic info
Explore the various aspects of your AWS Well-Architected Framework to identify the status, type, and ownership details. This can help you understand the current configuration and effectively manage your resources.

```sql+postgres
select
  lens_name,
  lens_alias,
  arn,
  lens_status,
  lens_type,
  owner
from
  aws_wellarchitected_lens;
```

```sql+sqlite
select
  lens_name,
  lens_alias,
  arn,
  lens_status,
  lens_type,
  owner
from
  aws_wellarchitected_lens;
```

### List unique lenses
Explore unique lenses within your AWS Well-Architected Tool to identify their names, statuses, and types, helping you gain insights into your architectural choices and manage your workloads more effectively.
When querying multiple regions in an AWS account, each region will return the AWS provided lenses. To only see unique lenses, please see the example below.


```sql+postgres
select distinct
  on(arn) arn,
  lens_name,
  lens_status,
  lens_type
from
  aws_wellarchitected_lens;
```

```sql+sqlite
select distinct
  arn,
  lens_name,
  lens_status,
  lens_type
from
  aws_wellarchitected_lens
group by
  arn;
```

### List custom lenses that are shared from other AWS accounts
Discover the segments that feature custom lenses shared from other AWS accounts. This can help to understand the distribution and status of shared custom lenses, providing valuable insights into the usage patterns and ownership within your AWS environment.

```sql+postgres
select
  lens_name,
  arn,
  lens_status,
  lens_type,
  owner,
  share_invitation_id
from
  aws_wellarchitected_lens
where
  lens_type = 'CUSTOM_SHARED';
```

```sql+sqlite
select
  lens_name,
  arn,
  lens_status,
  lens_type,
  owner,
  share_invitation_id
from
  aws_wellarchitected_lens
where
  lens_type = 'CUSTOM_SHARED';
```

### List deprecated lenses
Explore which lenses in your AWS Well-Architected Tool are deprecated. This is useful to ensure you're not using outdated tools in your architecture.

```sql+postgres
select
  lens_name,
  lens_status,
  lens_type,
  lens_version,
  owner
from
  aws_wellarchitected_lens
where
  lens_status = 'DEPRECATED';
```

```sql+sqlite
select
  lens_name,
  lens_status,
  lens_type,
  lens_version,
  owner
from
  aws_wellarchitected_lens
where
  lens_status = 'DEPRECATED';
```

### List lenses created in the last 30 days
Discover the segments that have been recently added in the last month. This is useful for keeping track of new additions and understanding their status, type, and version.

```sql+postgres
select
  lens_name,
  lens_status,
  lens_type,
  created_at,
  lens_version
from
  aws_wellarchitected_lens
where
  created_at <= now() - interval '30' day;
```

```sql+sqlite
select
  lens_name,
  lens_status,
  lens_type,
  created_at,
  lens_version
from
  aws_wellarchitected_lens
where
  created_at <= datetime('now', '-30 day');
```

### List lenses owned by this account
This query is used to identify the lenses that are owned by a specific account in the AWS Well-Architected Framework. It's useful for assessing ownership and managing resources within an AWS environment.

```sql+postgres
select
  lens_name,
  lens_status,
  lens_type,
  lens_version,
  owner,
  account_id
from
  aws_wellarchitected_lens
where
  owner = account_id;
```

```sql+sqlite
select
  lens_name,
  lens_status,
  lens_type,
  lens_version,
  owner,
  account_id
from
  aws_wellarchitected_lens
where
  owner = account_id;
```