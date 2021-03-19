package ioc

import (
	"context"
	"testing"
)

func TestRegister(t *testing.T) {

	Register(context.TODO(), "TestRegister", func(ctx context.Context) interface{} {
		return &MyObject{}
	})

	instance, err := Resolve(context.TODO(), "TestRegister")
	if err != nil {
		t.Error(err)
	}

	_, ok := instance.(*MyObject)
	if !ok {
		t.Error("Type is incorrect")
	}
}

func TestDuplicatedRegister(t *testing.T) {

	Register(context.TODO(), "TestDuplicatedRegister", func(ctx context.Context) interface{} {
		return &MyObject{}
	})

	err := Register(context.TODO(), "TestDuplicatedRegister", func(ctx context.Context) interface{} {
		return &MyObject{}
	})

	if err != ErrHasRegistered {
		t.Error("duplicated register should be blocked")
	}
}

func TestNotegister(t *testing.T) {

	_, err := Resolve(context.TODO(), "TestNotegister")

	if err != ErrNoRegistered {
		t.Error("ErrNoRegistered should be thrown.")
	}
}

func TestRegisterInstance(t *testing.T) {
	obj := &MyObject{
		Name: "myobject",
	}

	RegisterInstance(context.TODO(), "TestRegisterInstance", obj)

	instance, err := ResolveInstance(context.TODO(), "TestRegisterInstance")
	if err != nil {
		t.Error(err)
	}

	if instance != obj {
		t.Error("instance is not same")
	}
	resolvedObj, ok := instance.(*MyObject)
	if !ok {
		t.Error("Type is incorrect")
	}
	if resolvedObj.Name != obj.Name {
		t.Error("Name is not same")
	}
}

func TestDuplicatedRegisterInstance(t *testing.T) {

	obj := &MyObject{
		Name: "myobject",
	}

	RegisterInstance(context.TODO(), "TestDuplicatedRegisterInstance", obj)

	err := RegisterInstance(context.TODO(), "TestDuplicatedRegisterInstance", obj)
	if err != ErrHasRegistered {
		t.Error("ErrHasRegistered should be thrown")
	}

}

func TestNotRegisterInstance(t *testing.T) {

	_, err := ResolveInstance(context.TODO(), "TestNotRegisterInstance")
	if err != ErrNoRegistered {
		t.Error("ErrNoRegistered should be thrown.")
	}

}

type MyObject struct {
	Name string
}
