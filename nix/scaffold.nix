{ scaffold }:
scaffold {
  command = "resource";
  name = "firewall_rule_resource";
  scaffoldName = "FirewallRule";
  package = "resource_firewall_rule";
  # a2b scaffold does not pre-create $out; preRun hook does it
  env.preRun = "mkdir -p $out";
}
