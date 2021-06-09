select id, name
from aws.aws_guardduty_finding
where id = '{{ output.resource_id.value }}::dummy';
