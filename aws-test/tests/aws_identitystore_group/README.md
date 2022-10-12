# Tests: aws_identitystore_group

The tests for `aws_identitystore_group` require manual setup (Terraform is unable to create all required resources due to AWS API limitations):

* Enable AWS SSO: https://docs.aws.amazon.com/singlesignon/latest/userguide/step1.html
* Create an Identity Store group: https://docs.aws.amazon.com/singlesignon/latest/userguide/addgroups.html
  * Group name: TestGroup
