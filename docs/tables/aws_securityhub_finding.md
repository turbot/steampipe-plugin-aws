---
title: "Steampipe Table: aws_securityhub_finding - Query AWS Security Hub Findings using SQL"
description: "Allows users to query AWS Security Hub Findings to gather information about security issues identified within AWS resources. This includes details such as the severity, status, and description of the finding, the resources affected, and any remediation steps recommended."
folder: "Security Hub"
---

# Table: aws_securityhub_finding - Query AWS Security Hub Findings using SQL

The AWS Security Hub Findings is a feature within AWS Security Hub that aggregates, organizes, and prioritizes your security alerts, or findings, from multiple AWS services. It also consolidates findings from different AWS services such as Amazon GuardDuty, Amazon Inspector, and Amazon Macie, as well as from AWS Partner solutions. The findings are presented as potential security issues, such as insecure configurations or malicious activity.

## Table Usage Guide

The `aws_securityhub_finding` table in Steampipe provides you with information about security findings within AWS Security Hub. This table allows you as a security analyst or DevOps engineer to query details about identified security issues, including their severity, status, description, the resources affected, and any recommended remediation steps. You can utilize this table to gather insights on security vulnerabilities, such as open security groups, exposed access keys, and more. The schema outlines the various attributes of the security finding for you, including the finding ARN, ID, title, description, severity, and associated resources.

## Examples

### Basic info
Explore the critical security findings in your AWS environment to understand their severity and level of confidence. This can help prioritize remediation efforts and improve overall security posture.

```sql+postgres
select
  title,
  id,
  company_name,
  created_at,
  criticality,
  confidence
from
  aws_securityhub_finding;
```

```sql+sqlite
select
  title,
  id,
  company_name,
  created_at,
  criticality,
  confidence
from
  aws_securityhub_finding;
```

### List findings with high severity
Determine the areas in which high severity security findings occur within your AWS Security Hub to prioritize and address significant threats effectively.

```sql+postgres
select
  title,
  product_arn,
  product_name,
  severity ->> 'Original' as severity_original
from
  aws_securityhub_finding
where
  severity ->> 'Original' = 'HIGH';
```

```sql+sqlite
select
  title,
  product_arn,
  product_name,
  json_extract(severity, '$.Original') as severity_original
from
  aws_securityhub_finding
where
  json_extract(severity, '$.Original') = 'HIGH';
```

### Count the number of findings by severity
Determine the distribution of security issues based on their severity level, providing a quick overview of the security status and helping prioritize remediation efforts.

```sql+postgres
select
  severity ->> 'Original' as severity_original,
  count(severity ->> 'Original')
from
  aws_securityhub_finding
group by
  severity ->> 'Original'
order by
  severity ->> 'Original';
```

```sql+sqlite
select
  json_extract(severity, '$.Original') as severity_original,
  count(json_extract(severity, '$.Original'))
from
  aws_securityhub_finding
group by
  json_extract(severity, '$.Original')
order by
  json_extract(severity, '$.Original');
```

### List findings with failed compliance status
Determine the areas in which security findings have failed compliance checks. This is useful for identifying potential vulnerabilities and areas that require immediate attention in your AWS Security Hub.

```sql+postgres
select
  title,
  product_arn,
  product_name,
  compliance ->> 'Status' as compliance_status,
  compliance ->> 'StatusReasons' as compliance_status_reasons
from
  aws_securityhub_finding
where
  compliance ->> 'Status' = 'FAILED';
```

```sql+sqlite
select
  title,
  product_arn,
  product_name,
  json_extract(compliance, '$.Status') as compliance_status,
  json_extract(compliance, '$.StatusReasons') as compliance_status_reasons
from
  aws_securityhub_finding
where
  json_extract(compliance, '$.Status') = 'FAILED';
```

### List findings with malware
Discover the segments that have detected malware within your AWS SecurityHub findings. This allows for quick identification of potential security threats, helping to safeguard your cloud environment.

```sql+postgres
select
  title,
  product_arn,
  product_name,
  malware
from
  aws_securityhub_finding
where
  malware is not null;
```

```sql+sqlite
select
  title,
  product_arn,
  product_name,
  malware
from
  aws_securityhub_finding
where
  malware is not null;
```

### List critical findings from the last 10 days
Explore critical security issues from the past 10 days to understand potential threats. This query is beneficial in identifying and addressing major vulnerabilities swiftly to enhance overall security.

