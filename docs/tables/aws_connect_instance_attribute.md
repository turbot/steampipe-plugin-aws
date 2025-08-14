---
title: "Steampipe Table: aws_connect_instance_attribute - Query AWS Connect Instance Attributes using SQL"
description: "Allows users to query AWS Connect Instance Attributes for detailed information about each attribute configuration and value."
folder: "Connect"
---

# Table: aws_connect_instance_attribute - Query AWS Connect Instance Attributes using SQL

The AWS Connect Instance Attribute is a component of Amazon Connect that represents individual feature toggles and configuration settings for a Connect instance. Each attribute controls specific functionality such as Contact Lens, auto-resolve best voices, custom TTS voices, and other instance-level features. This table provides granular access to these configuration settings.

## Table Usage Guide

The `aws_connect_instance_attribute` table in Steampipe provides you with information about individual attributes within AWS Connect instances. This table allows you, as a DevOps engineer, to query attribute-specific details, including attribute types, values, and associated instance information. You can utilize this table to gather insights on Connect instance configurations, such as which features are enabled, attribute values, and more. The schema outlines the various attributes of the Connect instance attribute for you, including the instance ID, attribute type, and value.

## Examples

### Basic info
Explore which attributes are configured for your AWS Connect instances. This query can be particularly useful in understanding the feature configuration across your Connect instances.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute;
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute;
```

### List attributes for a specific instance
Identify all attributes configured for a specific Connect instance. This helps in understanding the complete configuration of a particular instance.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  instance_id = '3f276507-1576-4a30-bc23-f6388cec8893';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  instance_id = '3f276507-1576-4a30-bc23-f6388cec8893';
```

### Find instances with Contact Lens enabled
Identify Connect instances that have Contact Lens feature enabled. This is useful for understanding which instances have advanced analytics capabilities.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'CONTACT_LENS'
  and value = 'true';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'CONTACT_LENS'
  and value = 'true';
```

### Check auto-resolve best voices configuration
Analyze which Connect instances have auto-resolve best voices enabled. This feature helps in automatically selecting the best voice for text-to-speech.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'AUTO_RESOLVE_BEST_VOICES';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'AUTO_RESOLVE_BEST_VOICES';
```

### Find instances with custom TTS voices
Identify Connect instances that have custom text-to-speech voices enabled. This helps in understanding which instances use custom voice configurations.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'USE_CUSTOM_TTS_VOICES'
  and value = 'true';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'USE_CUSTOM_TTS_VOICES'
  and value = 'true';
```

### Check early media configuration
Analyze which Connect instances have early media enabled. This feature allows media to be sent before the call is answered.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'EARLY_MEDIA';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'EARLY_MEDIA';
```

### Find instances with multi-party conference enabled
Identify Connect instances that support multi-party conference calls. This helps in understanding which instances have advanced calling capabilities.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'MULTI_PARTY_CONFERENCE'
  and value = 'true';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'MULTI_PARTY_CONFERENCE'
  and value = 'true';
```

### Check contact flow logs configuration
Analyze which Connect instances have contact flow logs enabled. This feature provides detailed logging for contact flow execution.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'CONTACTFLOW_LOGS';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'CONTACTFLOW_LOGS';
```

### Find instances with high volume outbound enabled
Identify Connect instances that have high volume outbound calling enabled. This feature is useful for instances that need to make large numbers of outbound calls.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'HIGH_VOLUME_OUTBOUND'
  and value = 'true';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'HIGH_VOLUME_OUTBOUND'
  and value = 'true';
```

### Check enhanced contact monitoring configuration
Analyze which Connect instances have enhanced contact monitoring enabled. This feature provides additional monitoring capabilities for contacts.

```sql+postgres
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'ENHANCED_CONTACT_MONITORING';
```

```sql+sqlite
select
  instance_id,
  attribute_type,
  value
from
  aws_connect_instance_attribute
where
  attribute_type = 'ENHANCED_CONTACT_MONITORING';
```

### Get attribute information for all instances
Combine attribute data with instance information to get a comprehensive view of Connect instances and their configurations.

```sql+postgres
select
  a.instance_id,
  i.instance_alias,
  i.instance_status,
  a.attribute_type,
  a.value
from
  aws_connect_instance_attribute as a
  left join aws_connect_instance as i on a.instance_id = i.id
order by
  a.instance_id,
  a.attribute_type;
```

```sql+sqlite
select
  a.instance_id,
  i.instance_alias,
  i.instance_status,
  a.attribute_type,
  a.value
from
  aws_connect_instance_attribute as a
  left join aws_connect_instance as i on a.instance_id = i.id
order by
  a.instance_id,
  a.attribute_type;
```
