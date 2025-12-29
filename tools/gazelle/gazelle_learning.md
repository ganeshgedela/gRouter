# Gazelle Learning Guide

This document explains how to use Gazelle in the `gRouter` project to manage Bazel build files.

## What is Gazelle?

Gazelle is a Bazel build file generator for Go projects. It automatically creates and updates `BUILD.bazel` files by analyzing your Go source code and dependencies.

## Standard Usage in this Project

In the root of the `gRouter` workspace, you will find configuration directives in `BUILD.bazel` and a run target for Gazelle.

### Running Gazelle

To update all `BUILD.bazel` files in your project (e.g., after adding new files or imports), run:

```bash
bazel run //:gazelle
```

This command will:
1. Scan your Go files.
2. Identify imports and dependencies.
3. specific `BUILD.bazel` files in each package directory.
4. Auto-generate `go_library`, `go_test`, and `go_binary` targets.

### Updating Dependencies

If you have added new external dependencies to `go.mod`, you need to update Bazel's repository rules. Run:

```bash
bazel run //:gazelle-update-repos
```

This updates `deps.bzl` (or `WORKSPACE`) with the new Go modules defined in `go.mod`.

## Configuration

Gazelle is configured using special comments (directives) in your `BUILD.bazel` files.

### Common Directives

- **Define Module Name**: Tells Gazelle the module path (matches `go.mod`).
  ```starlark
  # gazelle:prefix grouter
  ```

- **Proto Mode**: disables or enables proto rule generation.
  ```starlark
  # gazelle:proto disable
  ```

- **Exclude Files/Directories**:
  ```starlark
  # gazelle:exclude node_modules
  ```

### Project Specifics

For this project `gRouter`, the main configuration is in the root `BUILD.bazel`.

```starlark
load("@bazel_gazelle//:def.bzl", "gazelle")

# gazelle:prefix grouter
gazelle(name = "gazelle")

gazelle(
    name = "gazelle-update-repos",
    args = [
        "-from_file=go.mod",
        "-to_macro=deps.bzl%go_dependencies",
        "-prune",
    ],
    command = "update-repos",
)
```

## Best Practices

1. **Run often**: Run `bazel run //:gazelle` whenever you add a file, change an import, or modify `go.mod`.
2. **Don't edit generated rules manually**: Gazelle will overwrite manual changes to generated rules. If you need to modify them, use directives or add `keep` comments (e.g., `# keep`).
3. **Commit generated files**: Check in the updated `BUILD.bazel` files to source control.

## Troubleshooting

- **Dependency not found**: If Bazel complains about a missing dependency, ensure it's in `go.mod`, run `go mod tidy`, then `bazel run //:gazelle-update-repos`.
- **Visibility issues**: Gazelle defaults to private visibility. To make a library public, use `# gazelle:default_visibility //visibility:public` in the package's `BUILD.bazel` file.
