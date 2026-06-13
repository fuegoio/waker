package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Get MAC address from environment variable
	macAddress := os.Getenv("MAC_ADDRESS")
	if macAddress == "" {
		log.Fatal("MAC_ADDRESS environment variable is required")
	}

	// Clean and validate MAC address
	macAddress = cleanMAC(macAddress)
	if !isValidMAC(macAddress) {
		log.Fatalf("Invalid MAC address: %s", macAddress)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Get broadcast address from environment variable
	broadcastAddr := os.Getenv("BROADCAST_IP")
	if broadcastAddr == "" {
		broadcastAddr = "255.255.255.255"
	}

	// Get WoL UDP port from environment variable
	wolPort := os.Getenv("WOL_PORT")
	if wolPort == "" {
		wolPort = "9"
	}

	// Serve static files from the web/dist directory
	fs := http.FileServer(http.Dir("../web/dist"))
	http.Handle("/", fs)

	// API endpoint to send Wake-on-LAN
	http.HandleFunc("/api/wake", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		log.Printf("Sending Wake-on-LAN to %s via %s:%s", macAddress, broadcastAddr, wolPort)
		err := sendWakeOnLan(macAddress, broadcastAddr, wolPort)
		if err != nil {
			log.Printf("Failed to send Wake-on-LAN: %v", err)
			http.Error(w, "Failed to send Wake-on-LAN", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Wake-on-LAN packet sent to %s", macAddress)
	})

	log.Printf("Server starting on port %s, MAC address: %s, broadcast: %s:%s", port, macAddress, broadcastAddr, wolPort)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// cleanMAC removes all non-hex characters and converts to uppercase
func cleanMAC(mac string) string {
	mac = strings.ToUpper(mac)
	var result strings.Builder
	for _, c := range mac {
		if (c >= '0' && c <= '9') || (c >= 'A' && c <= 'F') {
			result.WriteRune(c)
		}
	}
	return result.String()
}

// isValidMAC checks if the MAC address is valid (12 hex characters)
func isValidMAC(mac string) bool {
	if len(mac) != 12 {
		return false
	}
	for _, c := range mac {
		if !(c >= '0' && c <= '9') && !(c >= 'A' && c <= 'F') {
			return false
		}
	}
	return true
}

// sendWakeOnLan sends a Wake-on-LAN magic packet to the specified MAC address
func sendWakeOnLan(mac, broadcastIP, port string) error {
	// Create magic packet: 6 bytes of 0xFF followed by 16 repetitions of the MAC address
	magicPacket := make([]byte, 6+6*16)

	// First 6 bytes are 0xFF
	for i := 0; i < 6; i++ {
		magicPacket[i] = 0xFF
	}

	// Parse MAC address and add it 16 times
	macBytes := parseMAC(mac)
	for i := 0; i < 16; i++ {
		copy(magicPacket[6+i*6:], macBytes)
	}

	// Build broadcast address with port
	broadcastAddr := broadcastIP + ":" + port

	// Resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", broadcastAddr)
	if err != nil {
		return fmt.Errorf("failed to resolve UDP address: %v", err)
	}

	// Create UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("failed to dial UDP: %v", err)
	}
	defer conn.Close()

	// Enable broadcast
	if err := conn.SetWriteBuffer(0); err != nil {
		log.Printf("Warning: failed to set write buffer: %v", err)
	}

	// Send magic packet
	_, err = conn.Write(magicPacket)
	if err != nil {
		return fmt.Errorf("failed to write magic packet: %v", err)
	}

	return nil
}

// parseMAC converts a MAC address string to bytes
func parseMAC(mac string) []byte {
	if len(mac) != 12 {
		return nil
	}

	bytes := make([]byte, 6)
	for i := 0; i < 6; i++ {
		bytes[i] = hexToByte(mac[i*2 : i*2+2])
	}
	return bytes
}

// hexToByte converts two hex characters to a byte
func hexToByte(s string) byte {
	var result byte
	for i := 0; i < 2; i++ {
		c := s[i]
		var value byte
		if c >= '0' && c <= '9' {
			value = c - '0'
		} else if c >= 'A' && c <= 'F' {
			value = c - 'A' + 10
		}
		result = result*16 + value
	}
	return result
}
