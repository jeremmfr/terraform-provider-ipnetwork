default: install

.PHONY: install testunit cleanout changemd

# Install to use dev_overrides in provider_installation of Terraform
install:
	go install

# Run unit tests
testunit: 
	go test -race -v -coverprofile=coverage_unit.out ./...
	go tool cover -html=coverage_unit.out

# Cleanup out files from tests
cleanout:
	find . -maxdepth 1 -name "*.out" -type f -delete

changemd:
	cp .changes/.template.md .changes/$(shell git rev-parse --abbrev-ref HEAD).md
