// orchestrator/main.go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/xfrr/go-cqrsify/messaging"
	messagingnats "github.com/xfrr/go-cqrsify/messaging/nats"
	"github.com/xfrr/go-cqrsify/pkg/lock"
	"github.com/xfrr/go-cqrsify/saga"
)

const (
	streamName = "cqrsify_examples_sagas"
)

func main() {
	ctx, signalCancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer signalCancel()

	// Create NATS connection
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		panic(err)
	}

	_, err = js.CreateOrUpdateStream(
		ctx,
		jetstream.StreamConfig{
			Name:      streamName,
			Subjects:  []string{"cqrsify.examples.sagas.>"},
			Storage:   jetstream.MemoryStorage,
			Retention: jetstream.LimitsPolicy,
		},
	)
	if err != nil {
		panic(err)
	}

	cleanup := func() {
		_ = js.DeleteStream(ctx, streamName)
		nc.Drain()
		nc.Close()
	}

	// Create JetStream-based CommandBus
	jsonSerializer := createJsonSerializer()
	jsonDeserializer := createJsonDeserializer()

	publisher, err := messagingnats.NewJetStreamMessagePublisher(
		nc,
		streamName,
		jsonSerializer,
		jsonDeserializer,
	)
	if err != nil {
		cleanup()
		panic(err)
	}

	cmdConsumer, err := messagingnats.NewJetStreamMessageConsumer(
		nc,
		streamName,
		jsonSerializer,
		jsonDeserializer,
		messagingnats.WithConsumerConfig(jetstream.ConsumerConfig{
			Name:      "cqrsify_examples_sagas_consumer",
			AckPolicy: jetstream.AckExplicitPolicy,
			FilterSubjects: []string{
				// Listen to all saga command subjects for simplicity
				// In production, consider split consumers per bounded context
				"cqrsify.examples.sagas.order.reserve.cmd",
				"cqrsify.examples.sagas.order.reserve.compensate.cmd",
				"cqrsify.examples.sagas.payment.charge.cmd",
				"cqrsify.examples.sagas.payment.refund.cmd",
				"cqrsify.examples.sagas.inventory.reduce.cmd",
				"cqrsify.examples.sagas.inventory.restore.cmd",
				"cqrsify.examples.sagas.delivery.ship.cmd",
				"cqrsify.examples.sagas.delivery.cancel.cmd",
			},
		}),
	)
	if err != nil {
		cleanup()
		panic(err)
	}

	cmdRouter := newCommandRoutedHandler()

	unsub, err := cmdConsumer.SubscribeWithReply(ctx, cmdRouter)
	if err != nil {
		cleanup()
		panic(err)
	}
	defer unsub()

	cmdBus := messagingnats.NewJetStreamCommandBus(publisher, cmdConsumer)

	def := &saga.Definition{
		Name: "checkout",
		Steps: []saga.Step{
			{
				Name:          "create_order",
				Action:        saga.RemoteAction(cmdBus, saga.RemoteSubjects{Action: "cqrsify.examples.sagas.order.reserve.cmd", Compensate: "cqrsify.examples.sagas.order.reserve.compensate.cmd", Timeout: 5 * time.Second}),
				Compensate:    saga.RemoteCompensation(cmdBus, saga.RemoteSubjects{Action: "cqrsify.examples.sagas.order.reserve.cmd", Compensate: "cqrsify.examples.sagas.order.reserve.compensate.cmd", Timeout: 5 * time.Second}),
				Retry:         saga.RetryPolicy{MaxAttempts: 5, Backoff: 500 * time.Millisecond, MaxBackoff: 5 * time.Second},
				Timeout:       10 * time.Second,
				IdempotencyFn: func(ex *saga.Execution) string { return "order:" + ex.SagaID },
			},
			{
				Name:          "process_payment",
				Action:        saga.RemoteAction(cmdBus, saga.RemoteSubjects{Action: "cqrsify.examples.sagas.payment.charge.cmd", Compensate: "cqrsify.examples.sagas.payment.refund.cmd", Timeout: 15 * time.Second}),
				Compensate:    saga.RemoteCompensation(cmdBus, saga.RemoteSubjects{Action: "cqrsify.examples.sagas.payment.charge.cmd", Compensate: "cqrsify.examples.sagas.payment.refund.cmd", Timeout: 15 * time.Second}),
				Retry:         saga.RetryPolicy{MaxAttempts: 3, Backoff: time.Second, MaxBackoff: 8 * time.Second},
				Timeout:       20 * time.Second,
				IdempotencyFn: func(ex *saga.Execution) string { return "payment:" + ex.SagaID },
			},
			{
				Name:          "update_inventory",
				Action:        saga.RemoteAction(cmdBus, saga.RemoteSubjects{Action: "cqrsify.examples.sagas.inventory.reduce.cmd", Compensate: "cqrsify.examples.sagas.inventory.restore.cmd"}),
				Compensate:    saga.RemoteCompensation(cmdBus, saga.RemoteSubjects{Action: "cqrsify.examples.sagas.inventory.reduce.cmd", Compensate: "cqrsify.examples.sagas.inventory.restore.cmd"}),
				Retry:         saga.RetryPolicy{MaxAttempts: 4, Backoff: 400 * time.Millisecond, MaxBackoff: 6 * time.Second},
				IdempotencyFn: func(ex *saga.Execution) string { return "inventory:" + ex.SagaID },
			},
			{
				Name:          "deliver_order",
				Action:        saga.RemoteAction(cmdBus, saga.RemoteSubjects{Action: "cqrsify.examples.sagas.delivery.ship.cmd", Compensate: "cqrsify.examples.sagas.delivery.cancel.cmd"}),
				Compensate:    saga.RemoteCompensation(cmdBus, saga.RemoteSubjects{Action: "cqrsify.examples.sagas.delivery.ship.cmd", Compensate: "cqrsify.examples.sagas.delivery.cancel.cmd"}),
				Retry:         saga.RetryPolicy{MaxAttempts: 2, Backoff: time.Second, MaxBackoff: 5 * time.Second},
				IdempotencyFn: func(ex *saga.Execution) string { return "delivery:" + ex.SagaID },
			},
		},
	}

	coord := saga.NewCoordinator(def, saga.NewInMemoryStore(), lock.NewInMemoryLocker(), saga.CoordinatorConfig{
		LockTTL:      10 * time.Second,
		DefaultRetry: saga.RetryPolicy{MaxAttempts: 1, Backoff: 500 * time.Millisecond},
		UUID:         saga.DefaultUUIDProvider,
		Hooks: saga.Hooks{
			OnStepStart: func(ctx context.Context, si *saga.Instance, s saga.StepState) {
				fmt.Printf("→ start %q attempt %d\n", s.Name, s.Attempt)
			},
			OnStepSuccess: func(ctx context.Context, si *saga.Instance, s saga.StepState) {
				fmt.Printf("✓ success %q\n", s.Name)
			},
			OnStepFailure: func(ctx context.Context, si *saga.Instance, s saga.StepState, err error) {
				fmt.Printf("✗ failure %q: %v\n", s.Name, err)
			},
			OnStepCompensationOK: func(ctx context.Context, si *saga.Instance, s saga.StepState) {
				fmt.Printf("↶ compensated %q\n", s.Name)
			},
			OnStepCompensationKO: func(ctx context.Context, si *saga.Instance, s saga.StepState, err error) {
				fmt.Printf("↷ compensation failed for %q: %v\n", s.Name, err)
			},
			OnSagaCompensating: func(ctx context.Context, si *saga.Instance, from int) {
				fmt.Printf("↶ saga %q compensating from step %d\n", si.Name, from)
			},
			OnSagaCompensatingFinished: func(ctx context.Context, si *saga.Instance) {
				switch si.Status {
				case saga.StatusCompensateSuccess:
					fmt.Printf("✔ saga %q compensated successfully\n", si.Name)
				case saga.StatusCompensateFailed:
					fmt.Printf("⚠ saga %q compensated with errors\n", si.Name)
				}
			},
			OnSagaCompleted: func(ctx context.Context, si *saga.Instance) {
				fmt.Printf("🎉 saga %q completed\n", si.Name)
			},
		},
	})

	input := map[string]any{
		"user_id":          "u-42",
		"sku":              "SKU-RED-42",
		"qty":              2,
		"amount":           int64(2599),
		"currency":         "EUR",
		"shipping_address": "Av. Principal 123",
	}

	// Start a saga execution
	id, err := coord.Start(ctx, input, map[string]string{"tenant": "acme"})
	if err != nil {
		cleanup()
		panic(err)
	}

	// A scheduler or HTTP handler would call Run repeatedly until terminal:
	for range 20 {
		select {
		case <-ctx.Done():
			return
		default:
		}

		err = coord.Run(ctx, id)
		if err != nil {
			fmt.Println("✗ fatal error running saga:", err)
			break
		}

		time.Sleep(300 * time.Millisecond)
	}

	cleanup()
}

