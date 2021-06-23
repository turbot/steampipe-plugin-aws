select name, akas, tags, title
from aws_codepipeline_pipeline
where name = '{{ resourceName }}';