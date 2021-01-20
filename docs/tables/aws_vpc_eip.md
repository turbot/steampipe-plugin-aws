# Table: aws_vpc_eip

An Elastic IP address is a static, public IPv4 address designed for dynamic cloud computing.

## Examples

### List of unused elastic IPs

```sql
select
  public_ip,
  domain association_id
from
  aws_vpc_eip
where
  association_id is null;
```


### Count of elastic IPs by instance Ids

```sql
select
  public_ipv4_pool,
  count(public_ip) as elastic_ips
from
  aws_vpc_eip
group by
  public_ipv4_pool;
```
