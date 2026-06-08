{
  stdenvNoCC,
  fetchFromGitHub,
  php,
}:

let
  restapi-src = fetchFromGitHub {
    owner = "pfrest";
    repo = "pfSense-pkg-RESTAPI";
    rev = "4190c8ca40c73e5705e01f227a029b61dedf7b7f";
    hash = "sha256-zwa84mRQKhFqQ3d+pEkoZxmdBAPVrUz+IImA8oHXIHM=";
  };

  firebase-php-jwt = fetchFromGitHub {
    owner = "firebase";
    repo = "php-jwt";
    rev = "5645b43af647b6947daac1d0f659dd1fbe8d3b65";
    hash = "sha256-4KxEAR7xodzrmrUFckQSojsvqDNpPsJdXCN1JTSCmyo=";
  };

in

stdenvNoCC.mkDerivation {
  pname = "pfsense-openapi";
  version = "2.8.1";

  dontUnpack = true;

  nativeBuildInputs = [ php ];

  buildPhase = ''
    # Copy RESTAPI source; dir name must stay "RESTAPI" to match namespace→path mapping
    cp -r ${restapi-src}/pfSense-pkg-RESTAPI/files/usr/local/pkg/RESTAPI RESTAPI
    chmod -R u+w RESTAPI

    # Patch the one hardcoded /usr/local/pkg/ prefix in get_classes_from_namespace()
    sed -i "s|'/usr/local/pkg/'|'$PWD/'|" RESTAPI/Core/Tools.inc

    # Replace autoloader.inc with a no-op — bootstrap.php handles all loading;
    # individual .inc files re-include autoloader.inc as a safety net but it's redundant here
    php -r 'file_put_contents("RESTAPI/autoloader.inc", "<?php // replaced by openapi-bootstrap.php\n");'

    # Schema generator creates model instances without skip_init, causing internal_callable
    # functions to fire (which call pfSense system functions). Pass skip_init: true instead.
    # Use PHP for the replacement to avoid bash/sed dollar-sign escaping issues.
    php -r '$f="RESTAPI/Schemas/OpenAPISchema.inc"; file_put_contents($f, str_replace("\$model = new \$model();", "\$model = new \$model(skip_init: true);", file_get_contents($f)));'

    # Set up firebase/php-jwt without Composer
    mkdir -p RESTAPI/.resources/vendor/Firebase/JWT
    cp -r ${firebase-php-jwt}/src/. RESTAPI/.resources/vendor/Firebase/JWT/
    cat > RESTAPI/.resources/vendor/autoload.php << 'AUTOLOAD'
<?php
// Minimal GraphQL stubs — only what GraphQLResponse.inc needs during schema generation.
// Namespace blocks must come before global code in PHP.
namespace GraphQL\Executor {
    class ExecutionResult {
        public mixed $data = null;
        public array $errors = [];
        public array $extensions = [];
    }
}
namespace GraphQL\Type\Definition {
    class ObjectType { public function __construct(array $config = []) {} }
    class Type { public static function string(): mixed { return null; } }
}
namespace GraphQL {
    class GraphQL {
        public static function executeQuery(mixed ...$args): \GraphQL\Executor\ExecutionResult {
            return new \GraphQL\Executor\ExecutionResult();
        }
    }
}
namespace {
    // Firebase\JWT PSR-4 autoloader
    spl_autoload_register(function (string $class): void {
        if (str_starts_with($class, 'Firebase\\JWT\\')) {
            $file = __DIR__ . '/Firebase/JWT/' . str_replace('\\', '/', substr($class, 13)) . '.php';
            if (file_exists($file)) {
                require $file;
            }
        }
    });
}
AUTOLOAD

    php ${./openapi-bootstrap.php} "$PWD" > openapi.json
  '';

  installPhase = ''
    install -D openapi.json $out/openapi.json
  '';
}