func createJsonSerializer() messaging.MessageSerializer {
	jsonSerializer := messaging.NewJSONSerializer()
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.order.reserve.cmd",
		saga.RemotePayloadFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.order.reserve.cmd.result",
		saga.RemoteResultFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.order.reserve.compensate.cmd",
		saga.RemotePayloadFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.order.reserve.compensate.cmd.result",
		saga.RemoteResultFromJSONEncoder,
	)

	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.payment.charge.cmd",
		saga.RemotePayloadFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.payment.charge.cmd.result",
		saga.RemoteResultFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.payment.refund.cmd",
		saga.RemotePayloadFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.payment.refund.cmd.result",
		saga.RemoteResultFromJSONEncoder,
	)

	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.inventory.reduce.cmd",
		saga.RemotePayloadFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.inventory.reduce.cmd.result",
		saga.RemoteResultFromJSONEncoder,
	)

	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.inventory.restore.cmd",
		saga.RemotePayloadFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.inventory.restore.cmd.result",
		saga.RemoteResultFromJSONEncoder,
	)

	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.delivery.ship.cmd",
		saga.RemotePayloadFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.delivery.ship.cmd.result",
		saga.RemoteResultFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.delivery.cancel.cmd",
		saga.RemotePayloadFromJSONEncoder,
	)
	messaging.RegisterJSONMessageSerializer(
		jsonSerializer,
		"cqrsify.examples.sagas.delivery.cancel.cmd.result",
		saga.RemoteResultFromJSONEncoder,
	)
	return jsonSerializer
}

