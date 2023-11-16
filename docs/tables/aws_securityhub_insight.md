---
title: "Table: aws_securityhub_insight - Query AWS Security Hub Insights using SQL"
description: "Allows users to query AWS Security Hub Insights to retrieve information about the insights in AWS Security Hub. This includes details such as insight ARN, name, filters, group by attributes, and more."
---

# Table: aws_securityhub_insight - Query AWS Security Hub Insights using SQL

The `aws_securityhub_insight` table in Steampipe provides information about insights within AWS Security Hub. This table allows security analysts to query insight-specific details, including the insight ARN, name, filters, and group by attributes. Users can utilize this table to gather insights on the insights, such as insights with specific filters, the grouping of attributes, and more. The schema outlines the various attributes of the insight, including the insight ARN, name, filters, group by attributes, and associated metadata.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_securityhub_insight` table, you can use the `.inspect aws_securityhub_insight` command in Steampipe.

### Key columns:

- `arn`: The ARN of the insight. This is a unique identifier and can be used to join this table with other tables.
- `name`: The name of the insight. It can be useful to filter insights based on their names.
- `filters`: The filters of the insight. It can be used to understand the criteria used by the insight to identify findings.

## Examples

### Basic info

```sql
select
  name,
  arn,
  group_by_attribute,
  region
from
  aws_securityhub_insight;
```

### List insights by a particular attribute

```sql
select
  name,
  arn,
  group_by_attribute,
  region
from
  aws_securityhub_insight
where
  group_by_attribute = 'ResourceId';
```

### Get workflow status details for insights

```sql
select
  name,
  arn,
  filters ->> 'WorkflowStatus' as workflow_status
from
  aws_securityhub_insight;
```

### Get severity details for insights

```sql
select
  name,
  arn,
  filters ->> 'SeverityLabel' as severity_label
from
  aws_securityhub_insight;
```

### List insights that filter on critical severity labels 

```sql
select
  name,
  arn,
  filters ->> 'SeverityLabel' as severity
from
  aws_securityhub_insight
where
  filters ->> 'SeverityLabel' = '{"Comparison": "EQUALS", "Value": "CRITICAL"}'
```

### List insights that filter on ipv4_address threat intelligence type

```sql
select
  name,
  arn,
  filters ->> 'ThreatIntelIndicatorType' as threat_intelligence_details
from
  aws_securityhub_insight
where
  filters ->> 'ThreatIntelIndicatorType' = '{"Comparison": "EQUALS", "Value": "IPV4_ADDRESS"}'
```

### List insights that failed compliance

```sql
select
  name,
  arn,
  filters ->> 'ComplianceStatus' as compliance_status
from
  aws_securityhub_insight
where
  filters ->> 'ComplianceStatus' = '{"Comparison": "EQUALS", "Value": "FAILED"}'
```

### Get malware details for insights

```sql
select
  name,
  arn,
  filters ->> 'MalwareName' as malware_name,
  filters ->> 'MalwarePath' as malware_path,
  filters ->> 'MalwareType' as malware_type
from
  aws_securityhub_insight;
```

### Get network details for insights

```sql
select
  name,
  arn,
  filters ->> 'NetworkSourceDomain' as network_source_domain,
  filters ->> 'NetworkDestinationDomain' as network_destination_domain,
  filters ->> 'NetworkSourceIpV4' as network_source_ip_v4,
  filters ->> 'NetworkDestinationIpV4' as network_destination_ip_v4,
  filters ->> 'NetworkSourceIpV6' as network_source_ip_v6,
  filters ->> 'NetworkDestinationIpV6' as network_destination_ip_v6,
  filters ->> 'NetworkSourcePort' as network_source_port,
  filters ->> 'NetworkDestinationPort' as network_destination_port
from
  aws_securityhub_insight;
```

### Get record state details for a custom insight named 'sp'

```sql
select
  name,
  arn,
  filters ->> 'RecordState' as record_state
from
  aws_securityhub_insight
where
  name = 'sp';
```