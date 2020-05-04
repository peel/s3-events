{ buildGoModule
, nix-gitignore
}:

buildGoModule {
  pname = "s3-events";
  version = "0.0.1";
  src = nix-gitignore.gitignoreSource [] ./.;
  goPackagePath = "github.com/peel/s3-events";
  modSha256 = "07z0a7g5yc7dclzzqgpf6gp5xlyfpb2lhzmm38d61f18g4l95ncw";  
}
