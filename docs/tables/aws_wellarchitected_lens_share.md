---
title: "Steampipe Table: aws_wellarchitected_lens_share - Query AWS Well-Architected Tool Lens Shares using SQL"
description: "Allows users to query Lens Shares in AWS Well-Architected Tool, providing details about shared lenses including the share ARN, share status, and the AWS account ID of the lens owner."
folder: "Well-Architected"
---

# Table: aws_wellarchitected_lens_share - Query AWS Well-Architected Tool Lens Shares using SQL

The AWS Well-Architected Tool Lens Shares is a feature of the AWS Well-Architected Tool. It allows you to share your workload reviews with other AWS accounts or within your organization. This ensures that best practices and architectural recommendations are consistently applied across your workloads and teams.

## Table Usage Guide

The `aws_wellarchitected_lens_share` table in Steampipe provides you with information about Lens Shares within AWS Well-Architected Tool. This table allows you, as a cloud architect or developer, to query details about shared lenses, including the share ARN, share status, and the AWS account ID of the lens owner. You can utilize this table to gather insights on shared lenses, such as the status of shared lenses, the AWS account ID of the lens owner, and more. The schema outlines for you the various attributes of the Lens Share, including the share ARN, share status, and the AWS account ID of the lens owner.

## Examples

### Basic info
Explore the distribution of shared AWS Well-Architected Lens to understand which lenses are shared and with whom. This can help assess the spread of AWS best practices across your organization.

```sql+postgres
select
  lens_name,
  lens_alias,
  lens_arn,
  share_id,
  shared_with
from
  aws_wellarchitected_lens_share;
```

```sql+sqlite
select
  lens_name,
  lens_alias,
  lens_arn,
  share_id,
  shared_with
from
  aws_wellarchitected_lens_share;
```

### Get lens details of the shared lenses
Explore the shared lenses in your AWS Well-Architected tool to understand details like the lens status, type, owner, and share invitation ID. This can help in managing and tracking the shared lenses effectively.

```sql+postgres
select
  s.lens_name,
  l.arn,
  l.lens_status,
  l.lens_type,
  l.owner,
  l.share_invitation_id
from
  aws_wellarchitected_lens_share as s,
  aws_wellarchitected_lens as l
where
  s.lens_arn = l.arn;
```

```sql+sqlite
select
  s.lens_name,
  l.arn,
  l.lens_status,
  l.lens_type,
  l.owner,
  l.share_invitation_id
from
  aws_wellarchitected_lens_share as s,
  aws_wellarchitected_lens as l
where
  s.lens_arn = l.arn;
```

### List shared lenses that are in pending state
Identify shared lenses within the AWS Well-Architected Tool that are still pending approval. This can help you manage your shared resources and ensure timely responses to share requests.

```sql+postgres
select
  lens_name,
  lens_alias,
  lens_arn,
  share_id,
  shared_with,
  status
from
  aws_wellarchitected_lens_share
where
  status = 'PENDING';
```

```sql+sqlite
select
  lens_name,
  lens_alias,
  lens_arn,
  share_id,
  shared_with,
  status
from
  aws_wellarchitected_lens_share
where
  status = 'PENDING';
```