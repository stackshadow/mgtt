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

in
{
  # nix-build -A package --no-out-link
  package = binary;
  #module = import ./build/module.nix;
}
