select akas, version_label, region, tags, title
from aws_elastic_beanstalk_application_version
where version_label = '{{ output.version_label.value }}';