{ system ? builtins.currentSystem, binary, version ? "dev" }:
let
  pkgs = import <nixpkgs> { inherit system; };
  lib = pkgs.lib;
  inherit (pkgs) dockerTools;
  inherit (lib) sourceByRegex;

in
#buildLayeredImage
pkgs.dockerTools.streamLayeredImage {
  name = "mgtt";
  tag = version;
  maxLayers = 10;

  contents = with pkgs; [
    # basic stuff
    cacert

    # our binary
    binary
  ];

  created = "now";

  fakeRootCommands = ''
    ${pkgs.dockerTools.shadowSetup}
    groupadd mgtt
    useradd --gid 5005 mgtt

    mkdir -p /tmp
    chmod a+rwX /tmp
  '';

  config = {
    User = "5005";
    Env = [
      "HOME=/tmp"
    ];
    ExposedPorts = {
      "8000/tcp" = { };
    };
    Cmd = [
      "${binary}/bin/mgtt"
    ];
  };
}
