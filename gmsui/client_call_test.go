package gmsui_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/W3Tools/go-modules/gmsui"
	"github.com/W3Tools/go-modules/gmsui/client"
)

type clock struct {
	Id          gmsui.SuiMoveId `json:"id"`
	TimestampMs string          `json:"timestamp_ms"`
}

func TestClientCall_GetObjectAndUnmarshal(t *testing.T) {
	suiClient, err := client.NewSuiClient(context.Background(), client.GetFullNodeURL("mainnet"))
	if err != nil {
		t.Fatalf("failed to new sui client, msg: %v", err)
	}

	_, clockObject, err := gmsui.GetObjectAndUnmarshal[clock](suiClient, "0x6")
	if err != nil {
		t.Fatalf("failed to get object and unmarshal, msg: %v", err)
	}

	if reflect.DeepEqual(clockObject.Id.Id, "") {
		t.Errorf("expected id not none, but got %s", clockObject.Id.Id)
	}

	if reflect.DeepEqual(clockObject.TimestampMs, "") {
		t.Errorf("expected timestamp not none, but got %s", clockObject.TimestampMs)
	}

	_, _, err = gmsui.GetObjectAndUnmarshal[clock](suiClient, "0x2")
	if err == nil {
		t.Fatalf("unable to get 0x2 package, expected package but got unknown")
	}
}

func TestClientCall_GetObjectsAndUnmarshal(t *testing.T) {
	suiClient, err := client.NewSuiClient(context.Background(), client.GetFullNodeURL("mainnet"))
	if err != nil {
		t.Fatalf("failed to new sui client, msg: %v", err)
	}

	_, clockObjects, err := gmsui.GetObjectsAndUnmarshal[clock](suiClient, []string{"0x6", "0x6"})
	if err != nil {
		t.Fatalf("failed to get objects and unmarshal, msg: %v", err)
	}

	for _, obj := range clockObjects {
		if reflect.DeepEqual(obj.Id.Id, "") {
			t.Errorf("expected id not none, but got %s", obj.Id.Id)
		}

		if reflect.DeepEqual(obj.TimestampMs, "") {
			t.Errorf("expected timestamp not none, but got %s", obj.TimestampMs)
		}
	}
}

func TestClientCall_GetDynamicFieldObjectAndUnmarshal(t *testing.T) {
}
