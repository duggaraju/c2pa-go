module github.com/duggaraju/c2pa-go/c2pa

go 1.24.0

retract [v0.94.0, v0.94.1] // v0.94.0 was published accidentally; v0.94.1 contains retractions only.

require github.com/stretchr/testify v1.11.1

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
