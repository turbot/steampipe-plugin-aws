---
title: "Steampipe Table: aws_dx_customer_metadata - Query AWS Direct Connect Customer Metadata using SQL"
description: "Allows users to query AWS Direct Connect Customer Metadata for information about customer agreements and NNI partner types."
folder: "Direct Connect"
---

# Table: aws_dx_customer_metadata - Query AWS Direct Connect Customer Metadata using SQL

AWS Direct Connect Customer Metadata provides information about customer agreements and network-to-network interface (NNI) partner configurations. This metadata is essential for understanding the customer's agreement status and the type of network interconnections available through AWS Direct Connect.

## Table Usage Guide

The `aws_dx_customer_metadata` table in Steampipe provides you with information about AWS Direct Connect customer metadata. This table allows you, as a DevOps engineer, to query customer-specific details, including agreement names, status, and NNI partner types. You can utilize this table to gather insights on customer agreements, such as active agreements, partner types, and agreement status. The schema outlines the various attributes of the customer metadata for you, including the agreement name, status, and NNI partner type.

## Examples

### Basic customer metadata information
Explore the Direct Connect customer metadata to understand agreement status and partner configurations.

```sql+postgres
select
  agreement_name,
  status,
  nni_partner_type
from
  aws_dx_customer_metadata;
```

```sql+sqlite
select
  agreement_name,
  status,
  nni_partner_type
from
  aws_dx_customer_metadata;
```

### Group metadata by status
Analyze the distribution of customer agreements by their current status to understand the overall health of Direct Connect agreements.

```sql+postgres
select
  status,
  count(*) as agreement_count
from
  aws_dx_customer_metadata
group by
  status;
```

```sql+sqlite
select
  status,
  count(*) as agreement_count
from
  aws_dx_customer_metadata
group by
  status;
```

### NNI partner type distribution
Understand the distribution of different NNI partner types in your Direct Connect setup.

```sql+postgres
select
  nni_partner_type,
  count(*) as partner_count,
  array_agg(distinct status) as statuses
from
  aws_dx_customer_metadata
group by
  nni_partner_type;
```

```sql+sqlite
select
  nni_partner_type,
  count(*) as partner_count,
  group_concat(distinct status) as statuses
from
  aws_dx_customer_metadata
group by
  nni_partner_type;
```

### Active agreements with partner information
Find all active customer agreements and their associated NNI partner types for operational insights.

```sql+postgres
select
  agreement_name,
  status,
  nni_partner_type
from
  aws_dx_customer_metadata
where
  status = 'active';
```

```sql+sqlite
select
  agreement_name,
  status,
  nni_partner_type
from
  aws_dx_customer_metadata
where
  status = 'active';
```

### V2 partner agreements
Identify agreements with V2 NNI partners, which typically offer enhanced capabilities.

```sql+postgres
select
  agreement_name,
  status,
  nni_partner_type
from
  aws_dx_customer_metadata
where
  nni_partner_type = 'V2';
```

```sql+sqlite
select
  agreement_name,
  status,
  nni_partner_type
from
  aws_dx_customer_metadata
where
  nni_partner_type = 'V2';
```

### Customer metadata by region
Analyze customer metadata distribution across different AWS regions.

```sql+postgres
select
  region,
  count(*) as metadata_count,
  array_agg(distinct nni_partner_type) as partner_types
from
  aws_dx_customer_metadata
group by
  region
order by
  metadata_count desc;
```

```sql+sqlite
select
  region,
  count(*) as metadata_count,
  group_concat(distinct nni_partner_type) as partner_types
from
  aws_dx_customer_metadata
group by
  region
order by
  metadata_count desc;
```

### Cross-account metadata analysis
Identify customer metadata that spans multiple AWS accounts for complex organizational structures.

```sql+postgres
select
  account_id,
  count(*) as metadata_count,
  array_agg(distinct agreement_name) as agreements
from
  aws_dx_customer_metadata
group by
  account_id
having
  count(*) > 1;
```

```sql+sqlite
select
  account_id,
  count(*) as metadata_count,
  group_concat(distinct agreement_name) as agreements
from
  aws_dx_customer_metadata
group by
  account_id
having
  count(*) > 1;
```

### Non-partner agreements
Find agreements that are not associated with NNI partners for direct AWS connections.

```sql+postgres
select
  agreement_name,
  status,
  nni_partner_type
from
  aws_dx_customer_metadata
where
  nni_partner_type = 'nonPartner'
  or nni_partner_type is null;
```

```sql+sqlite
select
  agreement_name,
  status,
  nni_partner_type
from
  aws_dx_customer_metadata
where
  nni_partner_type = 'nonPartner'
  or nni_partner_type is null;
```