```sql+postgres
select
  title,
  product_arn,
  product_name,
  severity ->> 'Original' as severity_original
from
  aws_securityhub_finding
where
  severity ->> 'Original' = 'CRITICAL'
and
  created_at >= now() - interval '10' day;
```

```sql+sqlite
select
  title,
  product_arn,
  product_name,
  json_extract(severity, '$.Original') as severity_original
from
  aws_securityhub_finding
where
  json_extract(severity, '$.Original') = 'CRITICAL'
and
  created_at >= datetime('now', '-10 day');
```

### List findings ordered by criticality
Explore the security findings from AWS Security Hub, with a focus on identifying the most critical issues first. This can aid in prioritizing security measures and addressing the most severe vulnerabilities or threats promptly.

```sql+postgres
select
  title,
  product_arn,
  product_name,
  criticality
from
  aws_securityhub_finding
order by
  criticality desc nulls last;
```

```sql+sqlite
select
  title,
  product_arn,
  product_name,
  criticality
from
  aws_securityhub_finding
order by
  case when criticality is null then 1 else 0 end, criticality desc;
```

### List findings for Turbot company
Explore the security findings related to the company 'Turbot' to gain insights into potential vulnerabilities or breaches. This is crucial for maintaining robust security practices and addressing any potential threats promptly.

```sql+postgres
select
  title,
  id,
  product_arn,
  product_name,
  company_name
from
  aws_securityhub_finding
where
  company_name = 'Turbot';
```

```sql+sqlite
select
  title,
  id,
  product_arn,
  product_name,
  company_name
from
  aws_securityhub_finding
where
  company_name = 'Turbot';
```

### List findings updated in the last 30 days
Explore the recent updates made within the last month in your security findings. This can help identify any potential security threats or issues that have surfaced recently, thereby allowing for timely response and mitigation.

```sql+postgres
select
  title,
  product_arn,
  product_name,
  updated_at
from
  aws_securityhub_finding
where
   updated_at >= now() - interval '30' day;
```

```sql+sqlite
select
  title,
  product_arn,
  product_name,
  updated_at
from
  aws_securityhub_finding
where
   updated_at >= datetime('now','-30 day');
```

### List findings with workflow status NOTIFIED
Explore which security findings in your AWS environment have been flagged with a status of 'Notified'. This allows you to identify and focus on issues that have been brought to your attention, facilitating a more efficient security management process.

```sql+postgres
select
  title,
  id,
  product_arn,
  product_name,
  workflow_status
from
  aws_securityhub_finding
where
  workflow_status = 'NOTIFIED';
```

```sql+sqlite
select
  title,
  id,
  product_arn,
  product_name,
  workflow_status
from
  aws_securityhub_finding
where
  workflow_status = 'NOTIFIED';
```

### Get network detail for a particular finding
Explore the network details associated with a specific security finding to gain insights into potential security threats. This can assist in identifying the source and destination of any potential attacks, helping to bolster security measures.

```sql+postgres
select
  title,
  id,
  network ->> 'DestinationDomain' as network_destination_domain,
  network ->> 'DestinationIpV4' as network_destination_ip_v4,
  network ->> 'DestinationIpV6' as network_destination_ip_v6,
  network ->> 'DestinationPort' as network_destination_port,
  network ->> 'Protocol' as network_protocol,
  network ->> 'SourceIpV4' as network_source_ip_v4,
  network ->> 'SourceIpV6' as network_source_ip_v6,
  network ->> 'SourcePort' as network_source_port
from
  aws_securityhub_finding
where
  title = 'EC2 instance involved in SSH brute force attacks.';
```

```sql+sqlite
select
  title,
  id,
  json_extract(network, '$.DestinationDomain') as network_destination_domain,
  json_extract(network, '$.DestinationIpV4') as network_destination_ip_v4,
  json_extract(network, '$.DestinationIpV6') as network_destination_ip_v6,
  json_extract(network, '$.DestinationPort') as network_destination_port,
  json_extract(network, '$.Protocol') as network_protocol,
  json_extract(network, '$.SourceIpV4') as network_source_ip_v4,
  json_extract(network, '$.SourceIpV6') as network_source_ip_v6,
  json_extract(network, '$.SourcePort') as network_source_port
from
  aws_securityhub_finding
where
  title = 'EC2 instance involved in SSH brute force attacks.';
```

