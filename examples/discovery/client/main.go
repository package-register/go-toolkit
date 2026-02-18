package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"

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
	disc := discovery.NewDiscovery(hostname+"-Client", version, logger)

	// Register a handler for "ping" messages from the server
	disc.RegisterHandler("ping", func(from net.Addr, env discovery.MessageEnvelope) {
		logger.Info("Received PING from %s, responding with PONG", from.String())
		_ = disc.Send(discovery.MessageEnvelope{
			SendType: "response",
			SendTo:   from.String(),
			Command:  "pong",
			TaskID:   env.TaskID,
			Payload:  mustJSON("ok"),
		})
	})

	// Register a handler for "announce" messages (to prevent "unregistered handler" warnings)
	disc.RegisterHandler("announce", func(from net.Addr, env discovery.MessageEnvelope) {
		// The discovery package handles the device list update internally.
		// This handler just prevents the "unregistered command handler" log.
		logger.Info("Received announce from %s (UUID: %s)", from.String(), env.FromUUID)
	})

	// Register a handler for "exec_command" from the server
	disc.RegisterHandler("exec_command", func(from net.Addr, env discovery.MessageEnvelope) {
		var commandString string
		if err := json.Unmarshal(env.Payload, &commandString); err != nil {
			logger.Error("Failed to unmarshal command payload: %v", err)
			return
		}
		logger.Info("Received command to execute from %s: %s (TaskID: %s)", from.String(), commandString, env.TaskID)

		// Execute the command in a non-blocking way
		go func() {
			cmdParts := strings.Fields(commandString)
			if len(cmdParts) == 0 {
				logger.Error("Empty command string received.")
				return
			}

			cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
			output, err := cmd.CombinedOutput()

			result := fmt.Sprintf("Command '%s' executed.\nOutput:\n%s\n", commandString, string(output))
			if err != nil {
				result += fmt.Sprintf("Error: %v\n", err)
				logger.Error("Command execution failed: %v", err)
			} else {
				logger.Info("Command execution successful.")
			}

			// Send the result back to the server
			responseEnv := discovery.MessageEnvelope{
				SendType: "response", // Or a custom type like "command_result"
				SendTo:   from.String(),
				Command:  "command_result", // Command to indicate result
				TaskID:   env.TaskID,
				Payload:  mustJSON(result),
			}
			if err := disc.Send(responseEnv); err != nil {
				logger.Error("Failed to send command result: %v", err)
			}
		}()
	})

	if err := disc.Start(); err != nil {
		log.Fatal(err)
	}
	defer disc.Stop()

	logger.Info("Client started. Announcing presence...")

	// Keep the client running
	select {}
}

// mustJSON is a helper function, kept here for the main package's usage
func mustJSON(v any) json.RawMessage {
	b, _ := json.Marshal(v)
	return b
}
