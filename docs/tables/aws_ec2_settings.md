# Table: aws_ec2_settings

Settings associated with the AWS Account.

## Examples

### Basic settings info

```sql
select
  default_ebs_encryption_enabled,
  default_ebs_encryption_key,
  title,
  region
from
  aws_ec2_settings;
```


### Settings info for a particular region

```sql
select
  default_ebs_encryption_enabled,
  default_ebs_encryption_key,
  title,
  region
from
  aws_ec2_settings
where
  region='ap-south-1';
```


### List the regions along with the key where default EBS encryption is enabled

```sql
select
  region,
  default_ebs_encryption_enabled,
  default_ebs_encryption_key
from
  aws_ec2_settings
where
  default_ebs_encryption_enabled='true';
```
