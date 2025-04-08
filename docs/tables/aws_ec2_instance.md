---
title: "Steampipe Table: aws_ec2_instance - Query AWS EC2 Instances using SQL"
description: "Allows users to query AWS EC2 Instances for comprehensive data on each instance, including instance type, state, tags, and more."
folder: "EC2"
---

# Table: aws_ec2_instance - Query AWS EC2 Instances using SQL

The AWS EC2 Instance is a virtual server in Amazon's Elastic Compute Cloud (EC2) for running applications on the Amazon Web Services (AWS) infrastructure. It provides scalable computing capacity in the AWS cloud, eliminating the need to invest in hardware up front, so you can develop and deploy applications faster. With EC2, you can launch as many or as few virtual servers as you need, configure security and networking, and manage storage.

## Table Usage Guide

The `aws_ec2_instance` table in Steampipe provides you with information about EC2 Instances within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query instance-specific details, including instance state, launch time, instance type, and associated metadata. You can utilize this table to gather insights on instances, such as instances with specific tags, instances in a specific state, instances of a specific type, and more. The schema outlines the various attributes of the EC2 instance for you, including the instance ID, instance state, instance type, and associated tags.

## Examples

### Instance count in each availability zone
Discover the distribution of instances across different availability zones and types within your AWS EC2 service. This helps in understanding load balancing and can aid in optimizing resource utilization.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which detailed monitoring is not enabled for your AWS EC2 instances. This is useful for identifying potential blind spots in your system's monitoring coverage.

```sql+postgres
select
  instance_id,
  monitoring_state
from
  aws_ec2_instance
where
  monitoring_state = 'disabled';
```

```sql+sqlite
select
  instance_id,
  monitoring_state
from
  aws_ec2_instance
where
  monitoring_state = 'disabled';
```

### Count the number of instances by instance type
Determine the distribution of your virtual servers based on their configurations, allowing you to assess your resource allocation and optimize your infrastructure management strategy.

```sql+postgres
select
  instance_type,
  count(instance_type) as count
from
  aws_ec2_instance
group by
  instance_type;
```

```sql+sqlite
select
  instance_type,
  count(instance_type) as count
from
  aws_ec2_instance
group by
  instance_type;
```

### List instances stopped for more than 30 days
Determine the areas in which AWS EC2 instances have been stopped for over 30 days. This can be useful for identifying and managing instances that may be unnecessarily consuming resources or costing money.

```sql+postgres
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

```sql+sqlite
select
  instance_id,
  instance_state,
  launch_time,
  state_transition_time
from
  aws_ec2_instance
where
  instance_state = 'stopped'
  and state_transition_time <= date('now', '-30 day');
```

### List of instances without application tag key
Determine the areas in which EC2 instances are lacking the 'application' tag. This is useful to identify instances that may not be following your organization's tagging strategy, ensuring better resource management and cost tracking.

```sql+postgres
select
  instance_id,
  tags
from
  aws_ec2_instance
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  instance_id,
  tags
from
  aws_ec2_instance
where
  json_extract(tags, '$.application') is null;
```

### Get maintenance options for each instance
Determine the status of each instance's automatic recovery feature to plan for potential maintenance needs. This can help in understanding the instances' resilience and ensure uninterrupted services.

```sql+postgres
select
  instance_id,
  instance_state,
  launch_time,
  maintenance_options ->> 'AutoRecovery' as auto_recovery
from
  aws_ec2_instance;
```

```sql+sqlite
select
  instance_id,
  instance_state,
  launch_time,
  json_extract(maintenance_options, '$.AutoRecovery') as auto_recovery
from
  aws_ec2_instance;
```

### Get license details for each instance
Determine the license details associated with each of your instances to better manage and track your licensing agreements. This can help ensure compliance and avoid potential legal issues.

```sql+postgres
select
  instance_id,
  instance_type,
  instance_state,
  l ->> 'LicenseConfigurationArn' as license_configuration_arn
from
  aws_ec2_instance,
  jsonb_array_elements(licenses) as l;
```

```sql+sqlite
select
  instance_id,
  instance_type,
  instance_state,
  json_extract(l.value, '$.LicenseConfigurationArn') as license_configuration_arn
from
  aws_ec2_instance,
  json_each(licenses) as l;
```

### Get placement group details for each instance
This query can be used to gain insights into the geographic distribution and configuration of your AWS EC2 instances. It helps in managing resources efficiently by understanding their placement details such as affinity, availability zone, and tenancy.

```sql+postgres
select
  instance_id,
  instance_state,
  placement_affinity,
  placement_group_id,
  placement_group_name,
  placement_availability_zone,
  placement_host_id,
  placement_host_resource_group_arn,
  placement_partition_number,
  placement_tenancy
