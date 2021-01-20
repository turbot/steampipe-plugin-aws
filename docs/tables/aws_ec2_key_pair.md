# Table: aws_ec2_key_pair

A key pair, consisting of a private key and a public key, is a set of security credentials that is used to prove your identity when connecting to an instance.

## Examples

### Basic keypair info

```sql
select
  key_name,
  key_pair_id,
  region
from
  aws_ec2_key_pair;
```


### List of keypairs without owner tag key

```sql
select
  key_name,
  tags
from
  aws_ec2_key_pair
where
  not tags :: JSONB ? 'owner';
```