package main

import (
	"os"
	"testing"
)

func TestLoadConfigFromFiles(t *testing.T) {
	t.Run("a single config file should match the generated result",
		func(t *testing.T) {
			var (
				configs = []string{
					"database:\n  host: localhost\n  port: 5432\n  user: root\n  password: pass\n  database: testDb",
				}
				expected = &Config{
					Database: DatabaseConfig{
						Host:     "localhost",
						Port:     "5432",
						User:     "root",
						Password: "pass",
						Database: "testDb",
					},
				}
			)

			cfg, err := runTestLoadConfigFromFiles(configs)

			if err != nil {
				t.Fatalf("Test failed with error: %s", err)
			}

			if *cfg != *expected {
				t.Fatalf("expected %+v but got %+v", *expected, *cfg)
			}
		})
	t.Run("using multiple config files should apply their members in order",
		func(t *testing.T) {
			var (
				configs = []string{
					"database:\n  host: localhost\n  port: 1234\n  user: root\n  password: pass\n  database: testDb",
					"database:\n  database: databaseName\n  user: userName\n  password: userPassword",
				}
				expected = &Config{
					Database: DatabaseConfig{
						Host:     "localhost",
						Port:     "1234",
						User:     "userName",
						Password: "userPassword",
						Database: "databaseName",
					},
				}
			)

			cfg, err := runTestLoadConfigFromFiles(configs)

			if err != nil {
				t.Fatalf("Test failed with error: %s", err)
			}

			if *cfg != *expected {
				t.Fatalf("expected %+v but got %+v", *expected, *cfg)
			}
		})
	t.Run("providing no config should cause an error",
		func(t *testing.T) {
			var configs []string

			cfg, err := runTestLoadConfigFromFiles(configs)

			if err == nil {
				t.Fatalf("test was expected to error but instead returned %+v", *cfg)
			}
		})
	t.Run("providing an invalid file should error (json)",
		func(t *testing.T) {
			var configs = []string{
				"{'database': {'host':'localhost', 'port':5432, 'user':'root' 'password':'pass' 'database':'testDb'}}",
			}

			cfg, err := runTestLoadConfigFromFiles(configs)

			if err == nil {
				t.Fatalf("test was expected to error but instead returned %+v", *cfg)
			}
		})
	t.Run("if all files are missing there should only be an error returned",
		func(t *testing.T) {
			var configs = []string{"THIS_FILE_DOES_NOT_EXIST.yml"}

			cfg, err := LoadConfigFromFiles(configs...)

			if err == nil {
				t.Fatalf("test was expected to error but instead returned %+v", *cfg)
			}

			if cfg != nil {
				t.Fatalf("no config was expected to be returned")
			}
		})
	t.Run("if at least one config file is not missing, config should be generated but an error also returned",
		func(t *testing.T) {
			var (
				missingFile = "THIS_FILE_DOES_NOT_EXIST.yml"
				config      = "database:\n  host: localhost\n  port: 5432\n  user: root\n  password: pass\n  database: testDb"
				expected    = &Config{
					Database: DatabaseConfig{
						Host:     "localhost",
						Port:     "5432",
						User:     "root",
						Password: "pass",
						Database: "testDb",
					},
				}
				configFile, err = makeTempFile(config)
			)

			cfg, err := LoadConfigFromFiles(configFile, missingFile)

			if err == nil {
				t.Fatalf("test was expected to error")
			}

			if cfg == nil {
				t.Fatalf("config was expected to be returned")
			}

			if *cfg != *expected {
				t.Fatalf("expected %+v but got %+v", *expected, *cfg)
			}
		})
}

func runTestLoadConfigFromFiles(configs []string) (*Config, error) {
	fileNames := make([]string, len(configs))
	for i, config := range configs {
		fileName, err := makeTempFile(config)
		if err != nil {
			return nil, err
		}
		defer os.Remove(fileName)

		fileNames[i] = fileName
	}

	got, err := LoadConfigFromFiles(fileNames...)
	return got, err
}

func makeTempFile(config string) (string, error) {
	f, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(config)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
