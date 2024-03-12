---
title: "Steampipe Table: aws_eks_insight - Query AWS EKS Insights using SQL"
description: "Allows users to query AWS EKS Insight data, providing information about cluster insights within an AWS Elastic Kubernetes Service (EKS) cluster."
---

# Table: aws_eks_insight - Query AWS EKS Insights using SQL

The AWS EKS Insight is a resource within Amazon Elastic Kubernetes Service (EKS). 

## Table Usage Guide

The `aws_eks_insight` table in Steampipe provides you with information about each insight within an AWS Elastic Kubernetes Service (EKS) cluster. This table allows you, as a DevOps engineer, system administrator, or other technical professional, to gather insights into your EKS clusters and potential problems during an upgrade to next Kubernetes version. The schema outlines the various attributes of the EKS Cluster Insight for you, including the Kubernetes version, description, status and potential recommendations on how to remediate the insight.

## Examples

### Basic info
Explore the insights available for all your Amazon EKS clsuters. 

```sql+postgres
select
  id,
  name,
  cluster_name
  description,
  insight_status,
  recommendation
from
  aws_eks_insight;
```

```sql+sqlite
select
  id,
  name,
  cluster_name
  description,
  insight_status
from
  aws_eks_insight;
```

### List insights for specific Kubernetes version
Identify only the insights that are relevant to specific Kubernetes version.

```sql+postgres
select
  id,
  name,
  cluster_name
  description,
  insight_status
from
  aws_eks_insight
where
  kubernetes_version = '1.25';
```

```sql+sqlite
select
  id,
  name,
  cluster_name
  description,
  insight_status
from
  aws_eks_insight
where
  kubernetes_version = '1.25';
```

### Get insights for specific cluster
Get all insights for specific cluster.

```sql+postgres
select
  id,
  name,
  cluster_name
  description,
  insight_status
from
  aws_eks_insight
where
  cluster_name = 'eks.cluster.name';
```

```sql+sqlite
select
  id,
  name,
  cluster_name
  description,
  insight_status
from
  aws_eks_insight
where
  cluster_name = 'eks.cluster.name';
```