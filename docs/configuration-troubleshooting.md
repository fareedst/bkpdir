# Configuration Troubleshooting Guide

<!-- 🔻 CFG-006: Documentation - 📝 Troubleshooting guide -->

This guide helps you diagnose and resolve common configuration issues with the enhanced configuration inspection system (CFG-006). Use the powerful configuration visibility features to debug problems quickly.

## 🔍 Quick Diagnosis Commands

Before diving into specific issues, use these commands for quick diagnosis:

```bash
# Get complete configuration overview
bkpdir config --sources --format tree > config-debug.txt

# Show only problematic values  
bkpdir config --overrides-only --sources

# Check performance
time bkpdir config >/dev/null
```

## 🛠️ Configuration Value Issues

### Issue 1: Configuration Value Not Being Applied

**Symptoms:**
- Setting a configuration value but it's not taking effect
- Expected value doesn't match actual value

**Diagnosis:**
```bash
# Check source resolution for specific field
bkpdir config [field_name] --sources --format tree
```

**Example:**
```bash
$ bkpdir config archive_dir_path --sources --format tree
archive_dir_path:
├── Environment: BKPDIR_ARCHIVE_DIR_PATH → (not set)
├── Local Config: ./.bkpdir.yml → (file not found)
├── Inherited: ~/.bkpdir.yml → ~/Archives (USED)
└── Default: ./archives (fallback)
│
Final Value: ~/Archives (from ~/.bkpdir.yml)
```

**Common Causes & Solutions:**

1. **Higher priority source overriding value**
   ```bash
   # Check if environment variable is set
   env | grep BKPDIR_ARCHIVE_DIR_PATH
   
   # Solution: Unset environment variable or update it
   unset BKPDIR_ARCHIVE_DIR_PATH
   ```

2. **Typo in configuration file path or field name**
   ```bash
   # Verify config file exists and is readable
   ls -la ~/.bkpdir.yml ./.bkpdir.yml
   
   # Check field name spelling
   bkpdir config --format json | jq '.configuration[] | .name' | grep -i archive
   ```

3. **Configuration file syntax error**
   ```bash
   # Test YAML syntax
   python3 -c "import yaml; yaml.safe_load(open('~/.bkpdir.yml'))"
   ```

### Issue 2: Inheritance Chain Problems

**Symptoms:**
- Values from parent config not being inherited
- Merge strategies not working as expected

**Diagnosis:**
```bash
# View complete inheritance chain
bkpdir config --sources --format tree | grep -A 10 -B 2 "Inherited"

# Check specific field inheritance
bkpdir config exclude_patterns --sources --format tree
```

**Example:**
```bash
$ bkpdir config exclude_patterns --sources --format tree
exclude_patterns:
├── Environment: BKPDIR_EXCLUDE_PATTERNS → (not set)
├── Local Config: ./.bkpdir.yml → [node_modules/] (+ merge strategy)
├── Inherited: ~/.bkpdir.yml → [*.tmp, *.log] (base values)
└── Default: [] (fallback)
│
Final Value: [*.tmp, *.log, node_modules/] (merged result)
```

**Common Causes & Solutions:**

1. **Missing inherit declaration**
   ```yaml
   # Add to local config file (./.bkpdir.yml)
   inherit: "~/.bkpdir.yml"
   archive_dir_path: "./local-archives"
   ```

2. **Incorrect merge strategy syntax**
   ```yaml
   # Correct merge strategy syntax
   +exclude_patterns:  # Append to inherited values
     - "build/"
     - "dist/"
   ```

3. **Circular inheritance**
   ```bash
   # Check for circular dependencies
   bkpdir config --sources | grep -i "circular"
   ```

### Issue 3: Environment Variable Overrides Not Working

**Symptoms:**
- Environment variables not overriding config file values
- Inconsistent environment variable behavior

**Diagnosis:**
```bash
# Check current environment variables
env | grep BKPDIR_ | sort

# Verify field name to environment variable mapping
bkpdir config --sources | grep Environment
```

**Common Causes & Solutions:**

1. **Incorrect environment variable name**
   ```bash
   # Check correct environment variable name
   bkpdir config [field_name] --sources
   
   # Example: archive_dir_path uses BKPDIR_ARCHIVE_DIR_PATH
   export BKPDIR_ARCHIVE_DIR_PATH="/custom/path"
   ```

2. **Environment variable type mismatch**
   ```bash
   # Boolean values need true/false (not 1/0)
   export BKPDIR_INCLUDE_GIT_INFO=true  # Not BKPDIR_INCLUDE_GIT_INFO=1
   
   # Array values need proper format
   export BKPDIR_EXCLUDE_PATTERNS="[*.tmp, *.log]"
   ```

