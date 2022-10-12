# Tests: aws_identitystore_user

The tests for `aws_identitystore_user` require manual setup (Terraform is unable to create all required resources due to AWS API limitations):

* Enable AWS SSO: https://docs.aws.amazon.com/singlesignon/latest/userguide/step1.html
* Create an Identity Store user: https://docs.aws.amazon.com/singlesignon/latest/userguide/addusers.html
  * Username: TestUser
