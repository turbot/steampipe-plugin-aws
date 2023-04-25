# Table: aws_wellarchitected_lens

Lenses provide a way for you to consistently measure your architectures against best practices and identify areas for improvement. The AWS Well-Architected Framework Lens is automatically applied when a workload is defined.

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