3. **Environment variable not being passed to command**
   ```bash
   # Verify environment variable is in scope
   echo $BKPDIR_ARCHIVE_DIR_PATH
   
   # Use env to set temporarily
   env BKPDIR_INCLUDE_GIT_INFO=false bkpdir config include_git_info
   ```

## ⚡ Performance Issues

### Issue 4: Slow Config Command Response

**Symptoms:**
- Configuration inspection takes longer than expected (>1 second)
- Noticeable delay when running config command

**Diagnosis:**
```bash
# Measure config command performance
time bkpdir config --format table >/dev/null

# Check specific operations
time bkpdir config --overrides-only >/dev/null
time bkpdir config --sources >/dev/null
```

**Expected Performance:**
- First run: <100ms (with reflection)
- Cached runs: <50ms (with reflection caching)
- Filtered runs: <25ms (lazy evaluation)

**Common Causes & Solutions:**

1. **Reflection cache not working**
   ```bash
   # Clear and rebuild cache
   rm -rf ~/.cache/bkpdir/reflection-cache 2>/dev/null
   
   # Run twice and compare times
   time bkpdir config >/dev/null  # Should be slower (building cache)
   time bkpdir config >/dev/null  # Should be faster (using cache)
   ```

2. **Large configuration files**
   ```bash
   # Check config file sizes
   du -h ~/.bkpdir.yml ./.bkpdir.yml
   
   # Use filtering to reduce processing
   bkpdir config --filter "archive" --overrides-only
   ```

3. **Complex inheritance chains**
   ```bash
   # Identify inheritance depth
   bkpdir config --sources --format tree | grep "Inherited" | wc -l
   
   # Simplify inheritance if needed
   ```

### Issue 5: Memory Usage During Inspection

**Symptoms:**
- High memory usage when running config command
- System slowdown during configuration inspection

**Diagnosis:**
```bash
# Monitor memory usage
/usr/bin/time -v bkpdir config 2>&1 | grep "Maximum resident set size"

# Check for memory leaks
valgrind --tool=memcheck --leak-check=full bkpdir config 2>/dev/null
```

**Common Causes & Solutions:**

1. **Too many reflection operations**
   ```bash
   # Use targeted inspection instead of full scan
   bkpdir config --filter "specific_pattern"
   
   # Use lazy evaluation
   bkpdir config --overrides-only
   ```

2. **Large configuration structures**
   ```bash
   # Check configuration complexity
   bkpdir config --format json | jq '.configuration | length'
   
   # Use incremental access for specific fields
   bkpdir config specific_field_name
   ```

### Issue 6: Reflection Overhead

**Symptoms:**
- First-time config command is significantly slower
- Inconsistent performance across different systems

**Diagnosis:**
```bash
# Benchmark reflection operations
go test -bench=BenchmarkConfigReflectionOperations ./...

# Profile reflection performance
go test -cpuprofile=cpu.prof -bench=BenchmarkConfigReflectionOperations ./...
go tool pprof cpu.prof
```

**Solutions:**
```bash
# Warm up reflection cache
bkpdir config >/dev/null

# Use pre-warmed cache in scripts
if [ ! -f ~/.cache/bkpdir/reflection-cache ]; then
    bkpdir config >/dev/null 2>&1  # Warm up cache
fi
```

## 🎨 Display Issues

### Issue 7: Missing Configuration Fields

**Symptoms:**
- Expected configuration fields not appearing
- Incomplete field discovery

**Diagnosis:**
```bash
# Check total field count
bkpdir config --format json | jq '.configuration | length'

# Compare with expected count (should be 100+)
echo "Expected: 100+ fields, Got: $(bkpdir config --format json | jq '.configuration | length') fields"

# Check for reflection errors
bkpdir config 2>&1 | grep -i error
```

**Common Causes & Solutions:**

1. **Reflection not finding embedded fields**
   ```bash
   # Check for embedded struct fields
   bkpdir config --format json | jq '.configuration[] | select(.category == "")' 
   
   # Look for fields without categories (might indicate reflection issues)
   ```

2. **Type handling issues**
   ```bash
   # Check for unsupported types
   bkpdir config --format json | jq '.configuration[] | select(.type == "unknown")'
   ```

3. **Field filtering bug**
   ```bash
   # Test without filters
   bkpdir config --all
   
   # Test with different filters
   bkpdir config --filter ""
   ```

