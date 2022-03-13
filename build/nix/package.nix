{ system ? builtins.currentSystem
, sourceByRegex
, version
, vendorSha256 ? null
, homepage ? "https://gitlab.actaport.de/actaport/infrastructure/watchcat/"
}:
let
  pkgs = import <nixpkgs> { inherit system; };
  lib = pkgs.lib;
  inherit (pkgs) buildGoModule;

  localSource = sourceByRegex ../.. [
    "^go.mod"
    "^go.sum"
    "vendor"
    "vendor/.*"
    "pkg"
    "pkg/.*"
    "cmd"
    "cmd/.*"
    "internal"
    "internal/.*"
  ];



in
buildGoModule {
  pname = "mgtt";
  inherit version;
  src = localSource;

  # excludedPackages = "test";
  subPackages = [ "cmd/mgtt" ];

  inherit vendorSha256;
  proxyVendor = true;

  buildFlagsArray = ''
    -ldflags=
    -w -s
    -X gitlab.com/mgtt/internal/mgtt/cli.Version=${version}
    -extldflags=-static
  '';

  preBuild = ''
 
  '';

  meta = with lib; {
    description = "mgtt - an mqtt broker written in go";
    license = licenses.mit;
    homepage = homepage;
    maintainers = with maintainers; [ stackshadow ];
  };
}
