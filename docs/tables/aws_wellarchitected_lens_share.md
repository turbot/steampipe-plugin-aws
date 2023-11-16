---
title: "Table: aws_wellarchitected_lens_share - Query AWS Well-Architected Tool Lens Shares using SQL"
description: "Allows users to query Lens Shares in AWS Well-Architected Tool, providing details about shared lenses including the share ARN, share status, and the AWS account ID of the lens owner."
---

# Table: aws_wellarchitected_lens_share - Query AWS Well-Architected Tool Lens Shares using SQL

The `aws_wellarchitected_lens_share` table in Steampipe provides information about Lens Shares within AWS Well-Architected Tool. This table allows cloud architects and developers to query details about shared lenses, including the share ARN, share status, and the AWS account ID of the lens owner. Users can utilize this table to gather insights on shared lenses, such as the status of shared lenses, the AWS account ID of the lens owner, and more. The schema outlines the various attributes of the Lens Share, including the share ARN, share status, and the AWS account ID of the lens owner.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_wellarchitected_lens_share` table, you can use the `.inspect aws_wellarchitected_lens_share` command in Steampipe.

### Key columns:

- `share_arn`: The Amazon Resource Name (ARN) of the shared lens. It is a unique identifier and can be used to join this table with other tables.
- `share_status`: The status of the shared lens. It provides information about the current state of the lens share and can be useful to filter or sort the data based on the share status.
- `lens_owner_account_id`: The AWS account ID of the owner of the lens. It can be used to join this table with other AWS account-related tables.

## Examples

### Basic info

```sql
select
  lens_name,
  lens_alias,
  lens_arn,
  share_id,
  shared_with
from
  aws_wellarchitected_lens_share;
```

### Get lens details of the shared lenses

```sql
select
  s.lens_name,
  l.arn,
  l.lens_status,
  l.lens_type,
  l.owner,
  l.share_invitation_id
from
  aws_wellarchitected_lens_share as s,
  aws_wellarchitected_lens as l
where
  s.lens_arn = l.arn;
```

### List shared lenses that are in pending state

```sql
select
  lens_name,
  lens_alias,
  lens_arn,
  share_id,
  shared_with,
  status
from
  aws_wellarchitected_lens_share
where
  status = 'PENDING';
```
