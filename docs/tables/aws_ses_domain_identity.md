---
title: "Table: aws_ses_domain_identity - Query Amazon Simple Email Service Domain Identities using SQL"
description: "Allows users to query Amazon Simple Email Service Domain Identities. The aws_ses_domain_identity table in Steampipe provides information about domain identities within Amazon Simple Email Service (SES). This table allows DevOps engineers to query domain-specific details, including verification status, DKIM attributes, and associated metadata. Users can utilize this table to gather insights on domain identities, such as verification status, DKIM tokens, and more. The schema outlines the various attributes of the SES domain identity, including the identity name, verification status, DKIM enabled status, and DKIM tokens."
---

# Table: aws_ses_domain_identity - Query Amazon Simple Email Service Domain Identities using SQL

The `aws_ses_domain_identity` table in Steampipe provides information about domain identities within Amazon Simple Email Service (SES). This table allows DevOps engineers to query domain-specific details, including verification status, DKIM attributes, and associated metadata. Users can utilize this table to gather insights on domain identities, such as verification status, DKIM tokens, and more. The schema outlines the various attributes of the SES domain identity, including the identity name, verification status, DKIM enabled status, and DKIM tokens.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the aws_ses_domain_identity table, you can use the `.inspect aws_ses_domain_identity` command in Steampipe.

### Key columns:

- `identity_name`: The name of the identity. This is a key column that can be used to join this table with other tables to gather more specific information about a particular identity.
- `verification_status`: The verification status of the domain. This column is important as it helps to understand whether the domain identity is successfully verified or not.
- `dkim_enabled`: Indicates whether DKIM signing is enabled for the identity. This column is useful for security checks as it indicates if the domain identity is configured with DKIM, an email authentication method, which can help protect the domain from being used for phishing scams.

## Examples

### Basic info

```sql
select
  identity,
  arn,
  region,
  akas
from
  aws_ses_domain_identity;
```

### List domain identities which failed verification

```sql
select
  identity,
  region,
  verification_status
from
  aws_ses_domain_identity
where
  verification_status = 'Failed';
```
