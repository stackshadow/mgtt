{ system ? builtins.currentSystem
, pkgs ? import <nixpkgs> { inherit system; }
, version ? "development"
, vendorSha256 ? null
, homepage ? "https://gitlab.actaport.de/actaport/infrastructure/watchcat/"
}:
let
  lib = pkgs.lib;
  inherit (pkgs) buildGoModule;
  inherit (lib) sourceByRegex;


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
    -X gitlab.com/mgtt/internal/mgtt/config.Version=${version}
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
