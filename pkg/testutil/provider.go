// [IMPL-TESTING] [ARCH-TESTING_STRATEGY] [REQ-CODE_QUALITY]
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
package testutil

// DefaultTestUtilProvider provides access to all testing utilities.
// This is the main entry point for the testutil package.
type DefaultTestUtilProvider struct {
	configProvider  ConfigProvider
	envManager      EnvironmentManager
	fsHelper        FileSystemTestHelper
	cliHelper       CliTestHelper
	assertionHelper AssertionHelper
	fixtureManager  TestFixtureManager
}

// NewTestUtilProvider creates a new test utility provider with default implementations.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func NewTestUtilProvider() TestUtilProvider {
	return &DefaultTestUtilProvider{
		configProvider:  NewConfigProvider(nil),
		envManager:      NewEnvironmentManager(),
		fsHelper:        NewFileSystemTestHelper(),
		cliHelper:       NewCliTestHelper(),
		assertionHelper: NewAssertionHelper(),
		fixtureManager:  NewTestFixtureManager(""),
	}
}

// NewTestUtilProviderWithFixtureDir creates a test utility provider with a custom fixture directory.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func NewTestUtilProviderWithFixtureDir(fixtureDir string) TestUtilProvider {
	return &DefaultTestUtilProvider{
		configProvider:  NewConfigProvider(nil),
		envManager:      NewEnvironmentManager(),
		fsHelper:        NewFileSystemTestHelper(),
		cliHelper:       NewCliTestHelper(),
		assertionHelper: NewAssertionHelper(),
		fixtureManager:  NewTestFixtureManager(fixtureDir),
	}
}

// GetConfigProvider returns a configuration provider for testing.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (p *DefaultTestUtilProvider) GetConfigProvider() ConfigProvider {
	return p.configProvider
}

// GetEnvironmentManager returns an environment manager for testing.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (p *DefaultTestUtilProvider) GetEnvironmentManager() EnvironmentManager {
	return p.envManager
}

// GetFileSystemHelper returns a file system helper for testing.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (p *DefaultTestUtilProvider) GetFileSystemHelper() FileSystemTestHelper {
	return p.fsHelper
}

// GetCliHelper returns a CLI helper for testing.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (p *DefaultTestUtilProvider) GetCliHelper() CliTestHelper {
	return p.cliHelper
}

// GetAssertionHelper returns an assertion helper for testing.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (p *DefaultTestUtilProvider) GetAssertionHelper() AssertionHelper {
	return p.assertionHelper
}

// GetFixtureManager returns a fixture manager for testing.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (p *DefaultTestUtilProvider) GetFixtureManager() TestFixtureManager {
	return p.fixtureManager
}

// CreateScenario creates a new test scenario.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func (p *DefaultTestUtilProvider) CreateScenario(name, description string) TestScenario {
	return NewTestScenario(name, description)
}

// Package-level convenience functions

// GetDefaultProvider returns a default test utility provider.
// This is the most common way to access testutil functionality.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func GetDefaultProvider() TestUtilProvider {
	return NewTestUtilProvider()
}

// GetProviderWithFixtures returns a test utility provider with a custom fixture directory.
//
// EXTRACT-008: See architecture.md - Package Extraction [DECISION:maintenance]
func GetProviderWithFixtures(fixtureDir string) TestUtilProvider {
	return NewTestUtilProviderWithFixtureDir(fixtureDir)
}
