<?php

declare(strict_types=1);

// Send PHP warnings/notices to stderr so stdout contains only the JSON schema
ini_set('display_errors', 'stderr');
error_reporting(E_ALL);

// $argv[1] is the build directory containing the RESTAPI/ subdirectory
$root = $argv[1];
$restapi = $root . '/RESTAPI';

// Many RESTAPI files do `require_once 'RESTAPI/autoloader.inc'` as a safety net.
// We've replaced autoloader.inc with a no-op and load everything ourselves,
// so set include_path to $root so those require_once calls resolve without error.
ini_set('include_path', $root . ':' . get_include_path());

require_once $restapi . '/.resources/vendor/autoload.php';

// pfSense globals — static data extracted verbatim from pfSense CE source
global $config, $sysctls, $g, $pf_reserved_ifs, $filterent_db;
global $p1_ealgos, $p2_ealgos, $p1_halgos, $p2_halgos, $p1_dhgroups, $p2_dhgroups, $p2_pfskeygroups;
global $p1_authentication_methods, $p2_modes, $p2_protos;
global $ipsec_conid_prefix, $ipsec_conid_separator, $ipsec_reqid_base;
global $ipsec_loglevels, $ipsec_log_sevs, $ipsec_log_cats;
global $ipsec_identifier_list, $my_identifier_list, $peer_identifier_list;
global $ipsec_idhandling, $ipsec_preshared_key_type, $ipsec_startactions, $ipsec_closeactions;
global $openvpn_dh_lengths, $openvpn_tls_modes, $openvpn_dev_mode;
global $priv_list, $ipsec_descrs;
global $a_acmeserver, $acme_domain_validation_method;
global $current_openvpn_version, $current_openvpn_version_rev;
global $legacy_openvpn_version, $legacy_openvpn_version_rev;
global $previous_openvpn_version, $previous_openvpn_version_rev;
global $legacy_incompatible_ciphers;

// gettext stub — pfSense globals use it for i18n
if (!function_exists('gettext')) {
    function gettext(string $s): string { return $s; }
    function _($s): string { return $s; }
}

$config = [];
$sysctls = [];
$g = [];
$pf_reserved_ifs = [];
$filterent_db = [];
$priv_list = [];
$ipsec_descrs = [];
$a_acmeserver = [];
$acme_domain_validation_method = [];
$current_openvpn_version = '2.6';
$current_openvpn_version_rev = 0;
$legacy_openvpn_version = '2.5';
$legacy_openvpn_version_rev = 0;
$previous_openvpn_version = '2.4';
$previous_openvpn_version_rev = 0;
$legacy_incompatible_ciphers = [];
$openvpn_dh_lengths = [1024, 2048, 3072, 4096, 'none'];
$openvpn_tls_modes = ['tls' => 'TLS Authentication', 'tls-crypt' => 'TLS Encryption'];
$openvpn_dev_mode = ['tun' => 'tun', 'tap' => 'tap'];

// IPsec globals — verbatim from pfSense CE src/etc/inc/ipsec.inc
$ipsec_conid_prefix = 'con';
$ipsec_conid_separator = '_';
$ipsec_reqid_base = 5000;
$ipsec_loglevels = ['dmn' => 'Daemon', 'mgr' => 'SA Manager', 'ike' => 'IKE SA',
    'chd' => 'IKE Child SA', 'job' => 'Job Processing', 'cfg' => 'Configuration backend',
    'knl' => 'Kernel Interface', 'net' => 'Networking', 'asn' => 'ASN encoding',
    'enc' => 'Message encoding', 'esp' => 'IPsec traffic', 'lib' => 'StrongSwan Lib'];
$ipsec_log_sevs = ['-1' => 'Silent', '0' => 'Audit', '1' => 'Control',
    '2' => 'Diag', '3' => 'Raw', '4' => 'Highest'];
