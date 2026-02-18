package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/package-register/go-genius/discovery"
)

const (
	version = "1.0.0"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Panicln(err)
	}

	logger := &discovery.StdLogger{}
	disc := discovery.NewDiscovery(hostname+"-Server", version, logger)

	// Register a handler for "announce" messages (to prevent "unregistered handler" warnings)
	disc.RegisterHandler("announce", func(from net.Addr, env discovery.MessageEnvelope) {
		// The discovery package handles the device list update internally.
		// This handler just prevents the "unregistered command handler" log.
		logger.Info("Received announce from %s (UUID: %s)", from.String(), env.FromUUID)
	})

	// Register a handler for "pong" messages from clients
	disc.RegisterHandler("pong", func(from net.Addr, env discovery.MessageEnvelope) {
		logger.Info("Received PONG from %s (UUID: %s)", from.String(), env.FromUUID)
		// Here you could update client status, etc.
	})

	// Register a handler for command execution results from clients
	disc.RegisterHandler("command_result", func(from net.Addr, env discovery.MessageEnvelope) {
		logger.Info("Received command result from %s (TaskID: %s)", from.String(), env.TaskID)
		var result string
		_ = json.Unmarshal(env.Payload, &result)
		fmt.Printf("Command Result from %s: %s\n", env.FromUUID, result)
	})

	if err := disc.Start(); err != nil {
		log.Fatal(err)
	}
	defer disc.Stop()

	logger.Info("Server started. Waiting for clients and commands...")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command (e.g., 'exec <client_uuid> <command_string>'): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		parts := strings.SplitN(input, " ", 3)
		if len(parts) < 3 || parts[0] != "exec" {
			fmt.Println("Invalid command format. Use 'exec <client_uuid> <command_string>'")
			continue
		}

		clientUUID := parts[1]
		commandString := parts[2]

		// Find the client device
		var targetDevice *discovery.Device
		for _, dev := range disc.GetDevices() {
			if dev.UUID == clientUUID {
				targetDevice = dev
				break
			}
		}

		if targetDevice == nil {
			fmt.Printf("Client with UUID %s not found.\n", clientUUID)
			fmt.Println("Available devices:")
			for _, dev := range disc.GetDevices() {
				fmt.Printf("  UUID: %s, Name: %s, IP: %s:%d\n", dev.UUID, dev.Name, dev.IP, dev.Port)
			}
			continue
		}

		// Send the command to the client
		taskID := uuid.New().String()
		cmdEnv := discovery.MessageEnvelope{
			SendType: "spec", // Specific client
			SendTo:   fmt.Sprintf("%s:%d", targetDevice.IP, targetDevice.Port),
			Command:  "exec_command", // New command for client to execute
			TaskID:   taskID,
			Payload:  mustJSON(commandString),
		}

		logger.Info("Sending command '%s' to client %s (TaskID: %s)", commandString, clientUUID, taskID)
		if err := disc.Send(cmdEnv); err != nil {
			logger.Error("Failed to send command: %v", err)
		}
	}
}

// mustJSON is a helper function, kept here for the main package's usage
func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
