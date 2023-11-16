---
title: "Table: aws_wellarchitected_lens - Query AWS Well-Architected Lens using SQL"
description: "Allows users to query AWS Well-Architected Lens, providing details about each lens such as its name, description, and associated AWS Well-Architected Framework pillars."
---

# Table: aws_wellarchitected_lens - Query AWS Well-Architected Lens using SQL

The `aws_wellarchitected_lens` table in Steampipe provides information about each Well-Architected Lens within AWS Well-Architected Tool. This table allows DevOps engineers to query lens-specific details, including lens name, description, and associated AWS Well-Architected Framework pillars. Users can utilize this table to gather insights on lenses, such as their associated pillars and descriptions. The schema outlines the various attributes of the Well-Architected Lens, including the lens name, lens version, lens status, and associated pillars.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_lens` table, you can use the `.inspect aws_wellarchitected_lens` command in Steampipe.

### Key columns:

- `name`: The name of the lens. This can be used to join this table with other tables for more detailed analysis.
- `version`: The version of the lens. This is useful for tracking changes and updates to the lens over time.
- `status`: The status of the lens. This can be used to filter out lenses that are not currently active.

## Examples

### Basic info

```sql
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

When querying multiple regions in an AWS account, each region will return the AWS provided lenses. To only see unique lenses, please see the example below.

```sql
select distinct
  on(arn) arn,
  lens_name,
  lens_status,
  lens_type
from
  aws_wellarchitected_lens;
```

### List custom lenses that are shared from other AWS accounts

```sql
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

```sql
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

```sql
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

### List lenses owned by this account

```sql
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
