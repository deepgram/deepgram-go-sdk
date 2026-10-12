// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"errors"
	"testing"
	"time"

	listenv2api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket"
	msginterfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v2/websocket/interfaces"
)

const mockFluxWarning = `{"type":"Warning","code":"FORCE_END_TURN_NO_ACTIVE_TURN","description":"ForceEndTurn received with no active turn"}`

// baseFluxCallback implements msginterfaces.FluxMessageCallback with no-op methods
// and records UnhandledEvent payloads. It deliberately does NOT implement the
// optional FluxWarningCallback extension.
type baseFluxCallback struct {
	unhandled [][]byte
}

func (c *baseFluxCallback) Open(or *msginterfaces.OpenResponse) error           { return nil }
func (c *baseFluxCallback) Connected(cr *msginterfaces.ConnectedResponse) error { return nil }
func (c *baseFluxCallback) TurnInfo(tr *msginterfaces.TurnInfoResponse) error   { return nil }
func (c *baseFluxCallback) ConfigureSuccess(cs *msginterfaces.ConfigureSuccessResponse) error {
	return nil
}
func (c *baseFluxCallback) ConfigureFailure(cf *msginterfaces.ConfigureFailureResponse) error {
	return nil
}
func (c *baseFluxCallback) FatalError(fe *msginterfaces.FatalErrorResponse) error { return nil }
func (c *baseFluxCallback) Close(cr *msginterfaces.CloseResponse) error           { return nil }
func (c *baseFluxCallback) Error(er *msginterfaces.ErrorResponse) error           { return nil }
func (c *baseFluxCallback) UnhandledEvent(byData []byte) error {
	c.unhandled = append(c.unhandled, byData)
	return nil
}

// warnFluxCallback additionally implements the FluxWarningCallback extension.
type warnFluxCallback struct {
	baseFluxCallback
	warnings []*msginterfaces.WarningResponse
}

func (c *warnFluxCallback) Warning(wr *msginterfaces.WarningResponse) error {
	c.warnings = append(c.warnings, wr)
	return nil
}

// Test_FluxCallbackRouterWarning verifies the callback router treats server Warning
// messages as routine, non-fatal events (Greg's B1 on PR #350): a warning-aware
// handler receives them via Warning(), a plain handler receives them via
// UnhandledEvent(), and neither path reports ErrInvalidMessageType.
func Test_FluxCallbackRouterWarning(t *testing.T) {
	t.Run("handler implementing FluxWarningCallback receives the warning", func(t *testing.T) {
		handler := &warnFluxCallback{}
		router := listenv2api.NewCallbackRouter(handler)

		if err := router.Message([]byte(mockFluxWarning)); err != nil {
			t.Fatalf("Warning must be non-fatal, got error: %s", err)
		}
		if len(handler.warnings) != 1 {
			t.Fatalf("expected 1 warning delivered, got %d", len(handler.warnings))
		}
		wr := handler.warnings[0]
		if wr.Code != msginterfaces.WarningCodeForceEndTurnNoActiveTurn {
			t.Errorf("expected code %q, got %q", msginterfaces.WarningCodeForceEndTurnNoActiveTurn, wr.Code)
		}
		if wr.Description == "" {
			t.Error("expected a non-empty warning description")
		}
		if len(handler.unhandled) != 0 {
			t.Errorf("warning must not also hit UnhandledEvent, got %d calls", len(handler.unhandled))
		}
	})

	t.Run("handler without FluxWarningCallback gets UnhandledEvent, no error", func(t *testing.T) {
		handler := &baseFluxCallback{}
		router := listenv2api.NewCallbackRouter(handler)

		err := router.Message([]byte(mockFluxWarning))
		if err != nil {
			t.Fatalf("Warning must be non-fatal even without a Warning handler, got: %s", err)
		}
		if len(handler.unhandled) != 1 {
			t.Fatalf("expected warning forwarded to UnhandledEvent once, got %d calls", len(handler.unhandled))
		}
		if string(handler.unhandled[0]) != mockFluxWarning {
			t.Errorf("UnhandledEvent got altered payload: %s", handler.unhandled[0])
		}
	})

	t.Run("truly unknown message types still return ErrInvalidMessageType", func(t *testing.T) {
		handler := &baseFluxCallback{}
		router := listenv2api.NewCallbackRouter(handler)

		err := router.Message([]byte(`{"type":"SomeFutureEvent"}`))
		if !errors.Is(err, listenv2api.ErrInvalidMessageType) {
			t.Errorf("expected ErrInvalidMessageType for unknown types, got: %v", err)
		}
	})
}

