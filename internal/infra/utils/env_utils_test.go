package utils_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vs0uz4/labs-auction/internal/infra/utils"
)

func TestGetMaxBatchSizeInterval(t *testing.T) {
	// Variável de ambiente definida corretamente
	os.Setenv("BATCH_INSERT_INTERVAL", "1m")
	defer os.Unsetenv("BATCH_INSERT_INTERVAL")

	interval := utils.GetMaxBatchSizeInterval()
	assert.Equal(t, 1*time.Minute, interval)

	// Variável de ambiente inválida
	os.Setenv("BATCH_INSERT_INTERVAL", "invalid")
	interval = utils.GetMaxBatchSizeInterval()
	assert.Equal(t, 3*time.Minute, interval)

	// Variável de ambiente não definida
	os.Unsetenv("BATCH_INSERT_INTERVAL")
	interval = utils.GetMaxBatchSizeInterval()
	assert.Equal(t, 3*time.Minute, interval)
}

func TestGetMaxBatchSize(t *testing.T) {
	// Variável de ambiente definida corretamente
	os.Setenv("MAX_BATCH_SIZE", "10")
	defer os.Unsetenv("MAX_BATCH_SIZE")

	size := utils.GetMaxBatchSize()
	assert.Equal(t, 10, size)

	// Variável de ambiente inválida
	os.Setenv("MAX_BATCH_SIZE", "invalid")
	size = utils.GetMaxBatchSize()
	assert.Equal(t, 5, size)

	// Variável de ambiente não definida
	os.Unsetenv("MAX_BATCH_SIZE")
	size = utils.GetMaxBatchSize()
	assert.Equal(t, 5, size)
}
