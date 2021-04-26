select name, cache, source, title, akas
from aws.aws_codebuild_project
where name = 'dummy';