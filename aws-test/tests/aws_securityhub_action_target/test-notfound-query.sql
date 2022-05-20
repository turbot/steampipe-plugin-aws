select name, akas, region
from aws_securityhub_action_target
where arn = 'TestNotFound';