// baseFluxChan implements msginterfaces.FluxMessageChan with a single unhandled
// channel and nil for everything else. It deliberately does NOT implement the
// optional FluxWarningChan extension.
type baseFluxChan struct {
	unhandledChan chan *[]byte
}

func newBaseFluxChan() *baseFluxChan {
	return &baseFluxChan{unhandledChan: make(chan *[]byte, 4)}
}

func (h *baseFluxChan) GetOpen() []*chan *msginterfaces.OpenResponse           { return nil }
func (h *baseFluxChan) GetConnected() []*chan *msginterfaces.ConnectedResponse { return nil }
func (h *baseFluxChan) GetTurnInfo() []*chan *msginterfaces.TurnInfoResponse   { return nil }
func (h *baseFluxChan) GetConfigureSuccess() []*chan *msginterfaces.ConfigureSuccessResponse {
	return nil
}
func (h *baseFluxChan) GetConfigureFailure() []*chan *msginterfaces.ConfigureFailureResponse {
	return nil
}
func (h *baseFluxChan) GetFatalError() []*chan *msginterfaces.FatalErrorResponse { return nil }
func (h *baseFluxChan) GetClose() []*chan *msginterfaces.CloseResponse           { return nil }
func (h *baseFluxChan) GetError() []*chan *msginterfaces.ErrorResponse           { return nil }
func (h *baseFluxChan) GetUnhandled() []*chan *[]byte {
	return []*chan *[]byte{&h.unhandledChan}
}

// warnFluxChan additionally implements the FluxWarningChan extension.
type warnFluxChan struct {
	baseFluxChan
	warningChan chan *msginterfaces.WarningResponse
}

func newWarnFluxChan() *warnFluxChan {
	return &warnFluxChan{
		baseFluxChan: *newBaseFluxChan(),
		warningChan:  make(chan *msginterfaces.WarningResponse, 4),
	}
}

func (h *warnFluxChan) GetWarning() []*chan *msginterfaces.WarningResponse {
	return []*chan *msginterfaces.WarningResponse{&h.warningChan}
}

// Test_FluxChanRouterWarning is the channel-router counterpart of
// Test_FluxCallbackRouterWarning.
func Test_FluxChanRouterWarning(t *testing.T) {
	t.Run("handler implementing FluxWarningChan receives the warning", func(t *testing.T) {
		handler := newWarnFluxChan()
		router := listenv2api.NewChanRouter(handler)

		if err := router.Message([]byte(mockFluxWarning)); err != nil {
			t.Fatalf("Warning must be non-fatal, got error: %s", err)
		}

		select {
		case wr := <-handler.warningChan:
			if wr.Code != msginterfaces.WarningCodeForceEndTurnNoActiveTurn {
				t.Errorf("expected code %q, got %q", msginterfaces.WarningCodeForceEndTurnNoActiveTurn, wr.Code)
			}
		case <-time.After(time.Second):
			t.Fatal("no warning delivered on the warning channel")
		}

		select {
		case byData := <-handler.unhandledChan:
			t.Errorf("warning must not also hit the unhandled channel, got: %s", *byData)
		default:
		}
	})

	t.Run("handler without FluxWarningChan gets the raw bytes on unhandled, no error", func(t *testing.T) {
		handler := newBaseFluxChan()
		router := listenv2api.NewChanRouter(handler)

		if err := router.Message([]byte(mockFluxWarning)); err != nil {
			t.Fatalf("Warning must be non-fatal even without warning channels, got: %s", err)
		}

		select {
		case byData := <-handler.unhandledChan:
			if string(*byData) != mockFluxWarning {
				t.Errorf("unhandled channel got altered payload: %s", *byData)
			}
		case <-time.After(time.Second):
			t.Fatal("warning was not forwarded to the unhandled channel")
		}
	})

	t.Run("truly unknown message types still return ErrInvalidMessageType", func(t *testing.T) {
		handler := newBaseFluxChan()
		router := listenv2api.NewChanRouter(handler)

		err := router.Message([]byte(`{"type":"SomeFutureEvent"}`))
		if !errors.Is(err, listenv2api.ErrInvalidMessageType) {
			t.Errorf("expected ErrInvalidMessageType for unknown types, got: %v", err)
		}
		// drain the copy forwarded to the unhandled channel
		select {
		case <-handler.unhandledChan:
		case <-time.After(time.Second):
			t.Fatal("unknown message was not forwarded to the unhandled channel")
		}
	})
}
