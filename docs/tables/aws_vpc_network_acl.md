# Table: aws_vpc_network_acl

A network access control list (ACL) is an optional layer of security for your VPC that acts as a firewall for controlling traffic in and out of one or more subnets.

## Examples

### List the attached VPC IDs for each network ACL

```sql
select
  network_acl_id,
  arn,
  vpc_id
from
  aws_vpc_network_acl;
```


### List the default NACL associated with the VPCs

```sql
select
  network_acl_id,
  vpc_id,
  is_default
from
  aws_vpc_network_acl
where
  is_default = true;
```


### Subnet associated with each network ACL

```sql
select
  network_acl_id,
  vpc_id,
  association ->> 'SubnetId' as subnet_id,
  association ->> 'NetworkAclAssociationId' as network_acl_association_id
from
  aws_vpc_network_acl
  cross join jsonb_array_elements(associations) as association;
```