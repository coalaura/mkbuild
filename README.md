# mkbuild

A streamlined, opinionated GitHub Actions workflow generator for Go releases.

`mkbuild` creates ready-to-use `release.yml` files that handle cross-compilation, artifact packaging and GitHub releases. It supports static builds by default, CGO cross-compilation via Zig and interactive configuration, allowing you to generate a complete CI/CD pipeline in seconds.

## Installation

Download the latest binary from the [Releases Page](https://github.com/coalaura/mkbuild/releases) or install a prebuilt binary with one command:

```bash
curl -sL https://src.w2k.sh/mkbuild/install.sh | sh
```

## Usage

Run `mkbuild` in the root of your project directory.

```bash
# Generate workflow interactively
mkbuild my-app -i

# Preview generated file without writing
mkbuild my-app -n

# Generate with CGO/Zig support and archives
mkbuild my-app --cgo --archive
```

### Generated Artifacts

The tool creates a `.github/workflows/release.yml` file configured to:

1. **Build**: Cross-compile binaries for Linux, Windows and macOS (amd64/arm64).
2. **Package**: Optionally create `.tar.gz`/`.zip` archives.
3. **Release**: Create a GitHub Release with artifacts and SHA256 checksums.

## Configuration

`mkbuild` is designed to be flexible while providing sensible defaults.

* **`--cgo`**: Enables CGO and configures cross-compilation via Zig. Produces static `musl` binaries for Linux and `gnu` binaries for Windows.
* **`--main <path>`**: Path to the main package (default: `.`).
* **`--archive`**: Creates `.tar.gz` (Linux/macOS) and `.zip` (Windows) archives instead of raw binaries.
* **`--checksums`**: Generates a `checksums-sha256.txt` file in the release.
* **`--draft`**: Marks the GitHub Release as a draft.

## Workflow Features

* **Static by Default**: Builds with `CGO_ENABLED=0` for portable, dependency-free binaries.
* **Zig Cross-Compilation**: When CGO is enabled, uses Zig to produce fully static `musl` binaries for Linux, removing the need for specific build runners or `glibc` compatibility issues.
* **Optimized Output**: Uses `-trimpath` and `-ldflags "-s -w"` for reproducible, minimal binaries.
* **Version Injection**: Automatically injects the Git tag into `main.Version` via linker flags.
* **Smart Caching**: Caches Go modules and build caches for non-CGO builds to speed up workflows.
