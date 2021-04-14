# Table: aws_guardduty_threat_intel_set

ThreatIntelSet consists of known malicious IP addresses. GuardDuty generates findings based on the threatIntelSet when it is activated.

## Examples

### Basic info

```sql
select
  detector_id,
  threat_intel_set_id,
  name,
  format,
  location
from
  aws_guardduty_threat_intel_set;
```


### List disabled threat intel sets

```sql
select
  threat_intel_set_id,
  status
from
  aws_guardduty_threat_intel_set
where
  status = 'INACTIVE';
```
