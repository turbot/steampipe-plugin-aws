select name, description, artifacts, cache, source, title
from aws.aws_codebuild_project
where name = '{{ resourceName }}';
