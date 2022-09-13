select akas, name, region, title
from aws.aws_securityhub_standards_subscription
where name = 'PCI DSS v3.2.1' and region = '{{ output.aws_region.value }}';
