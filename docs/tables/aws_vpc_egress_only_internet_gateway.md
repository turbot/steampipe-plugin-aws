---
title: "Steampipe Table: aws_vpc_egress_only_internet_gateway - Query AWS VPC Egress Only Internet Gateways using SQL"
description: "Allows users to query AWS VPC Egress Only Internet Gateways, which provide egress only access for IPv6 traffic from the VPC to the internet."
folder: "VPC"
---

# Table: aws_vpc_egress_only_internet_gateway - Query AWS VPC Egress Only Internet Gateways using SQL

The AWS VPC Egress Only Internet Gateway is a resource that provides egress only access for IPv6 traffic from a Virtual Private Cloud (VPC) to the internet. It prevents inbound traffic from the internet, enhancing the security of your VPC. This feature is particularly useful when you want to allow outbound communication to the internet from your instances, but not allow any inbound traffic.

## Table Usage Guide

The `aws_vpc_egress_only_internet_gateway` table in Steampipe provides you with information about Egress Only Internet Gateways within Amazon Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer, to query gateway-specific details, including the gateway's attachments, creation time, and associated tags. You can utilize this table to gather insights on gateways, such as the gateways associated with a specific VPC, the state of the gateway's attachments, and more. The schema outlines the various attributes of the Egress Only Internet Gateway for you, including the gateway ID, VPC ID, and attachment state.

## Examples

### Egress only internet gateway basic info
Determine the status and associated Virtual Private Cloud (VPC) of your egress-only internet gateways across different regions. This is beneficial for managing network traffic and ensuring the secure flow of your outbound communication.

```sql+postgres
select
  id,
  att ->> 'State' as state,
  att ->> 'VpcId' as vpc_id,
  tags,
  region
from
  aws_vpc_egress_only_internet_gateway
  cross join jsonb_array_elements(attachments) as att;
```

```sql+sqlite
select
  id,
  json_extract(att.value, '$.State') as state,
  json_extract(att.value, '$.VpcId') as vpc_id,
  tags,
  region
from
  aws_vpc_egress_only_internet_gateway,
  json_each(attachments) as att;
```

### List unattached egress only gateways
Determine the areas in which egress-only internet gateways in your AWS VPC are unattached. This helps in identifying unused resources and potential cost savings.

```sql+postgres
select
  id,
  attachments
from
  aws_vpc_egress_only_internet_gateway
where
  attachments is null;
```

```sql+sqlite
select
  id,
  attachments
from
  aws_vpc_egress_only_internet_gateway
where
  attachments is null;
```


### List all the egress only gateways attached to default VPC
Determine the instances where egress-only internet gateways are connected to the default Virtual Private Cloud (VPC). This is useful for understanding the security configuration of your network and identifying potential areas of vulnerability.

```sql+postgres
select
  vig.id,
  vpc.is_default
from
  aws_vpc_egress_only_internet_gateway as vig
  cross join jsonb_array_elements(attachments) as i
  join aws_vpc vpc on i ->> 'VpcId' = vpc.vpc_id
where
  vpc.is_default = true;
```

```sql+sqlite
select
  vig.id,
  vpc.is_default
from
  aws_vpc_egress_only_internet_gateway as vig,
  json_each(attachments) as i
  join aws_vpc vpc on json_extract(i.value, '$.VpcId') = vpc.vpc_id
where
  vpc.is_default = 1;
```