{ pkgs ? import <nixpkgs> {} }:
  pkgs.mkShell {
    nativeBuildInputs = with pkgs.buildPackages; [ sqlite-interactive gcc ];
    shellHook = ''
      export CGO_ENABLED=1
    '';
}