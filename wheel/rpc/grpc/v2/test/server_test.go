package test

import (
	"fmt"
	"github.com/zhuanxuhit/go-in-practice/wheel/rpc/grpc/v2/method"
	"reflect"
	"testing"
)

func TestRegister(t *testing.T) {
	arith := new(method.Arith)
	var (
		typ   = reflect.TypeOf(arith)
		rcvr  = reflect.ValueOf(arith)
		sname = reflect.Indirect(rcvr).Type().Name()
	)
	//*method.Arith 0xc000016278 0 Arith ptr
	fmt.Printf("%+v %+v %+v %+v %+v", typ, rcvr, rcvr.Elem(), sname, typ.Kind())
}
