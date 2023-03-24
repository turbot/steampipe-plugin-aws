# Table: aws_networkfirewall_firewall

The firewall defines the configuration settings for an AWS Network Firewall firewall. The settings include the firewall policy, the subnets in your VPC to use for the firewall endpoints, and any tags that are attached to the firewall AWS resource.

## Examples

### Basic info

```sql
select
  arn,
  name,
  description,
  vpc_id,
  policy_arn,
  region,
  tags
from
  aws_networkfirewall_firewall;
```

### List firewalls using default encryption

```sql
select
  arn,
  name,
  description,
  encryption_configuration
from
  aws_networkfirewall_firewall
where 
  encryption_configuration ->> 'Type' = 'AWS_OWNED_KMS_KEY';
```

### List firewalls having deletion protection disabled

```sql
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not delete_protection;
```

### List firewalls having policy change protection disabled

```sql
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not policy_change_protection;
```

### List firewalls having subnet change protection disabled

```sql
select
  arn,
  name,
  description,
  vpc_id
from
  aws_networkfirewall_firewall
where
  not subnet_change_protection;
```

### Get subnet details for each firewall

```sql
select
  f.arn,
  f.name,
  f.vpc_id,
  s ->> 'SubnetId' as subnet_id,
  cidr_block,
  availability_zone,
  default_for_az
from
  aws_networkfirewall_firewall f,
  jsonb_array_elements(subnet_mappings) s,
  aws_vpc_subnet vs
where
  vs.subnet_id = s ->> 'SubnetId';
```

### Get KMS key details of firewalls encrypted with customer managed keys

```sql
select
  f.arn,
  f.name,
  f.vpc_id,
  k.arn as key_arn,
  key_rotation_enabled
from
  aws_networkfirewall_firewall f,
  aws_kms_key k
where
  k.id = encryption_configuration ->> 'KeyId'
  and not f.encryption_configuration ->> 'Type' = 'AWS_OWNED_KMS_KEY';
```

