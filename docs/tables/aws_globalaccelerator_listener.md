# Table: aws_globalaccelerator_listener

A listener processes inbound connections from clients to Global Accelerator, based on the port (or port range) and
protocol (or protocols) that you configure. A listener can be configured for TCP, UDP, or both TCP and UDP protocols.
Each listener has one or more endpoint groups associated with it, and traffic is forwarded to endpoints in one of
the groups. You associate endpoint groups with listeners by specifying the Regions that you want to distribute traffic
to. With a standard accelerator, traffic is distributed to optimal endpoints within the endpoint groups associated
with a listener.

## Examples

### Basic info

```sql
select
  title,
  client_affinity,
  port_ranges,
  protocol
from
  aws_globalaccelerator_listener;
```

### List listeners for a specific accelerator

```sql
select
  title,
  client_affinity,
  port_ranges,
  protocol
from
  aws_globalaccelerator_listener
where
  accelerator_arn = 'arn:aws:globalaccelerator::012345678901:accelerator/1234abcd';
```

### Basic info for all accelerators and listeners

```sql
select
  a.name as accelerator_name,
  a.status as accelerator_status,
  l.title as listener_title,
  l.client_affinity as listener_client_affinity,
  l.port_ranges as listener_port_ranges,
  l.protocol as listener_protocol
from
  aws_globalaccelerator_accelerator a,
  aws_globalaccelerator_listener l
where
  l.accelerator_arn = a.arn;
```

### List accelerators listening on TCP port 443

```sql
select
  a.name as accelerator_name,
  a.status as accelerator_status,
  l.protocol,
  port_range -> 'FromPort' as from_port,
  port_range -> 'ToPort' as to_port
from
  aws_globalaccelerator_accelerator a,
  aws_globalaccelerator_listener l,
  jsonb_array_elements(l.port_ranges) as port_range
where
  l.accelerator_arn = a.arn
  and l.protocol = 'TCP'
  and (port_range -> 'FromPort')::int <= 443
  and (port_range -> 'ToPort')::int >= 443;
```
