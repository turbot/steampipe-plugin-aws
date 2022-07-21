select name, description, artifacts, title, tags
from aws.aws_codebuild_project
where akas::text = '["{{ output.resource_aka.value }}"]';