### Issue 8: Incorrect Source Attribution

**Symptoms:**
- Wrong source shown for configuration values
- Source attribution inconsistencies

**Diagnosis:**
```bash
# Verify source attribution accuracy
bkpdir config --sources --format json | jq '.configuration[] | {name, value, source}'

# Cross-check specific field sources
bkpdir config archive_dir_path --sources --format tree
```

**Common Causes & Solutions:**

1. **Source tracking cache issue**
   ```bash
   # Clear source tracking cache
   rm -rf ~/.cache/bkpdir/source-cache 2>/dev/null
   
   # Rebuild source attribution
   bkpdir config --sources >/dev/null
   ```

2. **Inheritance resolution bug**
   ```bash
   # Check inheritance chain consistency
   bkpdir config --sources --format tree | grep -A 5 -B 5 "inconsistent\|error"
   ```

3. **Environment variable detection issue**
   ```bash
   # Manually verify environment source
   env | grep BKPDIR_ARCHIVE_DIR_PATH
   bkpdir config archive_dir_path --sources
   ```

### Issue 9: Formatting Problems

**Symptoms:**
- Broken table formatting
- JSON output errors
- Tree format display issues

**Diagnosis:**
```bash
# Test all formats
bkpdir config --format table | head -10
bkpdir config --format tree | head -10  
bkpdir config --format json | jq . >/dev/null && echo "JSON valid" || echo "JSON invalid"
```

**Common Causes & Solutions:**

1. **Terminal width issues**
   ```bash
   # Check terminal width
   tput cols
   
   # Use JSON format for narrow terminals
   bkpdir config --format json
   ```

2. **Special characters in values**
   ```bash
   # Check for problematic characters
   bkpdir config --format json | jq '.configuration[] | select(.value | contains("\\") or contains("\""))'
   
   # Use JSON format for complex values
   ```

3. **Unicode/encoding issues**
   ```bash
   # Check locale settings
   locale
   
   # Set UTF-8 encoding
   export LC_ALL=en_US.UTF-8
   ```

## 🔗 Integration Issues

### Issue 10: CFG-005 Inheritance Not Working

**Symptoms:**
- Inheritance features not available
- Error messages about missing inheritance support

**Diagnosis:**
```bash
# Check if CFG-005 is implemented
bkpdir config --help | grep -i inherit

# Test inheritance functionality
echo 'inherit: "~/.bkpdir.yml"' > test-inherit.yml
echo 'test_value: "inherited"' >> test-inherit.yml
```

**Solutions:**
```bash
# Verify CFG-005 implementation
grep -r "CFG-005" . | grep -v Binary

# Check for inheritance support in config loading
bkpdir config --sources | grep -i inherit
```

### Issue 11: EXTRACT-001 pkg/config Integration

**Symptoms:**
- Configuration loading errors
- Missing features from extracted package

**Diagnosis:**
```bash
# Check package integration
go list -m all | grep pkg/config

# Test extracted package functionality
go test ./pkg/config/... -v
```

**Solutions:**
```bash
# Verify package import paths
grep -r "pkg/config" *.go

# Check for version conflicts
go mod tidy
go mod verify
```

### Issue 12: CLI Framework Integration

**Symptoms:**
- Command-line flags not working
- Help text outdated or incorrect

**Diagnosis:**
```bash
# Check CLI framework integration
bkpdir config --help

# Test all command-line flags
bkpdir config --all --overrides-only --sources --format json --filter test
```

**Solutions:**
```bash
# Verify cobra integration
grep -r "cobra" main.go | grep config

# Check flag definitions
bkpdir config --help | grep -A 20 "Flags:"
```

## 🔧 Common Error Messages

### Error: "reflection: call of reflect.Value.Field on zero Value"

**Cause:** Reflection operating on nil or uninitialized struct

**Solution:**
```bash
# Check configuration loading
bkpdir config 2>&1 | grep -B 5 -A 5 "reflect"

# Verify config file syntax
python3 -c "import yaml; print(yaml.safe_load(open('~/.bkpdir.yml')))"
```

### Error: "field discovery failed"

**Cause:** Reflection unable to enumerate struct fields

**Solution:**
```bash
# Clear reflection cache and retry
rm -rf ~/.cache/bkpdir/reflection-cache 2>/dev/null
bkpdir config

# Check for struct definition issues
go vet ./...
```

### Error: "source attribution failed"

**Cause:** Source tracking unable to determine value origin

**Solution:**
```bash
# Reset source tracking
rm -rf ~/.cache/bkpdir/source-cache 2>/dev/null

# Check inheritance chain integrity
bkpdir config --sources --format tree 2>&1 | grep -i error
```

