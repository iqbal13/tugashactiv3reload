package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Struct untuk menyimpan data water dan wind
type WaterWind struct {
	Water       int    `json:"water"`
	Wind        int    `json:"wind"`
	WaterStatus string `json:"water_status"`
	WindStatus  string `json:"wind_status"`
}

func WaterStatus(water int) string {
	if water < 5 {
		return "Aman"
	} else if water >= 6 && water <= 8 {
		return "Siaga"
	} else {
		return "Bahaya"
	}
}

func WindStatus(wind int) string {
	if wind < 6 {
		return "Aman"
	} else if wind >= 7 && wind <= 15 {
		return "Siaga"
	} else {
		return "Bahaya"
	}
}

const REST_API = "http://localhost:8091/updatedata"

func SendData(data WaterWind) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}
	resp, err := http.Post(REST_API, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error sending data to API: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from API: %s", resp.Status)
	}

	//log.Println("Data successfully sent to API")
	return nil
}

func main() {
	// Seed untuk random number

	// K
	// Ticker untuk memperbarui data setiap 15 detik
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Generate angka acak untuk water dan wind
			water := rand.Intn(100) + 1
			wind := rand.Intn(100) + 1

			// Tentukan status
			waterStatus := WaterStatus(water)
			windStatus := WindStatus(wind)

			// Buat objek WeatherData untuk database

			data := WaterWind{
				Water:       water,
				Wind:        wind,
				WaterStatus: waterStatus,
				WindStatus:  windStatus,
			}

			// Kirim data ke API
			err := SendData(data)
			if err != nil {
				log.Printf("Error sending data to API: %v", err)
			}

			fmt.Printf(`{
				"water": %d,
				"wind": %d
				}
				status water: %s
				status wind: %s`, data.Water, data.Wind, data.WaterStatus, data.WindStatus)

		}
	}
}
