# Generic Documentation Update Prompts for Go Projects

This file contains generic prompts that can be used to instruct AI models to update project documentation based on changes in the codebase. These prompts are designed to be adaptable to different Go projects.

## Table of Contents
- [Update Component Documentation](#update-component-documentation)
- [Add New Component](#add-new-component)
- [Update Examples](#update-examples)
- [Update Dependencies](#update-dependencies)
- [Update Testing Information](#update-testing-information)
- [Complete Guide Review](#complete-guide-review)

## Update Component Documentation

Use this prompt when a specific component has been modified:

```
Review the current implementation of the {{component_name}} component in the current Go project by examining the relevant source files. Identify the files related to this component (typically named after the component).

Compare the implementation details with the {{component_name}} section in the main documentation guide. Update the documentation to reflect any changes in:
- Function signatures and parameters
- Configuration options
- Behavior or functionality
- Error handling
- Context usage
- Return values
- Interface definitions

Ensure code examples are accurate and reflect the current API. Pay attention to:
- Exported functions and methods
- Struct definitions and field descriptions
- Interface implementations
- Package-level variables and constants
```

## Add New Component

Use this prompt when a new component has been added to the package:

```
A new component named {{component_name}} has been added to the Go project. Please add a new section to the main documentation guide that includes:

1. A header for the new component following the existing format
2. Key Features section describing what the component does
3. Main Functions/Methods section with all exported functions:
   - Function signatures with parameter types
   - Parameter descriptions
   - Return value descriptions
   - Error conditions
4. Usage examples demonstrating how to use the component
5. Any configuration options or structs with detailed field descriptions
6. Context integration details if applicable
7. Dependencies specific to this component
8. Interface implementations if relevant

Review the relevant source files for implementation details to understand:
- How the component is intended to be used
- What problems it solves
- Its relationship to other components
- Any special considerations or limitations
```

## Update Examples

Use this prompt when examples need to be updated:

```
Review all code examples in the main documentation guide to ensure they accurately reflect the current API usage. Check:

1. Function signatures and parameter usage
2. Import statements and package references
3. Error handling patterns
4. Context usage (if applicable)
5. Method/Function call sequences
6. Interface usage
7. Package initialization and setup

Update any examples that no longer match the current implementation. Pay special attention to:
- Common usage patterns
- Best practices demonstrated in the codebase
- Error handling approaches
- Resource management (cleanup, defer statements)
- Concurrency patterns (if applicable)

Verify examples by cross-referencing with test files to ensure they follow established patterns in the project.
```

## Update Dependencies

Use this prompt when dependencies have changed:

```
Review the project's dependency management file (typically go.mod for Go projects) to identify any changes in dependencies. Update the Dependencies section in the main documentation guide to reflect:

1. New dependencies added
2. Dependencies that have been removed
3. Version updates for existing dependencies
4. Any changes in subpackage dependencies

For each dependency change, provide:
- The dependency name and version
- Its purpose in the project
- Any impact on functionality or API
- Migration considerations (if upgrading)

Ensure any version requirements (like minimum Go version) are also updated if needed.
```

## Update Testing Information

Use this prompt when testing procedures or coverage have changed:

```
Review the testing approach in the project by examining:
- Test files throughout the codebase
- Test utilities or helpers
- Any testing configuration files
- Scripts used to run tests

Update the Testing section in the main documentation guide to reflect changes in:
1. Test coverage requirements or goals
2. Testing patterns or methodologies used in the project
3. New testing utilities or frameworks adopted
4. Changes in how tests are executed
5. Updates to test structure or organization
6. Mocking strategies or test doubles used

Ensure the documentation accurately describes:
- How to run tests
- What testing frameworks are used
- How to write new tests that follow project conventions
- What to expect in terms of output or reporting
```

## Complete Guide Review

Use this prompt for a comprehensive update of the entire guide:

```
Perform a complete review of the Go project and update the main documentation guide to ensure it accurately reflects the current implementation. Follow these steps:

1. Review all source files to understand the current API and functionality
2. Review all test files to understand usage patterns and examples
3. Check dependency management files for dependency information
4. Review README and other documentation files for context

Update the main documentation guide with changes in:
- Project overview and core concepts
- Package structure and organization
- All component documentation with current function signatures
- Configuration options and structs
- Usage examples that reflect current best practices
- Dependencies and version requirements
- Testing information
- Development guidelines and conventions
- Important constants, variables, and types
- Interface definitions and implementations

Ensure all code examples are accurate and complete. Verify that cross-references between sections are correct and that the table of contents is up-to-date.

Pay special attention to:
- API changes and breaking modifications
- New patterns or conventions adopted
- Deprecated features or functions
- Performance considerations documented
- Security practices mentioned
- Error handling approaches
- Resource management patterns
```

## Update README.md

Use this prompt when the README.md file needs to be created or updated:

```
Review the Go project and create or update the README.md file to serve as a human-friendly introduction to the project. The README should be concise but informative, providing essential information without being overwhelming. Include:

1. Project name and brief description (1-2 sentences)
2. Key features (bullet points)
3. Installation instructions
4. Quick start example (a simple but complete working example)
5. Core components overview (brief descriptions)
6. Links to detailed documentation files
7. Testing information
8. Dependencies list with purposes
9. License information

Ensure the README:
- Is approachable for newcomers
- Accurately reflects the current API
- Contains working code examples
- Uses clear, simple language
- Is well-organized with appropriate headings
- Is not too long (aim for scannable content)
- Includes only the most important information
- Has up-to-date installation and usage instructions

Review the existing README.md and update it to match the current implementation. If it doesn't exist, create it following this structure.
```

## Usage Instructions

To use these prompts effectively:

1. Identify what has changed in the codebase
2. Select the appropriate prompt template
3. Fill in the placeholders with specific component names or details as needed
4. Provide the AI model with access to the relevant source files
5. Ask the AI to update the documentation using the prompt

These generic prompts are designed to work with most Go projects. The AI should be able to identify:
- The project name from the codebase context
- File naming conventions by examining existing files
- Component relationships by analyzing imports and usage
- Key constants and context values by searching the codebase
- Project-specific patterns and conventions by reviewing multiple files

This approach ensures that documentation can remain accurate and up-to-date as projects evolve, regardless of their specific structure or naming conventions.