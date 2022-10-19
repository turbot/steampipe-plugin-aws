# Table: aws_lightsail_instance

Amazon Lightsail is a service to provide easy virtual private servers with custom software already setup.

## Examples

### Instance count in each availability zone

```sql
select
  availability_zone as az,
  bundle_id,
  count(*)
from
  aws_lightsail_instance
group by
  availability_zone,
  bundle_id;
```

### List stopped instances created for more than 30 days

```sql
select
  name,
  state_name
from
  aws_lightsail_instance
where
  state_name = 'stopped'
  and created_at <= (current_date - interval '30' day);
```

### List public instances

```sql
select
  name,
  state_name,
  bundle_id,
  region
from
  aws_lightsail_instance
where
  public_ip_address is not null;
```

### List of instances without application tag key

```sql
select
  name,
  tags
from
  aws_lightsail_instance
where
  not tags :: JSONB ? 'application';
```

### Hardware specifications of the instances

```sql
select
  name,
  hardware ->> 'CpuCount' as "CPU Count",
  hardware ->> 'RamSizeInGb' as "RAM Size (in GB)"
from
  aws_lightsail_instance;
```