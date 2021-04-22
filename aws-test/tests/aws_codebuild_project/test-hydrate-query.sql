select name, description, akas, tags, title
from aws.aws_codebuild_project
where name = '{{ resourceName }}';