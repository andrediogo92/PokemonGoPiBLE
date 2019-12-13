package main

import (
	"PokemonGoPiBLE/examples/lib"
	"context"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"time"

	"github.com/go-ble/ble"
	"github.com/go-ble/ble/linux"
)

var (
	device = flag.String("device", "default", "implementation of ble")
	du     = flag.Duration("du", 5*time.Second, "advertising duration, 0 for indefinitely")
)

func main() {
	flag.Parse()

	d, err := NewDevice(*device)
	if err != nil {
		log.Fatalf("can't new device : %s", err)
	}
	ble.SetDefaultDevice(d)

	testSvc := ble.NewService(lib.EchoCharUUID)
	testSvc.AddCharacteristic(lib.NewEchoChar())

	if err := ble.AddService(testSvc); err != nil {
		log.Fatalf("can't add service: %s", err)
	}

	// Advertise for specified durantion, or until interrupted by user.
	fmt.Printf("Advertising for %s...\n", *du)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), *du))
	chkErr(ble.AdvertiseNameAndServices(ctx, "Gopher", testSvc.UUID))
}

func chkErr(err error) {
	switch errors.Cause(err) {
	case nil:
	case context.DeadlineExceeded:
		fmt.Printf("done\n")
	case context.Canceled:
		fmt.Printf("canceled\n")
	default:
		log.Fatalf(err.Error())
	}
}

// DefaultDevice ...
func DefaultDevice(opts ...ble.Option) (d ble.Device, err error) {
	return linux.NewDevice(opts...)
}

// NewDevice ...
func NewDevice(impl string, opts ...ble.Option) (d ble.Device, err error) {
	return DefaultDevice(opts...)
}
