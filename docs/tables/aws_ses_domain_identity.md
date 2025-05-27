---
title: "Steampipe Table: aws_ses_domain_identity - Query Amazon Simple Email Service Domain Identities using SQL"
description: "Allows users to query Amazon Simple Email Service Domain Identities. The aws_ses_domain_identity table in Steampipe provides information about domain identities within Amazon Simple Email Service (SES). This table allows DevOps engineers to query domain-specific details, including verification status, DKIM attributes, and associated metadata. Users can utilize this table to gather insights on domain identities, such as verification status, DKIM tokens, and more. The schema outlines the various attributes of the SES domain identity, including the identity name, verification status, DKIM enabled status, and DKIM tokens."
folder: "Simple Email Service (SES)"
---

# Table: aws_ses_domain_identity - Query Amazon Simple Email Service Domain Identities using SQL

The Amazon Simple Email Service (SES) Domain Identity is an entity that you use to send email. It represents the domain that you use for sending email, and is verified by Amazon SES. Once verified, you can send email from any address on the specified domain.

## Table Usage Guide

The `aws_ses_domain_identity` table in Steampipe provides you with information about domain identities within Amazon Simple Email Service (SES). This table allows you, as a DevOps engineer, to query domain-specific details, including verification status, DKIM attributes, and associated metadata. You can utilize this table to gather insights on domain identities, such as verification status, DKIM tokens, and more. The schema outlines the various attributes of the SES domain identity for you, including the identity name, verification status, DKIM enabled status, and DKIM tokens.

## Examples

### Basic info
Determine the areas in which your AWS Simple Email Service (SES) domain identities are located. This allows you to gain insights into the regional distribution of your SES domain identities, which can be useful for optimizing email delivery and managing regional compliance requirements.

```sql+postgres
select
  identity,
  arn,
  region,
  akas
from
  aws_ses_domain_identity;
```

```sql+sqlite
select
  identity,
  arn,
  region,
  akas
from
  aws_ses_domain_identity;
```

### List domain identities which failed verification
Discover the segments that have failed the verification process within your domain identities. This query is useful for identifying potential issues with your SES domain identities across different regions.

```sql+postgres
select
  identity,
  region,
  verification_status
from
  aws_ses_domain_identity
where
  verification_status = 'Failed';
```

```sql+sqlite
select
  identity,
  region,
  verification_status
from
  aws_ses_domain_identity
where
  verification_status = 'Failed';
```

### Retrieve tags for SES domain identities
Identify the tags associated with domain identities in Amazon SES. This query retrieves domain identity details along with their assigned tags, enabling better resource management and organization.

```sql+postgres
select 
  i.arn,
  i.identity,
  r.tags_src,
  r.tags
from 
  aws_ses_domain_identity as i
  join aws_tagging_resource as r on r.arn = i.arn;
```

```sql+sqlite
select 
  i.arn,
  i.identity,
  r.tags_src,
  r.tags
from 
  aws_ses_domain_identity as i
  join aws_tagging_resource as r on r.arn = i.arn;
```