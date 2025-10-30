---
title: "Steampipe Table: aws_eks_access_entry - Query AWS EKS Access Entries using SQL"
description: "Allows users to query AWS EKS Access Entries to retrieve information about IAM principals that have been granted access to Amazon EKS clusters."
folder: "EKS"
---

# Table: aws_eks_access_entry - Query AWS EKS Access Entries using SQL

AWS EKS Access Entries provide a way to grant IAM principals (users or roles) access to your Amazon EKS clusters. Access entries allow you to manage cluster authentication and authorization by associating IAM principals with Kubernetes groups and access policies, enabling fine-grained access control to your EKS cluster resources.

## Table Usage Guide

The `aws_eks_access_entry` table in Steampipe provides you with information about access entries associated with each Amazon EKS cluster. This table allows you, as a DevOps engineer or security administrator, to query access entry-specific details, including the IAM principal ARN, associated Kubernetes groups, access policies, and associated metadata. You can utilize this table to gather insights on cluster access management, such as which IAM principals have access to specific clusters, what type of access they have, and the Kubernetes groups they are associated with. The schema outlines the various attributes of the EKS access entry for you, including the principal ARN, access entry type, username, Kubernetes groups, and associated tags.

## Examples

### Basic info
Explore which IAM principals have access to your AWS EKS clusters and understand their access configuration.

```sql+postgres
select
  cluster_name,
  principal_arn,
  type,
  username,
  kubernetes_groups,
  created_at
from
  aws_eks_access_entry;
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  type,
  username,
  kubernetes_groups,
  created_at
from
  aws_eks_access_entry;
```

### List access entries for a specific cluster
Identify all IAM principals that have been granted access to a particular EKS cluster.

```sql+postgres
select
  cluster_name,
  principal_arn,
  type,
  username,
  kubernetes_groups
from
  aws_eks_access_entry
where
  cluster_name = 'my-eks-cluster';
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  type,
  username,
  kubernetes_groups
from
  aws_eks_access_entry
where
  cluster_name = 'my-eks-cluster';
```

### List access entries by type
Determine which access entries are of a specific type (e.g., STANDARD, EC2_LINUX, FARGATE_LINUX).

```sql+postgres
select
  cluster_name,
  principal_arn,
  type,
  username
from
  aws_eks_access_entry
where
  type = 'STANDARD';
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  type,
  username
from
  aws_eks_access_entry
where
  type = 'STANDARD';
```

### Get count of access entries by cluster
Determine the total number of access entries per cluster to understand cluster access distribution.

```sql+postgres
select
  cluster_name,
  count(*) as access_entry_count
from
  aws_eks_access_entry
group by
  cluster_name;
```

```sql+sqlite
select
  cluster_name,
  count(*) as access_entry_count
from
  aws_eks_access_entry
group by
  cluster_name;
```

### List access entries created in the last 30 days
Identify recently created access entries to monitor new cluster access grants.

```sql+postgres
select
  cluster_name,
  principal_arn,
  type,
  username,
  created_at
from
  aws_eks_access_entry
where
  created_at >= now() - interval '30 days';
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  type,
  username,
  created_at
from
  aws_eks_access_entry
where
  created_at >= datetime('now', '-30 days');
```

### List access entries with specific Kubernetes groups
Find access entries that grant access to specific Kubernetes groups for role-based access control.

```sql+postgres
select
  cluster_name,
  principal_arn,
  username,
  kubernetes_groups
from
  aws_eks_access_entry
where
  kubernetes_groups @> '["system:masters"]'::jsonb;
```

```sql+sqlite
select
  cluster_name,
  principal_arn,
  username,
  kubernetes_groups
from
  aws_eks_access_entry
where
  json_extract(kubernetes_groups, '$') like '%system:masters%';
```

### Get access entries with their associated IAM role details
Join access entries with IAM roles to get comprehensive access information.

```sql+postgres
select
  ae.cluster_name,
  ae.principal_arn,
  ae.type,
  ae.username,
  ae.kubernetes_groups,
  r.create_date as role_created_at,
  r.max_session_duration
from
  aws_eks_access_entry ae
  left join aws_iam_role r on ae.principal_arn = r.arn;
```

```sql+sqlite
select
  ae.cluster_name,
  ae.principal_arn,
  ae.type,
  ae.username,
  ae.kubernetes_groups,
  r.create_date as role_created_at,
  r.max_session_duration
from
  aws_eks_access_entry ae
  left join aws_iam_role r on ae.principal_arn = r.arn;
```

