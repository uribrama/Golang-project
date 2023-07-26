package config

import "github.com/uribrama/Golang-project/logger"

// Mock is a task-specific mock that can be built out as needed
type Mock struct{}

// Get mocks the Config get method
func (m Mock) Get() ProjectConfig {
	return ProjectConfig{
		// EnvVar
		EnvVar: EnvVar{
			Environment: "test",
			Debug:       true,
		},
	}
}

/*
// GetSecret mocks the Config GetSecret method
// Instead of putting actual secrets here, it's generally recommended to put
// tests that run with aws credentials in the acceptance/ directory
func (m Mock) GetSecret(path, key string) (string, error) {
	return "secretpassword", nil
}*/

// GetLoggingConfig mocks the method
func (m Mock) GetLogging() *logger.Logger {
	return logger.New(c.proyectConfig.Debug)
}
