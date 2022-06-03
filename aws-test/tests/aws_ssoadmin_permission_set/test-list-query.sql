select ps.arn, description, name, relay_state, session_duration
from aws.aws_ssoadmin_permission_set as ps
join aws.aws_ssoadmin_instance as i on ps.instance_arn = i.arn
where ps.name = '{{resourceName}}';
