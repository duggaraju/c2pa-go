module github.com/duggaraju/c2pa-go/c2pa

go 1.24.0

retract [v0.94.0, v0.94.1] // v0.94.0 was published accidentally; v0.94.1 contains retractions only.

require github.com/stretchr/testify v1.12.1

require go.yaml.in/yaml/v3 v3.0.5 // indirect
