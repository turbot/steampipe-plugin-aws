---
title: "Table: aws_ec2_instance_type - Query AWS EC2 Instance Type using SQL"
description: "Allows users to query AWS EC2 Instance Type data, including details about instance type name, current generation, vCPU, memory, storage, and network performance."
---

# Table: aws_ec2_instance_type - Query AWS EC2 Instance Type using SQL

The `aws_ec2_instance_type` table in Steampipe provides information about EC2 instance types within AWS Elastic Compute Cloud (EC2). This table allows DevOps engineers to query instance type-specific details, including its name, current generation, vCPU, memory, storage, and network performance. Users can utilize this table to gather insights on instance types, such as their capabilities, performance, and associated metadata. The schema outlines the various attributes of the EC2 instance type, including the instance type, current generation, vCPU, memory, storage, and network performance.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_instance_type` table, you can use the `.inspect aws_ec2_instance_type` command in Steampipe.

### Key columns:

- `instance_type`: This column contains the name of the instance type. It can be used to join this table with other tables that contain instance type data.
- `vcpu`: This column contains information about the number of virtual central processing units (vCPUs) that the instance type supports. It can be used to join this table with other tables that contain vCPU data to analyze compute capabilities.
- `memory`: This column contains information about the memory capacity of the instance type. It can be used to join this table with other tables that contain memory data to analyze memory capabilities.

## Examples

### List of instance types which supports dedicated host

```sql
select
  instance_type,
  dedicated_hosts_supported
from
  aws_ec2_instance_type
where
  dedicated_hosts_supported;
```


### List of instance types which does not support auto recovery

```sql
select
  instance_type,
  auto_recovery_supported
from
  aws_ec2_instance_type
where
  not auto_recovery_supported;
```


### List of instance types which have more than 24 cores

```sql
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


### List of instance types which does not support encryption to root volume

```sql
select
  instance_type,
  ebs_info ->> 'EncryptionSupport' as encryption_support
from
  aws_ec2_instance_type
where
  ebs_info ->> 'EncryptionSupport' = 'unsupported';
```


### List of instance types eligible for free tier

```sql
select
  instance_type,
  free_tier_eligible
from
  aws_ec2_instance_type
where
  free_tier_eligible;
```