# Table: aws_ec2_ami

An Amazon Machine Image is a special type of virtual appliance that is used to create a virtual machine within the Amazon Elastic Compute Cloud.

The `aws_ec2_ami` lists only private images. To list public images, or images that are shared with your account, use `aws_ec2_ami_shared`.

## Examples

### Basic AMI info

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

### List of public AMIs

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

### AWS AMI Volume info

```sql
select
  name,
  image_id,
  mapping -> 'Ebs' ->> 'VolumeSize' as volume_size,
  mapping -> 'Ebs' ->> 'VolumeType' as volume_type
from
  aws_ec2_ami
  cross join jsonb_array_elements(block_device_mappings) as mapping;
```

### AWS AMI Block Device Encryption status

```sql
select
  name,
  image_id,
  mapping -> 'Ebs' ->> 'Encrypted' as encryption_status,
  mapping -> 'Ebs' ->> 'KmsKeyId' as kms_key
from
  aws_ec2_ami
  cross join jsonb_array_elements(block_device_mappings) as mapping;
```

### AWS AMI Block Device Deletion Protection info

```sql
select
  name,
  image_id,
  mapping -> 'Ebs' ->> 'DeleteOnTermination' as delete_on_termination
from
  aws_ec2_ami
  cross join jsonb_array_elements(block_device_mappings) as mapping;
```
