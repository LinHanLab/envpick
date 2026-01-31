# E2E Testing for envpick

This directory contains end-to-end tests for envpick that verify all features work correctly from a user's perspective.

## Architecture

The E2E tests use an **in-process testing approach** where we test the application logic directly without compiling a binary. This provides:

- Fast execution (no compilation overhead)
- Easy mocking of external dependencies (fzf, browser, editor)
- Full control over the test environment
- Ability to inspect internal state
- Standard Go testing patterns

## Test Structure

```
test/e2e/
├── e2e_test.go      # Main E2E test scenarios
├── fixtures.go      # Test data and config templates
├── mocks.go         # Mock implementations for external dependencies
├── helpers.go       # Test helper functions
└── README.md        # This file
```

## Test Coverage

The E2E tests cover all major features:

### 1. Environment Variable Management
- `TestBasicEnvironmentSelection` - Load and export environment variables
- `TestDirectSelection` - Select config without persistence
- `TestSetCurrentConfig` - Change current configuration with persistence

### 2. Namespace Support
- `TestNamespaceIsolation` - Verify namespaces maintain separate state
- `TestMultipleNamespaceOperations` - Operations across multiple namespaces
- `TestNamespaceFiltering` - Namespace filtering logic

### 3. Configuration Management
- `TestConfigFileCreation` - Config directory and file creation
- `TestMetadataFiltering` - Metadata variables (starting with `_`) are not exported

### 4. State Management
- `TestStateMigration` - Migration from old to new state format
- `TestGetOptions` - Option generation for selection with active marking

### 5. Error Handling
- `TestErrorHandling` - Various error conditions (missing config, invalid TOML, non-existent config)

## Running Tests

```bash
# Run all E2E tests
go test ./test/e2e/...

# Run with verbose output
go test ./test/e2e/... -v

# Run specific test
go test ./test/e2e/... -run TestBasicEnvironmentSelection

# Run all tests including E2E
go test ./...

# Run quality checks (tests + formatting + linting)
make quality
```

## Test Isolation

Each test runs in complete isolation:

1. **Temporary HOME directory**: Each test gets its own temp directory via `t.TempDir()`
2. **Isolated config/state files**: Config and state files are created in the test's temp directory
3. **No interference**: Tests can run in parallel without affecting each other
4. **Automatic cleanup**: Temp directories are automatically cleaned up after tests

## Test Helpers

### TestEnv
The `TestEnv` struct provides an isolated test environment:

```go
env := NewTestEnv(t)
defer env.SetHome()()  // Set HOME to test dir and restore on cleanup

env.WriteConfig(BasicConfig)  // Write test config
env.WriteState(stateContent)  // Write test state
state := env.ReadState()      // Read current state
```

### Assertions
Helper functions for common assertions:

```go
AssertContains(t, output, "expected string")
AssertNotContains(t, output, "unexpected string")
env.AssertStateContains("expected in state")
```

## Test Fixtures

The `fixtures.go` file contains test configuration templates:

- `BasicConfig` - Simple dev/prod configuration
- `NamespaceConfig` - Configuration with multiple namespaces
- `MetadataConfig` - Configuration with metadata variables
- `MultiNamespaceConfig` - Configuration with 3 namespaces
- `InvalidConfig` - Invalid TOML for error testing
- `LegacyState` - Old state format for migration testing

## Mocking External Dependencies

The `mocks.go` file provides mock implementations for external dependencies:

- `MockSelector` - Simulates fzf selection
- `MockURLOpener` - Simulates browser opening
- `MockFileEditor` - Simulates editor opening

These mocks are currently defined but not yet integrated into the codebase. Future work will refactor the code to use dependency injection for these components.

## Pre-commit Hook

A pre-commit hook is installed at `.git/hooks/pre-commit` that runs:

1. All tests (including E2E tests)
2. Code formatting checks
3. Linting

This ensures all features work before code is committed.

## CI Integration

The E2E tests are automatically run in CI via the `make quality` target, which is typically called in the CI pipeline.

## Adding New Tests

To add a new E2E test:

1. Identify the user workflow to test
2. Create a new test function in `e2e_test.go`
3. Use `NewTestEnv(t)` to create an isolated environment
4. Write test config/state using fixtures or custom data
5. Execute the feature being tested
6. Verify outputs and side effects
7. Follow the existing test patterns for consistency

Example:

```go
func TestNewFeature(t *testing.T) {
    env := NewTestEnv(t)
    defer env.SetHome()()

    // Setup
    env.WriteConfig(BasicConfig)
    env.WriteState(NewStateDefault)

    // Action
    engine, err := core.NewEngine()
    if err != nil {
        t.Fatalf("Failed to create engine: %v", err)
    }

    // Verify
    result := engine.SomeMethod()
    if result != expected {
        t.Errorf("Expected %v, got %v", expected, result)
    }
}
```

## Design Principles

1. **Test behavior, not implementation** - Focus on what the feature does, not how it does it
2. **Test real workflows** - Each test should represent an actual user workflow
3. **Fast and reliable** - Tests should run quickly and never be flaky
4. **Easy to understand** - Tests should be readable and self-documenting
5. **Isolated** - Tests should not depend on each other or external state

## Future Enhancements

Potential improvements for the E2E testing architecture:

1. **Dependency Injection** - Refactor code to accept interfaces for fzf, browser, and editor
2. **Binary Testing** - Add supplementary tests that compile and test the actual binary
3. **Integration Tests** - Add tests that verify integration with actual fzf (when available)
4. **Performance Tests** - Add tests to ensure operations complete within acceptable time
5. **Platform-specific Tests** - Add tests for platform-specific behavior (macOS, Linux, Windows)
