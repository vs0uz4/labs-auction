package utils_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vs0uz4/labs-auction/internal/infra/utils"
)

func TestGetAuctionInterval(t *testing.T) {
	// Variável de ambiente definida corretamente
	os.Setenv("AUCTION_INTERVAL", "2m")
	defer os.Unsetenv("AUCTION_INTERVAL")

	interval := utils.GetAuctionInterval()
	assert.Equal(t, 2*time.Minute, interval)

	// Variável de ambiente inválida
	os.Setenv("AUCTION_INTERVAL", "invalid")
	interval = utils.GetAuctionInterval()
	assert.Equal(t, 5*time.Minute, interval)

	// Variável de ambiente não definida
	os.Unsetenv("AUCTION_INTERVAL")
	interval = utils.GetAuctionInterval()
	assert.Equal(t, 5*time.Minute, interval)
}

func TestIsAuctionExpired(t *testing.T) {
	os.Setenv("AUCTION_INTERVAL", "5m")
	defer os.Unsetenv("AUCTION_INTERVAL")

	// Timestamp expirado
	expiredTimestamp := time.Now().Add(-10 * time.Minute).Unix()
	assert.True(t, utils.IsAuctionExpired(expiredTimestamp))

	// Timestamp válido
	validTimestamp := time.Now().Add(-2 * time.Minute).Unix()
	assert.False(t, utils.IsAuctionExpired(validTimestamp))
}
