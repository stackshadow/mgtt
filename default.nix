{ system ? builtins.currentSystem }:
let
  # reproducable build
  nixpkgs = fetchTarball {
    url = "https://github.com/NixOS/nixpkgs/archive/9344233ab1cea59c3461c5cedde6a08abb89e6ea.tar.gz";
  };
  pkgs = import nixpkgs { };
  lib = pkgs.lib;
  callPackage = lib.callPackageWith (pkgs // pkgs.lib);
  inherit (lib) sourceByRegex;

  packageVersion = "0.36.0";


  # the binary
  binary = callPackage ./build/nix/package.nix {
    inherit pkgs;
    version = packageVersion;
  };

  dockerImage = callPackage ./build/nix/docker.nix {
    inherit pkgs;
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

  module = import ./build/nix/module.nix;
}
