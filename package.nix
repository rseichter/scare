{
  lib,
  buildGoModule,
  fetchFromGitHub,
}:
let
  myrev = "6c935bb9eec7757001a878fa129f4173acb6ff04";
  myhash = "sha256-b7ZqkYmOIWiV4Q/Kc6I9R2FnsI8TZ2wIddmdO/wfVf4=";
  # myhash = lib.fakeHash;
in
buildGoModule (finalAttrs: {
  pname = "scare";
  version = "0.6.dev1";

  # src = ./.;

  src = fetchFromGitHub {
    owner = "rseichter";
    repo = "scare";
    hash = myhash;
    rev = myrev;
    # tag = "v${finalAttrs.version}";
  };

  vendorHash = null;
  env.CGO_ENABLED = 0;

  meta = {
    homepage = "https://github.com/rseichter/scare";
    description = "An opinionated script care utility";
    maintainers = with lib.maintainers; [ rseichter ];
  };
})
