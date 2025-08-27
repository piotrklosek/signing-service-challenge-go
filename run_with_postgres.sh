go run . server \
  --db.type=postgres \
  --db.postgres.dsn="postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