$ipsec_log_cats = $ipsec_loglevels;
$ipsec_identifier_list = [
    'none' => ['desc' => 'Automatic based on content', 'mobile' => true],
    'email' => ['desc' => 'E-mail address', 'mobile' => true],
    'fqdn' => ['desc' => 'Fully qualified domain name', 'mobile' => true],
    'userfqdn' => ['desc' => 'User fully qualified domain name', 'mobile' => true],
    'keyid' => ['desc' => 'KeyID tag', 'mobile' => true],
    'asn1dn' => ['desc' => 'ASN.1 distinguished Name', 'mobile' => true],
];
$my_identifier_list = [
    'myaddress' => ['desc' => 'My IP address', 'mobile' => true],
    'address' => ['desc' => 'IP address', 'mobile' => true],
    'fqdn' => ['desc' => 'Fully qualified domain name', 'mobile' => true],
    'user_fqdn' => ['desc' => 'User fully qualified domain name / E-mail', 'mobile' => true],
    'asn1dn' => ['desc' => 'ASN.1 distinguished Name', 'mobile' => true],
    'keyid tag' => ['desc' => 'KeyID tag', 'mobile' => true],
    'dyn_dns' => ['desc' => 'Dynamic DNS', 'mobile' => true],
    'auto' => ['desc' => 'Automatic based on content', 'mobile' => true],
];
$peer_identifier_list = [
    'any' => ['desc' => 'Any', 'mobile' => true],
    'peeraddress' => ['desc' => 'Peer IP address', 'mobile' => false],
    'address' => ['desc' => 'IP address', 'mobile' => false],
    'fqdn' => ['desc' => 'Fully qualified domain name', 'mobile' => true],
    'user_fqdn' => ['desc' => 'User fully qualified domain name / E-mail', 'mobile' => true],
    'asn1dn' => ['desc' => 'ASN.1 distinguished Name', 'mobile' => true],
    'keyid tag' => ['desc' => 'KeyID tag', 'mobile' => true],
    'auto' => ['desc' => 'Automatic based on content', 'mobile' => true],
];
$ipsec_idhandling = ['replace' => 'Yes (Replace)', 'no' => 'No', 'never' => 'Never', 'keep' => 'Keep'];
$p1_ealgos = [
    'aes'            => ['name' => 'AES',              'keysel' => ['lo' => 128, 'hi' => 256, 'step' => 64]],
    'aes128gcm'      => ['name' => 'AES128-GCM',       'keysel' => ['lo' => 64,  'hi' => 128, 'step' => 32]],
    'aes192gcm'      => ['name' => 'AES192-GCM',       'keysel' => ['lo' => 64,  'hi' => 128, 'step' => 32]],
    'aes256gcm'      => ['name' => 'AES256-GCM',       'keysel' => ['lo' => 64,  'hi' => 128, 'step' => 32]],
    'chacha20poly1305' => ['name' => 'CHACHA20-POLY1305'],
];
$p2_ealgos = $p1_ealgos;
$p1_halgos = ['sha1' => 'SHA1', 'sha256' => 'SHA256', 'sha384' => 'SHA384',
    'sha512' => 'SHA512', 'aesxcbc' => 'AES-XCBC'];
$p2_halgos = ['hmac_sha1' => 'SHA1', 'hmac_sha256' => 'SHA256', 'hmac_sha384' => 'SHA384',
    'hmac_sha512' => 'SHA512', 'aesxcbc' => 'AES-XCBC'];
