with (import <nixpkgs> { });
let
  package = import ./default.nix { };
in
mkShell {
  buildInputs = [
    package.package
  ];
}
