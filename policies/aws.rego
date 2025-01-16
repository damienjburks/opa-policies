package aws.terraform

import input.plan as plan

# Enforces versioning to be enabled on S3 buckets
deny[msg] {
  resource := plan.resource_changes[_]
  resource.type == "aws_s3_bucket"
  not resource.change.after.versioning.enabled
  msg := sprintf("S3 bucket '%s' must have versioning enabled.", [resource.address])
}

# Denies all instance types other than t2.micro
deny[msg] {
  resource := plan.resource_changes[_]
  resource.type == "aws_instance"
  resource.change.after.instance_type != "t2.micro"
  msg := sprintf("EC2 instance '%s' must use instance type 't2.micro'.", [resource.address])
}

# Denies untagged resources
deny[msg] {
  resource := plan.resource_changes[_]
  not resource.change.after.tags
  msg := sprintf("Resource '%s' must have tags.", [resource.address])
}