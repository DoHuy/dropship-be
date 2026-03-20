data "external_schema" "gorm" {
  program = [
    "go",
    "run",
    "./tools/atlas_loader.go",
  ]
}

env "local" {
  src = data.external_schema.gorm.url

  # Database tạm để Dev (Bắt buộc)
  dev = "docker://postgres/15/dev?search_path=public"

  migration {
    dir = "file://dropship-deployment/sql/migrations?format=golang-migrate"
  }
}