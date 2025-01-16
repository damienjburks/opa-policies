package digitalocean.terraform

# Deny rule that aggregates messages from all sub-rules
deny[msg] {
  insufficient_droplet_size[msg]
}

deny[msg] {
  missing_tags[msg]
}

deny[msg] {
  public_ipv4_disallowed[msg]
}

# Rule: Enforces a minimum droplet size
insufficient_droplet_size[msg] {
  resource := input.resource_changes[_]
  resource.type == "digitalocean_droplet"
  sizes := {"s-1vcpu-1gb", "s-1vcpu-2gb"}  # Define unacceptable sizes
  not resource.change.after.size in sizes
  msg := sprintf("Droplet '%s' must have a size of either 's-1vcpu-1gb' or 's-1vcpu-2gb'.", [resource.address])
}

# Rule: Ensures tags are applied to droplets
missing_tags[msg] {
  resource := input.resource_changes[_]
  resource.type == "digitalocean_droplet"
  not resource.change.after.tags  # Check if tags are missing
  msg := sprintf("Droplet '%s' must have at least one tag.", [resource.address])
}

# Rule: Disallow public IPv4 addresses
public_ipv4_disallowed[msg] {
  resource := input.resource_changes[_]
  resource.type == "digitalocean_droplet"
  resource.change.after.ipv4_address != null  # Droplet has a public IPv4 address
  msg := sprintf("Droplet '%s' must not have a public IPv4 address.", [resource.address])
}
