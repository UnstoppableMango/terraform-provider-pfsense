# Resource paths at REST API v2.8.1

Every path in the pinned spec and what it maps to under [ADR 0002](adr/0002-resources-are-derived-from-endpoint-metadata.md).
Paths are shown without the `/api/v2` prefix; resource names are shown without the provider's `pfsense_` prefix.
This is a snapshot of the rule applied to v2.8.1, and has to be regenerated when the pin moves.

- **resource**: generated, with its Key from [ADR 0001](adr/0001-resources-are-identified-by-key-not-api-id.md).
- **Lookup for**: the Plural endpoint a resource's Lookup queries.
- **singleton resource**: GET and PATCH only.
- **not generated**: passes the rule, but has no Key yet and no `descr` or `description` that accepts a Marker.
- **excluded**: doesn't fit CRUD; the note says why.

| Path | Methods | Type | Model | Maps to | Notes |
| --- | --- | --- | --- | --- | --- |
| `/auth/jwt` | POST | Singular | RESTAPIJWT | excluded | auth token exchange |
| `/auth/key` | POST,DELETE | Singular | RESTAPIKey | excluded | no GET or PATCH; API key resource is fog on the map |
| `/auth/keys` | GET,DELETE | Plural | RESTAPIKey | excluded | plural of an excluded model |
| `/diagnostics/arp_table` | GET,DELETE | Plural | ARPTable | excluded | read-only state; data source candidate (later) |
| `/diagnostics/arp_table/entry` | GET,DELETE | Singular | ARPTable | excluded | read-only state; data source candidate (later) |
| `/diagnostics/command_prompt` | POST | Singular | CommandPrompt | excluded | action (POST only) |
| `/diagnostics/config_history/revision` | GET,DELETE | Singular | ConfigHistoryRevision | excluded | read-only state; data source candidate (later) |
| `/diagnostics/config_history/revisions` | GET,DELETE | Plural | ConfigHistoryRevision | excluded | read-only state; data source candidate (later) |
| `/diagnostics/halt_system` | POST | Singular | SystemHalt | excluded | action (POST only) |
| `/diagnostics/ping` | POST | Singular | Ping | excluded | action (POST only) |
| `/diagnostics/reboot` | POST | Singular | SystemReboot | excluded | action (POST only) |
| `/diagnostics/table` | GET,DELETE | Singular | Table | excluded | read-only state; data source candidate (later) |
| `/diagnostics/tables` | GET | Plural | Table | excluded | read-only state; data source candidate (later) |
| `/firewall/advanced_settings` | GET,PATCH | Singular | FirewallAdvancedSettings | singleton resource `firewall_advanced_settings` | Create and Update PATCH; Destroy forgets |
| `/firewall/alias` | GET,POST,PATCH,DELETE | Singular | FirewallAlias | resource `firewall_alias` | Key: `name` |
| `/firewall/aliases` | GET,PUT,DELETE | Plural | FirewallAlias | Lookup for `firewall_alias` |  |
| `/firewall/apply` | GET,POST | Singular | FirewallApply | excluded | apply endpoint; when to apply is #59 |
| `/firewall/nat/one_to_one/mapping` | GET,POST,PATCH,DELETE | Singular | OneToOneNATMapping | resource `firewall_nat_one_to_one_mapping` | Key: Marker |
| `/firewall/nat/one_to_one/mappings` | GET,PUT,DELETE | Plural | OneToOneNATMapping | Lookup for `firewall_nat_one_to_one_mapping` |  |
| `/firewall/nat/outbound/mapping` | GET,POST,PATCH,DELETE | Singular | OutboundNATMapping | resource `firewall_nat_outbound_mapping` | Key: Marker |
| `/firewall/nat/outbound/mappings` | GET,PUT,DELETE | Plural | OutboundNATMapping | Lookup for `firewall_nat_outbound_mapping` |  |
| `/firewall/nat/outbound/mode` | GET,PATCH | Singular | OutboundNATMode | singleton resource `firewall_nat_outbound_mode` | Create and Update PATCH; Destroy forgets |
| `/firewall/nat/port_forward` | GET,POST,PATCH,DELETE | Singular | PortForward | resource `firewall_nat_port_forward` | Key: Marker |
| `/firewall/nat/port_forwards` | GET,PUT,DELETE | Plural | PortForward | Lookup for `firewall_nat_port_forward` |  |
| `/firewall/rule` | GET,POST,PATCH,DELETE | Singular | FirewallRule | resource `firewall_rule` | Key: `tracker` |
| `/firewall/rules` | GET,PUT,DELETE | Plural | FirewallRule | Lookup for `firewall_rule` |  |
| `/firewall/schedule` | GET,POST,PATCH,DELETE | Singular | FirewallSchedule | resource `firewall_schedule` | Key: Marker |
| `/firewall/schedule/time_range` | GET,POST,PATCH,DELETE | Singular | FirewallScheduleTimeRange | not generated | no Key yet; would be `firewall_schedule_time_range`; child of FirewallSchedule |
| `/firewall/schedule/time_ranges` | GET,DELETE | Plural | FirewallScheduleTimeRange | not generated | plural of a model with no Key yet; child of FirewallSchedule |
| `/firewall/schedules` | GET,PUT,DELETE | Plural | FirewallSchedule | Lookup for `firewall_schedule` |  |
| `/firewall/state` | GET,DELETE | Singular | FirewallState | excluded | read-only state; data source candidate (later) |
| `/firewall/states` | GET,DELETE | Plural | FirewallState | excluded | read-only state; data source candidate (later) |
| `/firewall/states/size` | GET,PATCH | Singular | FirewallStatesSize | singleton resource `firewall_states_size` | Create and Update PATCH; Destroy forgets |
| `/firewall/traffic_shaper` | GET,POST,PATCH,DELETE | Singular | TrafficShaper | resource `firewall_traffic_shaper` | Key: `interface` |
| `/firewall/traffic_shaper/limiter` | GET,POST,PATCH,DELETE | Singular | TrafficShaperLimiter | resource `firewall_traffic_shaper_limiter` | Key: `name` |
| `/firewall/traffic_shaper/limiter/bandwidth` | GET,POST,PATCH,DELETE | Singular | TrafficShaperLimiterBandwidth | not generated | no Key yet; would be `firewall_traffic_shaper_limiter_bandwidth`; child of TrafficShaperLimiter |
| `/firewall/traffic_shaper/limiter/bandwidths` | GET,DELETE | Plural | TrafficShaperLimiterBandwidth | not generated | plural of a model with no Key yet; child of TrafficShaperLimiter |
| `/firewall/traffic_shaper/limiter/queue` | GET,POST,PATCH,DELETE | Singular | TrafficShaperLimiterQueue | resource `firewall_traffic_shaper_limiter_queue` | Key: `name`; child of TrafficShaperLimiter |
| `/firewall/traffic_shaper/limiter/queues` | GET,DELETE | Plural | TrafficShaperLimiterQueue | Lookup for `firewall_traffic_shaper_limiter_queue` | child of TrafficShaperLimiter |
| `/firewall/traffic_shaper/limiters` | GET,PUT | Plural | TrafficShaperLimiter | Lookup for `firewall_traffic_shaper_limiter` |  |
| `/firewall/traffic_shaper/queue` | GET,POST,PATCH,DELETE | Singular | TrafficShaperQueue | resource `firewall_traffic_shaper_queue` | Key: `name`; child of TrafficShaper |
| `/firewall/traffic_shaper/queues` | GET,DELETE | Plural | TrafficShaperQueue | Lookup for `firewall_traffic_shaper_queue` | child of TrafficShaper |
| `/firewall/traffic_shapers` | GET,PUT,DELETE | Plural | TrafficShaper | Lookup for `firewall_traffic_shaper` |  |
| `/firewall/virtual_ip` | GET,POST,PATCH,DELETE | Singular | VirtualIP | resource `firewall_virtual_ip` | Key: Marker |
| `/firewall/virtual_ip/apply` | GET,POST | Singular | VirtualIPApply | excluded | apply endpoint; when to apply is #59 |
| `/firewall/virtual_ips` | GET,DELETE | Plural | VirtualIP | Lookup for `firewall_virtual_ip` |  |
| `/graphql` | POST | Singular | GraphQL | excluded | GraphQL query endpoint, not a model |
| `/interface` | GET,POST,PATCH,DELETE | Singular | NetworkInterface | resource `interface` | Key: string API id |
| `/interface/apply` | GET,POST | Singular | InterfaceApply | excluded | apply endpoint; when to apply is #59 |
| `/interface/available_interfaces` | GET | Plural | AvailableInterface | excluded | plural only, read-only; data source candidate (later) |
| `/interface/bridge` | GET,POST,PATCH,DELETE | Singular | InterfaceBridge | resource `interface_bridge` | Key: Marker |
| `/interface/bridges` | GET | Plural | InterfaceBridge | Lookup for `interface_bridge` |  |
| `/interface/gre` | GET,POST,PATCH,DELETE | Singular | InterfaceGRE | resource `interface_gre` | Key: Marker |
| `/interface/gres` | GET,DELETE | Plural | InterfaceGRE | Lookup for `interface_gre` |  |
| `/interface/group` | GET,POST,PATCH,DELETE | Singular | InterfaceGroup | resource `interface_group` | Key: Marker |
| `/interface/groups` | GET,PUT,DELETE | Plural | InterfaceGroup | Lookup for `interface_group` |  |
| `/interface/lagg` | GET,POST,PATCH,DELETE | Singular | InterfaceLAGG | resource `interface_lagg` | Key: Marker |
| `/interface/laggs` | GET,DELETE | Plural | InterfaceLAGG | Lookup for `interface_lagg` |  |
| `/interface/vlan` | GET,POST,PATCH,DELETE | Singular | InterfaceVLAN | resource `interface_vlan` | Key: Marker |
| `/interface/vlans` | GET,DELETE | Plural | InterfaceVLAN | Lookup for `interface_vlan` |  |
| `/interfaces` | GET,DELETE | Plural | NetworkInterface | Lookup for `interface` |  |
| `/routing/apply` | GET,POST | Singular | RoutingApply | excluded | apply endpoint; when to apply is #59 |
| `/routing/gateway` | GET,POST,PATCH,DELETE | Singular | RoutingGateway | resource `routing_gateway` | Key: Marker |
| `/routing/gateway/default` | GET,PATCH | Singular | DefaultGateway | singleton resource `routing_gateway_default` | Create and Update PATCH; Destroy forgets |
| `/routing/gateway/group` | GET,POST,PATCH,DELETE | Singular | RoutingGatewayGroup | resource `routing_gateway_group` | Key: Marker |
| `/routing/gateway/group/priorities` | GET,DELETE | Plural | RoutingGatewayGroupPriority | Lookup for `routing_gateway_group_priority` | child of RoutingGatewayGroup |
| `/routing/gateway/group/priority` | GET,POST,PATCH,DELETE | Singular | RoutingGatewayGroupPriority | resource `routing_gateway_group_priority` | Key: `gateway`; child of RoutingGatewayGroup |
| `/routing/gateway/groups` | GET,DELETE | Plural | RoutingGatewayGroup | Lookup for `routing_gateway_group` |  |
| `/routing/gateways` | GET,DELETE | Plural | RoutingGateway | Lookup for `routing_gateway` |  |
| `/routing/static_route` | GET,POST,PATCH,DELETE | Singular | StaticRoute | resource `routing_static_route` | Key: Marker |
| `/routing/static_routes` | GET,DELETE | Plural | StaticRoute | Lookup for `routing_static_route` |  |
| `/services/acme/account_key` | GET,POST,PATCH,DELETE | Singular | ACMEAccountKey | resource `services_acme_account_key` | Key: Marker; needs pfSense-pkg-acme |
| `/services/acme/account_key/register` | POST | Singular | ACMEAccountKeyRegister | excluded | action (POST only); needs pfSense-pkg-acme |
| `/services/acme/account_key/registrations` | GET | Plural | ACMEAccountKeyRegister | excluded | plural of an excluded model; needs pfSense-pkg-acme |
| `/services/acme/account_keys` | GET,PUT,DELETE | Plural | ACMEAccountKey | Lookup for `services_acme_account_key` | needs pfSense-pkg-acme |
| `/services/acme/certificate` | GET,POST,PATCH,DELETE | Singular | ACMECertificate | resource `services_acme_certificate` | Key: Marker; needs pfSense-pkg-acme |
| `/services/acme/certificate/action` | GET,POST,PATCH,DELETE | Singular | ACMECertificateAction | excluded | no plural endpoint to Lookup on; child of ACMECertificate; needs pfSense-pkg-acme |
| `/services/acme/certificate/domain` | GET,POST,PATCH,DELETE | Singular | ACMECertificateDomain | excluded | no plural endpoint to Lookup on; child of ACMECertificate; needs pfSense-pkg-acme |
| `/services/acme/certificate/issuances` | GET | Plural | ACMECertificateIssue | excluded | plural of an excluded model; needs pfSense-pkg-acme |
| `/services/acme/certificate/issue` | POST | Singular | ACMECertificateIssue | excluded | action (POST only); needs pfSense-pkg-acme |
| `/services/acme/certificate/renew` | POST | Singular | ACMECertificateRenew | excluded | action (POST only); needs pfSense-pkg-acme |
| `/services/acme/certificate/renewals` | GET | Plural | ACMECertificateRenew | excluded | plural of an excluded model; needs pfSense-pkg-acme |
| `/services/acme/certificates` | GET,PUT,DELETE | Plural | ACMECertificate | Lookup for `services_acme_certificate` | needs pfSense-pkg-acme |
| `/services/acme/settings` | GET,PATCH | Singular | ACMESettings | singleton resource `services_acme_settings` | Create and Update PATCH; Destroy forgets; needs pfSense-pkg-acme |
| `/services/bind/access_list` | GET,POST,PATCH,DELETE | Singular | BINDAccessList | resource `services_bind_access_list` | Key: `name`; needs pfSense-pkg-bind |
| `/services/bind/access_list/entries` | GET,DELETE | Plural | BINDAccessListEntry | Lookup for `services_bind_access_list_entry` | child of BINDAccessList; needs pfSense-pkg-bind |
| `/services/bind/access_list/entry` | GET,POST,PATCH,DELETE | Singular | BINDAccessListEntry | resource `services_bind_access_list_entry` | Key: Marker in `description`; child of BINDAccessList; needs pfSense-pkg-bind |
| `/services/bind/access_lists` | GET,PUT,DELETE | Plural | BINDAccessList | Lookup for `services_bind_access_list` | needs pfSense-pkg-bind |
| `/services/bind/settings` | GET,PATCH | Singular | BINDSettings | singleton resource `services_bind_settings` | Create and Update PATCH; Destroy forgets; needs pfSense-pkg-bind |
| `/services/bind/sync/remote_host` | GET,POST,PATCH,DELETE | Singular | BINDSyncRemoteHost | not generated | no Key yet; would be `services_bind_sync_remote_host`; needs pfSense-pkg-bind |
| `/services/bind/sync/remote_hosts` | GET,PUT,DELETE | Plural | BINDSyncRemoteHost | not generated | plural of a model with no Key yet; needs pfSense-pkg-bind |
| `/services/bind/sync/settings` | GET,PATCH | Singular | BINDSyncSettings | singleton resource `services_bind_sync_settings` | Create and Update PATCH; Destroy forgets; needs pfSense-pkg-bind |
| `/services/bind/view` | GET,POST,PATCH,DELETE | Singular | BINDView | resource `services_bind_view` | Key: Marker; needs pfSense-pkg-bind |
| `/services/bind/views` | GET,PUT,DELETE | Plural | BINDView | Lookup for `services_bind_view` | needs pfSense-pkg-bind |
| `/services/bind/zone` | GET,POST,PATCH,DELETE | Singular | BINDZone | resource `services_bind_zone` | Key: `name`; needs pfSense-pkg-bind |
| `/services/bind/zone/record` | GET,POST,PATCH,DELETE | Singular | BINDZoneRecord | excluded | no plural endpoint to Lookup on; child of BINDZone; needs pfSense-pkg-bind |
| `/services/bind/zones` | GET,PUT,DELETE | Plural | BINDZone | Lookup for `services_bind_zone` | needs pfSense-pkg-bind |
| `/services/cron/job` | GET,POST,PATCH,DELETE | Singular | CronJob | not generated | no Key yet; would be `services_cron_job`; needs pfSense-pkg-Cron |
| `/services/cron/jobs` | GET,PUT,DELETE | Plural | CronJob | not generated | plural of a model with no Key yet; needs pfSense-pkg-Cron |
| `/services/dhcp_relay` | GET,PATCH | Singular | DHCPRelay | singleton resource `services_dhcp_relay` | Create and Update PATCH; Destroy forgets |
| `/services/dhcp_server` | GET,POST,PATCH,DELETE | Singular | DHCPServer | resource `services_dhcp_server` | Key: string API id |
| `/services/dhcp_server/address_pool` | GET,POST,PATCH,DELETE | Singular | DHCPServerAddressPool | not generated | no Key yet; would be `services_dhcp_server_address_pool`; child of DHCPServer |
| `/services/dhcp_server/address_pools` | GET,DELETE | Plural | DHCPServerAddressPool | not generated | plural of a model with no Key yet; child of DHCPServer |
| `/services/dhcp_server/apply` | GET,POST | Singular | DHCPServerApply | excluded | apply endpoint; when to apply is #59 |
| `/services/dhcp_server/backend` | PATCH | Singular | DHCPServerBackend | excluded | PATCH only, nothing to Read |
| `/services/dhcp_server/custom_option` | GET,POST,PATCH,DELETE | Singular | DHCPServerCustomOption | resource `services_dhcp_server_custom_option` | Key: `number`; child of DHCPServer |
| `/services/dhcp_server/custom_options` | GET,DELETE | Plural | DHCPServerCustomOption | Lookup for `services_dhcp_server_custom_option` | child of DHCPServer |
| `/services/dhcp_server/static_mapping` | GET,POST,PATCH,DELETE | Singular | DHCPServerStaticMapping | resource `services_dhcp_server_static_mapping` | Key: `parent_id` + `mac`; child of DHCPServer |
| `/services/dhcp_server/static_mappings` | GET,DELETE | Plural | DHCPServerStaticMapping | Lookup for `services_dhcp_server_static_mapping` | child of DHCPServer |
| `/services/dhcp_servers` | GET,PUT | Plural | DHCPServer | Lookup for `services_dhcp_server` |  |
| `/services/dns_forwarder/apply` | GET,POST | Singular | DNSForwarderApply | excluded | apply endpoint; when to apply is #59 |
| `/services/dns_forwarder/host_override` | GET,POST,PATCH,DELETE | Singular | DNSForwarderHostOverride | resource `services_dns_forwarder_host_override` | Key: Marker |
| `/services/dns_forwarder/host_override/alias` | GET,POST,PATCH,DELETE | Singular | DNSForwarderHostOverrideAlias | resource `services_dns_forwarder_host_override_alias` | Key: Marker in `description`; child of DNSForwarderHostOverride |
| `/services/dns_forwarder/host_override/aliases` | GET,DELETE | Plural | DNSForwarderHostOverrideAlias | Lookup for `services_dns_forwarder_host_override_alias` | child of DNSForwarderHostOverride |
| `/services/dns_forwarder/host_overrides` | GET,PUT,DELETE | Plural | DNSForwarderHostOverride | Lookup for `services_dns_forwarder_host_override` |  |
| `/services/dns_resolver/access_list` | GET,POST,PATCH,DELETE | Singular | DNSResolverAccessList | resource `services_dns_resolver_access_list` | Key: Marker in `description` |
| `/services/dns_resolver/access_list/network` | GET,POST,PATCH,DELETE | Singular | DNSResolverAccessListNetwork | resource `services_dns_resolver_access_list_network` | Key: Marker in `description`; child of DNSResolverAccessList |
| `/services/dns_resolver/access_list/networks` | GET,DELETE | Plural | DNSResolverAccessListNetwork | Lookup for `services_dns_resolver_access_list_network` | child of DNSResolverAccessList |
| `/services/dns_resolver/access_lists` | GET,PUT,DELETE | Plural | DNSResolverAccessList | Lookup for `services_dns_resolver_access_list` |  |
| `/services/dns_resolver/apply` | GET,POST | Singular | DNSResolverApply | excluded | apply endpoint; when to apply is #59 |
| `/services/dns_resolver/domain_override` | GET,POST,PATCH,DELETE | Singular | DNSResolverDomainOverride | resource `services_dns_resolver_domain_override` | Key: Marker |
| `/services/dns_resolver/domain_overrides` | GET,PUT,DELETE | Plural | DNSResolverDomainOverride | Lookup for `services_dns_resolver_domain_override` |  |
| `/services/dns_resolver/host_override` | GET,POST,PATCH,DELETE | Singular | DNSResolverHostOverride | resource `services_dns_resolver_host_override` | Key: Marker |
| `/services/dns_resolver/host_override/alias` | GET,POST,PATCH,DELETE | Singular | DNSResolverHostOverrideAlias | resource `services_dns_resolver_host_override_alias` | Key: Marker; child of DNSResolverHostOverride |
| `/services/dns_resolver/host_override/aliases` | GET,DELETE | Plural | DNSResolverHostOverrideAlias | Lookup for `services_dns_resolver_host_override_alias` | child of DNSResolverHostOverride |
| `/services/dns_resolver/host_overrides` | GET,PUT,DELETE | Plural | DNSResolverHostOverride | Lookup for `services_dns_resolver_host_override` |  |
| `/services/dns_resolver/settings` | GET,PATCH | Singular | DNSResolverSettings | singleton resource `services_dns_resolver_settings` | Create and Update PATCH; Destroy forgets |
| `/services/freeradius/client` | GET,POST,PATCH,DELETE | Singular | FreeRADIUSClient | resource `services_freeradius_client` | Key: `addr`; needs pfSense-pkg-freeradius3 |
| `/services/freeradius/clients` | GET,PUT,DELETE | Plural | FreeRADIUSClient | Lookup for `services_freeradius_client` | needs pfSense-pkg-freeradius3 |
| `/services/freeradius/eap` | GET,PATCH | Singular | FreeRADIUSEAP | singleton resource `services_freeradius_eap` | Create and Update PATCH; Destroy forgets; needs pfSense-pkg-freeradius3 |
| `/services/freeradius/interface` | GET,POST,PATCH,DELETE | Singular | FreeRADIUSInterface | not generated | no Key yet; would be `services_freeradius_interface`; its `description` rejects the Marker's brackets; needs pfSense-pkg-freeradius3 |
| `/services/freeradius/interfaces` | GET,PUT,DELETE | Plural | FreeRADIUSInterface | not generated | plural of a model with no Key yet; needs pfSense-pkg-freeradius3 |
| `/services/freeradius/ldap` | GET,PATCH | Singular | FreeRADIUSLDAP | singleton resource `services_freeradius_ldap` | Create and Update PATCH; Destroy forgets; needs pfSense-pkg-freeradius3 |
| `/services/freeradius/mac` | GET,POST,PATCH,DELETE | Singular | FreeRADIUSMAC | resource `services_freeradius_mac` | Key: `mac`; needs pfSense-pkg-freeradius3 |
| `/services/freeradius/macs` | GET,PUT,DELETE | Plural | FreeRADIUSMAC | Lookup for `services_freeradius_mac` | needs pfSense-pkg-freeradius3 |
| `/services/freeradius/user` | GET,POST,PATCH,DELETE | Singular | FreeRADIUSUser | resource `services_freeradius_user` | Key: `username`; needs pfSense-pkg-freeradius3 |
| `/services/freeradius/users` | GET,PUT,DELETE | Plural | FreeRADIUSUser | Lookup for `services_freeradius_user` | needs pfSense-pkg-freeradius3 |
| `/services/haproxy/apply` | GET,POST | Singular | HAProxyApply | excluded | apply endpoint; when to apply is #59; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend` | GET,POST,PATCH,DELETE | Singular | HAProxyBackend | resource `services_haproxy_backend` | Key: `name`; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend/acl` | GET,POST,PATCH,DELETE | Singular | HAProxyBackendACL | not generated | no Key yet; would be `services_haproxy_backend_acl`; child of HAProxyBackend; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend/acls` | GET,DELETE | Plural | HAProxyBackendACL | not generated | plural of a model with no Key yet; child of HAProxyBackend; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend/action` | GET,POST,PATCH,DELETE | Singular | HAProxyBackendAction | not generated | no Key yet; would be `services_haproxy_backend_action`; child of HAProxyBackend; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend/actions` | GET,DELETE | Plural | HAProxyBackendAction | not generated | plural of a model with no Key yet; child of HAProxyBackend; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend/error_file` | GET,POST,PATCH,DELETE | Singular | HAProxyBackendErrorFile | resource `services_haproxy_backend_error_file` | Key: `errorcode`; child of HAProxyBackend; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend/errorfiles` | GET,DELETE | Plural | HAProxyBackendErrorFile | Lookup for `services_haproxy_backend_error_file` | child of HAProxyBackend; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend/server` | GET,POST,PATCH,DELETE | Singular | HAProxyBackendServer | resource `services_haproxy_backend_server` | Key: `name`; child of HAProxyBackend; needs pfSense-pkg-haproxy |
| `/services/haproxy/backend/servers` | GET,DELETE | Plural | HAProxyBackendServer | Lookup for `services_haproxy_backend_server` | child of HAProxyBackend; needs pfSense-pkg-haproxy |
| `/services/haproxy/backends` | GET,PUT,DELETE | Plural | HAProxyBackend | Lookup for `services_haproxy_backend` | needs pfSense-pkg-haproxy |
| `/services/haproxy/file` | GET,POST,PATCH,DELETE | Singular | HAProxyFile | resource `services_haproxy_file` | Key: `name`; needs pfSense-pkg-haproxy |
| `/services/haproxy/files` | GET,PUT,DELETE | Plural | HAProxyFile | Lookup for `services_haproxy_file` | needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend` | GET,POST,PATCH,DELETE | Singular | HAProxyFrontend | resource `services_haproxy_frontend` | Key: Marker; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/acl` | GET,POST,PATCH,DELETE | Singular | HAProxyFrontendACL | not generated | no Key yet; would be `services_haproxy_frontend_acl`; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/acls` | GET,DELETE | Plural | HAProxyFrontendACL | not generated | plural of a model with no Key yet; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/action` | GET,POST,PATCH,DELETE | Singular | HAProxyFrontendAction | not generated | no Key yet; would be `services_haproxy_frontend_action`; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/actions` | GET,DELETE | Plural | HAProxyFrontendAction | not generated | plural of a model with no Key yet; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/address` | GET,POST,PATCH,DELETE | Singular | HAProxyFrontendAddress | not generated | no Key yet; would be `services_haproxy_frontend_address`; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/addresses` | GET,DELETE | Plural | HAProxyFrontendAddress | not generated | plural of a model with no Key yet; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/certificate` | GET,POST,PATCH,DELETE | Singular | HAProxyFrontendCertificate | not generated | no Key yet; would be `services_haproxy_frontend_certificate`; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/certificates` | GET,DELETE | Plural | HAProxyFrontendCertificate | not generated | plural of a model with no Key yet; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/error_file` | GET,POST,PATCH,DELETE | Singular | HAProxyFrontendErrorFile | resource `services_haproxy_frontend_error_file` | Key: `errorcode`; child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontend/error_files` | GET,DELETE | Plural | HAProxyFrontendErrorFile | Lookup for `services_haproxy_frontend_error_file` | child of HAProxyFrontend; needs pfSense-pkg-haproxy |
| `/services/haproxy/frontends` | GET,PUT,DELETE | Plural | HAProxyFrontend | Lookup for `services_haproxy_frontend` | needs pfSense-pkg-haproxy |
| `/services/haproxy/settings` | GET,PATCH | Singular | HAProxySettings | singleton resource `services_haproxy_settings` | Create and Update PATCH; Destroy forgets; needs pfSense-pkg-haproxy |
| `/services/haproxy/settings/dns_resolver` | GET,POST,PATCH,DELETE | Singular | HAProxyDNSResolver | resource `services_haproxy_settings_dns_resolver` | Key: `name`; child of HAProxySettings; needs pfSense-pkg-haproxy |
| `/services/haproxy/settings/dns_resolvers` | GET,DELETE | Plural | HAProxyDNSResolver | Lookup for `services_haproxy_settings_dns_resolver` | child of HAProxySettings; needs pfSense-pkg-haproxy |
| `/services/haproxy/settings/email_mailer` | GET,POST,PATCH,DELETE | Singular | HAProxyEmailMailer | resource `services_haproxy_settings_email_mailer` | Key: `name`; child of HAProxySettings; needs pfSense-pkg-haproxy |
| `/services/haproxy/settings/email_mailers` | GET,DELETE | Plural | HAProxyEmailMailer | Lookup for `services_haproxy_settings_email_mailer` | child of HAProxySettings; needs pfSense-pkg-haproxy |
| `/services/ntp/settings` | GET,PATCH | Singular | NTPSettings | singleton resource `services_ntp_settings` | Create and Update PATCH; Destroy forgets |
| `/services/ntp/time_server` | GET,POST,PATCH,DELETE | Singular | NTPTimeServer | not generated | no Key yet; would be `services_ntp_time_server` |
| `/services/ntp/time_servers` | GET,PUT,DELETE | Plural | NTPTimeServer | not generated | plural of a model with no Key yet |
| `/services/service_watchdog` | GET,POST,PATCH,DELETE | Singular | ServiceWatchdog | resource `services_service_watchdog` | Key: `name`; needs pfSense-pkg-Service_Watchdog |
| `/services/service_watchdogs` | GET,PUT,DELETE | Plural | ServiceWatchdog | Lookup for `services_service_watchdog` | needs pfSense-pkg-Service_Watchdog |
| `/services/ssh` | GET,PATCH | Singular | SSH | singleton resource `services_ssh` | Create and Update PATCH; Destroy forgets |
| `/services/wake_on_lan/send` | POST | Singular | WakeOnLANSend | excluded | action (POST only) |
| `/status/carp` | GET,PATCH | Singular | CARP | singleton resource `status_carp` | Create and Update PATCH; Destroy forgets |
| `/status/dhcp_server/leases` | GET,DELETE | Plural | DHCPServerLease | excluded | plural only, read-only; data source candidate (later) |
| `/status/gateways` | GET | Plural | RoutingGatewayStatus | excluded | plural only, read-only; data source candidate (later) |
| `/status/interfaces` | GET | Plural | InterfaceStats | excluded | plural only, read-only; data source candidate (later) |
| `/status/ipsec/child_sa` | GET | Singular | IPsecChildSAStatus | excluded | read-only state; data source candidate (later); child of IPsecSAStatus |
| `/status/ipsec/child_sas` | GET | Plural | IPsecChildSAStatus | excluded | read-only state; data source candidate (later); child of IPsecSAStatus |
| `/status/ipsec/sas` | GET | Plural | IPsecSAStatus | excluded | plural only, read-only; data source candidate (later) |
| `/status/logs/auth` | GET | Plural | AuthLog | excluded | plural only, read-only; data source candidate (later) |
| `/status/logs/dhcp` | GET | Plural | DHCPLog | excluded | plural only, read-only; data source candidate (later) |
| `/status/logs/firewall` | GET | Plural | FirewallLog | excluded | plural only, read-only; data source candidate (later) |
| `/status/logs/openvpn` | GET | Plural | OpenVPNLog | excluded | plural only, read-only; data source candidate (later) |
| `/status/logs/packages/restapi` | GET | Plural | RESTAPILog | excluded | plural only, read-only; data source candidate (later) |
| `/status/logs/settings` | GET,PATCH | Singular | LogSettings | singleton resource `status_logs_settings` | Create and Update PATCH; Destroy forgets |
| `/status/logs/system` | GET | Plural | SystemLog | excluded | plural only, read-only; data source candidate (later) |
| `/status/openvpn/clients` | GET | Plural | OpenVPNClientStatus | excluded | plural only, read-only; data source candidate (later) |
| `/status/openvpn/server/connection` | GET,DELETE | Singular | OpenVPNServerConnectionStatus | excluded | read-only state; data source candidate (later); child of OpenVPNServerStatus |
| `/status/openvpn/server/connections` | GET,DELETE | Plural | OpenVPNServerConnectionStatus | excluded | read-only state; data source candidate (later); child of OpenVPNServerStatus |
| `/status/openvpn/server/route` | GET | Singular | OpenVPNServerRouteStatus | excluded | read-only state; data source candidate (later); child of OpenVPNServerStatus |
| `/status/openvpn/server/routes` | GET | Plural | OpenVPNServerRouteStatus | excluded | read-only state; data source candidate (later); child of OpenVPNServerStatus |
| `/status/openvpn/servers` | GET | Plural | OpenVPNServerStatus | excluded | plural only, read-only; data source candidate (later) |
| `/status/service` | POST | Singular | Service | excluded | action (POST only) |
| `/status/services` | GET | Plural | Service | excluded | plural of an excluded model |
| `/status/system` | GET | Singular | SystemStatus | excluded | read-only state; data source candidate (later) |
| `/system/certificate` | GET,POST,PATCH,DELETE | Singular | Certificate | resource `system_certificate` | Key: Marker |
| `/system/certificate/generate` | POST | Singular | CertificateGenerate | excluded | action (POST only) |
| `/system/certificate/pkcs12/export` | POST | Singular | CertificatePKCS12Export | excluded | action (POST only) |
| `/system/certificate/renew` | POST | Singular | CertificateRenew | excluded | action (POST only) |
| `/system/certificate/signing_request` | POST | Singular | CertificateSigningRequest | excluded | action (POST only) |
| `/system/certificate/signing_request/sign` | POST | Singular | CertificateSigningRequestSign | excluded | action (POST only) |
| `/system/certificate_authorities` | GET,DELETE | Plural | CertificateAuthority | Lookup for `system_certificate_authority` |  |
| `/system/certificate_authority` | GET,POST,PATCH,DELETE | Singular | CertificateAuthority | resource `system_certificate_authority` | Key: Marker |
| `/system/certificate_authority/generate` | POST | Singular | CertificateAuthorityGenerate | excluded | action (POST only) |
| `/system/certificate_authority/renew` | POST | Singular | CertificateAuthorityRenew | excluded | action (POST only) |
| `/system/certificates` | GET,DELETE | Plural | Certificate | Lookup for `system_certificate` |  |
| `/system/console` | GET,PATCH | Singular | SystemConsole | singleton resource `system_console` | Create and Update PATCH; Destroy forgets |
| `/system/crl` | GET,POST,PATCH,DELETE | Singular | CertificateRevocationList | resource `system_crl` | Key: Marker |
| `/system/crl/revoked_certificate` | GET,POST,PATCH,DELETE | Singular | CertificateRevocationListRevokedCertificate | excluded | no plural endpoint to Lookup on; child of CertificateRevocationList |
| `/system/crls` | GET,DELETE | Plural | CertificateRevocationList | Lookup for `system_crl` |  |
| `/system/dns` | GET,PATCH | Singular | SystemDNS | singleton resource `system_dns` | Create and Update PATCH; Destroy forgets |
| `/system/enum` | POST | Singular | Enum | excluded | action (POST only) |
| `/system/hostname` | GET,PATCH | Singular | SystemHostname | singleton resource `system_hostname` | Create and Update PATCH; Destroy forgets |
| `/system/notifications/email_settings` | GET,PATCH | Singular | EmailNotificationSettings | singleton resource `system_notifications_email_settings` | Create and Update PATCH; Destroy forgets |
| `/system/package` | GET,POST,DELETE | Singular | Package | excluded | no PATCH; package install is fog on the map |
| `/system/package/available` | GET | Plural | AvailablePackage | excluded | plural only, read-only; data source candidate (later) |
| `/system/packages` | GET,DELETE | Plural | Package | excluded | plural of an excluded model |
| `/system/restapi/access_list` | GET,PUT,DELETE | Plural | RESTAPIAccessListEntry | Lookup for `system_restapi_access_list_entry` |  |
| `/system/restapi/access_list/entry` | GET,POST,PATCH,DELETE | Singular | RESTAPIAccessListEntry | resource `system_restapi_access_list_entry` | Key: Marker |
| `/system/restapi/settings` | GET,PATCH | Singular | RESTAPISettings | singleton resource `system_restapi_settings` | Create and Update PATCH; Destroy forgets |
| `/system/restapi/settings/sync` | POST | Singular | RESTAPISettingsSync | excluded | action (POST only) |
| `/system/restapi/version` | GET,PATCH | Singular | RESTAPIVersion | excluded | PATCH upgrades the REST API package |
| `/system/timezone` | GET,PATCH | Singular | SystemTimezone | singleton resource `system_timezone` | Create and Update PATCH; Destroy forgets |
| `/system/tunable` | GET,POST,PATCH,DELETE | Singular | SystemTunable | resource `system_tunable` | Key: Marker |
| `/system/tunables` | GET,PUT,DELETE | Plural | SystemTunable | Lookup for `system_tunable` |  |
| `/system/update` | POST | Singular | SystemUpdate | excluded | action (POST only) |
| `/system/version` | GET | Singular | SystemVersion | excluded | read-only state; data source candidate (later) |
| `/system/webgui/settings` | GET,PATCH | Singular | WebGUISettings | singleton resource `system_webgui_settings` | Create and Update PATCH; Destroy forgets |
| `/user` | GET,POST,PATCH,DELETE | Singular | User | resource `user` | Key: Marker |
| `/user/auth_server` | GET,POST,PATCH,DELETE | Singular | AuthServer | resource `user_auth_server` | Key: `name` |
| `/user/auth_servers` | GET,PUT,DELETE | Plural | AuthServer | Lookup for `user_auth_server` |  |
| `/user/group` | GET,POST,PATCH,DELETE | Singular | UserGroup | resource `user_group` | Key: `name` |
| `/user/groups` | GET,PUT,DELETE | Plural | UserGroup | Lookup for `user_group` |  |
| `/users` | GET,DELETE | Plural | User | Lookup for `user` |  |
| `/vpn/ipsec/apply` | GET,POST | Singular | IPsecApply | excluded | apply endpoint; when to apply is #59 |
| `/vpn/ipsec/phase1` | GET,POST,PATCH,DELETE | Singular | IPsecPhase1 | resource `vpn_ipsec_phase1` | Key: Marker |
| `/vpn/ipsec/phase1/encryption` | GET,POST,PATCH,DELETE | Singular | IPsecPhase1Encryption | not generated | no Key yet; would be `vpn_ipsec_phase1_encryption`; child of IPsecPhase1 |
| `/vpn/ipsec/phase1/encryptions` | GET,DELETE | Plural | IPsecPhase1Encryption | not generated | plural of a model with no Key yet; child of IPsecPhase1 |
| `/vpn/ipsec/phase1s` | GET,PUT,DELETE | Plural | IPsecPhase1 | Lookup for `vpn_ipsec_phase1` |  |
| `/vpn/ipsec/phase2` | GET,POST,PATCH,DELETE | Singular | IPsecPhase2 | resource `vpn_ipsec_phase2` | Key: Marker |
| `/vpn/ipsec/phase2/encryption` | GET,POST,PATCH,DELETE | Singular | IPsecPhase2Encryption | not generated | no Key yet; would be `vpn_ipsec_phase2_encryption`; child of IPsecPhase2 |
| `/vpn/ipsec/phase2/encryptions` | GET,DELETE | Plural | IPsecPhase2Encryption | not generated | plural of a model with no Key yet; child of IPsecPhase2 |
| `/vpn/ipsec/phase2s` | GET,PUT,DELETE | Plural | IPsecPhase2 | Lookup for `vpn_ipsec_phase2` |  |
| `/vpn/openvpn/client` | GET,POST,PATCH,DELETE | Singular | OpenVPNClient | resource `vpn_openvpn_client` | Key: `vpnid`, assigned by pfSense on create |
| `/vpn/openvpn/client_export` | POST | Singular | OpenVPNClientExport | excluded | action (POST only); needs pfSense-pkg-openvpn-client-export |
| `/vpn/openvpn/client_export/config` | GET,POST,PATCH,DELETE | Singular | OpenVPNClientExportConfig | not generated | no Key yet; would be `vpn_openvpn_client_export_config`; needs pfSense-pkg-openvpn-client-export |
| `/vpn/openvpn/client_export/configs` | GET,PUT,DELETE | Plural | OpenVPNClientExportConfig | not generated | plural of a model with no Key yet; needs pfSense-pkg-openvpn-client-export |
| `/vpn/openvpn/clients` | GET,DELETE | Plural | OpenVPNClient | Lookup for `vpn_openvpn_client` |  |
| `/vpn/openvpn/cso` | GET,POST,PATCH,DELETE | Singular | OpenVPNClientSpecificOverride | resource `vpn_openvpn_cso` | Key: Marker in `description` |
| `/vpn/openvpn/csos` | GET,DELETE | Plural | OpenVPNClientSpecificOverride | Lookup for `vpn_openvpn_cso` |  |
| `/vpn/openvpn/server` | GET,POST,PATCH,DELETE | Singular | OpenVPNServer | resource `vpn_openvpn_server` | Key: `vpnid`, assigned by pfSense on create |
| `/vpn/openvpn/servers` | GET,DELETE | Plural | OpenVPNServer | Lookup for `vpn_openvpn_server` |  |
| `/vpn/wireguard/apply` | GET,POST | Singular | WireGuardApply | excluded | apply endpoint; when to apply is #59 |
| `/vpn/wireguard/peer` | GET,POST,PATCH,DELETE | Singular | WireGuardPeer | resource `vpn_wireguard_peer` | Key: Marker; needs pfSense-pkg-WireGuard |
| `/vpn/wireguard/peer/allowed_ip` | GET,POST,PATCH,DELETE | Singular | WireGuardPeerAllowedIP | resource `vpn_wireguard_peer_allowed_ip` | Key: Marker; child of WireGuardPeer; needs pfSense-pkg-WireGuard |
| `/vpn/wireguard/peer/allowed_ips` | GET,DELETE | Plural | WireGuardPeerAllowedIP | Lookup for `vpn_wireguard_peer_allowed_ip` | child of WireGuardPeer; needs pfSense-pkg-WireGuard |
| `/vpn/wireguard/peers` | GET,PUT,DELETE | Plural | WireGuardPeer | Lookup for `vpn_wireguard_peer` | needs pfSense-pkg-WireGuard |
| `/vpn/wireguard/settings` | GET,PATCH | Singular | WireGuardSettings | singleton resource `vpn_wireguard_settings` | Create and Update PATCH; Destroy forgets; needs pfSense-pkg-WireGuard |
| `/vpn/wireguard/tunnel` | GET,POST,PATCH,DELETE | Singular | WireGuardTunnel | resource `vpn_wireguard_tunnel` | Key: Marker; needs pfSense-pkg-WireGuard |
| `/vpn/wireguard/tunnel/address` | GET,POST,PATCH,DELETE | Singular | WireGuardTunnelAddress | resource `vpn_wireguard_tunnel_address` | Key: Marker; child of WireGuardTunnel; needs pfSense-pkg-WireGuard |
| `/vpn/wireguard/tunnel/addresses` | GET,DELETE | Plural | WireGuardTunnelAddress | Lookup for `vpn_wireguard_tunnel_address` | child of WireGuardTunnel; needs pfSense-pkg-WireGuard |
| `/vpn/wireguard/tunnels` | GET,PUT,DELETE | Plural | WireGuardTunnel | Lookup for `vpn_wireguard_tunnel` | needs pfSense-pkg-WireGuard |
