package util

import (
	"fmt"
	"os"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

func Download(dst string, urlStr string) (*grab.Response, error) {
	// create client
	client := grab.NewClient()
	client.UserAgent = "wget/1.21.4"

	req, _ := grab.NewRequest(dst, urlStr)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %s / %s (%.2f%%)\n",
				HumanSize(float64(resp.BytesComplete())),
				HumanSize(float64(resp.Size())),
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	return resp, resp.Err()
}

func DownloadBatch(worker int, requests ...*grab.Request) error {
	client := grab.NewClient()
	respch := client.DoBatch(10, requests...)

	// start a ticker to update progress every 200ms
	t := time.NewTicker(500 * time.Millisecond)

	// monitor downloads
	completed := 0
	inProgress := 0
	responses := make([]*grab.Response, 0)

	hasErr := false
	for completed < len(requests) {
		select {
		case resp := <-respch:
			// a new response has been received and has started downloading
			// (nil is received once, when the channel is closed by grab)
			if resp != nil {
				responses = append(responses, resp)
			}

		case <-t.C:
			// clear lines
			if inProgress > 0 {
				fmt.Printf("\033[%dA\033[K", inProgress)
			}

			// update completed downloads
			for i, resp := range responses {
				if resp != nil && resp.IsComplete() {
					// print final result

					if err := resp.Err(); err != nil {
						fmt.Fprintf(os.Stderr, "Error downloading %s: %v\n", resp.Request.URL(), err)
						hasErr = true
					} else {
						fmt.Printf(
							"Finished %s %s / %s (%.2f%%)\n",
							resp.Filename,
							HumanSize(float64(resp.BytesComplete())),
							HumanSize(float64(resp.Size())),
							100*resp.Progress())
					}

					// mark completed
					responses[i] = nil
					completed++
				}
			}

			// update downloads in progress
			inProgress = 0
			for _, resp := range responses {
				if resp != nil {
					inProgress++
					fmt.Printf(
						"Downloading %s %s / %s (%.2f%%)\033[K\n",
						resp.Filename,
						HumanSize(float64(resp.BytesComplete())),
						HumanSize(float64(resp.Size())),
						100*resp.Progress())
				}
			}
		}
	}

	t.Stop()

	if hasErr {
		return fmt.Errorf("Some requests failed.")
	}

	fmt.Printf("%d files successfully downloaded.\n", len(requests))
	return nil
}
