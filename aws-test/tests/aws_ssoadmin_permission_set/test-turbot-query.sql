select ps.title, ps.tags, ps.akas, ps.partition, ps.region, ps.account_id
from aws.aws_ssoadmin_permission_set as ps
join aws.aws_ssoadmin_instance as i on ps.instance_arn = i.arn
where ps.name = '{{resourceName}}';
