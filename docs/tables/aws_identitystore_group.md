# Table: aws_identitystore_group

Contains a specified group’s metadata and attributes. Queries to this table must include the `identity_store_id` and either the `name` or `id` columns.

## Examples

### Get group by ID

```sql
select
  id,
  name
from
  aws_identitystore_group
where identity_store_id = 'd-1234567890' and id = '1234567890-12345678-abcd-abcd-abcd-1234567890ab';
```

### List groups by name

```sql
select
  id,
  name
from
  aws_identitystore_group
where identity_store_id = 'd-1234567890' and name = 'test';
```
