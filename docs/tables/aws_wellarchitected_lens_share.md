# Table: aws_wellarchitected_lens_share

AWS Well-Architected Framework can share a custom lens with other AWS accounts, users, AWS Organizations, and organization units (OUs)..

## Examples

## Basic info

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

## Get lens details for the lens shares

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

## List lens shares that are in pending state

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