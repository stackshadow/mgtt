# nix-instantiate '<nixpkgs/nixos>' -I nixos-config=./configuration.example.nix -A system
# nix-build '<nixpkgs/nixos>' -I nixos-config=./configuration.example.nix -A system
{ config, lib, pkgs, modulesPath, ... }:
let

  mgttPackage = import ../../default.nix { };


  #codeservernew = callPackage ../packages/code-server-bin.nix { };
in
{

  imports = [
    # Include the results of the hardware scan.
    mgttPackage.module
  ];

  boot.loader.grub.device = "/dev/sda";
  fileSystems."/".device = "/dev/sda1";

  services.mgtt = {
    enable = true;
    config = {
      level = "info";
      json = false;

      url = "tcp://0.0.0.0:8883";

      plugins = {
        acl = {
          enable = true;
          rules = {
            firstuser = [
              {
                direction = "w";
                route = "device1/#";
                allow = true;
              }
              {
                direction = "r";
                route = "device1/#";
                allow = true;
              }
            ];
          };
        };
        auth = {
          enable = true;
          new = [
            {
              username = "firstuser";
              password = "password";
            }
          ];
        };
      };

      tls = {
        cert = {
          file = "./mgtt.cert";
        };
      };

      db = "./messages.db";
    };
  };



}
