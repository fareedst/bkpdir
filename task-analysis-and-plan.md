# Task Analysis and Implementation Plan

## 🎯 Task Overview

**TASK**: Enhance backup commands to output formatted file information
- All commands must output to stdout a single line about any files generated or existing to satisfy a backup request
- That output line must be formatted using a format string
- Provide `stat`-like info for the file as named replacements for the format string
- Currently: "inc" command displays the file name, "full" command displays nothing

## 🔍 Current State Analysis

### Command Output Analysis

**"inc" Command (Incremental Archive)**:
- Location: `archive.go:777` in `createAndVerifyIncrementalArchive()`
- Current behavior: Calls `formatter.PrintIncrementalCreated(cfg.Path)`
- Format string used: `format_incremental_created: "Created incremental archive: %s\n"`
- **Status**: ✅ Already outputs to stdout, but only shows path

**"full" Command (Full Archive)**:
- Location: `archive.go:416` in `CreateFullArchiveWithContext()`
- Current behavior: No output for successful creation (only dry-run output)
- **Status**: ❌ Missing output for successful archive creation

### Current Format String Infrastructure

The system already has comprehensive format string support:
- Configuration in `config.go` with format strings like `FormatCreatedArchive`
- OutputFormatter in `formatter.go` with methods like `FormatCreatedArchive()` and `PrintCreatedArchive()`
- Template-based formatting alongside printf-style formatting

### Gap Analysis

1. **Full command missing output**: No success message when archive is created
2. **Limited file information**: Only path is included, no stat-like info
3. **Format string needs enhancement**: Need named replacements for stat info

## 📋 Implementation Plan

### Phase 1: Analysis and Setup (CRITICAL)

#### 1.1 CRITICAL VALIDATION (MANDATORY - Execute First)
- [x] Read AI Assistant Context Documentation Index 
- [x] Follow PHASE 1 CRITICAL VALIDATION requirements
- [x] Identify task as MODIFICATION Protocol (modifying existing inc/full commands)
- [ ] Find or create Feature ID in feature-tracking.md
- [ ] Check for conflicts with immutable.md

#### 1.2 Technical Analysis
- [ ] Examine stat-like information available for files in Go
- [ ] Design format string structure with named replacements  
- [ ] Identify insertion points for new output in full command
- [ ] Plan consistent approach for both inc and full commands

### Phase 2: Design Enhanced Format Strings (HIGH)

#### 2.1 Stat-like Information Design
Design named replacements for format strings:
- `{path}` - Full file path
- `{name}` - File name only
- `{size}` - File size in bytes
- `{size_human}` - Human-readable size (1.2MB, 455KB, etc.)
- `{mtime}` - Modification time
- `{mtime_unix}` - Modification time as unix timestamp
- `{mode}` - File permissions/mode
- `{type}` - File type (regular, directory, symlink)

#### 2.2 Format String Configuration
Add new configuration options:
- `format_archive_created_detailed` - Enhanced format with stat info
- `format_incremental_created_detailed` - Enhanced format for incremental
- Maintain backward compatibility with existing format strings

### Phase 3: Implementation (HIGH)

#### 3.1 Add File Stat Gathering Function
Create utility function to gather stat-like info:
```go
type FileStatInfo struct {
    Path       string
    Name       string  
    Size       int64
    SizeHuman  string
    MTime      time.Time
    MTimeUnix  int64
    Mode       os.FileMode
    Type       string
}

func GatherFileStatInfo(path string) (*FileStatInfo, error)
```

#### 3.2 Enhance OutputFormatter
- Add methods for detailed formatting with stat info
- Update FormatCreatedArchive and FormatIncrementalCreated to support stat info
- Add template-based formatting for named replacements

#### 3.3 Update Archive Creation Functions
- Modify `CreateFullArchiveWithContext` to output success message
- Modify `createAndVerifyIncrementalArchive` to use enhanced formatting
- Ensure consistent behavior between inc and full commands

### Phase 4: Integration and Testing (MEDIUM)

#### 4.1 Configuration Integration
- Add new format strings to DefaultConfig()
- Update configuration template generation
- Add configuration documentation

#### 4.2 Testing
- Add unit tests for new stat gathering function
- Test format string template processing
- Test both inc and full command output
- Verify backward compatibility

### Phase 5: Documentation Updates (MEDIUM)

#### 5.1 Context Documentation (MANDATORY)
- Update feature-tracking.md with task completion
- Update specification.md with new output behavior
- Update architecture.md with format string enhancements
- Update requirements.md with stat info requirements

#### 5.2 User Documentation
- Update configuration examples
- Document available named replacements
- Add examples of customizing output format

## 🔧 Technical Implementation Details

### File Stat Information Gathering

```go
import (
    "os"
    "path/filepath"
    "time"
)

func GatherFileStatInfo(path string) (*FileStatInfo, error) {
    info, err := os.Stat(path)
    if err != nil {
        return nil, err
    }
    
    return &FileStatInfo{
        Path:      path,
        Name:      filepath.Base(path),
        Size:      info.Size(),
        SizeHuman: formatHumanSize(info.Size()),
        MTime:     info.ModTime(),
        MTimeUnix: info.ModTime().Unix(),
        Mode:      info.Mode(),
        Type:      getFileType(info),
    }, nil
}
```

### Format String Enhancement

Current format string:
```yaml
format_created_archive: "Created archive: %s\n"
```

Enhanced format string with named replacements:
```yaml
format_created_archive: "Created archive: {path} (size: {size_human}, modified: {mtime})\n"
```

### Integration Points

1. **CreateFullArchiveWithContext** - Add success output after `createAndVerifyArchive`
2. **createAndVerifyIncrementalArchive** - Enhance existing `PrintIncrementalCreated` call
3. **Configuration** - Add new format string options
4. **OutputFormatter** - Add stat-aware formatting methods

## 🎯 Success Criteria

1. **Full command outputs**: Single line to stdout when archive created successfully
2. **Inc command enhanced**: Output includes stat-like information, not just path
3. **Format strings**: Support named replacements for file information
4. **Backward compatibility**: Existing format strings continue to work
5. **Consistent behavior**: Both inc and full commands behave similarly
6. **Configuration**: Users can customize output format via configuration

## 🔄 Risk Assessment

**Low Risk**:
- Adding new functionality alongside existing code
- Backward compatible changes
- Well-defined scope

**Mitigation Strategies**:
- Maintain existing format string behavior as fallback
- Add comprehensive unit tests
- Follow existing architectural patterns

## 📝 Implementation Tokens

All code changes will use the format:
```go
// OUTPUT-FORMAT-001: Enhanced command output with stat information
```

## 🏁 Completion Criteria

- [ ] Both inc and full commands output single line to stdout with file info
- [ ] Format strings support named replacements for stat-like info
- [ ] All tests pass
- [ ] Documentation updated per MODIFICATION Protocol
- [ ] Feature marked complete in feature-tracking.md 