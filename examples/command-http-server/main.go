package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/xfrr/go-cqrsify/messaging"
	messaginghttp "github.com/xfrr/go-cqrsify/messaging/http"
	"github.com/xfrr/go-cqrsify/pkg/apix"
)

const sampleCommandType = "cqrsify.commands.print_text"

type sampleCommand struct {
	messaging.BaseCommand

	sampleCommandAttributes
}

type sampleCommandAttributes struct {
	Text string `json:"text"`
}

type sampleCommandHandler struct{}

func (h *sampleCommandHandler) Handle(ctx context.Context, cmd sampleCommand) error {
	fmt.Println("")
	fmt.Printf("Received command: %s\n", cmd.MessageType())
	fmt.Printf("Command ID: %s\n", cmd.CommandID())
	fmt.Printf("Command Text: %s\n", cmd.Text)
	fmt.Println("")
	return nil
}

func main() {
	ctx, cancelSignal := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancelSignal()

	// Create inmemory commandBus
	cmdbus := messaging.NewInMemoryCommandBus(
		messaging.ConfigureInMemoryMessageBusSubjects(sampleCommandType),
		messaging.ConfigureInMemoryMessageBusErrorHandler(
			func(messageType string, err error) {
				fmt.Printf("Error handling message of type %s: %v\n", messageType, err)
			},
		))

	// Register command handler
	unsub, err := messaging.SubscribeCommand(
		ctx,
		cmdbus,
		&sampleCommandHandler{},
	)
	if err != nil {
		panic(err)
	}
	defer unsub()

	// Create HTTP Command Bus HTTP Handler
	cmdHandler := messaginghttp.NewCommandHandler(cmdbus)
	messaginghttp.RegisterJSONAPICommandDecoder(
		cmdHandler,
		sampleCommandType,
		func(sd apix.SingleDocument[sampleCommandAttributes]) (messaging.Command, error) {
			cmd := sampleCommand{
				BaseCommand:             messaginghttp.CreateBaseCommandFromSingleDocument(sampleCommandType, sd),
				sampleCommandAttributes: sd.Data.Attributes,
			}
			return cmd, nil
		})

	// Create Gin engine
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	// Create HTTP Command Server
	commandHTTPServer := messaginghttp.NewCommandGINServer(cmdHandler, r)
	commandWSServer := messaginghttp.NewCommandWebsocketServer(cmdbus)
	r.Any("/ws/commands", gin.WrapH(commandWSServer))

	// Start HTTP Command Server
	go func() {
		_ = commandHTTPServer.ListenAndServe(":8091")
	}()

	defer func() {
		_ = commandHTTPServer.Close()
		_ = commandWSServer.Close()
	}()

	fmt.Println("HTTP Command Server is running on :8091")
	fmt.Println("Press Ctrl+C to stop.")

	// curl example
	fmt.Println(" ")
	fmt.Println("Example curl command to send a command:")
	fmt.Println("> websocat ws://localhost:8091/ws/commands")
	fmt.Println(`> curl -X POST http://localhost:8091/commands -H "Content-Type: application/vnd.api+json" -d '{"data": {"type": "cqrsify.commands.print_text", "id": "cmd-123", "attributes": {"text": "Hello, CQRS!"}}}'`)
	<-ctx.Done()
	fmt.Println("Shutting down HTTP Command Server...")
}
