---
title: "Table: aws_eks_identity_provider_config - Query Amazon EKS Identity Provider Configurations using SQL"
description: "Allows users to query Amazon EKS Identity Provider Configurations for detailed information about the identity provider configurations for Amazon EKS clusters."
---

# Table: aws_eks_identity_provider_config - Query Amazon EKS Identity Provider Configurations using SQL

The `aws_eks_identity_provider_config` table in Steampipe provides information about the identity provider configurations for Amazon EKS clusters. This table allows DevOps engineers to query configuration-specific details, including the type of identity provider, client ID, issuer URL, and associated metadata. Users can utilize this table to gather insights on configurations, such as the type of identity provider, the client ID, and the issuer URL. The schema outlines the various attributes of the identity provider configuration, including the cluster name, creation time, tags, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_eks_identity_provider_config` table, you can use the `.inspect aws_eks_identity_provider_config` command in Steampipe.

**Key columns**:

- `cluster_name`: The name of the Amazon EKS cluster associated with the identity provider configuration. This can be used to join with other tables that contain information about EKS clusters.
- `identity_provider_config_name`: The name of the identity provider configuration. This can be used to join with other tables that contain information about identity provider configurations.
- `identity_provider_config_type`: The type of the identity provider configuration. This column is important as it allows users to filter results based on the type of identity provider configuration.

## Examples

### Basic info

```sql
select
  name,
  arn,
  cluster_name,
  tags,
  status
from
  aws_eks_identity_provider_config;
```

### List OIDC type Identity provider config

```sql
select
  name,
  arn,
  cluster_name,
  type
from
  aws_eks_identity_provider_config
where 
  type = 'oidc';
```