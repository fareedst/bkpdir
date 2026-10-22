# [IMPL-TIED_FILES] [ARCH-TIED_STRUCTURE] [REQ-TIED_SETUP]

## Summary contract

Bootstrap script copy_files.sh installs tied-project/ methodology layout, project indexes, detail directories, and root agent loader files for new client repositories.

INPUT: target project path, TIED template tree
OUTPUT: tied-project/ directory with indexes and guides, root AGENTS.md and .cursorrules
DATA: templates/ sources, methodology vs project YAML split

## TIED_FILE_CREATION

SPEC-ID: IMPL-TIED_FILES::TIED_FILE_CREATION
STEP T001: create tied-project/, copy YAML indexes and guides from templates, create detail subdirs, and copy AGENTS.md plus .cursorrules to project root

- [IMPL-TIED_FILES] [ARCH-TIED_STRUCTURE] [REQ-TIED_SETUP] — How: create tied-project/, copy YAML indexes and guides from templates, create detail subdirs, and copy AGENTS.md plus .cursorrules to project root.

PROCEDURE tied_file_creation(target_project):
  # [IMPL-TIED_FILES][ARCH-TIED_STRUCTURE] [REQ-TIED_SETUP] How: create tied-project/, copy YAML indexes and guides from templates, create detail subdirs, and copy AGENTS.md plus .cursorrules to project root.
  Contract:
    INPUT: target project path, TIED template tree
    OUTPUT: tied-project/ directory with indexes and guides, root AGENTS.md and .cursorrules
    PRE: prerequisites for this process step are satisfied
    POST: step outputs match the documented How intent
    EFFECTS: process documentation; runtime unchanged unless a separate LEAP slice applies


  CREATE tied-project/ at project root
  COPY requirements.yaml, architecture-decisions.yaml, implementation-decisions.yaml, semantic-tokens.yaml
  COPY semantic-tokens.md, processes.md, commit-guidelines.md
  CREATE requirements/, architecture-decisions/, implementation-decisions/ under tied-project/
  COPY .cursorrules and AGENTS.md to project root
