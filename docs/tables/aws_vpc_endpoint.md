# Table: aws_vpc_endpoint

A VPC endpoint enables private connections between your VPC and supported AWS services and VPC endpoint services powered by AWS PrivateLink.

## Examples

### List of VPC endpoint and the corresponding services

```sql
select
  vpc_endpoint_id,
  vpc_id,
  service_name
from
  aws_vpc_endpoint;
```

### Subnet Id count for each VPC endpoints

```sql
select
  vpc_endpoint_id,
  jsonb_array_length(subnet_ids) as subnet_id_count
from
  aws_vpc_endpoint;
```

### Network details for each VPC endpoint

```sql
select
  vpc_endpoint_id,
  vpc_id,
  jsonb_array_elements(subnet_ids) as subnet_ids,
  jsonb_array_elements(network_interface_ids) as network_interface_ids,
  jsonb_array_elements(route_table_ids) as route_table_ids,
  sg ->> 'GroupName' as sg_name
from
  aws_vpc_endpoint
  cross join jsonb_array_elements(groups) as sg;
```

### DNS information for the VPC endpoints

```sql
select
  vpc_endpoint_id,
  private_dns_enabled,
  dns ->> 'DnsName' as dns_name,
  dns ->> 'HostedZoneId' as hosted_zone_id
from
  aws_vpc_endpoint
  cross join jsonb_array_elements(dns_entries) as dns;
```

### VPC endpoint count by VPC ID

```sql
select
  vpc_id,
  count(vpc_endpoint_id) as vpc_endpoint_count
from
  aws_vpc_endpoint
group by
  vpc_id;
```

### Count endpoints by endpoint type

```sql
select
  vpc_endpoint_type,
  count(vpc_endpoint_id)
from
  aws_vpc_endpoint
group by
  vpc_endpoint_type;
```

### List 'interface' type VPC Endpoints

```sql
select
  vpc_endpoint_id,
  service_name,
  vpc_id,
  vpc_endpoint_type
from
  aws_vpc_endpoint
where
  vpc_endpoint_type = 'Interface';
```