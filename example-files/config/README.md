

Done:
* CloudFront logging must be enabled

* ELB must access logging should be enabled

* EBS Volumes must be encrypted

* IamPolicyNotActionRule
* IamPolicyNotResourceRule
* IamPolicyWildcardActionRule
* IamPolicyWildcardResourceRule

* IamRoleNotActionOnPermissionsPolicyRule
* IamRoleNotResourceOnPermissionsPolicyRule

* S3BucketPolicyNotActionRule
* S3BucketPolicyNotPrincipalRule
* S3BucketPolicyWildcardActionRule
* S3BucketPolicyWildcardPrincipalRule

* SecurityGroupEgressOpenToWorldRule
* SecurityGroupEgressPortRangeRule
* SecurityGroupIngressCidrNon32Rule - but what is standalong ingress in cfn_nag?
* SecurityGroupIngressOpenToWorldRule
* SecurityGroupIngressPortRangeRule
* SecurityGroupMissingEgressRule

* SnsTopicPolicyNotActionRule
* SnsTopicPolicyNotPrincipalRule
* SnsTopicPolicyWildcardPrincipalRule

* SqsQueuePolicyNotActionRule
* SqsQueuePolicyNotPrincipalRule
* SqsQueuePolicyWildcardActionRule
* SqsQueuePolicyWildcardPrincipalRule

TODO
* CloudFront resource !Metadata['AWS::CloudFront::Authentication'].nil?  How to specify in Terraform?

* IamManagedPolicyNotActionRule  - How is this different than a plain IamPolicy?
* IamManagedPolicyNotResourceRule
* IamManagedPolicyWildcardActionRule
* IamManagedPolicyWildcardResourceRule

* IamRoleNotActionOnTrustPolicyRule
* IamRoleNotPrincipalOnTrustPolicyRule
* IamRoleWildcardActionOnPermissionsPolicyRule
* IamRoleWildcardActionOnTrustPolicyRule
* IamRoleWildcardResourceOnPermissionsPolicyRule

* LambdaPermissionInvokeFunctionActionRule
* LambdaPermissionWildcardPrincipalRule

* ManagedPolicyOnUserRule
* PolicyOnUserRule

* S3BucketPublicReadAclRule
* S3BucketPublicReadWriteAclRule


* UserHasInlinePolicyRule
* UserMissingGroupRule

* WafWebAclDefaultActionRule
