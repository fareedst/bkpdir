# Delayed Output Management

The delayed output functionality allows you to collect output messages and control when they are displayed, rather than printing them immediately. This is useful for operations that need to validate success before showing output, or for batching output messages.

## Basic Usage

### Creating a Delayed Output Formatter

```go
// Create a config (normally loaded from file)
cfg := DefaultConfig()

// Create an output collector
collector := NewOutputCollector()

// Create a formatter with delayed output enabled
formatter := NewOutputFormatterWithCollector(cfg, collector)
```

### Collecting Messages

```go
// These calls will collect messages instead of printing immediately
formatter.PrintCreatedArchive("/path/to/archive.zip")
formatter.PrintConfigValue("archive_dir", "/archives", "default")
formatter.PrintVerificationSuccess("backup.zip")

// Messages are now stored in the collector, not printed
```

### Displaying Collected Messages

```go
// Display all collected messages at once
collector.FlushAll()

// Or display only stdout messages
collector.FlushStdout()

// Or display only stderr messages  
collector.FlushStderr()

// Or discard all messages without displaying
collector.Clear()
```

## Advanced Usage

### Conditional Output

```go
func performOperation() error {
    cfg := DefaultConfig()
    collector := NewOutputCollector()
    formatter := NewOutputFormatterWithCollector(cfg, collector)
    
    // Perform operation and collect output
    formatter.PrintDryRunFilesHeader()
    formatter.PrintDryRunFileEntry("file1.txt")
    formatter.PrintCreatedArchive("/archives/result.zip")
    
    // Check if operation was successful
    if operationSuccessful() {
        // Show output only if successful
        collector.FlushAll()
        return nil
    } else {
        // Discard output if failed
        collector.Clear()
        return errors.New("operation failed")
    }
}
```

### Inspecting Messages

```go
// Get all collected messages for inspection
messages := collector.GetMessages()

for _, msg := range messages {
    fmt.Printf("Destination: %s, Type: %s, Content: %s\n", 
        msg.Destination, msg.Type, msg.Content)
}
```

### Switching Between Modes

```go
formatter := NewOutputFormatter(cfg)

// Enable delayed output
collector := NewOutputCollector()
formatter.SetCollector(collector)

// Check if in delayed mode
if formatter.IsDelayedMode() {
    fmt.Println("Messages will be collected")
}

// Disable delayed output (return to immediate printing)
formatter.SetCollector(nil)
```

## Message Types

The OutputCollector categorizes messages by type:

- `"info"` - General information messages
- `"error"` - Error messages  
- `"warning"` - Warning messages
- `"config"` - Configuration-related messages
- `"dry-run"` - Dry-run operation messages

## Backward Compatibility

The delayed output functionality is fully backward compatible. Existing code that uses `NewOutputFormatter()` will continue to work unchanged, printing output immediately as before.

Only code that explicitly uses `NewOutputFormatterWithCollector()` or `SetCollector()` will use delayed output. 