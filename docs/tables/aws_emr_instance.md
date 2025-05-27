---
title: "Steampipe Table: aws_emr_instance - Query AWS EMR Instances using SQL"
description: "Allows users to query AWS EMR Instances for detailed information about the status, configuration, and other metadata of each instance."
folder: "EMR"
---

# Table: aws_emr_instance - Query AWS EMR Instances using SQL

The AWS Elastic MapReduce (EMR) Instance is a component of the Amazon EMR service, which provides a managed Hadoop framework to process vast amounts of data across dynamically scalable Amazon EC2 instances. It enables businesses, researchers, data analysts, and developers to easily and cost-effectively process vast amounts of data. It uses Hadoop processing combined with several AWS products to do tasks such as web indexing, data mining, log file analysis, machine learning, scientific simulation, and data warehousing.

## Table Usage Guide

The `aws_emr_instance` table in Steampipe provides you with information about instances within AWS Elastic MapReduce (EMR). This table allows you, as a DevOps engineer, to query instance-specific details, including instance status, instance group ID, and associated metadata. You can utilize this table to gather insights on instances, such as instance health status, instance configuration details, and more. The schema outlines the various attributes of the EMR instance for you, including the instance ID, EBS volumes, and associated tags.

## Examples

### Basic info
Explore which instances are associated with a specific cluster in your AWS Elastic Map Reduce service. This query can be particularly useful in understanding the distribution of resources and optimizing cluster management.

```sql+postgres
select
  id,
  cluster_id,
  ec2_instance_id,
  instance_type,
  private_dns_name,
  private_ip_address
from
  aws_emr_instance;
```

```sql+sqlite
select
  id,
  cluster_id,
  ec2_instance_id,
  instance_type,
  private_dns_name,
  private_ip_address
from
  aws_emr_instance;
```

### List instances by type
Identify instances where the type is specified as 'm2.4xlarge' in the AWS EMR service. This can be useful in understanding the distribution and usage of specific instance types within your cloud infrastructure.

```sql+postgres
select
  id,
  ec2_instance_id,
  instance_type
from
  aws_emr_instance
where
  instance_type = 'm2.4xlarge';
```

```sql+sqlite
select
  id,
  ec2_instance_id,
  instance_type
from
  aws_emr_instance
where
  instance_type = 'm2.4xlarge';
```

### List instances for a cluster
This query is useful to identify the specific instances associated with a particular cluster in a cloud-based environment. It aids in understanding the composition and configuration of the cluster for better resource management and optimization.

```sql+postgres
select
  id,
  ec2_instance_id,
  instance_type
from
  aws_emr_instance
where
  cluster_id = 'j-21HIX5R2NZMXJ';
```

```sql+sqlite
select
  id,
  ec2_instance_id,
  instance_type
from
  aws_emr_instance
where
  cluster_id = 'j-21HIX5R2NZMXJ';
```

### Get volume details for an instance
This example allows you to discover the details of the volume attached to a specific instance in your AWS Elastic MapReduce (EMR) service. It can be used to better understand the storage configuration of your EMR instances, which can aid in optimizing storage usage and costs.

```sql+postgres
select
  id,
  ec2_instance_id,
  instance_type,
  v -> 'Device' as device,
  v -> 'VolumeId' as volume_id
from
  aws_emr_instance as ei,
  jsonb_array_elements(ebs_volumes) as v
where
  ei.id = 'ci-ULCFS2ZN0FK7';
```

```sql+sqlite
select
  aws_emr_instance.id,
  ec2_instance_id,
  instance_type,
  json_extract(v.value, '$.Device') as device,
  json_extract(v.value, '$.VolumeId') as volume_id
from
  aws_emr_instance,
  json_each(ebs_volumes) as v
where
  aws_emr_instance.id = 'ci-ULCFS2ZN0FK7';
```