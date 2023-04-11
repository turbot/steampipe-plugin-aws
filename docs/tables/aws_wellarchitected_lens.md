# Table: aws_wellarchitected_lens

Lenses provide a way for you to consistently measure your architectures against best practices and identify areas for improvement. The AWS Well-Architected Framework Lens is automatically applied when a workload is defined.

## Examples

## Basic info

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

## List lenses that are of type custom shared

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

## List deprecated lenses

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

## List lenses that are created in the last 30 days

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
