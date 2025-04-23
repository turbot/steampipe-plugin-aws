---
title: "Steampipe Table: aws_emr_security_configuration - Query AWS EMR Security Configurations using SQL"
description: "Allows users to query AWS EMR (Amazon Elastic MapReduce) Security Configurations. This table provides information about security settings and configurations that can be applied to EMR clusters, managing encryption, authentication, and authorization. These configurations are crucial for ensuring the secure handling of data, protecting sensitive information, and complying with various data security standards and regulations."
folder: "Config"
---

# Table: aws_emr_security_configuration - Query AWS EMR Security Configurations using SQL

AWS EMR (Amazon Elastic MapReduce) Security Configuration is a set of security settings and configurations that can be applied to EMR clusters to manage encryption, authentication, and authorization. These configurations are crucial for ensuring that your EMR clusters handle data securely, protecting sensitive information, and complying with various data security standards and regulations.

## Table Usage Guide

The `aws_emr_security_configuration` table in Steampipe allows users to query information about AWS EMR Security Configurations. These configurations are essential for securing EMR clusters, managing encryption, and ensuring compliance with data security standards. Users can retrieve details such as the configuration name, creation date and time, encryption configuration, instance metadata service configuration, and the overall security configuration.

## Examples

### Basic info
Retrieve basic information about AWS EMR Security Configurations, including their names, creation date and time, encryption configurations, instance metadata service configurations, and security configurations. This query provides an overview of the security configurations in your AWS EMR environment.

```sql+postgres
select
  name,
  creation_date_time,
  encryption_configuration,
  instance_metadata_service_configuration,
  security_configuration
from
  aws_emr_security_configuration;
```

```sql+sqlite
select
  name,
  creation_date_time,
  encryption_configuration,
  instance_metadata_service_configuration,
  security_configuration
from
  aws_emr_security_configuration;
```

### List security configurations created in the last 30 days
Identify AWS EMR Security Configurations created within the last 30 days. This query helps you keep track of recently created security configurations.

```sql+postgres
select
  name,
  creation_date_time,
  security_configuration
from
  aws_emr_security_configuration
where
  creation_date_time >= now() - interval '30' day;
```

```sql+sqlite
select
  name,
  creation_date_time,
  security_configuration
from
  aws_emr_security_configuration
where
  creation_date_time >= datetime('now', '-30 day');
```

### Get encryption configuration details for security configurations
Retrieve detailed encryption configuration information for AWS EMR Security Configurations. This includes information such as AWS KMS keys, EBS encryption settings, encryption key provider types, S3 encryption configurations, and more. This query allows you to inspect the encryption settings in your security configurations.

```sql+postgres
select
  name,
  creation_date_time,
  encryption_configuration -> 'AtRestEncryptionConfiguration' -> 'LocalDiskEncryptionConfiguration' ->> 'AwsKmsKey' as aws_kms_key,
  encryption_configuration -> 'AtRestEncryptionConfiguration' -> 'LocalDiskEncryptionConfiguration' ->> 'EnableEbsEncryption' as enable_ebs_encryption,
  encryption_configuration -> 'AtRestEncryptionConfiguration' -> 'LocalDiskEncryptionConfiguration' ->> 'EncryptionKeyProviderType' as encryption_key_provider_type,
  encryption_configuration -> 'S3EncryptionConfiguration' as s3_encryption_configuration,
  encryption_configuration ->> 'EnableAtRestEncryption' as enable_at_rest_encryption,
  encryption_configuration ->> 'EnableInTransitEncryption' as enable_in_transit_encryption,
  encryption_configuration -> 'InTransitEncryptionConfiguration' as in_transit_encryption_configuration
from
  aws_emr_security_configuration;
```

```sql+sqlite
select
  name,
  creation_date_time,
  json_extract(encryption_configuration, '$.AtRestEncryptionConfiguration.LocalDiskEncryptionConfiguration.AwsKmsKey') as aws_kms_key,
  json_extract(encryption_configuration, '$.AtRestEncryptionConfiguration.LocalDiskEncryptionConfiguration.EnableEbsEncryption') as enable_ebs_encryption,
  json_extract(encryption_configuration, '$.AtRestEncryptionConfiguration.LocalDiskEncryptionConfiguration.EncryptionKeyProviderType') as encryption_key_provider_type,
  json_extract(encryption_configuration, '$.S3EncryptionConfiguration') as s3_encryption_configuration,
  json_extract(encryption_configuration, '$.EnableAtRestEncryption') as enable_at_rest_encryption,
  json_extract(encryption_configuration, '$.EnableInTransitEncryption') as enable_in_transit_encryption,
  json_extract(encryption_configuration, '$.InTransitEncryptionConfiguration') as in_transit_encryption_configuration
from
  aws_emr_security_configuration;
```