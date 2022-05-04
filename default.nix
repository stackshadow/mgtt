{ system ? builtins.currentSystem
, pkgs ? import <nixpkgs> { inherit system; }
}:
let
  lib = pkgs.lib;
  callPackage = lib.callPackageWith (pkgs // pkgs.lib);
  inherit (lib) sourceByRegex;

  packageVersion = "0.34.0";


  # the binary
  binary = callPackage ./build/nix/package.nix {
    version = packageVersion;
  };

  dockerImage = callPackage ./build/nix/docker.nix {
    inherit binary;
    version = "dev";
  };

in
{
  # nix-build -A package --no-out-link
  package = binary;

  packageAarch64 = callPackage ./build/nix/package.nix {
    pkgs = pkgs.pkgsCross.aarch64-multiplatform;
    version = packageVersion;
  };

  # nix-build -A docker --no-out-link

  # $(nix-build -A docker) | docker load
  docker = dockerImage;

}
