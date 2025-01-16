package gcp.terraform

# Deny rule that aggregates all policy violations
deny[msg] {
  insecure_cloud_function[msg]
}

# Rule: Cloud Functions must not allow unauthenticated access
insecure_cloud_function[msg] {
  resource := input.resource_changes[_]
  resource.type == "google_cloudfunctions_function"
  resource.change.after.entry_point != null
  resource.change.after.https_trigger.security_level != "SECURE_ALWAYS"
  msg := sprintf("Cloud Function '%s' must enforce secure HTTPS connections.", [resource.address])
}