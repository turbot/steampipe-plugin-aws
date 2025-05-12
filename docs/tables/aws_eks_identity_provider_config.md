---
title: "Steampipe Table: aws_eks_identity_provider_config - Query Amazon EKS Identity Provider Configurations using SQL"
description: "Allows users to query Amazon EKS Identity Provider Configurations for detailed information about the identity provider configurations for Amazon EKS clusters."
folder: "Config"
---

# Table: aws_eks_identity_provider_config - Query Amazon EKS Identity Provider Configurations using SQL

The Amazon EKS Identity Provider Configurations is a feature of Amazon Elastic Kubernetes Service (EKS). It allows you to integrate and manage third-party identity providers for authentication with your EKS clusters. This ensures secure access and identity management for your Kubernetes workloads.

## Table Usage Guide

The `aws_eks_identity_provider_config` table in Steampipe provides you with information about the identity provider configurations for Amazon EKS clusters. This table allows you, as a DevOps engineer, to query configuration-specific details, including the type of identity provider, client ID, issuer URL, and associated metadata. You can utilize this table to gather insights on configurations, such as the type of identity provider, the client ID, and the issuer URL. The schema outlines the various attributes of the identity provider configuration, including the cluster name, creation time, tags, and status for you.

## Examples

### Basic info
Explore which AWS EKS identity provider configurations are in use and their current status. This can help you manage and monitor your AWS EKS resources more effectively.

```sql+postgres
select
  name,
  arn,
  cluster_name,
  tags,
  status
from
  aws_eks_identity_provider_config;
```

```sql+sqlite
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
Determine the areas in which OpenID Connect (OIDC) type identity provider configurations are used within your AWS Elastic Kubernetes Service (EKS) clusters. This is useful for understanding your security setup and ensuring that it aligns with your organization's policies.

```sql+postgres
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

```sql+sqlite
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