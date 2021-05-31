package fastlystackdriver

import (
	"reflect"

	"go.uber.org/zap"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
)

type valuer func(v reflect.Value) *monitoringpb.TypedValue

var int64Valuer valuer = func(v reflect.Value) *monitoringpb.TypedValue {
	var vv int64
	switch v.Kind() {
	case reflect.Int:
		vv = int64(v.Interface().(int))
	case reflect.Int8:
		vv = int64(v.Interface().(int8))
	case reflect.Int16:
		vv = int64(v.Interface().(int16))
	case reflect.Int32:
		vv = int64(v.Interface().(int32))
	case reflect.Int64:
		vv = int64(v.Interface().(int64))
	case reflect.Uint:
		vv = int64(v.Interface().(uint))
	case reflect.Uint8:
		vv = int64(v.Interface().(uint8))
	case reflect.Uint16:
		vv = int64(v.Interface().(uint16))
	case reflect.Uint32:
		vv = int64(v.Interface().(uint32))
	case reflect.Uint64:
		vv = int64(v.Interface().(uint64))
	default:
		zap.S().Warnf("bad value for 'int64' with kind %v", v.Kind())
		return nil
	}
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_Int64Value{
			Int64Value: vv,
		},
	}
}

var doubleValuer valuer = func(v reflect.Value) *monitoringpb.TypedValue {
	var vv float64
	switch v.Kind() {
	case reflect.Float32:
		vv = float64(v.Interface().(float32))
	case reflect.Float64:
		vv = float64(v.Interface().(float64))
	default:
		zap.S().Warnf("bad value for 'double' with kind %v", v.Kind())
		return nil
	}
	return &monitoringpb.TypedValue{
		Value: &monitoringpb.TypedValue_DoubleValue{
			DoubleValue: vv,
		},
	}
}
