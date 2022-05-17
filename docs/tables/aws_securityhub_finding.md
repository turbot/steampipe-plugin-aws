# Table: aws_securityhub_finding

AWS Security Hub eliminates the complexity of addressing large volumes of findings from multiple providers. It reduces the effort required to manage and improve the security of all of your AWS accounts, resources, and workloads.

## Examples

### Basic info

```sql
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

```sql
select
  title,
  product_arn,
  product_name,
  severity ->> 'Original' as severity_original
from
  aws_securityhub_finding
where
  severity ->> 'Original' = 'HIGH'
```

### List findings with failed compliance

```sql
select
  title,
  product_arn,
  product_name,
  compliance ->> 'Status' as compliance_status,
  compliance ->> 'StatusReasons'as compliance_status_reason
from
  aws_securityhub_finding
where
  compliance ->> 'Status' = 'FAILED'
```

### List findings with malware

```sql
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

### List last 10 days critical findings

```sql
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
   date_part('day', now() - created_at) <= 10;
```

### List findings with highest criticality

```sql
select
  title,
  product_arn,
  product_name,
  criticality
from
  aws_securityhub_finding
order by criticality desc;
```

### List findings for company Turbot

```sql
select
  title,
  id,
  product_arn,
  product_name,
  company_name
from
  aws_securityhub_finding
where company_name = 'Turbot';
```

### List findings which are updated in last 30 days

```sql
select
  title,
  product_arn,
  product_name,
  updated_at
from
  aws_securityhub_finding
where
   date_part('day', now() - updated_at) <= 30;
```

### List findings with assigned workflow state

```sql
select
  title,
  id,
  product_arn,
  product_name,
  workflow_state
from
  aws_securityhub_finding
where workflow_state = 'ASSIGNED';
```

### Get network detail for a particular finding

```sql
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
where title = 'EC2 instance involved in SSH brute force attacks.';
```

### Get patch summary for a particular finding

```sql
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
where title = 'EC2 instance involved in SSH brute force attacks.';
```

### Get vulnerabilities for a finding

```sql
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
where title = 'EC2 instance involved in SSH brute force attacks.';
```