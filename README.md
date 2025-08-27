# github-org-repos-sync

[![Build and Test](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/build-test.yaml/badge.svg)](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/build-test.yaml) [![golangci-lint](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/golangci-lint.yml) [![CodeQL](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/codeql-analysis.yml) [![Release](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/create-release.yaml/badge.svg)](https://github.com/xbglowx/github-org-repos-sync/actions/workflows/create-release.yaml)

A command-line tool to efficiently sync all repositories from a GitHub organization to your local machine. Perfect for developers who need to maintain local copies of multiple repositories for backup, analysis, or bulk operations. 

## Features

- **Bulk Repository Sync**: Clone or update all repositories from a GitHub organization in one command
- **Smart Updates**: Automatically switches to default branch and pulls latest changes
- **Dirty Repository Handling**: Safely stashes uncommitted changes before updating
- **Parallel Processing**: Configurable concurrency for faster synchronization
- **Repository Filtering**: Include/exclude specific repositories by name pattern
- **Archived Repository Support**: Option to skip archived repositories
- **Permission Aware**: Automatically skips repositories without pull permissions
- **Empty Repository Handling**: Gracefully handles empty repositories without errors

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/xbglowx/github-org-repos-sync/releases).

### Build from Source

1. **Clone the repository**:
   ```bash
   git clone git@github.com:xbglowx/github-org-repos-sync.git
   cd github-org-repos-sync
   ```

2. **Build the binary**:
   ```bash
   go build .
   ```

   **Optional**: Build with version from git:
   ```bash
   VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev-$(git rev-parse --short HEAD)")
   go build -ldflags "-X github.com/xbglowx/github-org-repos-sync/cmd.Version=$VERSION" .
   ```

   **Optional**: Build with custom version:
   ```bash
   go build -ldflags "-X github.com/xbglowx/github-org-repos-sync/cmd.Version=1.0.0" .
   ```

## Prerequisites

### Authentication

1. **GitHub Personal Access Token**: Create a token with `repo` scope
   ```bash
   export GITHUB_TOKEN=<your-token>
   ```

2. **Git CLI Authentication**: Ensure `git` is authenticated with GitHub (SSH keys or credential helper)

## Usage

### Basic Syntax

```bash
./github-org-repos-sync <organization> -d <destination-path> [options]
```

### Examples

#### Sync all repositories from an organization
```bash
./github-org-repos-sync myorg -d ~/repos/myorg
```

#### Sync with custom parallelism (faster for large organizations)
```bash
./github-org-repos-sync myorg -d ~/repos -p 5
```

#### Skip archived repositories
```bash
./github-org-repos-sync myorg -d ~/repos --skip-archived
```

#### Filter repositories (include only those containing "api")
```bash
./github-org-repos-sync myorg -d ~/repos --include "api"
```

#### Exclude specific repositories (skip test repos)
```bash
./github-org-repos-sync myorg -d ~/repos --exclude "test"
```

#### Backup an entire organization
```bash
./github-org-repos-sync mycompany -d ~/backups/mycompany-repos
```

#### Sync only service repositories
```bash
./github-org-repos-sync mycompany -d ~/services --include "service"
```

#### Show help
```bash
./github-org-repos-sync --help
```

#### Check version
```bash
# Using version subcommand
./github-org-repos-sync version

# Using version flag
./github-org-repos-sync --version
./github-org-repos-sync -v
```

## How It Works

1. **Repository Discovery**: Fetches a list of all repositories in the specified organization that you have access to
2. **Local Check**: Determines if each repository already exists locally
3. **Clone or Update**:
   - **New repositories**: Clones them to the destination path
   - **Existing repositories**: 
     - Fetches latest changes from remote
     - Stashes uncommitted changes if the working directory is dirty
     - Switches to the repository's default branch
     - Pulls latest changes with rebase
4. **Parallel Processing**: Processes multiple repositories concurrently for improved performance
5. **Error Handling**: Gracefully handles edge cases like empty repositories, missing branches, and permission issues

## Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `-d, --dest` | Destination directory for repositories | Required |
| `-p, --parallelism` | Number of concurrent operations | 10 |
| `--skip-archived` | Skip archived repositories | false |
| `--include` | Include only repositories containing this string | "" |
| `--exclude` | Exclude repositories containing this string | "" |
| `-h, --help` | Show help message | - |

## Troubleshooting

### Permission Issues
- Ensure your GitHub token has `repo` scope
- Verify you have access to the organization
- Check that git CLI is properly authenticated

### Network Issues
- The tool will retry failed operations
- Use lower parallelism (`-p 1`) for unstable connections

### Large Organizations
- Consider using filters (`--include`/`--exclude`) for large organizations
- Monitor disk space when syncing many repositories

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.
