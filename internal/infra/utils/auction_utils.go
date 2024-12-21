package utils

import (
	"os"
	"time"
)

func GetAuctionInterval() time.Duration {
	auctionInterval := os.Getenv("AUCTION_INTERVAL")
	duration, err := time.ParseDuration(auctionInterval)
	if err != nil {
		return time.Minute * 5
	}

	return duration
}

func IsAuctionExpired(timestamp int64) bool {
	interval := GetAuctionInterval()
	elapsed := time.Since(time.Unix(timestamp, 0))
	return elapsed > interval
}
