---
title: "Table: aws_ec2_launch_template_version - Query AWS EC2 Launch Template Versions using SQL"
description: "Allows users to query AWS EC2 Launch Template Versions, providing details about each version of an Amazon EC2 launch template."
---

# Table: aws_ec2_launch_template_version - Query AWS EC2 Launch Template Versions using SQL

The `aws_ec2_launch_template_version` table in Steampipe provides information about each version of an Amazon EC2 launch template. This table allows DevOps engineers, system administrators, and other IT professionals to query version-specific details, including the template ID, version number, and associated metadata. Users can utilize this table to gather insights on EC2 launch template versions, such as tracking changes between versions, verifying configuration details, and more. The schema outlines the various attributes of the EC2 launch template version, including the template ID, version number, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_launch_template_version` table, you can use the `.inspect aws_ec2_launch_template_version` command in Steampipe.

### Key columns:

- `launch_template_id`: This is the ID of the launch template. It is useful for joining with other tables that reference EC2 launch templates.
- `version_number`: This is the version number of the launch template. It is essential for tracking changes and differentiating between different versions of the same template.
- `created_by`: This column indicates the AWS account that created the launch template version. It can be used to track and manage access and usage across different AWS accounts.

## Examples

### Basic info

```sql
select
  launch_template_name,
  launch_template_id,
  created_time,
  created_by,
  default_version,
  version_description,
  version_number
from
  aws_ec2_launch_template_version;
```

### List launch template versions created by a user

```sql
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

```sql
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

### List default version launch templates

```sql
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

### Count versions by launch template

```sql
select
  launch_template_id,
  count(version_number) as number_of_versions
from
  aws_ec2_launch_template_version
group by
  launch_template_id;
```

### Get launch template data details of each version

```sql
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

### List launch template versions where instance is optimized for Amazon EBS I/O

```sql
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

### List launch template versions where instance termination is restricted via console, CLI, or API

```sql
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

### List template versions where instance stop protection is enabled

```sql
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