func createJsonDeserializer() messaging.MessageDeserializer {
	jsonDeserializer := messaging.NewJSONDeserializer()
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.order.reserve.cmd",
		saga.RemotePayloadFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.order.reserve.cmd.result",
		saga.RemoteResultFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.order.reserve.compensate.cmd",
		saga.RemotePayloadFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.order.reserve.compensate.cmd.result",
		saga.RemoteResultFromJSONDecoder,
	)

	// Payment
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.payment.charge.cmd",
		saga.RemotePayloadFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.payment.charge.cmd.result",
		saga.RemoteResultFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.payment.refund.cmd",
		saga.RemotePayloadFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.payment.refund.cmd.result",
		saga.RemoteResultFromJSONDecoder,
	)

	// Inventory
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.inventory.reduce.cmd",
		saga.RemotePayloadFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.inventory.reduce.cmd.result",
		saga.RemoteResultFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.inventory.restore.cmd",
		saga.RemotePayloadFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.inventory.restore.cmd.result",
		saga.RemoteResultFromJSONDecoder,
	)

	// Delivery
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.delivery.ship.cmd",
		saga.RemotePayloadFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.delivery.ship.cmd.result",
		saga.RemoteResultFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.delivery.cancel.cmd",
		saga.RemotePayloadFromJSONDecoder,
	)
	messaging.RegisterJSONMessageDeserializer(
		jsonDeserializer,
		"cqrsify.examples.sagas.delivery.cancel.cmd.result",
		saga.RemoteResultFromJSONDecoder,
	)
	return jsonDeserializer
}

