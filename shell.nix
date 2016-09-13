with import <nixpkgs> {};

stdenv.mkDerivation {
  name = "linecfg-go-dev";
  buildInputs = [
    go
  ];
}
