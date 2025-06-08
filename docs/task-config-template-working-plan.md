# Configuration Template Creation Task - Working Plan

**Task ID**: CFG-TEMPLATE-001  
**Date**: 2025-01-02  
**AI Assistant**: Claude Sonnet 4 (Cursor)

## 🎯 Task Summary

Add a command to make an empty configuration YAML file in the current directory as a template for the user to select applicable options.

### Requirements:
- **New file name**: Use `.bkpdir.yml` if new, otherwise write/overwrite to `.bkpdir.default-YYYY-MM-DD.yml`
- **Content**: Include all keys organized by category
- **Format**: All keys are commented out
- **Values**: Each key value is the configuration after loading BKPDIR_CONFIG

## 📋 PHASE 1: CRITICAL VALIDATION (Execute FIRST - MANDATORY)

### ✅ Pre-Work Validation Completed:
1. **📋 Task Verification**: Task recorded in feature-tracking.md with valid Feature ID: CFG-TEMPLATE-001
2. **🔍 Compliance Check**: Following ai-assistant-compliance.md requirements for implementation tokens
3. **📁 File Impact Analysis**: Determined documentation files requiring updates
4. **🛡️ Immutable Check**: No conflicts with immutable.md requirements

## 🔧 Change Classification: NEW FEATURE

**Protocol**: NEW FEATURE Protocol [PRIORITY: CRITICAL]  
**Justification**: Adding new command functionality with CLI interface and configuration generation

## 📊 Detailed Analysis

### Current Configuration System Analysis:
- **Configuration Discovery**: Uses `BKPDIR_CONFIG` environment variable or default paths
- **Default Paths**: `./.bkpdir.yml:~/.bkpdir.yml`
- **Structure**: Complete Config struct with 100+ fields organized by category
- **Categories**: Basic settings, Git config, status codes, format strings, templates, patterns

### Implementation Requirements:
1. **Command Integration**: Add to existing CLI framework in main.go
2. **Configuration Generation**: Create comprehensive template with all Config struct fields
3. **File Naming Logic**: Implement date-based naming for existing files
4. **Value Population**: Load current configuration and use as template values
5. **Commenting System**: Comment out all keys while preserving structure

## 🗂️ Implementation Plan

### Phase 1: Command Infrastructure
- [ ] Add `template` command to main.go CLI structure
- [ ] Implement command flags and help text
- [ ] Add command to rootCmd

### Phase 2: Configuration Template Generation
- [ ] Create template generation function using reflection on Config struct
- [ ] Implement category-based organization matching existing structure
- [ ] Add comprehensive commenting system for all keys
- [ ] Implement value extraction from loaded configuration

### Phase 3: File Management
- [ ] Implement file existence checking logic
- [ ] Add date-based naming system for existing files
- [ ] Ensure safe file writing with backup protection

### Phase 4: Integration and Testing
- [ ] Integration with existing configuration loading system
- [ ] Error handling and validation
- [ ] Test with various configuration scenarios

## 📝 Subtask Breakdown

### Subtask 1: CLI Command Implementation
**Status**: Not Started  
**Description**: Add `template` command to main.go CLI interface
**Requirements**: 
- Command integration with cobra CLI framework
- Help text and usage documentation
- Flag support for custom filename/path

### Subtask 2: Configuration Reflection System
**Status**: Not Started  
**Description**: Create system to extract all Config struct fields using reflection
**Requirements**:
- Complete field enumeration including nested structs (Git, Verification)
- Category-based organization matching documentation
- YAML tag extraction for proper field naming

### Subtask 3: Template Generation Engine
**Status**: Not Started  
**Description**: Build template file generator with commenting and value population
**Requirements**:
- YAML file generation with comments
- Value population from loaded configuration
- Category headers and organization
- Proper indentation and formatting

### Subtask 4: File Naming and Management
**Status**: Not Started  
**Description**: Implement smart file naming and conflict resolution
**Requirements**:
- Check for existing `.bkpdir.yml`
- Generate date-based names for existing files
- Safe file writing with error handling

