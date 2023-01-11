# Table: aws_ec2_ami_shared

An Amazon Machine Image is a special type of virtual appliance that is used to create a virtual machine within the Amazon Elastic Compute Cloud.

The `aws_ec2_ami_shared` table only lists public and shared images. To list private images, use the `aws_ec2_ami` table.

**You must specify an owner ID or Image ID** in a `where` clause (`where owner_id='`), (`where image_id='`).

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
