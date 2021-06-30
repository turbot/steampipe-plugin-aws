select name, arn, title, tags
from aws.aws_codepipeline_pipeline
where akas::text = '["{{ output.resource_aka.value }}"]';