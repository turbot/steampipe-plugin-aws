# Table: aws_vpc_verified_access_group

An AWS Verified Access group is a collection of Verified Access endpoints and a group-level Verified Access policy. Each endpoint within a group shares the Verified Access policy. You can use groups to gather together endpoints that have common security requirements. This can help simplify policy administration by using one policy for the security needs of multiple applications.

## Examples

### Basic info

```sql
select
  verified_access_group_id,
  arn,
  verified_access_instance_id,
  creation_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group;
```

### List groups that older than 30 days

```sql
select
  verified_access_group_id,
  creation_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group
where
  creation_time <= now() - interval '30' day;
```

### List active groups

```sql
select
  verified_access_group_id,
  creation_time,
  deletion_time,
  description,
  last_updated_time
from
  aws_vpc_verified_access_group
where
  deletion_time is null;
```

### Get trusted provider details for each group

```sql
select
  g.verified_access_group_id,
  g.creation_time,
  i.creation_time as instance_create_time,
  i.verified_access_instance_id,
  jsonb_pretty(i.verified_access_trust_providers) as verified_access_trust_providers
from
  aws_vpc_verified_access_group as g,
  aws_vpc_verified_access_instance as i
where
  g.verified_access_instance_id = i.verified_access_instance_id;
```