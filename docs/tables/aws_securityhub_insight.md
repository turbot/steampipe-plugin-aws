---
title: "Steampipe Table: aws_securityhub_insight - Query AWS Security Hub Insights using SQL"
description: "Allows users to query AWS Security Hub Insights to retrieve information about the insights in AWS Security Hub. This includes details such as insight ARN, name, filters, group by attributes, and more."
folder: "Security Hub"
---

# Table: aws_securityhub_insight - Query AWS Security Hub Insights using SQL

The AWS Security Hub Insight is a feature of AWS Security Hub that provides a summary of a specific security issue. It aggregates security findings across accounts, services, and supported AWS partners to provide a comprehensive view of your security posture. This allows you to quickly identify and react to potential security threats.

## Table Usage Guide

The `aws_securityhub_insight` table in Steampipe provides you with information about insights within AWS Security Hub. This table enables you, as a security analyst, to query insight-specific details, including the insight ARN, name, filters, and group by attributes. You can utilize this table to gather insights on the insights, such as insights with specific filters, the grouping of attributes, and more. The schema outlines the various attributes of the insight for you, including the insight ARN, name, filters, group by attributes, and associated metadata.

## Examples

### Basic info
Explore which security insights are grouped by specific attributes across different regions in your AWS SecurityHub, helping you manage and understand your security posture better.

```sql+postgres
select
  name,
  arn,
  group_by_attribute,
  region
from
  aws_securityhub_insight;
```

```sql+sqlite
select
  name,
  arn,
  group_by_attribute,
  region
from
  aws_securityhub_insight;
```

### List insights by a particular attribute
Discover the segments that are grouped by a specific attribute in the AWS Security Hub. This can help in identifying patterns or anomalies based on that attribute, enhancing your security management strategy.

```sql+postgres
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

```sql+sqlite
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
Identify instances where you can gain insights into the status of various workflows within your AWS SecurityHub. This can be useful in monitoring the progress and status of different security insights, aiding in efficient security management.

```sql+postgres
select
  name,
  arn,
  filters ->> 'WorkflowStatus' as workflow_status
from
  aws_securityhub_insight;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(filters, '$.WorkflowStatus') as workflow_status
from
  aws_securityhub_insight;
```

### Get severity details for insights
Gain insights into the severity level of security insights, which can help prioritize responses to potential security threats. This query is useful in identifying and understanding the criticality of the insights for efficient threat management.

```sql+postgres
select
  name,
  arn,
  filters ->> 'SeverityLabel' as severity_label
from
  aws_securityhub_insight;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(filters, '$.SeverityLabel') as severity_label
from
  aws_securityhub_insight;
```

### List insights that filter on critical severity labels 
Determine the areas in which critical security threats have been identified in your AWS Security Hub. This query allows you to focus on high-risk issues, enabling more efficient security management and response.

```sql+postgres
select
  name,
  arn,
  filters ->> 'SeverityLabel' as severity
from
  aws_securityhub_insight
where
  filters ->> 'SeverityLabel' = '{"Comparison": "EQUALS", "Value": "CRITICAL"}'
```

```sql+sqlite
select
  name,
  arn,
  json_extract(filters, '$.SeverityLabel') as severity
from
  aws_securityhub_insight
where
  json_extract(filters, '$.SeverityLabel') = '{"Comparison": "EQUALS", "Value": "CRITICAL"}'
```

### List insights that filter on ipv4_address threat intelligence type
This query allows you to identify potential security threats by pinpointing insights that are specifically filtering on IPv4 address threat intelligence type. This can be particularly useful in enhancing your cybersecurity measures by focusing on areas where your system may be vulnerable to IP-based threats.

```sql+postgres
select
  name,
  arn,
  filters ->> 'ThreatIntelIndicatorType' as threat_intelligence_details
from
  aws_securityhub_insight
where
  filters ->> 'ThreatIntelIndicatorType' = '{"Comparison": "EQUALS", "Value": "IPV4_ADDRESS"}'
```

```sql+sqlite
select
  name,
  arn,
  json_extract(filters, '$.ThreatIntelIndicatorType') as threat_intelligence_details
from
  aws_securityhub_insight
where
  json_extract(filters, '$.ThreatIntelIndicatorType') = '{"Comparison": "EQUALS", "Value": "IPV4_ADDRESS"}'
```

### List insights that failed compliance
Determine the areas in which security insights have failed to meet compliance standards, enabling you to focus your efforts on addressing these specific vulnerabilities.

```sql+postgres
select
  name,
  arn,
  filters ->> 'ComplianceStatus' as compliance_status
from
  aws_securityhub_insight
where
  filters ->> 'ComplianceStatus' = '{"Comparison": "EQUALS", "Value": "FAILED"}'
```

```sql+sqlite
select
  name,
  arn,
  json_extract(filters, '$.ComplianceStatus') as compliance_status
from
  aws_securityhub_insight
where
  json_extract(filters, '$.ComplianceStatus') = '{"Comparison": "EQUALS", "Value": "FAILED"}'
```

### Get malware details for insights
Explore potential security threats by identifying the instances of malware in your system. This query will help you gain insights into the name, path, and type of malware, aiding in your cybersecurity measures.

```sql+postgres
select
  name,
  arn,
  filters ->> 'MalwareName' as malware_name,
  filters ->> 'MalwarePath' as malware_path,
  filters ->> 'MalwareType' as malware_type
from
  aws_securityhub_insight;
```

```sql+sqlite
select
  name,
  arn,
  json_extract(filters, '$.MalwareName') as malware_name,
  json_extract(filters, '$.MalwarePath') as malware_path,
  json_extract(filters, '$.MalwareType') as malware_type
from
  aws_securityhub_insight;
```

### Get network details for insights
Discover the segments that are crucial for understanding your network's security. This query provides insights into network details such as source and destination domains, IPv4 and IPv6 addresses, and ports, which can be extremely useful in identifying potential security threats or areas of vulnerability.

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  json_extract(filters, '$.NetworkSourceDomain') as network_source_domain,
  json_extract(filters, '$.NetworkDestinationDomain') as network_destination_domain,
  json_extract(filters, '$.NetworkSourceIpV4') as network_source_ip_v4,
  json_extract(filters, '$.NetworkDestinationIpV4') as network_destination_ip_v4,
  json_extract(filters, '$.NetworkSourceIpV6') as network_source_ip_v6,
  json_extract(filters, '$.NetworkDestinationIpV6') as network_destination_ip_v6,
  json_extract(filters, '$.NetworkSourcePort') as network_source_port,
  json_extract(filters, '$.NetworkDestinationPort') as network_destination_port
from
  aws_securityhub_insight;
```

### Get record state details for a custom insight named 'sp'
Discover the status of a custom security insight within your AWS Security Hub. This is particularly useful for tracking and managing the state of your security insights.

```sql+postgres
select
  name,
  arn,
  filters ->> 'RecordState' as record_state
from
  aws_securityhub_insight
where
  name = 'sp';
```

```sql+sqlite
select
  name,
  arn,
  json_extract(filters, '$.RecordState') as record_state
from
  aws_securityhub_insight
where
  name = 'sp';
```