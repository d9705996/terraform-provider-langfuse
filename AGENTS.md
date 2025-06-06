# Repository Guidelines

## Commit Messages
- Use **Conventional Commits** for all commit messages (e.g., `feat: add feature`, `fix: correct bug`).

## Development Workflow
- After code changes run:
  - `go mod tidy`
  - `gofmt -w $(git ls-files '*.go')`
  - `go vet ./...`
  - `golangci-lint run`
  - `go test ./...`
- Ensure all checks pass before committing.

## Changelog Automation
- A GitHub Actions workflow updates `CHANGELOG.md` when PRs are merged into `main`. Do not manually edit the changelog.
