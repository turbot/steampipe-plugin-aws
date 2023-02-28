# Table: aws_ebs_volume

An Amazon EBS volume is a durable, block-level storage device that you can attach to your instances.

## Examples

### List of unencrypted EBS volumes

```sql
select
  volume_id,
  encrypted
from
  aws_ebs_volume
where
  not encrypted;
```

### List of unattached EBS volumes

```sql
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  jsonb_array_length(attachments) = 0;
```

### List of Provisioned IOPS SSD (io1) volumes

```sql
select
  volume_id,
  volume_type
from
  aws_ebs_volume
where
  volume_type = 'io1';
```

### List of EBS volumes with size more than 100GiB

```sql
select
  volume_id,
  size
from
  aws_ebs_volume
where
  size > '100';
```

### Count the number of EBS volumes by volume type

```sql
select
  volume_type,
  count(volume_type) as count
from
  aws_ebs_volume
group by
  volume_type;
```

### Find EBS Volumes Attached To Stopped EC2 Instances

```sql
select
  volume_id,
  size,
  att ->> 'InstanceId' as instance_id
from
  aws_ebs_volume
  cross join jsonb_array_elements(attachments) as att
  join aws_ec2_instance as i on i.instance_id = att ->> 'InstanceId'
where
  instance_state = 'stopped';
```
