{
  lib,
  buildGoModule,
}:
buildGoModule (finalAttrs: {
  pname = "scare";
  version = "0.6.dev1";

  src = ./.;
  # src = fetchFromGitHub {
  #   owner = "rseichter";
  #   repo = "scare";
  #   tag = "v${finalAttrs.version}";
  #   hash = lib.fakeHash;
  # };
  vendorHash = null;
  env.CGO_ENABLED = 0;

  meta = {
    homepage = "https://github.com/rseichter/scare";
    description = "An opinionated script care utility";
    maintainers = with lib.maintainers; [ rseichter ];
  };
})
