select name, title, tags, akas, account_id
from aws.aws_eks_cluster
where name = 'dummy-{{ resourceName }}';