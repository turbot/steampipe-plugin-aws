select
  name,
  jsonb_pretty(policy) as policy
from
  aws_s3_bucket;

select
  jsonb_pretty(policy),
  is_public,
  access_level,
  allowed_organization_ids,
  allowed_principal_account_ids,
  allowed_principal_federated_identities,
  allowed_principal_services,
  allowed_principals,
  public_access_levels,
  public_statement_ids
from
  aws_resource_policy_analysis
where
  policy = '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"OrganizationAccess","Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test","Condition":{"StringEquals":{"aws:PrincipalOrgID":["o-123456"]}}},{"Sid":"AccountPrincipals","Effect":"Allow","Principal":{"AWS":["arn:aws:iam::123456789012:user/victor@xyz.com","arn:aws:iam::111122223333:root"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"FederatedPrincipals","Effect":"Allow","Principal":{"Federated":"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"ServicePrincipals","Effect":"Allow","Principal":{"Service":["ecs.amazonaws.com","elasticloadbalancing.amazonaws.com"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"PublicAccess","Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"}]}';

select
  jsonb_pretty(policy),
  is_public,
  access_level,
  allowed_organization_ids,
  allowed_principal_account_ids,
  allowed_principal_federated_identities,
  allowed_principal_services,
  allowed_principals,
  public_access_levels,
  public_statement_ids
from
  aws_resource_policy_analysis
where
  policy = '{"Statement":[{"Action":"s3:GetBucketAcl","Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::aws-cloudtrail-logs-013122550996-1247fc6c","Sid":"AWSCloudTrailAclCheck20150319"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"AWS:SourceArn":"arn:aws:cloudtrail:ap-south-1:013122550996:trail/management-events","s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::aws-cloudtrail-logs-013122550996-1247fc6c/AWSLogs/013122550996/*","Sid":"AWSCloudTrailWrite20150319"},{"Action":"s3:PutObject","Condition":{"StringEquals":{"AWS:SourceArn":"arn:aws:cloudtrail:us-east-1:013122550996:trail/test-sd-pci","s3:x-amz-acl":"bucket-owner-full-control"}},"Effect":"Allow","Principal":{"Service":"cloudtrail.amazonaws.com"},"Resource":"arn:aws:s3:::aws-cloudtrail-logs-013122550996-1247fc6c/AWSLogs/013122550996/*","Sid":"AWSCloudTrailWrite20150319"}],"Version":"2012-10-17"}';

-- allowed_organization_ids,
-- allowed_principal_account_ids,
-- allowed_principal_federated_identities,
-- allowed_principal_services,
-- allowed_principals,
-- public_access_levels,
select
  s3.name as bucket,
  rpa.is_public,
  access_level,
  public_statement_ids
from
  _aws_s3_bucket as s3
  left join aws_resource_policy_analysis as rpa on s3.policy = rpa.policy;

select is_public, access_level, allowed_organization_ids, allowed_principal_account_ids, allowed_principal_federated_identities, allowed_principal_services, allowed_principals, public_access_levels, public_statement_ids from aws_resource_policy_analysis where policy = '{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"OrganizationAccess","Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test","Condition":{"StringEquals":{"aws:PrincipalOrgID":["o-123456"]}}},{"Sid":"AccountPrincipals","Effect":"Allow","Principal":{"AWS":["arn:aws:iam::123456789012:user/victor@xyz.com","arn:aws:iam::111122223333:root"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"FederatedPrincipals","Effect":"Allow","Principal":{"Federated":"arn:aws:iam::111011101110:saml-provider/AWSSSO_DO_NOT_DELETE"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"ServicePrincipals","Effect":"Allow","Principal":{"Service":["ecs.amazonaws.com","elasticloadbalancing.amazonaws.com"]},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"},{"Sid":"PublicAccess","Effect":"Allow","Principal":{"AWS":"*"},"Action":["s3:GetBucketLocation","s3:ListBucket"],"Resource":"arn:aws:s3:::test"}]}'