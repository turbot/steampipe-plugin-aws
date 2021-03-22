select hub_arn, auto_enable_controls, akas, region
from aws_securityhub_hub
where akas = '["arn:aws:securityhub:us-east-1:352685002396:hub/default"]';
