select name, arn, akas
from aws.aws_codepipeline_pipeline
where name = 'dummy';