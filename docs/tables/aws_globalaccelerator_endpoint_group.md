# Table: aws_globalaccelerator_endpoint_group

Each endpoint group is associated with a specific AWS Region. Endpoint groups include one or more endpoints in the
Region. With a standard accelerator, you can increase or reduce the percentage of traffic that would be otherwise
directed to an endpoint group by adjusting a setting called a traffic dial. The traffic dial lets you easily do
performance testing or blue/green deployment testing, for example, for new releases across different AWS Regions.

## Examples

### Basic info

```sql
select
  title,
  endpoint_descriptions,
  endpoint_group_region,
  traffic_dial_percentage,
  port_overrides,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  threshold_count
from
  aws_globalaccelerator_endpoint_group;
```

### List endpoint groups for a specific listener

```sql
select
  title,
  endpoint_descriptions,
  endpoint_group_region,
  traffic_dial_percentage,
  port_overrides,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  threshold_count
from
  aws_globalaccelerator_endpoint_group
where
  listener_arn = 'arn:aws:globalaccelerator::012345678901:accelerator/1234abcd-abcd-1234-abcd-1234abcdefgh/listener/abcdef1234';
```

### Get basic info for all accelerators, listeners, and endpoint groups

```sql
select
  a.name as accelerator_name,
  l.client_affinity as listener_client_affinity,
  l.port_ranges as listener_port_ranges,
  l.protocol as listener_protocol,
  eg.endpoint_descriptions,
  eg.endpoint_group_region,
  eg.traffic_dial_percentage,
  eg.port_overrides,
  eg.health_check_interval_seconds,
  eg.health_check_path,
  eg.health_check_port,
  eg.health_check_protocol,
  eg.threshold_count
from
  aws_globalaccelerator_accelerator a,
  aws_globalaccelerator_listener l,
  aws_globalaccelerator_endpoint_group eg
where
  eg.listener_arn = l.arn
  and l.accelerator_arn = a.arn;
```