### Get patch summary details for a particular finding
Determine the status and details of patch installations for a specific security finding, especially useful when assessing the severity and impact of potential security threats.

```sql+postgres
select
  title,
  id,
  patch_summary ->> 'Id' as patch_id,
  patch_summary ->> 'FailedCount' as failed_count,
  patch_summary ->> 'InstalledCount' as installed_count,
  patch_summary ->> 'InstalledOtherCount' as installed_other_count,
  patch_summary ->> 'InstalledPendingReboot' as installed_pending_reboot,
  patch_summary ->> 'InstalledRejectedCount' as installed_rejected_count,
  patch_summary ->> 'MissingCount' as missing_count,
  patch_summary ->> 'Operation' as operation,
  patch_summary ->> 'OperationEndTime' as operation_end_time,
  patch_summary ->> 'OperationStartTime' as operation_start_time,
  patch_summary ->> 'RebootOption' as reboot_option
from
  aws_securityhub_finding
where
  title = 'EC2 instance involved in SSH brute force attacks.';
```

```sql+sqlite
select
  title,
  id,
  json_extract(patch_summary, '$.Id') as patch_id,
  json_extract(patch_summary, '$.FailedCount') as failed_count,
  json_extract(patch_summary, '$.InstalledCount') as installed_count,
  json_extract(patch_summary, '$.InstalledOtherCount') as installed_other_count,
  json_extract(patch_summary, '$.InstalledPendingReboot') as installed_pending_reboot,
  json_extract(patch_summary, '$.InstalledRejectedCount') as installed_rejected_count,
  json_extract(patch_summary, '$.MissingCount') as missing_count,
  json_extract(patch_summary, '$.Operation') as operation,
  json_extract(patch_summary, '$.OperationEndTime') as operation_end_time,
  json_extract(patch_summary, '$.OperationStartTime') as operation_start_time,
  json_extract(patch_summary, '$.RebootOption') as reboot_option
from
  aws_securityhub_finding
where
  title = 'EC2 instance involved in SSH brute force attacks.';
```

### Get vulnerabilities for a finding
Discover the details of potential security vulnerabilities linked to a specific finding. This query is useful in understanding the severity, vendor information, and associated vulnerabilities tied to a particular security incident, such as an SSH brute force attack on an EC2 instance.

```sql+postgres
select
  title,
  v ->> 'Id' as vulnerabilitie_id,
  v -> 'Vendor' ->> 'Name' as vendor_name,
  v -> 'Vendor' ->> 'Url' as vendor_url,
  v -> 'Vendor' ->> 'VendorCreatedAt' as vendor_created_at,
  v -> 'Vendor' ->> 'VendorSeverity' as vendor_severity,
  v -> 'Vendor' ->> 'VendorUpdatedAt' as vendor_updated_at,
  v ->> 'Cvss' as cvss,
  v ->> 'ReferenceUrls' as reference_urls,
  v ->> 'RelatedVulnerabilities' as related_vulnerabilities,
  v ->> 'VulnerablePackages' as vulnerable_packages
from
  aws_securityhub_finding,
  jsonb_array_elements(vulnerabilities) as v
where
  title = 'EC2 instance involved in SSH brute force attacks.';
```

```sql+sqlite
select
  title,
  json_extract(v.value, '$.Id') as vulnerabilitie_id,
  json_extract(v.value, '$.Vendor.Name') as vendor_name,
  json_extract(v.value, '$.Vendor.Url') as vendor_url,
  json_extract(v.value, '$.Vendor.VendorCreatedAt') as vendor_created_at,
  json_extract(v.value, '$.Vendor.VendorSeverity') as vendor_severity,
  json_extract(v.value, '$.Vendor.VendorUpdatedAt') as vendor_updated_at,
  json_extract(v.value, '$.Cvss') as cvss,
  json_extract(v.value, '$.ReferenceUrls') as reference_urls,
  json_extract(v.value, '$.RelatedVulnerabilities') as related_vulnerabilities,
  json_extract(v.value, '$.VulnerablePackages') as vulnerable_packages
from
  aws_securityhub_finding,
  json_each(vulnerabilities) as v
where
  title = 'EC2 instance involved in SSH brute force attacks.';
```

### List EC2 instances with failed compliance status
Determine the areas in which EC2 instances have failed compliance checks. This is useful for identifying potential security risks and rectifying them to maintain the integrity of your AWS environment.

