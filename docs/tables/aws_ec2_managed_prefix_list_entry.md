# Table: aws_ec2_managed_prefix_list_entry

AWS EC2 Managed Prefix List Entry is a configuration entry within a managed prefix list in Amazon EC2. A managed prefix list is a virtual firewall that allows or denies traffic based on IP prefixes (CIDR blocks). Each entry within a managed prefix list specifies a specific CIDR block and an associated action (allow or deny) to control traffic flow.

There are two types of prefix lists:

## Examples

### Basic Info

```sql
select
  prefix_list_id,
  cidr,
  description
from
  aws_ec2_managed_prefix_list_entry;
```

### List customer-managed prefix lists entries

```sql
select
  l.name,
  l.id,
  e.cidr,
  e.description,
  l.state,
  l.owner_id
from
  aws_ec2_managed_prefix_list_entry as e,
  aws_ec2_managed_prefix_list as l
where
  l.owner_id <> 'AWS';
```

### Count prefix list entries by prefix list

```sql
select
  prefix_list_id,
  count(cidr) as numbers_of_entries
from
  aws_ec2_managed_prefix_list_entry
group by
  prefix_list_id;
```
