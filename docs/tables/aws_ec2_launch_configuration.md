# Table: aws_ec2_launch_configuration

A launch configuration is a template that an EC2 Auto Scaling group uses to launch EC2 instances

## Examples

### Basic launch configuration info

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

### Get IAM role attached to each launch configuration

```sql
select
  name,
  iam_instance_profile
from
  aws_ec2_launch_configuration;
```

### List launch configurations with public IPs

```sql
select
  name,
  associate_public_ip_address
from
  aws_ec2_launch_configuration
where
  associate_public_ip_address;
```

### Security groups attached to each launch configuration

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as security_groups
from
  aws_ec2_launch_configuration;
```

### List launch configurations with secrets in user data

```sql
select
  name,
  user_data
from
  aws_ec2_launch_configuration
where
  user_data like any (array ['%pass%', '%secret%','%token%','%key%'])
  or user_data ~ '(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]';
```