func newCommandRoutedHandler() *messaging.MessageHandlerWithReplyTypedRouter {
	cmdRouter := messaging.NewMessageHandlerWithReplyTypedRouter()
	cmdRouter.Register(
		"cqrsify.examples.sagas.order.reserve.cmd",
		messaging.NewCommandHandlerWithReplyFn(func(ctx context.Context, cmd saga.RemotePayload) (messaging.CommandReply, error) {
			remoteResultOK := saga.RemoteResultOK(cmd, map[string]any{
				"order_id": "order-12345",
			})
			return remoteResultOK, nil
		}),
	)
	cmdRouter.Register(
		"cqrsify.examples.sagas.order.reserve.compensate.cmd",
		messaging.NewCommandHandlerWithReplyFn(func(ctx context.Context, cmd saga.RemotePayload) (messaging.CommandReply, error) {
			remoteResultOK := saga.RemoteResultOK(cmd, map[string]any{
				"compensated": true,
			})
			return remoteResultOK, nil
		}),
	)

	// Payment
	cmdRouter.Register(
		"cqrsify.examples.sagas.payment.charge.cmd",
		messaging.NewCommandHandlerWithReplyFn(func(ctx context.Context, cmd saga.RemotePayload) (messaging.CommandReply, error) {
			remoteResultOK := saga.RemoteResultOK(cmd, map[string]any{
				"payment_id": "payment-67890",
			})
			return remoteResultOK, nil
		}),
	)
	cmdRouter.Register(
		"cqrsify.examples.sagas.payment.refund.cmd",
		messaging.NewCommandHandlerWithReplyFn(func(ctx context.Context, cmd saga.RemotePayload) (messaging.CommandReply, error) {
			remoteResultOK := saga.RemoteResultOK(cmd, map[string]any{
				"refunded": true,
			})
			return remoteResultOK, nil
		}),
	)
	cmdRouter.Register(
		"cqrsify.examples.sagas.inventory.reduce.cmd",
		messaging.NewCommandHandlerWithReplyFn(func(ctx context.Context, cmd saga.RemotePayload) (messaging.CommandReply, error) {
			remoteResultOK := saga.RemoteResultOK(cmd, map[string]any{
				"inventory_updated": true,
			})
			return remoteResultOK, nil
		}),
	)
	cmdRouter.Register(
		"cqrsify.examples.sagas.inventory.restore.cmd",
		messaging.NewCommandHandlerWithReplyFn(func(ctx context.Context, cmd saga.RemotePayload) (messaging.CommandReply, error) {
			remoteResultOK := saga.RemoteResultOK(cmd, map[string]any{
				"inventory_restored": true,
			})
			return remoteResultOK, nil
		}),
	)
	// Disabled shipping to simulate failure and trigger compensation
	// cmdRouter.Register(
	// 	"cqrsify.examples.sagas.delivery.ship.cmd",
	// 	messaging.CommandHandlerWithReplyFn(func(ctx context.Context, cmd saga.RemotePayload) (messaging.CommandReply, error) {
	// 		remoteResultOK := saga.RemoteResultOK(cmd, map[string]any{
	// 			"shipment_id": "shipment-54321",
	// 		})
	// 		return remoteResultOK, nil
	// 	}),
	// )
	cmdRouter.Register(
		"cqrsify.examples.sagas.delivery.cancel.cmd",
		messaging.NewCommandHandlerWithReplyFn(func(ctx context.Context, cmd saga.RemotePayload) (messaging.CommandReply, error) {
			remoteResultOK := saga.RemoteResultOK(cmd, map[string]any{
				"delivery_canceled": true,
			})
			return remoteResultOK, nil
		}),
	)
	return cmdRouter
}
