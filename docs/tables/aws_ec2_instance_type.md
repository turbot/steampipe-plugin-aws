# Table: aws_ec2_instance_type

Amazon EC2 provides a wide selection of instance types optimized to fit different use cases. Instance types comprise varying combinations of CPU, memory, storage, and networking capacity and give you the flexibility to choose the appropriate mix of resources for your applications.

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