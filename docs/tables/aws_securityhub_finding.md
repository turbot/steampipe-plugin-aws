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
  severity ->> 'Original' = 'HIGH';
```

### Count the number of findings by severity

```sql
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

### List findings with failed compliance status

```sql
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

### List critical findings from the last 10 days

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
  created_at >= now() - interval '10' day;
```

### List findings ordered by criticality

```sql
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

### List findings for Turbot company

```sql
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

```sql
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

### DEPRECATED, List findings with assigned workflow state

```sql
select
  title,
  id,
  product_arn,
  product_name,
  workflow_state
from
  aws_securityhub_finding
where 
  workflow_state = 'ASSIGNED';
```

### List findings with NOTIFIED workflow status

```sql
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
where 
  title = 'EC2 instance involved in SSH brute force attacks.';
```

### Get patch summary details for a particular finding

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
where 
  title = 'EC2 instance involved in SSH brute force attacks.';
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
where 
  title = 'EC2 instance involved in SSH brute force attacks.';
```

### List EC2 instances with failed compliance status

```sql
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

### Count resources with failed compliance status

```sql
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

### List findings for CIS AWS foundations benchmark

```sql
 select 
  title,
  id,
  company_name,
  created_at,
  criticality,
  confidence
from 
  aws_redhood.aws_securityhub_finding
where 
  standards_control_arn like '%cis-aws-foundations-benchmark%';
```

### List findings for a particular standard control (Config.1)

```sql
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

```sql
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

### List findings for production resources

```sql
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

### Count finding resources by environment tag

```sql
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