$p1_dhgroups = [
    1 => '1 (768 bit)', 2 => '2 (1024 bit)', 5 => '5 (1536 bit)',
    14 => '14 (2048 bit)', 15 => '15 (3072 bit)', 16 => '16 (4096 bit)',
    17 => '17 (6144 bit)', 18 => '18 (8192 bit)',
    19 => '19 (nist ecp256)', 20 => '20 (nist ecp384)', 21 => '21 (nist ecp521)',
    22 => '22 (1024(sub 160) bit)', 23 => '23 (2048(sub 224) bit)',
    24 => '24 (2048(sub 256) bit)', 25 => '25 (nist ecp192)', 26 => '26 (nist ecp224)',
    27 => '27 (brainpool ecp224)', 28 => '28 (brainpool ecp256)',
    29 => '29 (brainpool ecp384)', 30 => '30 (brainpool ecp512)',
    31 => '31 (Elliptic Curve 25519, 256 bit)', 32 => '32 (Elliptic Curve 448, 448 bit)',
];
$p2_dhgroups = $p1_dhgroups;
$p2_pfskeygroups = array_merge([0 => 'off'], $p1_dhgroups);
$p1_authentication_methods = [
    'hybrid_cert_server' => ['name' => 'Hybrid Certificate + Xauth', 'mobile' => true],
    'xauth_cert_server'  => ['name' => 'Mutual Certificate + Xauth', 'mobile' => true],
    'xauth_psk_server'   => ['name' => 'Mutual PSK + Xauth', 'mobile' => true],
    'eap-tls'            => ['name' => 'EAP-TLS', 'mobile' => true],
    'eap-radius'         => ['name' => 'EAP-RADIUS', 'mobile' => true],
    'eap-mschapv2'       => ['name' => 'EAP-MSChapv2', 'mobile' => true],
    'cert'               => ['name' => 'Mutual Certificate', 'mobile' => false],
    'pkcs11'             => ['name' => 'Mutual Certificate (PKCS#11)', 'mobile' => false],
    'pre_shared_key'     => ['name' => 'Mutual PSK', 'mobile' => false],
];
$ipsec_preshared_key_type = ['PSK' => 'PSK', 'EAP' => 'EAP'];
$ipsec_startactions = ['' => 'Default', 'none' => 'None (Responder Only)',
    'start' => 'Initiate at start', 'trap' => 'Initiate on demand'];
$ipsec_closeactions = ['' => 'Default', 'none' => 'Close connection and clear SA',
    'start' => 'Restart/Reconnect', 'trap' => 'Close connection and reconnect on demand'];
$p2_modes = ['tunnel' => 'Tunnel IPv4', 'tunnel6' => 'Tunnel IPv6',
    'transport' => 'Transport', 'vti' => 'Routed (VTI)'];
$p2_protos = ['esp' => 'ESP', 'ah' => 'AH'];

// Synthetic pfSense config — provides enough structure so Model constructors don't crash.
// Values mirror the field defaults defined in RESTAPISettings.
const PFSENSE_STUB_CONFIG = [
    'installedpackages' => [
        'package' => [
            [
                'name' => 'RESTAPI',
                'conf' => [
                    'enabled' => true,
                    'read_only' => false,
                    'keep_backup' => true,
                    'login_protection' => true,
                    'log_successful_auth' => false,
                    'log_level' => 'LOG_WARNING',
                    'allow_pre_releases' => false,
                    'allow_development_packages' => false,
                    'hateoas' => false,
                    'expose_sensitive_fields' => false,
                    'override_sensitive_fields' => [],
                    'allowed_interfaces' => [],
                    'represent_interfaces_as' => 'descr',
                    'auth_methods' => ['BasicAuth'],
                    'jwt_exp' => 3600,
                    'ha_sync' => false,
                    'ha_sync_validate_certs' => true,
                    'ha_sync_hosts' => [],
                    'ha_sync_username' => '',
                    'ha_sync_password' => '',
                ],
            ],
        ],
    ],
];

// pfSense function stubs — only what schema generation actually calls
function config_get_path(string $path, mixed $default = null): mixed
{
    $parts = array_values(array_filter(explode('/', $path)));
    $current = PFSENSE_STUB_CONFIG;
    foreach ($parts as $part) {
        if (!is_array($current) || !array_key_exists($part, $current)) {
            return $default ?? [];
        }
        $current = $current[$part];
    }
    return $current;
}

function config_set_path(string $path, mixed $value, mixed $default = null): mixed
{
    return $default;
}

function config_init_path(string $path): void {}

function config_path_enabled(string $path, string $enable_key = 'enable'): bool
{
    return false;
}

function write_config(string $desc = 'Unknown', mixed $custom_config = null, bool $allow_empty = false): void {}

