# Table: aws_ec2_launch_configuration

A launch configuration is a template that an EC2 Auto Scaling group uses to launch EC2 instances

## Examples

### Basic AMI info

```sql
select
  name,
  created_time,
  associate_public_ip_address,
  ebs_optimized,
  image_id,
  instance_monitoring_enabled,
  instance_type,
  key_name
from
  aws_ec2_launch_configuration;
```


### IAM role attached to each Launch Configurations

```sql
select
  name,
  iam_instance_profile
from
  aws_ec2_launch_configuration;
```


### List Launch Configurations with public IPs

```sql
select
  name,
  associate_public_ip_address
from
  aws_ec2_launch_configuration
where
  associate_public_ip_address;
```


### Security groups attached to each Launch Configuration

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as security_groups
from
  aws_ec2_launch_configuration;
```