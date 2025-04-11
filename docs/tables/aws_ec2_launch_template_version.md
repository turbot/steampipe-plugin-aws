---
title: "Steampipe Table: aws_ec2_launch_template_version - Query AWS EC2 Launch Template Versions using SQL"
description: "Allows users to query AWS EC2 Launch Template Versions, providing details about each version of an Amazon EC2 launch template."
folder: "EC2"
---

# Table: aws_ec2_launch_template_version - Query AWS EC2 Launch Template Versions using SQL

An AWS EC2 Launch Template Version is a configuration template that helps you avoid the trouble of specifying the same instance configuration details every time you launch an instance. It includes information like the ID of the Amazon Machine Image (AMI), the instance type, key pair, security groups, and the other parameters you typically provide when launching an instance. By using versions of launch templates, you can create different configurations changeable over time without altering the original template.

## Table Usage Guide

The `aws_ec2_launch_template_version` table in Steampipe provides you with information about each version of an Amazon EC2 launch template. This table allows you, as a DevOps engineer, system administrator, or other IT professional, to query version-specific details, including the template ID, version number, and associated metadata. You can utilize this table to gather insights on EC2 launch template versions, such as tracking changes between versions, verifying configuration details, and more. The schema outlines the various attributes of the EC2 launch template version for you, including the template ID, version number, creation date, and associated tags.

## Examples

### Basic info
Explore which EC2 launch templates are being used in your AWS environment, including details such as who created them and the default versions. This can help you gain insights into your AWS EC2 usage patterns and streamline your resource management.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  created_by,
  default_version,
  version_description,
  version_number
from
  aws_ec2_launch_template_version;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  created_by,
  default_version,
  version_description,
  version_number
from
  aws_ec2_launch_template_version;
```

### List launch template versions created by a user
Determine the instances where a specific user has created versions of a launch template. This can be useful for understanding user activity and maintaining security and consistency within your AWS EC2 environment.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  create_time,
  created_by,
  version_description,
  version_number
from
  aws_ec2_launch_template_version
where
  created_by like '%turbot';
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  create_time,
  created_by,
  version_description,
  version_number
from
  aws_ec2_launch_template_version
where
  created_by like '%turbot';
```

### List launch template versions created in the last 30 days
Explore which launch template versions have been created in the past 30 days to maintain a current understanding of your AWS EC2 environment. This query is useful for tracking recent changes and ensuring the latest configurations are being utilized.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  create_time,
  default_version,
  version_number
from
  aws_ec2_launch_template_version
where
  create_time >= now() - interval '30' day;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  create_time,
  default_version,
  version_number
from
  aws_ec2_launch_template_version
where
  create_time >= datetime('now', '-30 day');
```

### List default version launch templates
Determine the default versions of your launch templates to understand which configurations are set as standard when launching new instances. This can be helpful for maintaining consistency across your deployments.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  create_time,
  default_version,
  version_number
from
  aws_ec2_launch_template_version
where
  default_version;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  create_time,
  default_version,
  version_number
from
  aws_ec2_launch_template_version
where
  default_version = 1;
```

### Count versions by launch template
Assess the elements within each AWS EC2 launch template to understand the total number of versions associated with each. This can be helpful in managing and tracking the evolution of your launch templates.

```sql+postgres
select
  launch_template_id,
  count(version_number) as number_of_versions
from
  aws_ec2_launch_template_version
group by
  launch_template_id;
```

```sql+sqlite
select
  launch_template_id,
  count(version_number) as number_of_versions
from
  aws_ec2_launch_template_version
group by
  launch_template_id;
```

