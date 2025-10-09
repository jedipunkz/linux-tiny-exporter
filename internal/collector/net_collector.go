package collector

import (
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	netDeviceHeaderLines = 2 // Number of header lines to skip in /proc/net/dev
)

type NetCollector struct {
	receivedBytesDiff   *prometheus.Desc
	transmitBytesDiff   *prometheus.Desc
	receivedPacketsDiff *prometheus.Desc
	transmitPacketsDiff *prometheus.Desc
	lastReceivedBytes   map[string]float64
	lastTransmitBytes   map[string]float64
	lastReceivedPackets map[string]float64
	lastTransmitPackets map[string]float64
}

func NewNetCollector() *NetCollector {
	return &NetCollector{
		receivedBytesDiff: prometheus.NewDesc("net_received_bytes_diff",
			"Received bytes by network interface since the last scrape.",
			[]string{"interface"}, nil,
		),
		transmitBytesDiff: prometheus.NewDesc("net_transmit_bytes_diff",
			"Transmitted bytes by network interface since the last scrape.",
			[]string{"interface"}, nil,
		),
		receivedPacketsDiff: prometheus.NewDesc("net_received_packets_diff",
			"Received packets by network interface since the last scrape.",
			[]string{"interface"}, nil,
		),
		transmitPacketsDiff: prometheus.NewDesc("net_transmit_packets_diff",
			"Transmitted packets by network interface since the last scrape.",
			[]string{"interface"}, nil,
		),
		lastReceivedBytes:   make(map[string]float64),
		lastTransmitBytes:   make(map[string]float64),
		lastReceivedPackets: make(map[string]float64),
		lastTransmitPackets: make(map[string]float64),
	}
}

func (collector *NetCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.receivedBytesDiff
	ch <- collector.transmitBytesDiff
	ch <- collector.receivedPacketsDiff
	ch <- collector.transmitPacketsDiff
}

func (collector *NetCollector) Collect(ch chan<- prometheus.Metric) {
	data, err := getNetDevice()
	if err != nil {
		log.Printf("Error reading network device data: %v", err)
		return
	}

	lines := strings.Split(data, "\n")

	for i, line := range lines {
		if i < netDeviceHeaderLines {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}

		interfaceName := strings.Trim(fields[0], ":")
		receivedBytes, _ := strconv.ParseFloat(fields[1], 64)
		transmitBytes, _ := strconv.ParseFloat(fields[9], 64)
		receivedPackets, _ := strconv.ParseFloat(fields[2], 64)
		transmitPackets, _ := strconv.ParseFloat(fields[10], 64)

		collector.emitMetricDiff(ch, collector.receivedBytesDiff, collector.lastReceivedBytes, interfaceName, receivedBytes)
		collector.emitMetricDiff(ch, collector.transmitBytesDiff, collector.lastTransmitBytes, interfaceName, transmitBytes)
		collector.emitMetricDiff(ch, collector.receivedPacketsDiff, collector.lastReceivedPackets, interfaceName, receivedPackets)
		collector.emitMetricDiff(ch, collector.transmitPacketsDiff, collector.lastTransmitPackets, interfaceName, transmitPackets)

		collector.lastReceivedBytes[interfaceName] = receivedBytes
		collector.lastTransmitBytes[interfaceName] = transmitBytes
		collector.lastReceivedPackets[interfaceName] = receivedPackets
		collector.lastTransmitPackets[interfaceName] = transmitPackets
	}
}

func (collector *NetCollector) emitMetricDiff(ch chan<- prometheus.Metric, desc *prometheus.Desc, lastValues map[string]float64, interfaceName string, currentValue float64) {
	if lastValue, ok := lastValues[interfaceName]; ok {
		diff := currentValue - lastValue
		ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, diff, interfaceName)
	}
}

func getNetDevice() (string, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return "", err
	}
	defer file.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, file)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}
