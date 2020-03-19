with import <nixpkgs> { };
let
  GOPATH = toString ./gopath;
  bootstrap = pkgs.writeShellScript "bootstrap" ''
    mkdir -p ${GOPATH}
  '';
in
pkgs.mkShell rec {
  GOROOT = "${go}/share/go";
  inherit bootstrap GOPATH;
  buildInputs = [
    go
    goimports
    go-langserver
    vgo2nix
    gotools
  ];
}
