package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"time"
)

type Ping struct {
	Latency float64
}

type Download struct {
	Bandwidth float64
}

type Upload struct {
	Bandwidth float64
}

type Result struct {
	Url string
}

type Speedtest struct {
	Timestamp string
	Ping      Ping
	Download  Download
	Upload    Upload
	Result    Result
}

func main() {

Cmd:
	cmd := exec.Command("speedtest", "-fjson")
	stdout, err := cmd.Output()

	if err != nil {
		goto Cmd
	}

	var speedtest Speedtest
	json.Unmarshal([]byte(stdout), &speedtest)
	ping := math.Round(speedtest.Ping.Latency)
	dlMbps := math.Round((speedtest.Download.Bandwidth/125000)*100) / 100
	ulMbps := math.Round((speedtest.Upload.Bandwidth/125000)*100) / 100

	timestamp, _ := time.Parse("2006-01-02T15:04:05Z", speedtest.Timestamp)
	loc, _ := time.LoadLocation("Europe/Budapest")

	url := speedtest.Result.Url
	csv := fmt.Sprintf(`"%.0f", "%.2f Mbps", "%.2f Mbps", "%s", "%s"`+"\n", ping, dlMbps, ulMbps, timestamp.In(loc).Format("2006-01-02 3:04 PM"), url)

	io.WriteString(os.Stdout, csv)
}
