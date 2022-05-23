{ config, lib, pkgs, ... }:

# see
# https://github.com/microsoft/vscode/blob/27a61bd9852cb8b808af99f0acedd3b5d3b9afd5/src/vs/server/serverEnvironmentService.ts#L12-L59

with lib;
let

  cfg = config.services.mgtt;
  defaultUser = "mgtt";
  defaultGroup = defaultUser;
  defaultPackage = pkgs.callPackage ../../default.nix { };



in
{
  ###### interface
  options = {
    services.mgtt = {
      enable = mkEnableOption "mgtt";

      package = mkOption {
        default = defaultPackage.package;
        description = "The package to use";
        type = types.package;
      };

      user = mkOption {
        type = types.str;
        default = defaultUser;
        example = "yourUser";
        description = ''
          The user to run this.
          By default, a user named <literal>${defaultUser}</literal> will be created.
        '';
      };

      group = mkOption {
        type = types.str;
        default = defaultGroup;
        example = "yourGroup";
        description = ''
          The group to run this.
          By default, a group named <literal>${defaultGroup}</literal> will be created.
        '';
      };

      dataDir = mkOption {
        type = types.path;
        default = "/var/lib/mgtt";
        example = "/home/yourUser";
      };

      config = lib.mkOption {
        default = {
          level = "info";
          json = false;
          
          url = "tcp://0.0.0.0:8883";

          plugins = "auth,acl";

          tls = {
            cert = {
              file = "./mgtt.cert";
            };
          };

          db = "./messages.db";
        };
        defaultText = "";
        description = "The literal config of dyndb";

      };
    };
  };

  ###### implementation

  config =
    let
      configFile = pkgs.writeTextFile
        {
          name = "mgtt-config";
          text = lib.generators.toYAML { } cfg.config;
        };

    in
    mkIf cfg.enable {

      users.users = mkIf (cfg.user == defaultUser) {
        "${defaultUser}" =
          {
            isSystemUser = true;
            group = cfg.group;
            home = cfg.dataDir;
            createHome = true;
            description = "mgtt user";
          };
      };

      users.groups = mkIf (cfg.group == defaultGroup) {
        "${defaultGroup}" = {
          name = "mgtt";
        };
      };

      # system.activationScripts = {
      #   dyndbsession = ''
      #     if [ ! -d ${dirOf cfg.config.session.storage.bolt} ]; then
      #       mkdir -vp ${dirOf cfg.config.session.storage.bolt}
      #     fi
      #     chown -R ${cfg.user} ${dirOf cfg.config.session.storage.bolt}
      #     chmod -R u=rwX ${dirOf cfg.config.session.storage.bolt}
      #   '';
      # };

      systemd.services.mgtt = {
        description = "mgtt";
        wantedBy = [ "multi-user.target" ];
        after = [ "network-online.target" ];
        environment = { };
        serviceConfig = {
          ExecStart = "${cfg.package}/bin/mgtt -c ${configFile}";
          ExecReload = "${pkgs.coreutils}/bin/kill -HUP $MAINPID";
          User = cfg.user;
          Restart = "on-failure";

          # Used as root directory
          RuntimeDirectory = "mgtt";

          # This currently not work :(
          #RootDirectory = "/run/dyndb";
          #BindReadOnlyPaths = [
          #  "/etc/ssl"
          #  "/etc/static/ssl"
          #  "/etc/resolv.conf"
          #  "${configFile}"
          #  "${cfg.package}/bin/dyndb"
          #];
          ReadWritePaths = [
            (dirOf cfg.config.tls.cert.file)
            (dirOf cfg.config.db)
          ];

          # Some security
          CapabilityBoundingSet = [ "CAP_NET_BIND_SERVICE" ];
          DevicePolicy = "closed";
          LockPersonality = true;
          MemoryDenyWriteExecute = true;
          NoNewPrivileges = true;
          ProtectHome = "read-only";
          PrivateDevices = true;
          PrivateMounts = true;
          PrivateTmp = true;
          PrivateUsers = true;
          ProtectClock = true;
          ProtectControlGroups = true;
          ProtectHostname = true;
          ProtectKernelLogs = true;
          ProtectKernelModules = true;
          ProtectKernelTunables = true;
          ProtectSystem = "full";
          RestrictAddressFamilies = [ "AF_INET" "AF_INET6" "AF_UNIX" ];
          RestrictNamespaces = true;
          RestrictRealtime = true;
          RestrictSUIDSGID = true;
          SystemCallArchitectures = "native";
          SystemCallFilter = [ "@system-service" "~@resources" ];
          UMask = "0077";
        };

      };

    };

  meta.maintainers = with maintainers; [ stackshadow ];
}