from
  aws_ec2_instance;
```

```sql+sqlite
select
  instance_id,
  instance_state,
  placement_affinity,
  placement_group_id,
  placement_group_name,
  placement_availability_zone,
  placement_host_id,
  placement_host_resource_group_arn,
  placement_partition_number,
  placement_tenancy
from
  aws_ec2_instance;
```

### List of EC2 instances provisioned with undesired(for example t2.large and m3.medium is desired) instance type(s).
Identify instances where EC2 instances have been provisioned with types other than the desired ones, such as t2.large and m3.medium. This can help you manage your resources more effectively by spotting any instances that may not meet your specific needs or standards.

```sql+postgres
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

```sql+sqlite
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
Identify instances where the termination protection safety feature is enabled in EC2 instances. This is beneficial for preventing accidental terminations and ensuring system stability.

```sql+postgres
select
  instance_id,
  disable_api_termination
from
  aws_ec2_instance
where
  not disable_api_termination;
```

```sql+sqlite
select
  instance_id,
  disable_api_termination
from
  aws_ec2_instance
where
  disable_api_termination = 0;
```

### Find instances which have default security group attached
Discover the segments that have the default security group attached to them in order to identify potential security risks. This is useful for maintaining optimal security practices and ensuring that instances are not using default settings, which may be more vulnerable.

```sql+postgres
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

```sql+sqlite
select
  instance_id,
  json_extract(sg.value, '$.GroupId') as group_id,
  json_extract(sg.value, '$.GroupName') as group_name
from
  aws_ec2_instance,
  json_each(aws_ec2_instance.security_groups) as sg
where
  json_extract(sg.value, '$.GroupName') = 'default';
```

### List the unencrypted volumes attached to the instances
Identify instances where data storage volumes attached to cloud-based virtual servers are not encrypted. This is useful for enhancing security measures by locating potential vulnerabilities where sensitive data might be exposed.

```sql+postgres
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

```sql+sqlite
select
  i.instance_id,
  json_extract(vols.value, '$.Ebs.VolumeId') as vol_id,
  vol.encrypted
from
  aws_ec2_instance as i,
  json_each(i.block_device_mappings) as vols
  join aws_ebs_volume as vol on vol.volume_id = json_extract(vols.value, '$.Ebs.VolumeId')
where
  not vol.encrypted;
```

### List instances with secrets in user data
Discover the instances that might contain sensitive information in their user data. This is beneficial in identifying potential security risks and ensuring data privacy compliance.

```sql+postgres
select
  instance_id,
  user_data
from
  aws_ec2_instance
where
  user_data like any (array ['%pass%', '%secret%','%token%','%key%'])
  or user_data ~ '(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]';
```

```sql+sqlite
select
  instance_id,
  user_data
from
  aws_ec2_instance
where
  user_data like '%pass%'
  or user_data like '%secret%'
  or user_data like '%token%'
  or user_data like '%key%'
  or (user_data REGEXP '[a-z]' and user_data REGEXP '[A-Z]' and user_data REGEXP '\d' and user_data REGEXP '[@$!%*?&]');
```

### Get launch template data for the instances
Analyze the settings to understand the configuration and specifications of your cloud instances. This can help you assess the elements within your instances, such as network interfaces and capacity reservation specifications, which can be useful for optimizing resource usage and management.

```sql+postgres
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

```sql+sqlite
select
  instance_id,
  json_extract(launch_template_data, '$.ImageId') as image_id,
  json_extract(launch_template_data, '$.Placement') as placement,
  json_extract(launch_template_data, '$.DisableApiStop') as disable_api_stop,
  json_extract(launch_template_data, '$.MetadataOptions') as metadata_options,
  json_extract(launch_template_data, '$.NetworkInterfaces') as network_interfaces,
  json_extract(launch_template_data, '$.BlockDeviceMappings') as block_device_mappings,
  json_extract(launch_template_data, '$.CapacityReservationSpecification') as capacity_reservation_specification
from
  aws_ec2_instance;
```

### Get subnet details for each instance
Explore the association between instances and subnets in your AWS environment. This can be helpful in understanding how resources are distributed and for planning infrastructure changes or improvements.

```sql+postgres
select 
  i.instance_id, 
  i.vpc_id, 
  i.subnet_id, 
  s.tags ->> 'Name' as subnet_name
from 
  aws_ec2_instance as i, 
  aws_vpc_subnet as s 
where 
  i.subnet_id = s.subnet_id;
```

```sql+sqlite
select 
  i.instance_id, 
  i.vpc_id, 
  i.subnet_id, 
  json_extract(s.tags, '$.Name') as subnet_name
from 
  aws_ec2_instance as i, 
  aws_vpc_subnet as s 
where 
  i.subnet_id = s.subnet_id;
```