### Error: "performance optimization failed"

**Cause:** Caching or lazy evaluation error

**Solution:**
```bash
# Disable optimization temporarily
BKPDIR_DISABLE_CACHE=true bkpdir config

# Clear all caches
rm -rf ~/.cache/bkpdir/ 2>/dev/null
```

## 📊 Diagnostic Scripts

### Complete Configuration Health Check

```bash
#!/bin/bash
# config-health-check.sh

echo "=== Configuration Health Check ==="
echo "Generated: $(date)"
echo

echo "=== Basic Functionality ==="
if bkpdir config >/dev/null 2>&1; then
    echo "✅ Basic config command works"
else
    echo "❌ Basic config command failed"
    exit 1
fi

echo "=== Field Discovery ==="
FIELD_COUNT=$(bkpdir config --format json | jq '.configuration | length' 2>/dev/null)
if [ "$FIELD_COUNT" -gt 50 ]; then
    echo "✅ Field discovery working ($FIELD_COUNT fields found)"
else
    echo "⚠️  Field discovery may have issues ($FIELD_COUNT fields found, expected 100+)"
fi

echo "=== Source Attribution ==="
if bkpdir config --sources >/dev/null 2>&1; then
    echo "✅ Source attribution working"
else
    echo "❌ Source attribution failed"
fi

echo "=== Performance ==="
PERF_TIME=$(time ( bkpdir config >/dev/null ) 2>&1 | grep real | awk '{print $2}')
echo "⏱️  Performance: $PERF_TIME"

echo "=== Output Formats ==="
for format in table tree json; do
    if bkpdir config --format $format >/dev/null 2>&1; then
        echo "✅ Format '$format' working"
    else
        echo "❌ Format '$format' failed"
    fi
done

echo "=== Inheritance Support ==="
if bkpdir config --sources | grep -q "Inherit"; then
    echo "✅ Inheritance support detected"
else
    echo "⚠️  Inheritance support not detected"
fi

echo
echo "=== Health Check Complete ==="
```

### Performance Benchmark

```bash
#!/bin/bash
# config-performance-benchmark.sh

echo "=== Configuration Performance Benchmark ==="

echo "Testing first run (cold cache)..."
rm -rf ~/.cache/bkpdir/ 2>/dev/null
time bkpdir config >/dev/null

echo "Testing second run (warm cache)..."
time bkpdir config >/dev/null

echo "Testing filtered run..."
time bkpdir config --filter "archive" >/dev/null

echo "Testing overrides only..."
time bkpdir config --overrides-only >/dev/null

echo "Testing JSON format..."
time bkpdir config --format json >/dev/null

echo "Testing with sources..."
time bkpdir config --sources >/dev/null
```

## 📝 Prevention Best Practices

### 1. Regular Configuration Audits

```bash
# Weekly configuration review
bkpdir config --overrides-only --sources > weekly-config-$(date +%Y%m%d).txt
```

### 2. Environment Variable Documentation

```bash
# Document current environment variables
env | grep BKPDIR_ > env-vars-$(date +%Y%m%d).txt
```

### 3. Configuration Backup

```bash
# Backup configuration state
bkpdir config --format json > config-backup-$(date +%Y%m%d-%H%M%S).json
```

### 4. Performance Monitoring

```bash
# Monitor performance trends
echo "$(date): $(time bkpdir config >/dev/null 2>&1 | grep real)" >> config-performance.log
```

## 🆘 Getting Help

If you encounter issues not covered in this guide:

1. **Collect diagnostic information:**
   ```bash
   bkpdir config --sources --format tree > config-debug.txt
   bkpdir version >> config-debug.txt
   go version >> config-debug.txt
   ```

2. **Check for known issues:**
   ```bash
   grep -r "CFG-006" docs/context/ | grep -i "issue\|problem\|bug"
   ```

3. **Test with minimal configuration:**
   ```bash
   mv ~/.bkpdir.yml ~/.bkpdir.yml.backup
   bkpdir config  # Test with defaults only
   ```

4. **Enable debug mode:**
   ```bash
   BKPDIR_DEBUG=true bkpdir config
   ```

Remember: The CFG-006 system is designed to be self-diagnosing. Use the powerful configuration inspection capabilities to understand and resolve issues quickly.

---

*For comprehensive feature documentation, see the [Configuration Inspection Guide](configuration-inspection-guide.md).*
*For practical usage patterns, see the [Configuration Examples](configuration-examples.md).* 