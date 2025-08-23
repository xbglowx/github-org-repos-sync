# GitHub Copilot Instructions for github-org-repos-sync

## Repository Overview

This is a Go CLI application that synchronizes GitHub organization repositories to local storage. It provides functionality to clone and update multiple repositories from a GitHub organization in parallel, with options for filtering and customization.

**Repository**: `xbglowx/github-org-repos-sync`  
**Language**: Go 1.23+  
**Primary Framework**: Cobra CLI framework  
**License**: Apache License 2.0

## Project Purpose & Architecture

### Core Functionality
- **Repository Discovery**: Fetches all repositories from a specified GitHub organization using GitHub API
- **Parallel Operations**: Performs git clone/update operations in parallel for efficiency
- **Smart Updates**: Switches to default branch, stashes if dirty, fetches and rebases
- **Filtering Options**: Include/exclude repositories by name patterns
- **Archive Handling**: Option to skip archived repositories

### Key Components

1. **`main.go`** - Entry point that calls Cobra command execution
2. **`cmd/root.go`** - Cobra CLI configuration, flags, and command structure
3. **`cmd/github-org-sync.go`** - Core business logic for GitHub API integration and git operations

### Architecture Pattern
- **CLI Layer**: Cobra-based command interface with flags and validation
- **Service Layer**: GitHub API client integration and repository operations
- **Git Operations**: Direct git command execution via `os/exec`

## Development Guidelines

### Go Standards & Patterns

- **Go Version**: Use Go 1.23+ features and idioms
- **Error Handling**: Always handle errors explicitly, no panics in user-facing code
- **Context Usage**: Use `context.Context` for API calls and cancelable operations
- **Concurrency**: Use goroutines with semaphores for controlled parallelism
- **Resource Management**: Proper cleanup with defer statements

### Code Style & Conventions

- **Package Structure**: Follow standard Go project layout with `cmd/` for applications
- **Naming**: Use descriptive names, follow Go naming conventions (camelCase for private, PascalCase for public)
- **Documentation**: Use Go doc comments for exported functions and types
- **Formatting**: Code must pass `go fmt` and `golangci-lint`

### Dependencies & Frameworks

#### Core Dependencies
- **CLI Framework**: `github.com/spf13/cobra` - Command line interface
- **GitHub API**: `github.com/google/go-github` - GitHub API client
- **OAuth2**: `golang.org/x/oauth2` - Authentication for GitHub API

#### Build Tools
- **Linting**: `golangci-lint` (configured via GitHub Actions)
- **Testing**: Standard Go testing with `go test`
- **Building**: Standard Go build with `go build`

## Configuration & Environment

### Required Environment Variables
- **`GITHUB_TOKEN`**: Personal access token with `repo` scope for GitHub API access
- Token must have access to the target organization's repositories

### CLI Flags & Options
```go
--destination-path, -d    // Destination path for repositories (default: current directory)
--parallelism, -p        // Number of parallel git operations (default: 1)
--skip-archived          // Skip archived repositories
--exclude-repos         // Exclude repositories containing specified string
--include-repos         // Include only repositories containing specified string
```

### Validation Rules
- GitHub token must be present in environment
- Git CLI must be available in PATH
- Cannot use both `--exclude-repos` and `--include-repos` simultaneously

## Git Operations & Workflow

### Repository States Handled
1. **New Repository**: Clone from GitHub to local destination
2. **Existing Clean Repository**: Checkout default branch, fetch, and rebase
3. **Existing Dirty Repository**: Stash changes, then proceed with update
4. **Empty Repository**: Skip update operations

### Git Commands Used
- `git clone <url> <destination>` - Clone new repositories
- `git -C <path> diff-index --quiet --cached HEAD --` - Check for staged changes
- `git -C <path> diff-files --quiet` - Check for unstaged changes
- `git -C <path> stash push` - Stash dirty changes
- `git -C <path> fetch origin` - Fetch latest changes
- `git -C <path> checkout <branch>` - Switch to default branch
- `git -C <path> pull --rebase` - Update with rebase

## Testing & Quality Assurance

### Current State
- **No unit tests currently exist** - This is an area for improvement
- **Manual testing**: Build and run against real GitHub organizations
- **CI/CD**: GitHub Actions for build validation and linting

### Testing Guidelines (for future development)
- Add unit tests for core functions in `cmd/github-org-sync.go`
- Mock GitHub API responses for consistent testing
- Test git operations with temporary repositories
- Integration tests with test GitHub organizations

### Quality Checks
- **Build Validation**: `go build` must succeed
- **Linting**: `golangci-lint` checks for code quality
- **Code Analysis**: CodeQL analysis for security issues

## GitHub Actions Workflows

