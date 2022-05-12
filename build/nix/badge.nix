# run it with nix-shell build/nix/badge.nix

{ pkgs ? import <nixpkgs> { }
, version ? "development"
}:
let
  lib = pkgs.lib;
  inherit (lib) sourceByRegex;

  gobadgeSource =
    builtins.fetchGit
      {
        name = "dyndb";
        url = "https://gitlab.com/stackshadow/gobadge-cli.git";
        ref = "refs/heads/main";
        rev = "2e939eb04d5d0dd5e047e0b80919442a79eeeb12";
      };
  gobadge = import gobadgeSource { };

  # we need a newer version for gocyclo
  nixosGoCyclo = import
    (pkgs.fetchzip {
      url = "https://github.com/NixOS/nixpkgs/archive/bf01537f0c9deccf7906b51e101d05c039390bb8.zip";
      sha256 = "sha256-fgPiS1heTNSi5i+22pMxoj7t/iOg42zRZJqxeTCJPjU=";
    })
    { };
  gocyclo = nixosGoCyclo.gocyclo;

in
pkgs.mkShell {
  buildInputs = with pkgs;
    [
      go
      gnumake
      gocyclo
      gosec
      gobadge.package
    ];
  shellHook = ''
    export FONT_FILE=${pkgs.freefont_ttf}/share/fonts/truetype/FreeSans.ttf

    clean() {
      nix-shell ${gobadgeSource}/shell.nix --command clean
    }
    version() {
      gobadge-cli --label=version --text --value="${version}" --file-name=./version.svg
    }
    coverage() {
      set +e
      nix-shell ${gobadgeSource}/shell.nix --command cover
    }
    cyclo() {
      nix-shell ${gobadgeSource}/shell.nix --command cyclo
    }
    sec() {
      gosec -color=false -no-fail -severity medium ./... > gosec.out.gobadge
      value=$(cat gosec.out.gobadge | grep Issues | cut -d':' -f2 | sed -e 's/\s*//')

      gobadge-cli --value-min=10.00 --value-max=0.01 --label=gosec --value=$value --file-name=./gosec.svg
    }
    lastbuild(){
      nix-shell ${gobadgeSource}/shell.nix --command lastbuild
    }
  '';
}
