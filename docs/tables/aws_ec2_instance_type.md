---
title: "Steampipe Table: aws_ec2_instance_type - Query AWS EC2 Instance Type using SQL"
description: "Allows users to query AWS EC2 Instance Type data, including details about instance type name, current generation, vCPU, memory, storage, and network performance."
folder: "EC2"
---

# Table: aws_ec2_instance_type - Query AWS EC2 Instance Type using SQL

The AWS EC2 Instance Type is a component of Amazon's Elastic Compute Cloud (EC2), which provides scalable computing capacity in the Amazon Web Services (AWS) cloud. It defines the hardware of the host computer used for the instance. Different instance types offer varying combinations of CPU, memory, storage, and networking capacity, giving you the flexibility to choose the appropriate mix of resources for your applications.

## Table Usage Guide

The `aws_ec2_instance_type` table in Steampipe provides you with information about EC2 instance types within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query instance type-specific details, including its name, current generation, vCPU, memory, storage, and network performance. You can utilize this table to gather insights on instance types, such as their capabilities, performance, and associated metadata. The schema outlines the various attributes of the EC2 instance type for you, including the instance type, current generation, vCPU, memory, storage, and network performance.

**Important Notes**
- This table supports the optional quals `instance_type` and `instance_type_pattern`.
- Queries with optional quals are optimised to use additional filtering provided by the AWS API function.
- To filter by a specific `instance_type`, you need to include it in the WHERE clause, such as `where instance_type = 't2.small'`, to retrieve a single instance type.
- If you want to fetch instance types using a wildcard pattern, you can use `instance_type_pattern` in the WHERE clause, like `where instance_type_pattern = 't2*'`.

## Examples

### List of instance types which supports dedicated host
Explore which AWS EC2 instance types support a dedicated host. This is useful for identifying the types of instances that can be used for tasks requiring dedicated resources, enhancing performance and security.

```sql+postgres
select
  instance_type,
  dedicated_hosts_supported
from
  aws_ec2_instance_type
where
  dedicated_hosts_supported;
```

```sql+sqlite
select
  instance_type,
  dedicated_hosts_supported
from
  aws_ec2_instance_type
where
  dedicated_hosts_supported = 1;
```

### List of instance types which does not support auto recovery
Discover the segments of AWS EC2 instances that do not support auto-recovery. This is useful to identify potential risk areas in your infrastructure that may require manual intervention in case of system failures.

```sql+postgres
select
  instance_type,
  auto_recovery_supported
from
  aws_ec2_instance_type
where
  not auto_recovery_supported;
```

```sql+sqlite
select
  instance_type,
  auto_recovery_supported
from
  aws_ec2_instance_type
where
  auto_recovery_supported = 0;
```

### List of instance types which have more than 24 cores
Determine the areas in which AWS EC2 instance types support dedicated hosts and have more than 24 default cores. This can be useful for identifying high-performance instances suitable for resource-intensive applications.

```sql+postgres
select
  instance_type,
  dedicated_hosts_supported,
  v_cpu_info -> 'DefaultCores' as default_cores,
  v_cpu_info -> 'DefaultThreadsPerCore' as default_threads_per_core,
  v_cpu_info -> 'DefaultVCpus' as default_vcpus,
  v_cpu_info -> 'ValidCores' as valid_cores,
  v_cpu_info -> 'ValidThreadsPerCore' as valid_threads_per_core
from
  aws_ec2_instance_type
where
  v_cpu_info ->> 'DefaultCores' > '24';
```

```sql+sqlite
select
  instance_type,
  dedicated_hosts_supported,
  json_extract(v_cpu_info, '$.DefaultCores') as default_cores,
  json_extract(v_cpu_info, '$.DefaultThreadsPerCore') as default_threads_per_core,
  json_extract(v_cpu_info, '$.DefaultVCpus') as default_vcpus,
  json_extract(v_cpu_info, '$.ValidCores') as valid_cores,
  json_extract(v_cpu_info, '$.ValidThreadsPerCore') as valid_threads_per_core
from
  aws_ec2_instance_type
where
  CAST(json_extract(v_cpu_info, '$.DefaultCores') AS INTEGER) > 24;
```


### List of instance types which does not support encryption to root volume
Identify instances where the type of Amazon EC2 instance does not support encryption for the root volume. This is beneficial for maintaining security standards and ensuring sensitive data is adequately protected.

```sql+postgres
select
  instance_type,
  ebs_info ->> 'EncryptionSupport' as encryption_support
from
  aws_ec2_instance_type
where
  ebs_info ->> 'EncryptionSupport' = 'unsupported';
```

```sql+sqlite
select
  instance_type,
  json_extract(ebs_info, '$.EncryptionSupport') as encryption_support
from
  aws_ec2_instance_type
where
  json_extract(ebs_info, '$.EncryptionSupport') = 'unsupported';
```

### List of instance types eligible for free tier
Determine the types of instances that are eligible for the free tier in AWS EC2, aiding in cost-efficient resource allocation.

```sql+postgres
select
  instance_type,
  free_tier_eligible
from
  aws_ec2_instance_type
where
  free_tier_eligible;
```

```sql+sqlite
select
  instance_type,
  free_tier_eligible
from
  aws_ec2_instance_type
where
  free_tier_eligible = 1;
```