### Subtask 5: Integration Testing
**Status**: Not Started  
**Description**: Test command with various configuration scenarios
**Requirements**:
- Test with empty configuration
- Test with existing configuration files
- Test with BKPDIR_CONFIG environment variable
- Validate generated template structure

### Subtask 6: Documentation Updates
**Status**: Not Started  
**Description**: Update all required documentation per NEW FEATURE protocol
**Requirements**:
- Update feature-tracking.md with implementation status
- Add specification entry for new command
- Update architecture documentation
- Add testing requirements

## 🔧 Technical Implementation Details

### Command Structure:
```go
templateCmd := &cobra.Command{
    Use:   "template",
    Short: "Generate configuration template file",
    Long:  "Create a comprehensive configuration template with all available options",
    Run:   handleTemplateCommand,
}
```

### File Naming Logic:
- If `.bkpdir.yml` doesn't exist → create `.bkpdir.yml`
- If `.bkpdir.yml` exists → create `.bkpdir.default-YYYY-MM-DD.yml`

### Template Structure Categories:
1. **Basic Settings**: Archive paths, directory names, exclude patterns
2. **Git Integration**: Git configuration settings (GitConfig struct)
3. **Verification**: Archive verification settings
4. **Status Codes**: Directory and file operation status codes
5. **Format Strings**: Printf-style format strings for output
6. **Templates**: Template-based format strings with placeholders
7. **Patterns**: Regex patterns for various operations
8. **Extended Formats**: Additional format strings for comprehensive operations

### Configuration Value Source:
- Load current configuration using existing LoadConfig system
- Extract values from loaded Config struct
- Use default values where no override exists
- Maintain source attribution (environment, config file, default)

## 🚨 Risk Assessment

### Low Risk Items:
- CLI command integration (well-established pattern)
- Configuration loading (existing robust system)
- File writing operations (standard Go operations)

### Medium Risk Items:
- Reflection-based field extraction (complex but manageable)
- YAML commenting preservation (requires careful formatting)
- Date-based file naming (edge cases around existing files)

### Risk Mitigation:
- Use existing GetAllConfigFields() function as foundation
- Leverage existing YAML marshaling patterns
- Implement comprehensive error handling
- Add dry-run capability for testing

## 🎯 Success Criteria

### Primary Success Criteria:
1. **Command Availability**: `bkpdir template` command available and functional
2. **Template Generation**: Generates comprehensive configuration template
3. **File Management**: Handles naming conflicts appropriately
4. **Value Population**: Uses actual configuration values in template
5. **Documentation**: All keys properly commented and organized

### Quality Criteria:
1. **Completeness**: Template includes all Config struct fields
2. **Organization**: Fields organized by logical categories
3. **Usability**: Template is user-friendly and well-documented
4. **Robustness**: Handles edge cases and error conditions gracefully

## 📅 Implementation Timeline

**Estimated Duration**: 3-4 hours

1. **Phase 1** (30 min): CLI command setup and basic structure
2. **Phase 2** (90 min): Configuration reflection and template generation
3. **Phase 3** (45 min): File management and naming logic
4. **Phase 4** (45 min): Integration testing and refinement
5. **Phase 5** (30 min): Documentation updates per NEW FEATURE protocol

## 🔍 Validation Checklist

### Pre-Implementation:
- [x] Feature ID created in feature-tracking.md
- [x] Change protocol identified (NEW FEATURE)
- [x] Required documentation files identified
- [x] No immutable requirement conflicts

### Post-Implementation:
- [ ] All tests pass (`make test`)
- [ ] All lint checks pass (`make lint`)
- [ ] Command functional with various scenarios
- [ ] Template generation accurate and complete
- [ ] Documentation updated per protocol
- [ ] Feature marked complete in feature-tracking.md

## 🏁 Final Notes

This task implements a valuable user experience improvement by providing a comprehensive configuration template. The implementation leverages existing robust configuration infrastructure while adding meaningful new functionality.

**Key Implementation Token**: `// CFG-TEMPLATE-001: Configuration template generation`

---

**Next Action**: Begin Subtask 1 - CLI Command Implementation 