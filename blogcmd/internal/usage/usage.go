package usage

const Text string = `使い方:
  ./blog {dev|stg|prd} {up|down|restart|recreate} [docker compose args...]
  ./blog {dev|stg|prd} [docker compose args...]
  ./blog {dev|stg|prd} mysql [mysql args...]
  ./blog {dev|stg|prd} migrate {up|down|refresh}
  ./blog {dev|stg|prd} seed <seed_name>
  ./blog {stg|prd} deploy [docker compose args...]
  ./blog {stg|prd} image {build|push|build-push|pull}
  ./blog back fmt
  ./blog back mod tidy
  ./blog back test [package path]
  ./blog ent generate
  ./blog front install
  ./blog build
`
