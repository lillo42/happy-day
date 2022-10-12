//go:build wireinject
// +build wireinject

package main

import (
	"fmt"
	"happy_day/apis"
	"happy_day/application"
	"happy_day/infrastructure"
	"reflect"

	"github.com/google/uuid"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	uuidType    = reflect.TypeOf(uuid.UUID{})
	uuidSubtype = byte(0x04)

	ProviderSet = wire.NewSet(
		ProvideMongoDbOptions,

		apis.ProvideSet,
		application.ProvideSet,
		infrastructure.ProvideSet,
	)
)

func ProvideMongoDbOptions() *options.ClientOptions {
	connectionString := viper.GetString("connectingStrings.mongo")

	return options.Client().
		ApplyURI(connectionString).
		SetRegistry(bson.NewRegistryBuilder().
			RegisterTypeEncoder(uuidType, bsoncodec.ValueEncoderFunc(uuidEncodeValue)).
			RegisterTypeDecoder(uuidType, bsoncodec.ValueDecoderFunc(uuidDecodeValue)).
			Build())

}

func uuidEncodeValue(_ bsoncodec.EncodeContext, writer bsonrw.ValueWriter, value reflect.Value) error {
	if !value.IsValid() || value.Type() != uuidType {
		return bsoncodec.ValueEncoderError{
			Name:     "uuidEncodeValue",
			Types:    []reflect.Type{uuidType},
			Received: value,
		}
	}

	b := value.Interface().(uuid.UUID)
	return writer.WriteBinaryWithSubtype(b[:], uuidSubtype)
}

func uuidDecodeValue(_ bsoncodec.DecodeContext, reader bsonrw.ValueReader, value reflect.Value) error {
	if !value.CanSet() || value.Type() != uuidType {
		return bsoncodec.ValueDecoderError{
			Name:     "uuidDecodeValue",
			Types:    []reflect.Type{uuidType},
			Received: value,
		}
	}

	var data []byte
	var subtype byte
	var err error

	switch valueType := reader.Type(); valueType {
	case bsontype.Binary:
		data, subtype, err = reader.ReadBinary()
		if subtype != uuidSubtype {
			return fmt.Errorf("unsupported binary subtype %v for UUID", subtype)
		}
	case bsontype.Null:
		err = reader.ReadNull()
	case bsontype.Undefined:
		err = reader.ReadUndefined()
	default:
		return fmt.Errorf("cannot decode %v into a UUID", valueType)

	}

	if err != nil {
		return err
	}

	uuid2, err := uuid.FromBytes(data)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(uuid2))
	return nil
}

func initializeReservationController() apis.ReservationController {
	wire.Build(ProviderSet)
	return apis.ReservationController{}
}

func initializeProductController() apis.ProductController {
	wire.Build(ProvideSet)
	return apis.ProductController{}
}
