# Waker - Wake-on-LAN Server

A simple Go server with a Vite/Vue frontend to send Wake-on-LAN magic packets.

## Setup

### Local Development

1. **Install dependencies:**
   ```bash
   go mod download
   cd web
   npm install
   cd ..
   ```

2. **Configure MAC address:**
   - Copy `.env.example` to `.env`
   - Edit `.env` and set `MAC_ADDRESS` to your device's MAC address
   - Supported formats: `00:11:22:33:44:55`, `00-11-22-33-44-55`, `0011.2233.4455`, or `001122334455`

3. **Build the frontend:**
   ```bash
   cd web
   npm run build
   cd ..
   ```

4. **Run the server:**
   ```bash
   go run server/main.go
   ```

5. **Access the web interface:**
   - Open http://localhost:8080 in your browser
   - Click "Wake Up Device" to send the magic packet

### Docker

Build and run with Docker:

```bash
# Build the image
docker build -t waker .

# Run the container
docker run -d \
  --name waker \
  -p 8080:8080 \
  -e MAC_ADDRESS=00:11:22:33:44:55 \
  -e BROADCAST_IP=192.168.1.255 \
  -e WOL_PORT=9 \
  waker

# Access at http://localhost:8080
```

Pull from GitHub Container Registry:
```bash
docker run -d \
  --name waker \
  -p 8080:8080 \
  -e MAC_ADDRESS=00:11:22:33:44:55 \
  -e BROADCAST_IP=192.168.1.255 \
  -e WOL_PORT=9 \
  ghcr.io/<your-username>/waker:latest
```

## Development

To run the frontend in development mode with hot reload:
```bash
cd web
npm run dev
```

The dev server runs on port 3000 and proxies API requests to the Go server.

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `MAC_ADDRESS` | MAC address of the device to wake | Required |
| `PORT` | HTTP server port | 8080 |
| `BROADCAST_IP` | UDP broadcast IP address | 255.255.255.255 |
| `WOL_PORT` | Wake-on-LAN UDP port | 9 |

## How it works

The server sends a Wake-on-LAN magic packet (UDP broadcast) to the local network. The magic packet consists of:
- 6 bytes of 0xFF (broadcast address)
- 16 repetitions of the target MAC address

This is sent to the broadcast address `255.255.255.255` on port 9 (standard WoL port).

## Notes

- The target device must have Wake-on-LAN enabled in its BIOS/UEFI settings
- The device's network interface must be configured to respond to WoL
- Both the server and target device must be on the same local network
- Some networks may block broadcast packets, which would prevent WoL from working
