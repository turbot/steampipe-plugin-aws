# Table: aws_ec2_ami_shared

An Amazon Machine Image is a special type of virtual appliance that is used to create a virtual machine within the Amazon Elastic Compute Cloud.

The `aws_ec2_ami` lists only private images. To list public images, or images that are shared with your account, use `aws_ec2_ami_shared`.

This table requires an '=' qualifier for column: `owner_id`.

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
  aws_ec2_ami_shared
where
  owner_id = '137112412989';
```

### AWS instances using images from a specific AWS account

```sql
select
  i.title,
  i.instance_id,
  i.image_id,
  ami.name,
  ami.description,
  ami.platform_details
from
  aws_ec2_instance as i
 join aws_ec2_ami_shared as ami on i.image_id = ami.image_id
where
  ami.owner_id = '137112412989'
```
