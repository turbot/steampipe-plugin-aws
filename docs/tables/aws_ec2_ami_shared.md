# Table: aws_ec2_ami_shared

An Amazon Machine Image is a special type of virtual appliance that is used to create a virtual machine within the Amazon Elastic Compute Cloud.

**You must specify an Owner ID or Image ID** in a `where` clause (`where owner_id='`), (`where image_id='`).

The `aws_ec2_ami_shared` table can list any image but you must specify `owner_id` or `image_id`.
If you want to list all of the images in your account then you can use the `aws_ec2_ami` table.

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
  aws_ec2_ami_shared
where
  owner_id = '137112412989';
```

### List EC2 instances using AMIs owned by a specific AWS account

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
  ami.owner_id = '137112412989';
```
