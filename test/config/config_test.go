package config_test

import (
	"burakozkan138/questionanswerapi/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	t.Run("Test Load Config", func(t *testing.T) {
		err := config.LoadConfig(".env.test", "env", "../../config")
		assert.NoError(t, err)
	})

	t.Run("Test Initialize Config", func(t *testing.T) {
		err := config.LoadConfig(".env.test", "env", "../../config")
		assert.NoError(t, err)
		cfg, err := config.InitializeConfig()
		assert.NoError(t, err)
		assert.NotNil(t, cfg)
	})
}
