# Table: aws_emr_security_configuration

AWS EMR (Amazon Elastic MapReduce) Security Configuration is a set of security settings and configurations that can be applied to EMR clusters to manage encryption, authentication, and authorization. These configurations are crucial for ensuring that your EMR clusters handle data securely, protecting sensitive information and complying with various data security standards and regulations.

## Examples

### Basic info

```sql
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

```sql
select
  name,
  creation_date_time,
  security_configuration
from
  aws_emr_security_configuration
where
  creation_date_time >= now() - interval '30' day;
```

### Get encryption configuration details for security configurations

```sql
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