```sql+postgres
select
  distinct i.instance_id,
  i.instance_state,
  i.instance_type,
  f.title,
  f.compliance_status,
  f.severity ->> 'Original' as severity_original
from
  aws_ec2_instance as i,
  aws_securityhub_finding as f,
  jsonb_array_elements(resources) as r
where
  compliance_status = 'FAILED'
and
  r ->> 'Type' = 'AwsEc2Instance'
and
  i.arn = r ->> 'Id';
```

```sql+sqlite
select
  distinct i.instance_id,
  i.instance_state,
  i.instance_type,
  f.title,
  f.compliance_status,
  json_extract(f.severity, '$.Original') as severity_original
from
  aws_ec2_instance as i,
  aws_securityhub_finding as f,
  json_each(f.resources) as r
where
  f.compliance_status = 'FAILED'
and
  json_extract(r.value, '$.Type') = 'AwsEc2Instance'
and
  i.arn = json_extract(r.value, '$.Id');
```

### Count resources with failed compliance status
Determine areas in your AWS environment where resources have failed compliance checks. This allows you to identify potential security risks and take corrective action.

```sql+postgres
select
  r ->> 'Type' as resource_type,
  count(r ->> 'Type')
from
  aws_securityhub_finding,
  jsonb_array_elements(resources) as r
group by
  r ->> 'Type'
order by
  count desc;
```

```sql+sqlite
select
  json_extract(r.value, '$.Type') as resource_type,
  count(json_extract(r.value, '$.Type'))
from
  aws_securityhub_finding,
  json_each(resources) as r
group by
  json_extract(r.value, '$.Type')
order by
  count(json_extract(r.value, '$.Type')) desc;
```

### List findings for CIS AWS foundations benchmark
Explore which findings are associated with the CIS AWS foundations benchmark in your AWS Security Hub. This can assist in identifying potential security risks or non-compliance issues for your company.

```sql+postgres
 select
  title,
  id,
  company_name,
  created_at,
  criticality,
  confidence
from
  aws_securityhub_finding
where
  standards_control_arn like '%cis-aws-foundations-benchmark%';
```

```sql+sqlite
select
  title,
  id,
  company_name,
  created_at,
  criticality,
  confidence
from
  aws_securityhub_finding
where
  standards_control_arn like '%cis-aws-foundations-benchmark%';
```

### List findings for a particular standard control (Config.1)
Identify instances where specific security findings are associated with a particular standard control. This is beneficial in understanding the security posture of your organization by analyzing the criticality and confidence of the findings.

```sql+postgres
 select
  f.title,
  f.id,
  f.company_name,
  f.created_at,
  f.criticality,
  f.confidence
from
  aws_securityhub_finding as f,
  aws_securityhub_standards_control as c
where
  c.arn = f.standards_control_arn
and
  c.control_id = 'Config.1';
```

```sql+sqlite
select
  f.title,
  f.id,
  f.company_name,
  f.created_at,
  f.criticality,
  f.confidence
from
  aws_securityhub_finding as f,
  aws_securityhub_standards_control as c
where
  c.arn = f.standards_control_arn
and
  c.control_id = 'Config.1';
```

### List resources with a failed compliance status for CIS AWS foundations benchmark
Discover the segments that have failed to comply with the CIS AWS foundations benchmark. This allows you to identify and rectify areas in your AWS resources that are not adhering to the recommended security controls, thereby enhancing your overall security posture.

```sql+postgres
select
  distinct r ->> 'Id' as resource_arn,
  r ->> 'Type' as resource_type,
  f.title,
  f.compliance_status,
  f.severity ->> 'Original' as severity_original
from
  aws_securityhub_finding as f,
  jsonb_array_elements(resources) as r
where
  f.compliance_status = 'FAILED'
and
  standards_control_arn like '%cis-aws-foundations-benchmark%';
```

```sql+sqlite
select
  distinct json_extract(r.value, '$.Id') as resource_arn,
  json_extract(r.value, '$.Type') as resource_type,
  f.title,
  f.compliance_status,
  json_extract(f.severity, '$.Original') as severity_original
from
  aws_securityhub_finding as f,
  json_each(f.resources) as r
where
  f.compliance_status = 'FAILED'
and
  f.standards_control_arn like '%cis-aws-foundations-benchmark%';
```

### List findings for production resources
Uncover the details of potential security issues within your production resources. This query aids in identifying non-compliant resources, assessing the severity of the issues, and understanding the areas that require immediate attention to enhance the security posture.

