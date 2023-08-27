package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	bitcoinBlocks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_blocks",
		Help: "Block height",
	})
	bitcoinDifficulty = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_difficulty",
		Help: "Difficulty",
	})
	bitcoinPeers = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_peers",
		Help: "Number of peers",
	})
	bitcoinConnIn = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_conn_in",
		Help: "Number of connections in",
	})
	bitcoinConnOut = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "bitcoin_conn_out",
		Help: "Number of connections out",
	})

	bitcoinHashPS = make(map[int]prometheus.Gauge)
	bitcoinEstSmartFee = make(map[int]prometheus.Gauge)

	retryExceptions = []string{
		"btcrpcclient.InWarmupError",
		"btcrpcclient.JSONRPCError",
		"btcrpcclient.BTCDRPCError",
	}

	bitcoinRPCUser   string
	bitcoinRPCPass   string
	bitcoinRPCServer string
)

func init() {
	// Register all metrics with the Prometheus client library
	prometheus.MustRegister(bitcoinBlocks)
	prometheus.MustRegister(bitcoinDifficulty)
	prometheus.MustRegister(bitcoinPeers)
	prometheus.MustRegister(bitcoinConnIn)
	prometheus.MustRegister(bitcoinConnOut)

	for _, numBlocks := range []int{-1, 1, 120} {
		bitcoinHashPS[numBlocks] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("bitcoin_hashps_%d", numBlocks),
			Help: fmt.Sprintf("Estimated network hash rate per second for the last %d blocks", numBlocks),
		})
		prometheus.MustRegister(bitcoinHashPS[numBlocks])
	}

	for _, numBlocks := range []int{2, 3, 5, 20} {
		bitcoinEstSmartFee[numBlocks] = prometheus.NewGauge(prometheus.GaugeOpts{
			Name: fmt.Sprintf("bitcoin_est_smart_fee_%d", numBlocks),
			Help: fmt.Sprintf("Estimated smart fee per kilobyte for confirmation in %d blocks", numBlocks),
		})
		prometheus.MustRegister(bitcoinEstSmartFee[numBlocks])
	}

	// Add other metrics as per the Python code
	// ...
}

func main() {
	// Parse command-line flags for Bitcoin RPC settings
	flag.StringVar(&bitcoinRPCUser, "rpcuser", "", "Bitcoin RPC username")
	flag.StringVar(&bitcoinRPCPass, "rpcpass", "", "Bitcoin RPC password")
	flag.StringVar(&bitcoinRPCServer, "rpcserver", "localhost:8332", "Bitcoin RPC server address")
	flag.Parse()

	// Set up Prometheus metrics handler
	http.Handle("/metrics", promhttp.Handler())

	// Gracefully handle SIGTERM to shut down the HTTP server
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM)
	go func() {
		<-signalChannel
		log.Println("Received SIGTERM. Exiting...")
		os.Exit(0)
	}()

	// Start the HTTP server to serve metrics
	log.Println("Starting HTTP server on :9332")
	log.Fatal(http.ListenAndServe(":9332", nil))
}

// Add other metric collection functions as per the Python code
// ...
