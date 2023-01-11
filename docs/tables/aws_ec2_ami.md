# Table: aws_ec2_ami

An Amazon Machine Image is a special type of virtual appliance that is used to create a virtual machine within the Amazon Elastic Compute Cloud.

The `aws_ec2_ami` table only lists images in your account. To list other images shared with you, please use the `aws_ec2_ami_shared` table.

## Examples

### Basic info

```sql
select
  name,
  image_id,
  state,
  image_location,
  creation_date,
  public,
  root_device_name
from
  aws_ec2_ami;
```

### List public AMIs

```sql
select
  name,
  image_id,
  public
from
  aws_ec2_ami
where
  public;
```

### Get volume info for each AMI

```sql
select
  name,
  image_id,
  mapping -> 'Ebs' ->> 'VolumeSize' as volume_size,
  mapping -> 'Ebs' ->> 'VolumeType' as volume_type,
  mapping -> 'Ebs' ->> 'Encrypted' as encryption_status,
  mapping -> 'Ebs' ->> 'KmsKeyId' as kms_key,
  mapping -> 'Ebs' ->> 'DeleteOnTermination' as delete_on_termination
from
  aws_ec2_ami
  cross join jsonb_array_elements(block_device_mappings) as mapping;
```
