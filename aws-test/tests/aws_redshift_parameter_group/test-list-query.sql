select name , partition, region, tags, title
from aws.aws_redshift_parameter_group
where akas::text = '["{{ output.resource_aka.value }}"]'