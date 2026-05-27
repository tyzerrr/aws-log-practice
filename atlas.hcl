variable "url" {
  type    = string
  default = getenv("DATABASE_URL")
}

env "local" {
  src = "file://db/schema.sql"
  url = var.url
  dev = "docker://postgres/17/dev?search_path=public"
  migration {
    dir = "file://db/migrations"
  }
}
