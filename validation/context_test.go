package validation

import (
	"testing"

	"github.com/yaitoo/sparrow/micro"
)

var form = `
#Section 分组节点
[login]
"passwd"=[
	{Rule="Required",Message={zh_CN="passwd参数必填",en_US="passwd is required"}},  
	{Rule="MaxSize:6",Message={zh_CN="passwd参数长度要少于6位",en_US="passwd shouldn't be longer than 6 chars"}}
]

#栏位
    [[login.login]]
        Rule="Required"
        Message={zh_CN="login参数必填",en_US="login is required"}
    [[login.login]]      
        Rule="MinSize:6"
        Message={zh_CN="login参数长度至少要6位",en_US="login should be longer than 6 chars"}


`

func TestValidateGroup(t *testing.T) {

	vg := NewContext(micro.NewContext(), WithForms(form))

	vg.ValidateGroup("login").
		WithNames("login", "passwd").
		WithValues("testuser1234", "")

	if !vg.IsValid() {
		t.Fatal("IsValid failed:", vg.Errors())
	}

}

func TestValidate(t *testing.T) {
	vg := NewContext(micro.NewContext(), WithForms(form))

	vg.Validate("login", "passwd", "123")
	if vg.IsValid() {
		t.Fatal("test login passwd MinSize(6) failed")
	}
	vg.Clear()
	vg.Validate("login", "passwd", "123456")
	if !vg.IsValid() {
		t.Fatal("test login passwd MinSize(6) failed")
	}

	vg.Clear()
	vg.Validate("login", "passwd", "1234567890123")
	if vg.IsValid() {
		t.Fatal("test login passwd MaxSize(12) failed")
	}
	vg.Clear()
	vg.Validate("login", "passwd", "123456789012")
	if !vg.IsValid() {
		t.Fatal("test login passwd MaxSize(12) failed")
	}
}

func TestErrorLang(t *testing.T) {
	vg := NewContext(micro.NewContext())
	vg.Validate("login", "passwd", "123")
	if vg.IsValid() {
		t.Fatal("test login passwd MinSize(6) failed")
	}

	errors := vg.Errors()
	if len(errors) != 1 {
		t.Fatalf("test login passwd failed: error len %d != 1", len(errors))
	}
	if errors[0] != "passwd参数长度至少要6位" {
		t.Fatalf("Test default error lang failed")
	}

	vg = NewContext(micro.WithValues(micro.NewContext(), map[string]string{"lang": "en-US"}))
	vg.Validate("login", "passwd", "123")
	if vg.IsValid() {
		t.Fatal("test login passwd MinSize(6) failed")
	}

	errors = vg.Errors()
	if len(errors) != 1 {
		t.Fatalf("test login passwd failed: error len %d != 1", len(errors))
	}
	if errors[0] != "passwd should be longer than 6 chars" {
		t.Fatalf("Test error lang en-US failed")
	}
}
