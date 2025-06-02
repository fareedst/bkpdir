# Context File Responsibilities

## 📁 Feature Tracking Matrix (feature-tracking.md)

- **ALWAYS UPDATE**: Must be updated for every code change
- **Content**: Feature registry, implementation tokens, status tracking, decision records
- **Triggers**: Any code modification, new features, status changes, architectural decisions

## 📄 Specification Document (specification.md)

- **UPDATE WHEN**: User-facing behavior changes, new commands, configuration options, output format changes
- **Content**: User interface specifications, command behaviors, configuration schemas
- **Triggers**: New CLI commands, configuration changes, output modifications, user-visible behavior changes

## 📋 Requirements Document (requirements.md)

- **UPDATE WHEN**: Implementation requirements change, new functional requirements, non-functional requirement modifications
- **Content**: Implementation requirements, traceability matrices, acceptance criteria
- **Triggers**: New features, requirement modifications, acceptance criteria changes

## 🏗️ Architecture Document (architecture.md)

- **UPDATE WHEN**: Technical implementation changes, new components, interface modifications, design decisions
- **Content**: System architecture, component interactions, design patterns, technical decisions
- **Triggers**: Code structure changes, new components, interface modifications, architectural decisions

## 🧪 Testing Document (testing.md)

- **UPDATE WHEN**: New tests added, test strategies change, coverage requirements modified
- **Content**: Test strategies, coverage requirements, test descriptions, validation approaches
- **Triggers**: New test files, test strategy changes, coverage modifications

## 🔒 Immutable Requirements (immutable.md)

- **UPDATE WHEN**: Never (immutable requirements), but must CHECK for violations
- **Content**: Unchangeable requirements, backward compatibility constraints
- **Triggers**: Check only - never modify