### Existing Workflows
1. **`build-test.yaml`** - Go build and test validation
2. **`golangci-lint.yml`** - Code linting and style checks
3. **`codeql-analysis.yml`** - Security analysis
4. **`create-release.yaml`** - Automated releases

### Workflow Standards
- Triggered on pull requests and pushes (excluding version tags)
- Uses latest Ubuntu runners
- Go version specified in `go.mod` file
- All workflows must pass for PR approval

## Error Handling & Logging

### Error Patterns
- Environment validation errors are returned early from `checkRequirements()`
- Git operation errors are logged but don't stop other repositories
- API errors should be handled gracefully with informative messages

### Logging Style
- Use `fmt.Printf()` for user-facing output
- Include repository names in error messages for context
- Distinguish between INFO, ERROR, and success messages
- Example: `fmt.Printf("ERROR: Repo %s failed to clone: %v\n", repo.GetName(), err)`

## Common Development Patterns

### Goroutine with Semaphore
```go
func (gh *GhOrgSync) processRepo(sem chan struct{}, repo *github.Repository) {
    defer gh.wg.Done()
    
    sem <- struct{}{}
    defer func() { <-sem }()
    
    // Repository processing logic
}
```

### Git Command Execution
```go
cmd := exec.Command("git", "-C", repoPath, "status")
err := cmd.Run()
if err != nil {
    // Handle error appropriately
}
```

### GitHub API Usage
```go
ctx := context.Background()
repos, _, err := client.Repositories.ListByOrg(ctx, org, &github.RepositoryListByOrgOptions{
    ListOptions: github.ListOptions{PerPage: 100},
})
```

## Security Considerations

### Authentication
- Use GitHub personal access tokens (never hardcode credentials)
- Tokens should have minimal required scopes (`repo` for private repos, public access for public)
- Validate token presence before making API calls

### Git Operations
- All git commands use `-C` flag to specify working directory (safer than `cd`)
- Repository paths are constructed safely using `fmt.Sprintf`
- No shell injection risks as commands use `exec.Command` with separate arguments

## Release & Versioning

### Release Process
- Uses `release-please` for automated releases
- Follows conventional commits for changelog generation
- Version managed in `.release-please-manifest.json`
- Releases triggered by merging release PRs

### Versioning Strategy
- Semantic versioning (SemVer)
- Current version: 0.0.8 (as of latest release)
- Breaking changes increment minor version (pre-1.0)

## Contribution Guidelines

### Making Changes
1. **Environment Setup**: Ensure Go 1.23+, git, and GitHub token configured
2. **Code Quality**: Run `go fmt`, `golangci-lint`, and `go build` before committing
3. **Testing**: Manual testing with real or test GitHub organizations
4. **Documentation**: Update README.md for user-facing changes

### Commit Standards
- Follow conventional commits format
- Examples: `feat: add new filtering option`, `fix: handle empty repositories`
- Use imperative mood in commit messages

## Common Tasks & Troubleshooting

### Local Development
```bash
# Setup
export GITHUB_TOKEN=<your-token>
go mod download

# Build
go build .

# Run
./github-org-repos-sync <org-name>

# Test with options
./github-org-repos-sync --parallelism 3 --skip-archived <org-name>
```

### Debugging
- Check GitHub token permissions if API calls fail
- Verify git is in PATH if git commands fail
- Use `--parallelism 1` for easier debugging of git operations
- Check repository permissions if clone/fetch operations fail

### Performance
- Increase `--parallelism` for faster operation (default: 1, recommended: 3-5)
- GitHub API has rate limits (5000 requests/hour for authenticated users)
- Large organizations may require pagination handling

## Future Enhancement Areas

### Potential Improvements
1. **Testing**: Add comprehensive unit and integration tests
2. **Logging**: Implement structured logging with levels
3. **Configuration**: Support for config files beyond CLI flags  
4. **Git Integration**: Use git libraries instead of shelling out to git commands
5. **Progress**: Add progress bars for long-running operations
6. **Dry Run**: Add option to preview what would be done without executing

### Architecture Considerations
- Consider separating GitHub API logic into separate package
- Abstract git operations behind an interface for testing
- Add configuration validation beyond environment checks
- Consider using channels for better error aggregation across goroutines

## Contact & Maintenance

- **Primary Maintainer**: Based on repository ownership (xbglowx)
- **Issue Tracking**: GitHub Issues on the repository
- **CI Status**: Check GitHub Actions for build and lint status
- **Dependencies**: Renovate.json configured for automated dependency updates

## Important Notes

⚠️ **Repository Access**: Ensure your GitHub token has access to the target organization's repositories

⚠️ **Disk Space**: Large organizations can require significant disk space for all repositories

⚠️ **Network**: Initial cloning may take considerable time depending on repository sizes and network speed