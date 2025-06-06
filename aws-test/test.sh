 #!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 YELLOW="\e[93m"
 ENDCOLOR="\e[0m"

# Define your function here
run_test () {
   echo -e "${YELLOW}Running $1 ${ENDCOLOR}"
 if ! ./tint.js $1 >> output.txt
   then
    echo -e "${RED}Failed -> $1 ${ENDCOLOR}"
    echo $1 >> failed_tests.txt
  else
    echo -e "${BOLDGREEN}Passed -> $1 ${ENDCOLOR}"
    echo $1 >> passed_tests.txt
   fi
 }

 # output.txt - store output of each test
 # failed_tests.txt - names of failed test
 # passed_tests.txt names of passed test

 # removes files from previous test
rm -rf output.txt failed_tests.txt passed_tests.txt
date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt

run_test aws_accessanalyzer_analyzer
run_test aws_acm_certificate
run_test aws_amplify_app
run_test aws_api_gateway_api_key
run_test aws_api_gateway_authorizer
run_test aws_api_gateway_rest_api
run_test aws_api_gateway_stage
run_test aws_api_gateway_usage_plan
run_test aws_api_gatewayv2_api
run_test aws_api_gatewayv2_domain_name
run_test aws_api_gatewayv2_integration
run_test aws_api_gatewayv2_route
run_test aws_api_gatewayv2_stage
run_test aws_appautoscaling_policy
run_test aws_appautoscaling_target
run_test aws_appconfig_application
run_test aws_appstream_fleet
run_test aws_auditmanager_assessment
run_test aws_auditmanager_control
run_test aws_auditmanager_framework
run_test aws_backup_framework
run_test aws_backup_plan
run_test aws_backup_report_plan
run_test aws_backup_selection
run_test aws_backup_vault
run_test aws_cloudformation_stack
run_test aws_cloudformation_stack_resource
run_test aws_cloudformation_stack_set
run_test aws_cloudfront_cache_policy
run_test aws_cloudfront_distribution
run_test aws_cloudfront_function
run_test aws_cloudfront_origin_access_identity
run_test aws_cloudfront_origin_request_policy
run_test aws_cloudfront_response_headers_policy
run_test aws_cloudsearch_domain
run_test aws_cloudtrail_event_data_store
run_test aws_cloudtrail_trail
run_test aws_cloudwatch_alarm
run_test aws_cloudwatch_log_group
run_test aws_cloudwatch_log_resource_policy
run_test aws_cloudwatch_log_stream
run_test aws_codeartifact_domain
run_test aws_codeartifact_repository
run_test aws_codebuild_project
run_test aws_codebuild_source_credential
run_test aws_codecommit_repository
run_test aws_codedeploy_app
run_test aws_codedeploy_deployment_config
run_test aws_codedeploy_deployment_group
run_test aws_codepipeline_pipeline
run_test aws_config_aggregate_authorization
run_test aws_config_configuration_recorder
run_test aws_config_conformance_pack
run_test aws_config_rule


date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt