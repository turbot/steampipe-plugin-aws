# Table: aws_vpc_dhcp_options

The Dynamic Host Configuration Protocol (DHCP) provides a standard for passing configuration information to hosts on a TCP/IP network.

## Examples

### DHCP options configuration parameters info

```sql
select
  dhcp_options_id,
  domain_name,
  domain_name_servers,
  netbios_name_servers,
  netbios_node_type,
  ntp_servers
from
  aws_vpc_dhcp_options;
```


### List DHCP options which are not using AWS provided DNS

```sql
select
  dhcp_options_id,
  domain_name,
  domain_name_servers
from
  aws_vpc_dhcp_options
where
  domain_name_servers ? 'AmazonProvidedDNS';
```


### List all DHCP options without desired netbios (for example 2 - P-node is desired) node type

```sql
select
  dhcp_options_id,
  netbios_node_type
from
  aws_vpc_dhcp_options
  cross join jsonb_array_elements_text(netbios_node_type) as i
where
  not i.value :: int in (2);
```
