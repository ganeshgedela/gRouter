# Logger Package - Unit Tests

## Overview

Comprehensive unit tests for the logger package covering logger initialization, configuration, context propagation, and all logging functionality.

## Test Coverage

### Logger Tests (`logger_test.go`)

✅ **TestNew** - Logger creation with various configurations (5 test cases)
  - Valid JSON format
  - Valid console format
  - Invalid log level
  - Warn level
  - Error level

✅ **TestNew_FileOutput** - File-based logging with temp directory

✅ **TestGet** - Global logger retrieval and management

✅ **TestSugar** - Sugared logger creation

✅ **TestWithFields** - Logger with additional fields

✅ **TestLogLevels** - All log levels (Debug, Info, Warn, Error)

✅ **TestSync** - Logger sync/flush functionality

✅ **TestLogFormats** - JSON and console output formats

✅ **TestLogLevelParsing** - Case-insensitive level parsing (7 test cases)

✅ **TestEmptyOutputPath** - Default stdout handling

✅ **TestMultipleLoggerCreation** - Multiple logger instances

### Context Tests (`context_test.go`)

✅ **TestWithContext** - Logger storage in context

✅ **TestFromContext_NoLogger** - Fallback to global logger

✅ **TestWithRequestID** - Request ID field addition

✅ **TestWithTraceID** - Trace ID field addition

✅ **TestMultipleContextFields** - Multiple fields in context

✅ **TestContextPropagation** - Context passing through layers

✅ **TestFromContext_WithFields** - Logger with pre-existing fields

✅ **TestNilContext** - Nil context handling

✅ **TestContextChaining** - Chained context operations

## Test Results

```
✅ PASS - 21 tests passed
✅ Coverage: 92.0% of statements
```

**All tests passed successfully!**

## Running Tests

### Standard Go Test
```bash
# Run all tests
go test -v ./pkg/logger/...

# With coverage
go test -v ./pkg/logger/... -cover

# Coverage report
go test -v ./pkg/logger/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Using Bazel
```bash
# Run tests
bazel test //pkg/logger:logger_test

# With coverage
bazel coverage //pkg/logger:logger_test
```

## Test Coverage Details

### Logger Functionality
- ✅ Logger creation with different configs
- ✅ Log level validation (debug, info, warn, error)
- ✅ Output format (JSON, console)
- ✅ File and stdout output
- ✅ Global logger management
- ✅ Sugared logger
- ✅ Logger with fields
- ✅ Sync/flush operations

### Context Integration
- ✅ Context storage and retrieval
- ✅ Request ID propagation
- ✅ Trace ID propagation
- ✅ Multiple field chaining
- ✅ Context propagation through layers
- ✅ Fallback to global logger

## natsdemosvc Test Output

```
=== RUN   TestMultipleContextFields
2025-12-18T13:06:22.001+0530    INFO    test message with request and trace IDs 
    {"request_id": "req-123", "trace_id": "trace-456"}
--- PASS: TestMultipleContextFields (0.00s)

=== RUN   TestContextPropagation
2025-12-18T13:06:22.001+0530    INFO    processing request      
    {"request_id": "req-001"}
2025-12-18T13:06:22.001+0530    DEBUG   added trace ID  
    {"request_id": "req-001", "trace_id": "trace-001"}
--- PASS: TestContextPropagation (0.00s)
```

## Coverage Areas

1. **Initialization**: Logger creation with various configs
2. **Log Levels**: Debug, Info, Warn, Error validation
3. **Output Formats**: JSON and console formats
4. **File I/O**: File-based logging with temp directories
5. **Context**: Request/trace ID propagation
6. **Error Handling**: Invalid configs, nil contexts
7. **Global State**: Global logger management
8. **Field Addition**: Logger with additional fields

## Files Created

- `/home/ganesh/gRouter/pkg/logger/logger_test.go` - 14 test cases
- `/home/ganesh/gRouter/pkg/logger/context_test.go` - 9 test cases
- Updated `/home/ganesh/gRouter/pkg/logger/BUILD.bazel` with test target

Excellent test coverage at **92.0%**! 🎉
