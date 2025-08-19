---
title: "Steampipe Table: aws_servicecatalog_portfolio_share - Query AWS Service Catalog Portfolio Shares using SQL"
description: "Allows users to query AWS Service Catalog Portfolio Shares, providing information about portfolio sharing configurations and access permissions."
folder: "ServiceCatalog"
---

# Table: aws_servicecatalog_portfolio_share - Query AWS Service Catalog Portfolio Shares using SQL

AWS Service Catalog Portfolio Shares define how portfolios are shared with different entities such as accounts, organizations, organizational units, or organization member accounts. These shares control access to portfolios and their associated products, allowing organizations to manage portfolio distribution across their AWS environment.

## Table Usage Guide

The `aws_servicecatalog_portfolio_share` table in Steampipe provides you with information about portfolio shares within your AWS Service Catalog. This table allows you, as a cloud administrator or DevOps engineer, to query portfolio share details including share types, recipient entities, acceptance status, and sharing permissions. You can utilize this table to gather insights on portfolio sharing, such as which portfolios are shared with which entities, whether shares have been accepted, and what permissions are granted through the shares.

**Important notes:**
- In order to query portfolio shares, the `type` column must be specified in the WHERE clause. The table requires a specific share type to retrieve the corresponding portfolio share information. Valid values are:
  - `ACCOUNT` - External account to account shares
  - `ORGANIZATION` - Shares to an organization
  - `ORGANIZATIONAL_UNIT` - Shares to organizational units
  - `ORGANIZATION_MEMBER_ACCOUNT` - Shares to organization member accounts

## Examples

### List all portfolio shares by type
Analyze portfolio shares for a specific type to understand sharing configurations across your Service Catalog portfolios. This is useful for auditing portfolio access and permissions.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  share_principals,
  share_tag_options,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ACCOUNT';
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  share_principals,
  share_tag_options,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ACCOUNT';
```

### Find accepted portfolio shares
Identify portfolio shares that have been accepted by recipient entities. This helps track which shares are actively being used.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ACCOUNT'
  and accepted = true;
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ACCOUNT'
  and accepted = 1;
```

### Check organization shares
Review portfolio shares that are shared with organizations. This is useful for understanding organization-wide portfolio access.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  share_principals,
  share_tag_options,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATION';
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  share_principals,
  share_tag_options,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATION';
```

### Find shares with principal sharing enabled
Identify portfolio shares that have principal sharing enabled, allowing recipients to share the portfolio with other entities.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  share_principals,
  share_tag_options,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATION'
  and share_principals = true;
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  share_principals,
  share_tag_options,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATION'
  and share_principals = 1;
```

### Check account-specific shares
Review portfolio shares that are shared with specific AWS accounts. This helps track cross-account portfolio access.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ACCOUNT';
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ACCOUNT';
```

### Find shares with tag options sharing
Identify portfolio shares that have tag options sharing enabled, allowing recipients to access and use tag options associated with the portfolio.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  share_tag_options,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATIONAL_UNIT'
  and share_tag_options = true;
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  share_tag_options,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATIONAL_UNIT'
  and share_tag_options = 1;
```

### Get specific portfolio share
Retrieve details for a specific portfolio share using the portfolio ID, type, and principal ID.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  share_principals,
  share_tag_options,
  portfolio_arn,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ACCOUNT'
  and portfolio_id = 'port-1234567890abcdef'
  and principal_id = '123456789012';
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  share_principals,
  share_tag_options,
  portfolio_arn,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ACCOUNT'
  and portfolio_id = 'port-1234567890abcdef'
  and principal_id = '123456789012';
```

### Check organizational unit shares
Review portfolio shares that are shared with organizational units. This helps understand OU-level portfolio access.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATIONAL_UNIT';
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATIONAL_UNIT';
```

### Find pending portfolio shares
Identify portfolio shares that have not been accepted yet. This helps track shares that are awaiting recipient approval.

```sql+postgres
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATION_MEMBER_ACCOUNT'
  and accepted = false;
```

```sql+sqlite
select
  portfolio_display_name,
  portfolio_id,
  type,
  principal_id,
  accepted,
  region
from
  aws_servicecatalog_portfolio_share
where
  type = 'ORGANIZATION_MEMBER_ACCOUNT'
  and accepted = 0;
```
