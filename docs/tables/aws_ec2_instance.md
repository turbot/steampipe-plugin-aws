# Table: aws_ec2_instance

An AWS EC2 instance is a virtual server in the AWS cloud.

## Examples

### Instance count in each availability zone

```sql
select
  placement_availability_zone as az,
  instance_type,
  count(*)
from
  aws_ec2_instance
group by
  placement_availability_zone,
  instance_type;
```

### List instances whose detailed monitoring is not enabled

```sql
select
  instance_id,
  monitoring_state
from
  aws_ec2_instance
where
  monitoring_state = 'disabled';
```

### Count the number of instances by instance type

```sql
select
  instance_type,
  count(instance_type) as count
from
  aws_ec2_instance
group by
  instance_type;
```

### List instances stopped for more than 30 days

```sql
select
  instance_id,
  instance_state,
  launch_time,
  state_transition_time
from
  aws_ec2_instance
where
  instance_state = 'stopped'
  and state_transition_time <= (current_date - interval '30' day);
```

### List of instances without application tag key

```sql
select
  instance_id,
  tags
from
  aws_ec2_instance
where
  not tags :: JSONB ? 'application';
```

### Get maintenance option for each instance

```sql
select
  instance_id,
  instance_state,
  launch_time,
  maintenance_options ->> 'AutoRecovery' as auto_recovery
from
  aws_ec2_instance;
```

### Get license details for each instance

```sql
select
  instance_id,
  instance_type,
  instance_state,
  l ->> 'LicenseConfigurationArn' as license_configuration_arn
from
  aws_ec2_instance,
  jsonb_array_elements(licenses) as l;
```

### List instances by placement group

```sql
select
  instance_id,
  instance_state,
  launch_time,
  placement_group_id,
  placement_group_name
from
  aws_ec2_instance
where
  placement_group_id = 'pg-015f18b75049eaa0d';
```

### List of EC2 instances provisioned with undesired(for example t2.large and m3.medium is desired) instance type(s).

```sql
select
  instance_type,
  count(*) as count
from
  aws_ec2_instance
where
  instance_type not in ('t2.large', 'm3.medium')
group by
  instance_type;
```

### List EC2 instances having termination protection safety feature enabled

```sql
select
  instance_id,
  disable_api_termination
from
  aws_ec2_instance
where
  not disable_api_termination;
```

### Find instances which have default security group attached

```sql
select
  instance_id,
  sg ->> 'GroupId' as group_id,
  sg ->> 'GroupName' as group_name
from
  aws_ec2_instance
  cross join jsonb_array_elements(security_groups) as sg
where
  sg ->> 'GroupName' = 'default';
```

### List the unencrypted volumes attached to the instances

```sql
select
  i.instance_id,
  vols -> 'Ebs' ->> 'VolumeId' as vol_id,
  vol.encrypted
from
  aws_ec2_instance as i
  cross join jsonb_array_elements(block_device_mappings) as vols
  join aws_ebs_volume as vol on vol.volume_id = vols -> 'Ebs' ->> 'VolumeId'
where
  not vol.encrypted;
```

### List instances with secrets in user data

```sql
select
  instance_id,
  user_data
from
  aws_ec2_instance
where
  user_data like any (array ['%pass%', '%secret%','%token%','%key%'])
  or user_data ~ '(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]';
```

### Get launch template data for the instances

```sql
select
  instance_id,
  launch_template_data -> 'ImageId' as image_id,
  launch_template_data -> 'Placement' as placement,
  launch_template_data -> 'DisableApiStop' as disable_api_stop,
  launch_template_data -> 'MetadataOptions' as metadata_options,
  launch_template_data -> 'NetworkInterfaces' as network_interfaces,
  launch_template_data -> 'BlockDeviceMappings' as block_device_mappings,
  launch_template_data -> 'CapacityReservationSpecification' as capacity_reservation_specification
from
  aws_ec2_instance;
```