```sql+postgres
select
  distinct r ->> 'Id' as resource_arn,
  r ->> 'Type' as resource_type,
  f.title,
  f.compliance_status,
  f.severity ->> 'Original' as severity_original
from
  aws_securityhub_finding as f,
  jsonb_array_elements(resources) as r
where
  r -> 'Tags' ->> 'Environment' = 'PROD';
```

```sql+sqlite
select
  distinct json_extract(r.value, '$.Id') as resource_arn,
  json_extract(r.value, '$.Type') as resource_type,
  f.title,
  f.compliance_status,
  json_extract(f.severity, '$.Original') as severity_original
from
  aws_securityhub_finding as f,
  json_each(f.resources) as r
where
  json_extract(json_extract(r.value, '$.Tags'), '$.Environment') = 'PROD';
```

### Count finding resources by environment tag
This query helps identify the number of security findings associated with different environments in your AWS infrastructure. It's useful for understanding the distribution of potential security issues across various operational contexts.

```sql+postgres
select
  r -> 'Tags' ->> 'Environment' as environment,
  count(r ->> 'Tags')
from
  aws_securityhub_finding as f,
  jsonb_array_elements(resources) as r
group by
  r -> 'Tags' ->> 'Environment'
order by
  count desc;
```

```sql+sqlite
select
  json_extract(r.value, '$.Tags.Environment') as environment,
  count(json_extract(r.value, '$.Tags'))
from
  aws_securityhub_finding as f,
  json_each(f.resources) as r
group by
  json_extract(r.value, '$.Tags.Environment')
order by
  count(*) desc;
```

### List all findings for affected account 0123456789012
Determine the areas in which security issues have been identified for a specific account in AWS Security Hub. This is beneficial for pinpointing and addressing vulnerabilities in a targeted manner.

```sql+postgres
select
  title,
  f.severity ->> 'Original' as severity,
  r ->> 'Type' as resource_type,
  source_account_id
from
  aws_securityhub_finding as f,
  jsonb_array_elements(resources) r
where
  source_account_id = '0123456789012';
```

```sql+sqlite
select
  title,
  json_extract(f.severity, '$.Original') as severity,
  json_extract(r.value, '$.Type') as resource_type,
  source_account_id
from
  aws_securityhub_finding as f,
  json_each(resources) as r
where
  source_account_id = '0123456789012';
```

### Count the number of findings by affected account
Discover the segments that have security findings in your AWS accounts. This query helps you identify which accounts have the most findings, allowing you to prioritize your security efforts.

```sql+postgres
select
  source_account_id,
  count(*) as finding_count
from
  aws_securityhub_finding
group by
  source_account_id
order by
  source_account_id;
```

```sql+sqlite
select
  source_account_id,
  count(*) as finding_count
from
  aws_securityhub_finding
group by
  source_account_id
order by
  source_account_id;
```

### Retrieve findings updated within a specific time range
This query retrieves AWS Security Hub findings that were updated within a specified time interval. It provides details such as the ID, company name, the first time the finding was observed, the last update time, criticality, and verification state.

```sql+postgres
select
  id,
  company_name,
  first_observed_at,
  updated_at,
  criticality,
  verification_state
from
  aws_securityhub_finding
where
  updated_at between '2023-06-26T13:00:21+05:30' and '2024-07-04T14:45:00+05:30';
```

```sql+sqlite
select
  id,
  company_name,
  first_observed_at,
  updated_at,
  criticality,
  verification_state
from
  aws_securityhub_finding
where
  updated_at between '2023-06-26T13:00:21+05:30' and '2024-07-04T14:45:00+05:30';
```

### List findings that are created in the last month
This query is useful for retrieving security findings from AWS Security Hub that were created within the last 30 days. It allows you to filter and monitor findings based on their creation date, helping identify recent security issues, assess compliance status, and track product-specific incidents. This can be particularly valuable for auditing, incident response, or compliance reporting, ensuring you're working with the most recent data.

```sql+postgres
select
  id,
  company_name,
  created_at,
  confidence,
  compliance_status,
  product_name,
  product_arn
from
  aws_securityhub_finding
where
  created_at >= now() - interval '30d';
```

```sql+sqlite
select
  id,
  company_name,
  created_at,
  confidence,
  compliance_status,
  product_name,
  product_arn
from
  aws_securityhub_finding
where
  created_at >= datetime('now', '-30 days');
```