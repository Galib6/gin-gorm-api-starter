data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./cmd/atlas_loader",
  ]
}

env "local" {
  src = data.external_schema.gorm.url
  
  // Dev database to calculate diffs (requires Docker)
  dev = "docker://postgres/18/dev?search_path=public"
  
  migration {
    dir = "file://database/migrations"
    format = goose
  }
}