function mark_subsystem_dirty(string $subsystem = '', bool $force_sync = false): void {}

function clear_subsystem_dirty(string $subsystem = ''): void {}

function is_subsystem_dirty(string $subsystem = ''): bool
{
    return false;
}

function log_error(string $message): void
{
    fwrite(STDERR, "[pfsense-stub] $message\n");
}

function authenticate_user(string $username, string $password, array &$attributes = []): bool
{
    return false;
}

function is_user_enabled(mixed $user = null): bool
{
    return true;
}

function get_package_id(string $pkg_name): int
{
    return -1;
}

function pfsense_default_state_size(): int
{
    return 1000000;
}

function get_pfstate(): string
{
    return '0/1000000';
}

function get_next_bridgeif(): string
{
    return 'bridge0';
}

function get_next_number(string $prefix = '', array $existing = []): int
{
    return 0;
}

function ipsec_ikeid_next(): int
{
    return 1;
}

function ipsec_new_reqid(): int
{
    return 1;
}

function get_temp(): float|false
{
    return false;
}

function get_load_average(): array
{
    return [0.0, 0.0, 0.0];
}

function get_carp_interface_status(string $if = ''): string
{
    return '';
}

function get_carp_internal(): array
{
    return [];
}

function get_backups(): array
{
    return [];
}

function get_networks(string $if = ''): array
{
    return [];
}

function g_get(string $var): mixed
{
    return null;
}

function g_set(string $var, mixed $value): void {}

function get_config(string $path = '', mixed $default = null): mixed
{
    return $default ?? [];
}

function get_single_sysctl(string $name): mixed
{
    return null;
}

function set_single_sysctl(string $name, mixed $value): void {}

function get_sysctl(array $names): array
{
    return [];
}

function set_sysctl(array $names): void {}

function get_os_version(): string
{
    return '2.7.0';
}

function get_php_version(): string
{
    return PHP_VERSION;
}

function get_freenas_version(): string
{
    return '';
}

function get_unbound_status(): string
{
    return '';
}

function openvpn_get_active_clients(): array
{
    return [];
}

function openvpn_get_active_servers(): array
{
    return [];
}

function ipsec_get_status(): array
{
    return [];
}

function get_timezone_list(): array
{
    return [];
}

function get_country_codes(): array
{
    return [];
}

function get_rsync_status(): string
{
    return '';
}

function format_bytes(int $bytes, int $precision = 2): string
{
    return $bytes . 'B';
}

function format_uptime(int $seconds): string
{
    return $seconds . 's';
}

function get_uptime_s(): int
{
    return 0;
}

function get_mbuf_stats(): array
{
    return [];
}

function read_dummynet_config(): void {}

function dnpipe_find_nextnumber(): int
{
    return 1;
}

function dnqueue_find_nextnumber(): int
{
    return 1;
}

function read_altq_config(): void {}

function altq_get_interface_list(): array
{
    return [];
}

function get_real_interface(string $interface = '', string $family = '', bool $realv6iface = false): string
{
    return $interface;
}

function does_interface_exist(string $interface = ''): bool
{
    return false;
}

function get_interface_ip(string $interface = ''): string
{
    return '';
}

function get_interface_ipv6(string $interface = ''): string
{
    return '';
}

function get_interface_subnet(string $interface = ''): string
{
    return '';
}

function get_interface_subnetv6(string $interface = ''): string
{
    return '';
}

function get_interface_mtu(string $if = ''): int
{
    return 1500;
}

function get_interface_stats(string $if = ''): array
{
    return [];
}

function get_interface_info(string $if = ''): array
{
    return [];
}

function get_interface_ports(string $if = ''): array
{
    return [];
}

function get_configured_interface_with_descr(bool $only_opt = false): array
{
    return [];
}

function get_supported_media(string $if = ''): array
{
    return [];
}

function get_pkg_info(string $pkg, bool $metadata = false): array
{
    return [];
}

function get_services(): array
{
    return [];
}

function get_service_status(array $service): bool
{
    return false;
}

