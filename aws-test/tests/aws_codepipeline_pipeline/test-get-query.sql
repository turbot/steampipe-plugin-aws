select name, arn, encryption_key, role_arn, tags_src, title
from aws_codepipeline_pipeline
where name = '{{ resourceName }}';
