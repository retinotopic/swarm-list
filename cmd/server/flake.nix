{
  description = "flake for dockerfile";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs?ref=nixos-25.05";
  };

  outputs = { self, nixpkgs }:
    let
      allSystems = [
        "x86_64-linux" # 64-bit Intel/AMD Linux
        "aarch64-linux" # 64-bit ARM Linux
        "x86_64-darwin" # 64-bit Intel macOS
        "aarch64-darwin" # 64-bit ARM macOS
      ];
      forAllSystems = f: nixpkgs.lib.genAttrs allSystems (system: f {
        pkgs = import nixpkgs { inherit system; };
      });
    in
    {
      packages = forAllSystems ({ pkgs }: {
        default = pkgs.buildGo124Module {
          pname = "server";
          version = "0.1.0";
          src = ./.;
          env.CGO_ENABLED = 0;
          GOFLAGS = [
            "-buildmode=exe"
          ];
          ldflags = [
            "-s -w"
          ];
          
          vendorHash = "sha256-0iHh1F1VH0aX+R8b0a+1KPi7QsWRpYRoyFmW+O8PB5I=";
        };
      });
    };
}