### Get launch template data details of each version
Identify instances where detailed information about each launch template version is required. This can be useful for understanding and managing the different settings and configurations associated with each version.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  version_number,
  launch_template_data -> 'BlockDeviceMappings' as block_device_mappings,
  launch_template_data -> 'CapacityReservationSpecification' as capacity_reservation_specification,
  launch_template_data -> 'CpuOptions' as cpu_options,
  launch_template_data -> 'CreditSpecification' as credit_specification,
  launch_template_data -> 'DisableApiStop' as disable_api_stop,
  launch_template_data -> 'DisableApiTermination' as disable_api_termination,
  launch_template_data -> 'EbsOptimized' as ebs_optimized,
  launch_template_data -> 'ElasticGpuSpecifications' as elastic_gpu_specifications,
  launch_template_data -> 'ElasticInferenceAccelerators' as elastic_inference_accelerators,
  launch_template_data -> 'EnclaveOptions' as enclave_options,
  launch_template_data -> 'IamInstanceProfile' as iam_instance_profile,
  launch_template_data -> 'ImageId' as image_id,
  launch_template_data -> 'InstanceInitiatedShutdownBehavior' as instance_initiated_shutdown_behavior,
  launch_template_data -> 'InstanceRequirements' as instance_requirements,
  launch_template_data -> 'InstanceType' as instance_type,
  launch_template_data -> 'KernelId' as kernel_id,
  launch_template_data -> 'LicenseSpecifications' as license_specifications,
  launch_template_data -> 'MaintenanceOptions' as maintenance_options,
  launch_template_data -> 'MetadataOptions' as metadata_options,
  launch_template_data -> 'Monitoring' as monitoring,
  launch_template_data -> 'NetworkInterfaces' as network_interfaces,
  launch_template_data -> 'PrivateDnsNameOptions' as private_dns_name_options,
  launch_template_data -> 'RamDiskId' as ram_disk_id,
  launch_template_data -> 'SecurityGroupIds' as security_group_ids,
  launch_template_data -> 'SecurityGroups' as security_groups,
  launch_template_data -> 'TagSpecifications' as tag_specifications,
  launch_template_data -> 'UserData' as user_data
from
  aws_ec2_launch_template_version;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  version_number,
  json_extract(launch_template_data, '$.BlockDeviceMappings') as block_device_mappings,
  json_extract(launch_template_data, '$.CapacityReservationSpecification') as capacity_reservation_specification,
  json_extract(launch_template_data, '$.CpuOptions') as cpu_options,
  json_extract(launch_template_data, '$.CreditSpecification') as credit_specification,
  json_extract(launch_template_data, '$.DisableApiStop') as disable_api_stop,
  json_extract(launch_template_data, '$.DisableApiTermination') as disable_api_termination,
  json_extract(launch_template_data, '$.EbsOptimized') as ebs_optimized,
  json_extract(launch_template_data, '$.ElasticGpuSpecifications') as elastic_gpu_specifications,
  json_extract(launch_template_data, '$.ElasticInferenceAccelerators') as elastic_inference_accelerators,
from
  aws_ec2_launch_template_version;
```

### List launch template versions where instance is optimized for Amazon EBS I/O
Determine the versions of launch templates that are optimized for Amazon EBS I/O. This is useful for identifying instances that are designed for high-performance EBS operations.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  version_number,
  version_description,
  ebs_optimized
from
  aws_ec2_launch_template_version
where
  ebs_optimized;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  version_number,
  version_description,
  ebs_optimized
from
  aws_ec2_launch_template_version
where
  ebs_optimized = 1;
```

### List launch template versions where instance termination is restricted via console, CLI, or API
Determine the areas in which instance termination is restricted for various versions of launch templates. This is useful to ensure that vital instances are safeguarded from accidental termination via console, CLI, or API.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  version_number,
  version_description,
  disable_api_termination
from
  aws_ec2_launch_template_version
where
  disable_api_termination;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  version_number,
  version_description,
  disable_api_termination
from
  aws_ec2_launch_template_version
where
  disable_api_termination = 1;
```

### List template versions where instance stop protection is enabled
Identify versions of launch templates where the protection against instance stops is activated. This is useful for managing and safeguarding critical instances within your AWS EC2 environment.

```sql+postgres
select
  launch_template_name,
  launch_template_id,
  version_number,
  disable_api_stop
from
  aws_ec2_launch_template_version
where
  disable_api_stop;
```

```sql+sqlite
select
  launch_template_name,
  launch_template_id,
  version_number,
  disable_api_stop
from
  aws_ec2_launch_template_version
where
  disable_api_stop = 1;
```