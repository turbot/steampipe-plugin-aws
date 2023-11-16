---
title: "Table: aws_vpc_dhcp_options - Query AWS VPC DHCP Options using SQL"
description: "Allows users to query DHCP Options associated with Virtual Private Cloud (VPC) in AWS."
---

# Table: aws_vpc_dhcp_options - Query AWS VPC DHCP Options using SQL

The `aws_vpc_dhcp_options` table in Steampipe provides information about DHCP Options associated with Virtual Private Cloud (VPC) within Amazon Web Services (AWS). This table allows network administrators and DevOps engineers to query DHCP Options specific details, including domain name servers, domain name, NTP servers, and associated metadata. Users can utilize this table to gather insights on DHCP Options, such as the configured domain name servers, NTP servers, and NetBIOS name servers. The schema outlines the various attributes of the DHCP Options, including the DHCP Options ID, domain name, domain name servers, NTP servers, NetBIOS name servers, NetBIOS node type, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_dhcp_options` table, you can use the `.inspect aws_vpc_dhcp_options` command in Steampipe.

### Key columns:

- `dhcp_options_id`: The ID of the set of DHCP options. It can be used to join this table with other tables to get more detailed information.
- `domain_name`: The domain name that instances that are associated with this set of DHCP options are given. It is useful to identify the domain associated with the DHCP options.
- `domain_name_servers`: The IP addresses of up to four domain name servers, or AmazonProvidedDNS. It is useful to identify the DNS servers associated with the DHCP options.

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
