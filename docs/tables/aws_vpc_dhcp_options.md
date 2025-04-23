---
title: "Steampipe Table: aws_vpc_dhcp_options - Query AWS VPC DHCP Options using SQL"
description: "Allows users to query DHCP Options associated with Virtual Private Cloud (VPC) in AWS."
folder: "VPC"
---

# Table: aws_vpc_dhcp_options - Query AWS VPC DHCP Options using SQL

The AWS VPC DHCP Options is a feature within Amazon's Virtual Private Cloud (VPC) that allows you to configure Domain Name System (DNS) settings for your instances that get their IP addresses from a DHCP set. You can specify DNS servers and domain names that Amazon EC2 instances use when they're launched in your VPC. DHCP options sets provide a simple way to manage DNS settings consistently across your entire VPC, enhancing the overall network management.

## Table Usage Guide

The `aws_vpc_dhcp_options` table in Steampipe provides you with information about DHCP Options associated with Virtual Private Cloud (VPC) within Amazon Web Services (AWS). This table allows you, as a network administrator or DevOps engineer, to query DHCP Options specific details, including domain name servers, domain name, NTP servers, and associated metadata. You can utilize this table to gather insights on DHCP Options, such as the configured domain name servers, NTP servers, and NetBIOS name servers. The schema outlines the various attributes of the DHCP Options for you, including the DHCP Options ID, domain name, domain name servers, NTP servers, NetBIOS name servers, NetBIOS node type, and associated tags.

## Examples

### DHCP options configuration parameters info
Analyze the settings to understand the configuration of DHCP options within your AWS VPC. This is useful for maintaining network stability and optimizing domain and server configurations.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that are not utilizing the DNS provided by AWS in their DHCP options. This could be useful to ensure compliance with company policies or to identify potential areas for optimization.

```sql+postgres
select
  dhcp_options_id,
  domain_name,
  domain_name_servers
from
  aws_vpc_dhcp_options
where
  domain_name_servers ? 'AmazonProvidedDNS';
```

```sql+sqlite
select
  dhcp_options_id,
  domain_name,
  domain_name_servers
from
  aws_vpc_dhcp_options
where
  json_extract(domain_name_servers, '$.AmazonProvidedDNS') is not null;
```


### List all DHCP options without desired netbios (for example 2 - P-node is desired) node type
This query is used to identify the DHCP options that are not configured with the desired netbios node type, in this case, P-node. This can help in managing network settings efficiently by pinpointing those options that may need to be updated to ensure optimal network performance.

```sql+postgres
select
  dhcp_options_id,
  netbios_node_type
from
  aws_vpc_dhcp_options
  cross join jsonb_array_elements_text(netbios_node_type) as i
where
  not i.value :: int in (2);
```

```sql+sqlite
select
  dhcp_options_id,
  json_extract(i.value, '$') as netbios_node_type
from
  aws_vpc_dhcp_options,
  json_each(netbios_node_type) as i
where
  cast(json_extract(i.value, '$') as int) not in (2);
```