function is_service_enabled(string $service): bool
{
    return false;
}

function gen_subnet(string $ip, int $bits): string
{
    return '0.0.0.0';
}

function gen_subnet_max(string $ip, int $bits): string
{
    return '255.255.255.255';
}

function get_unique_id(): string
{
    return uniqid('', true);
}

function lookup_gateway_ip_by_name(string $name): string
{
    return '';
}

function return_gateways_array(bool $disabled = false, bool $localhost = false, bool $inactive = false): array
{
    return [];
}

function get_ipsec_sa_statuses(): array
{
    return [];
}

function ipsec_list_sa(): array
{
    return [];
}

function interface_ipsec_vti_list_all(): array
{
    return [];
}

function interface_lagg_configure(string $if): void {}

function interface_gre_configure(string $if): void {}

function ca_create(array &$ca, int $keylen, int $lifetime, string $dn, string $digest): bool
{
    return true;
}

function ca_inter_create(array &$ca, int $keylen, int $lifetime, string $dn, string $sign, string $digest): bool
{
    return true;
}

function ca_get_all_services(array $ca): array
{
    return [];
}

function cert_create(array &$cert, string $caref, int $keylen, int $lifetime, string $dn, string $type, string $digest): bool
{
    return true;
}

function cert_import(array &$cert, string $crt, string $key): bool
{
    return true;
}

function cert_get_dates(string $str): array
{
    return ['start' => '', 'end' => ''];
}

function cert_get_serial(string $str): string
{
    return '';
}

function cert_get_purpose(string $str): array
{
    return [];
}

function cert_in_use(string $ref): bool
{
    return false;
}

function cert_renew(array &$cert): bool
{
    return true;
}

function cert_restart_services(array $cert): void {}

function cert_get_all_services(array $cert): array
{
    return [];
}

function lookup_cert(string $refid): array
{
    return [];
}

function crl_in_use(string $ref): bool
{
    return false;
}

function crl_update(array &$crl, string $certs, int $lifetime): bool
{
    return true;
}

function csr_generate(array &$cert, int $keylen, string $dn, string $type, string $digest): bool
{
    return true;
}

function csr_sign(array &$cert, array $ca, int $lifetime, string $type, string $digest): bool
{
    return true;
}

function get_configured_interface_list(bool $only_opt = false, bool $construct = false): array
{
    return [];
}

function get_configured_interface_list_by_realif(bool $only_opt = false): array
{
    return [];
}

function convert_real_interface_to_friendly_interface_name(string $realif): string
{
    return $realif;
}

function get_interface_list(): array
{
    return [];
}

function get_alias_list(string $type = ''): array
{
    return [];
}

function get_certificates(bool $keytype = false): array
{
    return [];
}

function get_ca_list(): array
{
    return [];
}

function get_cert_list(): array
{
    return [];
}

function get_user_settings(): array
{
    return [];
}

function get_system_timezone(): string
{
    return 'UTC';
}

function mwexec(string $cmd, bool $nologging = false): int
{
    return 0;
}

function exec_safe(string $fmt, string ...$args): string
{
    return '';
}

function shell_safe(string $fmt, string ...$args): string
{
    return '';
}

// BaseTraits must be first — everything else depends on it
require_once $restapi . '/Core/BaseTraits.inc';

// Load in the same order as RESTAPI_LIBRARIES; skip GraphQL/Forms/Tests (not needed for OpenAPI)
$dirs = [
    '/Core/',
    '/Dispatchers/',
    '/Caches/',
    '/Responses/',
    '/Validators/',
    '/Fields/',
    '/ModelTraits/',
    '/Models/',
    '/QueryFilters/',
    '/ContentHandlers/',
    '/Schemas/',
    '/Auth/',
    '/Endpoints/',
];

foreach ($dirs as $dir) {
    foreach (glob($restapi . $dir . '*.inc') as $file) {
        require_once $file;
    }
}

$schema = new \RESTAPI\Schemas\OpenAPISchema();
echo $schema->get_schema_str();
