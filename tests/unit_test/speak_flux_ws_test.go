// Copyright 2026 Deepgram SDK contributors. All Rights Reserved.
// Use of this source code is governed by a MIT license that can be found in the LICENSE file.
// SPDX-License-Identifier: MIT

package deepgram_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dvonthenen/websocket"

	speakmsg "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/speak/v2/websocket/interfaces"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	interfacesv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces/v2"
	speakv2 "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2"
	speakclientws "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/speak/v2/websocket"
)

const wireWaitTimeout = 5 * time.Second

// speakTestServerScript reacts to one received text frame on the mock server.
type speakTestServerScript func(conn *websocket.Conn, frame string)

// newFluxSpeakTestServer starts a TLS WebSocket server for /v2/speak wire tests.
// Every received text frame is recorded on frames and handed to script; pings are
// counted on pings.
func newFluxSpeakTestServer(t *testing.T, script speakTestServerScript) (host string, frames chan string, pings chan struct{}, shutdown func()) {
	t.Helper()

	frames = make(chan string, 32)
	pings = make(chan struct{}, 32)
	upgrader := websocket.Upgrader{}

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			t.Logf("test server upgrade failed: %s", err)
			return
		}
		defer conn.Close()

		conn.SetPingHandler(func(appData string) error {
			pings <- struct{}{}
			return conn.WriteControl(websocket.PongMessage, []byte(appData), time.Now().Add(time.Second))
		})

		for {
			msgType, msg, err := conn.ReadMessage()
			if err != nil {
				return // client closed the connection
			}
			if msgType != websocket.TextMessage {
				continue
			}
			frames <- string(msg)
			if script != nil {
				script(conn, string(msg))
			}
		}
	}))

	return strings.TrimPrefix(srv.URL, "https://"), frames, pings, srv.Close
}

// drainAndClose is the graceful-close script: on Close it simulates a long turn
// still synthesizing — final audio and SessionMetadata delayed well beyond the
// 100ms the old Stop() waited — then closes the socket like the live server.
func drainAndClose(delay time.Duration) speakTestServerScript {
	return func(conn *websocket.Conn, frame string) {
		switch {
		case strings.Contains(frame, `"Speak"`):
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"SpeechStarted","speech_id":"`+testSpeechID+`"}`))
			_ = conn.WriteMessage(websocket.BinaryMessage, []byte{0x01, 0x02, 0x03, 0x04})
		case strings.Contains(frame, `"Flush"`):
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"Flushed","speech_id":"`+testSpeechID+`"}`))
		case strings.Contains(frame, `"Close"`):
			time.Sleep(delay)
			_ = conn.WriteMessage(websocket.BinaryMessage, []byte{0x05, 0x06, 0x07, 0x08})
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"SpeechMetadata","speech_id":"`+testSpeechID+`","audio_duration_ms":1875,"input_character_count":47,"billable_character_count":47,"controls_applied":{"pronunciations_applied":0,"breaks_applied":0,"pronunciation_warnings":0}}`))
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"SessionMetadata","total_audio_duration_ms":1875,"total_input_character_count":47,"total_billable_character_count":47}`))
			_ = conn.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
		}
	}
}

// speakWireHandler is a FluxSpeakMessageCallback that records every event.
type speakWireHandler struct {
	connected   chan *speakmsg.ConnectedResponse
	binary      chan []byte
	started     chan *speakmsg.SpeechStartedResponse
	speechMeta  chan *speakmsg.SpeechMetadataResponse
	interrupted chan *speakmsg.SpeechInterruptedResponse
	flushed     chan *speakmsg.FlushedResponse
	sessionMeta chan *speakmsg.SessionMetadataResponse
	cfgSuccess  chan *speakmsg.ConfigureSuccessResponse
	cfgFailure  chan *speakmsg.ConfigureFailureResponse
	warnings    chan *speakmsg.WarningResponse
	fatals      chan *speakmsg.FatalErrorResponse
	unhandled   chan []byte
}

func newSpeakWireHandler() *speakWireHandler {
	return &speakWireHandler{
		connected:   make(chan *speakmsg.ConnectedResponse, 8),
		binary:      make(chan []byte, 64),
		started:     make(chan *speakmsg.SpeechStartedResponse, 8),
		speechMeta:  make(chan *speakmsg.SpeechMetadataResponse, 8),
		interrupted: make(chan *speakmsg.SpeechInterruptedResponse, 8),
		flushed:     make(chan *speakmsg.FlushedResponse, 8),
		sessionMeta: make(chan *speakmsg.SessionMetadataResponse, 8),
		cfgSuccess:  make(chan *speakmsg.ConfigureSuccessResponse, 8),
		cfgFailure:  make(chan *speakmsg.ConfigureFailureResponse, 8),
		warnings:    make(chan *speakmsg.WarningResponse, 8),
		fatals:      make(chan *speakmsg.FatalErrorResponse, 8),
		unhandled:   make(chan []byte, 8),
	}
}

func (h *speakWireHandler) Open(or *speakmsg.OpenResponse) error { return nil }
func (h *speakWireHandler) Connected(cr *speakmsg.ConnectedResponse) error {
	h.connected <- cr
	return nil
}
func (h *speakWireHandler) Binary(byMsg []byte) error {
	h.binary <- byMsg
	return nil
}
func (h *speakWireHandler) SpeechStarted(ss *speakmsg.SpeechStartedResponse) error {
	h.started <- ss
	return nil
}
func (h *speakWireHandler) SpeechMetadata(sm *speakmsg.SpeechMetadataResponse) error {
	h.speechMeta <- sm
	return nil
}
func (h *speakWireHandler) SpeechInterrupted(si *speakmsg.SpeechInterruptedResponse) error {
	h.interrupted <- si
	return nil
}
func (h *speakWireHandler) Flushed(fl *speakmsg.FlushedResponse) error {
	h.flushed <- fl
	return nil
}
func (h *speakWireHandler) SessionMetadata(sm *speakmsg.SessionMetadataResponse) error {
	h.sessionMeta <- sm
	return nil
}
func (h *speakWireHandler) ConfigureSuccess(cs *speakmsg.ConfigureSuccessResponse) error {
	h.cfgSuccess <- cs
	return nil
}
func (h *speakWireHandler) ConfigureFailure(cf *speakmsg.ConfigureFailureResponse) error {
	h.cfgFailure <- cf
	return nil
}
func (h *speakWireHandler) Warning(wr *speakmsg.WarningResponse) error {
	h.warnings <- wr
	return nil
}
func (h *speakWireHandler) FatalError(fe *speakmsg.FatalErrorResponse) error {
	h.fatals <- fe
	return nil
}
func (h *speakWireHandler) Close(cr *speakmsg.CloseResponse) error { return nil }
func (h *speakWireHandler) Error(er *speakmsg.ErrorResponse) error { return nil }
func (h *speakWireHandler) UnhandledEvent(byData []byte) error {
	h.unhandled <- byData
	return nil
}

func speakWireClientOptions(host string) *interfaces.ClientOptionsV2 {
	return &interfaces.ClientOptionsV2{
		Host:           host,
		SkipServerAuth: true, // the test server uses a self-signed certificate
	}
}

func speakWireOptions() *interfaces.SpeakV2WSOptions {
	return &interfaces.SpeakV2WSOptions{
		Model:      testFluxModel,
		Encoding:   "linear16",
		SampleRate: 48000,
	}
}

func awaitFrame(t *testing.T, frames chan string, want string) {
	t.Helper()
	select {
	case frame := <-frames:
		if frame != want {
			t.Errorf("frame mismatch:\n got %s\nwant %s", frame, want)
		}
	case <-time.After(wireWaitTimeout):
		t.Fatalf("server never received frame %s", want)
	}
}

// Test_FluxSpeakWireCallback drives the public callback client end to end against
// a mock WebSocket server (Greg's S3 on PR #351): every outbound method emits its
// exact wire frame, inbound text and binary events route to the handler in order,
// Finish drains a delayed turn (B1), and write errors propagate after shutdown.
func Test_FluxSpeakWireCallback(t *testing.T) {
	script := drainAndClose(300 * time.Millisecond)
	host, frames, _, shutdown := newFluxSpeakTestServer(t, func(conn *websocket.Conn, frame string) {
		switch {
		case strings.Contains(frame, `"Configure"`):
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"ConfigureSuccess","applied":{"speed":1.05}}`))
		case strings.Contains(frame, `"Interrupt"`):
			_ = conn.WriteMessage(websocket.TextMessage, []byte(`{"type":"SpeechInterrupted","audio_played_ms":2340,"metadata":{"speech_id":"`+testSpeechID+`","audio_duration_ms":4200,"input_character_count":47,"billable_character_count":47,"controls_applied":{"pronunciations_applied":0,"breaks_applied":0,"pronunciation_warnings":0}}}`))
		default:
			script(conn, frame)
		}
	})
	defer shutdown()

	handler := newSpeakWireHandler()
	dgClient, err := speakv2.NewWSUsingCallback(
		context.Background(), MockAPIKey, speakWireClientOptions(host), speakWireOptions(), handler)
	if err != nil {
		t.Fatalf("NewWSUsingCallback failed: %s", err)
	}
	if !dgClient.Connect() {
		t.Fatal("Connect to the test server failed")
	}

	// outbound wire formats, one frame per method
	if err := dgClient.Speak("Sure, I can help."); err != nil {
		t.Fatalf("Speak failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Speak","text":"Sure, I can help."}`)

	// inbound events triggered by Speak: SpeechStarted, then binary audio
	select {
	case <-handler.started:
	case <-time.After(wireWaitTimeout):
		t.Fatal("SpeechStarted never reached the handler")
	}
	select {
	case audio := <-handler.binary:
		if len(audio) != 4 {
			t.Errorf("expected the 4-byte audio frame, got %d bytes", len(audio))
		}
	case <-time.After(wireWaitTimeout):
		t.Fatal("binary audio never reached the handler")
	}

	if err := dgClient.Flush(); err != nil {
		t.Fatalf("Flush failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Flush"}`)
	select {
	case fl := <-handler.flushed:
		if fl.SpeechID != testSpeechID {
			t.Errorf("unexpected Flushed speech_id %q", fl.SpeechID)
		}
	case <-time.After(wireWaitTimeout):
		t.Fatal("Flushed never reached the handler")
	}

	if err := dgClient.Interrupt(); err != nil {
		t.Fatalf("Interrupt failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Interrupt"}`)
	select {
	case si := <-handler.interrupted:
		if si.AudioPlayedMs != 2340 {
			t.Errorf("unexpected SpeechInterrupted: %+v", si)
		}
	case <-time.After(wireWaitTimeout):
		t.Fatal("SpeechInterrupted never reached the handler")
	}

	if err := dgClient.InterruptWithOffset(2340); err != nil {
		t.Fatalf("InterruptWithOffset failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Interrupt","playback_offset":{"type":"time_ms","value":2340}}`)
	<-handler.interrupted

	if err := dgClient.Configure(&interfaces.SpeakV2ConfigureOptions{Speed: float64Ptr(1.05)}); err != nil {
		t.Fatalf("Configure failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Configure","speed":1.05}`)
	select {
	case cs := <-handler.cfgSuccess:
		if cs.Applied.Speed != 1.05 {
			t.Errorf("unexpected ConfigureSuccess: %+v", cs)
		}
	case <-time.After(wireWaitTimeout):
		t.Fatal("ConfigureSuccess never reached the handler")
	}

	// B1: graceful close — the server delays the turn's final audio and
	// SessionMetadata 300ms past the old fixed 100ms shutdown window
	finishCtx, cancel := context.WithTimeout(context.Background(), wireWaitTimeout)
	defer cancel()
	if err := dgClient.Finish(finishCtx); err != nil {
		t.Fatalf("Finish failed: %s", err)
	}
	awaitedTail := false
	for !awaitedTail {
		select {
		case audio := <-handler.binary:
			if len(audio) == 4 && audio[0] == 0x05 {
				awaitedTail = true
			}
		default:
			t.Fatal("the delayed final audio frame was not delivered before Finish returned")
		}
	}
	select {
	case <-handler.sessionMeta:
	default:
		t.Fatal("SessionMetadata was not delivered before Finish returned")
	}

	// write errors propagate after shutdown
	if err := dgClient.Speak("too late"); err == nil {
		t.Error("Speak after Finish must propagate a write error, got nil")
	}
}

// Test_FluxSpeakConfigureValidation covers Greg's S4 on PR #351: nil options and
// invalid speeds are rejected client-side; omitted speed cannot masquerade as an
// empty-but-successful update.
func Test_FluxSpeakConfigureValidation(t *testing.T) {
	host, frames, _, shutdown := newFluxSpeakTestServer(t, nil)
	defer shutdown()

	dgClient, err := speakv2.NewWSUsingCallback(
		context.Background(), MockAPIKey, speakWireClientOptions(host), speakWireOptions(), newSpeakWireHandler())
	if err != nil {
		t.Fatalf("NewWSUsingCallback failed: %s", err)
	}
	if !dgClient.Connect() {
		t.Fatal("Connect to the test server failed")
	}
	defer dgClient.Stop()

	t.Run("nil options return an error", func(t *testing.T) {
		if err := dgClient.Configure(nil); !errors.Is(err, interfacesv2.ErrOptionsRequired) {
			t.Errorf("expected ErrOptionsRequired, got: %v", err)
		}
	})

	t.Run("omitted speed is rejected, not sent as an empty update", func(t *testing.T) {
		err := dgClient.Configure(&interfaces.SpeakV2ConfigureOptions{})
		if !errors.Is(err, interfacesv2.ErrConfigureNoFields) {
			t.Errorf("expected ErrConfigureNoFields, got: %v", err)
		}
	})

	for _, invalid := range []float64{0, -1, 0.45, 1.55, 1.07, 3.5} {
		invalid := invalid
		t.Run("invalid speed is rejected", func(t *testing.T) {
			err := dgClient.Configure(&interfaces.SpeakV2ConfigureOptions{Speed: float64Ptr(invalid)})
			if !errors.Is(err, interfacesv2.ErrSpeedOutOfRange) {
				t.Errorf("speed %v: expected ErrSpeedOutOfRange, got: %v", invalid, err)
			}
		})
	}

	for _, valid := range []float64{0.5, 0.85, 1.0, 1.15, 1.5} {
		valid := valid
		t.Run("valid speed is sent", func(t *testing.T) {
			if err := dgClient.Configure(&interfaces.SpeakV2ConfigureOptions{Speed: float64Ptr(valid)}); err != nil {
				t.Fatalf("speed %v: Configure failed: %s", valid, err)
			}
			select {
			case frame := <-frames:
				if !strings.Contains(frame, `"Configure"`) {
					t.Errorf("expected a Configure frame, got %s", frame)
				}
			case <-time.After(wireWaitTimeout):
				t.Fatal("Configure frame never reached the server")
			}
		})
	}
}

// speakWireChanHandler is a FluxSpeakMessageChan with buffered channels for the
// events the channel-client wire test asserts on.
type speakWireChanHandler struct {
	binaryCh      chan *[]byte
	flushedCh     chan *speakmsg.FlushedResponse
	sessionMetaCh chan *speakmsg.SessionMetadataResponse
	warningCh     chan *speakmsg.WarningResponse
}

func newSpeakWireChanHandler() *speakWireChanHandler {
	return &speakWireChanHandler{
		binaryCh:      make(chan *[]byte, 64),
		flushedCh:     make(chan *speakmsg.FlushedResponse, 8),
		sessionMetaCh: make(chan *speakmsg.SessionMetadataResponse, 8),
		warningCh:     make(chan *speakmsg.WarningResponse, 8),
	}
}

func (h *speakWireChanHandler) GetOpen() []*chan *speakmsg.OpenResponse           { return nil }
func (h *speakWireChanHandler) GetConnected() []*chan *speakmsg.ConnectedResponse { return nil }
func (h *speakWireChanHandler) GetBinary() []*chan *[]byte {
	return []*chan *[]byte{&h.binaryCh}
}
func (h *speakWireChanHandler) GetSpeechStarted() []*chan *speakmsg.SpeechStartedResponse {
	return nil
}
func (h *speakWireChanHandler) GetSpeechMetadata() []*chan *speakmsg.SpeechMetadataResponse {
	return nil
}
func (h *speakWireChanHandler) GetSpeechInterrupted() []*chan *speakmsg.SpeechInterruptedResponse {
	return nil
}
func (h *speakWireChanHandler) GetFlushed() []*chan *speakmsg.FlushedResponse {
	return []*chan *speakmsg.FlushedResponse{&h.flushedCh}
}
func (h *speakWireChanHandler) GetSessionMetadata() []*chan *speakmsg.SessionMetadataResponse {
	return []*chan *speakmsg.SessionMetadataResponse{&h.sessionMetaCh}
}
func (h *speakWireChanHandler) GetConfigureSuccess() []*chan *speakmsg.ConfigureSuccessResponse {
	return nil
}
func (h *speakWireChanHandler) GetConfigureFailure() []*chan *speakmsg.ConfigureFailureResponse {
	return nil
}
func (h *speakWireChanHandler) GetWarning() []*chan *speakmsg.WarningResponse {
	return []*chan *speakmsg.WarningResponse{&h.warningCh}
}
func (h *speakWireChanHandler) GetFatalError() []*chan *speakmsg.FatalErrorResponse { return nil }
func (h *speakWireChanHandler) GetClose() []*chan *speakmsg.CloseResponse           { return nil }
func (h *speakWireChanHandler) GetError() []*chan *speakmsg.ErrorResponse           { return nil }
func (h *speakWireChanHandler) GetUnhandled() []*chan *[]byte                       { return nil }

// Test_FluxSpeakWireChannel drives the public channel client end to end: outbound
// methods emit their exact frames, binary and text events land on the handler
// channels, and Finish drains the delayed turn before returning (B1).
func Test_FluxSpeakWireChannel(t *testing.T) {
	host, frames, _, shutdown := newFluxSpeakTestServer(t, drainAndClose(300*time.Millisecond))
	defer shutdown()

	handler := newSpeakWireChanHandler()
	dgClient, err := speakv2.NewWSUsingChan(
		context.Background(), MockAPIKey, speakWireClientOptions(host), speakWireOptions(), handler)
	if err != nil {
		t.Fatalf("NewWSUsingChan failed: %s", err)
	}
	if !dgClient.Connect() {
		t.Fatal("Connect to the test server failed")
	}

	if err := dgClient.Speak("Sure, I can help."); err != nil {
		t.Fatalf("Speak failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Speak","text":"Sure, I can help."}`)
	select {
	case audio := <-handler.binaryCh:
		if len(*audio) != 4 {
			t.Errorf("expected the 4-byte audio frame, got %d bytes", len(*audio))
		}
	case <-time.After(wireWaitTimeout):
		t.Fatal("binary audio never reached the handler channel")
	}

	if err := dgClient.Flush(); err != nil {
		t.Fatalf("Flush failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Flush"}`)
	select {
	case <-handler.flushedCh:
	case <-time.After(wireWaitTimeout):
		t.Fatal("Flushed never reached the handler channel")
	}

	if err := dgClient.InterruptWithOffset(1200); err != nil {
		t.Fatalf("InterruptWithOffset failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Interrupt","playback_offset":{"type":"time_ms","value":1200}}`)

	if err := dgClient.Configure(&interfaces.SpeakV2ConfigureOptions{Speed: float64Ptr(0.9)}); err != nil {
		t.Fatalf("Configure failed: %s", err)
	}
	awaitFrame(t, frames, `{"type":"Configure","speed":0.9}`)

	finishCtx, cancel := context.WithTimeout(context.Background(), wireWaitTimeout)
	defer cancel()
	if err := dgClient.Finish(finishCtx); err != nil {
		t.Fatalf("Finish failed: %s", err)
	}
	// the delayed tail audio and SessionMetadata were dispatched before Finish
	// returned; the buffered handler channels must already hold them
	select {
	case audio := <-handler.binaryCh:
		if (*audio)[0] != 0x05 {
			t.Errorf("expected the delayed tail frame, got % x", *audio)
		}
	default:
		t.Fatal("the delayed final audio frame was not dispatched before Finish returned")
	}
	select {
	case <-handler.sessionMetaCh:
	default:
		t.Fatal("SessionMetadata was not dispatched before Finish returned")
	}

	if err := dgClient.Speak("too late"); err == nil {
		t.Error("Speak after Finish must propagate a write error, got nil")
	}
}

// Test_FluxSpeakFinishEdgeCases covers the bounded-close contract: ctx expiry when
// the server never drains, and a peer close without the final SessionMetadata.
func Test_FluxSpeakFinishEdgeCases(t *testing.T) {
	t.Run("Finish times out when the server never drains", func(t *testing.T) {
		host, _, _, shutdown := newFluxSpeakTestServer(t, nil) // ignores Close
		defer shutdown()

		dgClient, err := speakv2.NewWSUsingCallback(
			context.Background(), MockAPIKey, speakWireClientOptions(host), speakWireOptions(), newSpeakWireHandler())
		if err != nil {
			t.Fatalf("NewWSUsingCallback failed: %s", err)
		}
		if !dgClient.Connect() {
			t.Fatal("Connect to the test server failed")
		}

		finishCtx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancel()
		if err := dgClient.Finish(finishCtx); !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("expected context.DeadlineExceeded, got: %v", err)
		}
	})

	t.Run("peer close without SessionMetadata is reported", func(t *testing.T) {
		host, _, _, shutdown := newFluxSpeakTestServer(t, func(conn *websocket.Conn, frame string) {
			if strings.Contains(frame, `"Close"`) {
				_ = conn.WriteControl(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now().Add(time.Second))
			}
		})
		defer shutdown()

		dgClient, err := speakv2.NewWSUsingCallback(
			context.Background(), MockAPIKey, speakWireClientOptions(host), speakWireOptions(), newSpeakWireHandler())
		if err != nil {
			t.Fatalf("NewWSUsingCallback failed: %s", err)
		}
		if !dgClient.Connect() {
			t.Fatal("Connect to the test server failed")
		}

		finishCtx, cancel := context.WithTimeout(context.Background(), wireWaitTimeout)
		defer cancel()
		if err := dgClient.Finish(finishCtx); !errors.Is(err, speakclientws.ErrGracefulCloseIncomplete) {
			t.Errorf("expected ErrGracefulCloseIncomplete, got: %v", err)
		}
	})
}

// Test_FluxSpeakKeepAlive covers Greg's B2 on PR #351: the default keepalive
// interval leaves a comfortable margin inside the server's 60s NET-0004 idle
// deadline, and an enabled client actually sends protocol pings on its interval
// (verified with a short interval so the test stays fast).
func Test_FluxSpeakKeepAlive(t *testing.T) {
	t.Run("default interval is at most half the server idle timeout", func(t *testing.T) {
		if speakclientws.DefaultKeepAlivePeriod > speakclientws.ServerIdleTimeout/2 {
			t.Errorf("DefaultKeepAlivePeriod %v must be at most half of the %v idle timeout",
				speakclientws.DefaultKeepAlivePeriod, speakclientws.ServerIdleTimeout)
		}
	})

	t.Run("pings are sent on the configured interval", func(t *testing.T) {
		host, _, pings, shutdown := newFluxSpeakTestServer(t, nil)
		defer shutdown()

		cOptions := speakWireClientOptions(host)
		cOptions.EnableKeepAlive = true
		cOptions.KeepAlivePeriod = 50 * time.Millisecond

		dgClient, err := speakv2.NewWSUsingCallback(
			context.Background(), MockAPIKey, cOptions, speakWireOptions(), newSpeakWireHandler())
		if err != nil {
			t.Fatalf("NewWSUsingCallback failed: %s", err)
		}
		if !dgClient.Connect() {
			t.Fatal("Connect to the test server failed")
		}
		defer dgClient.Stop()

		for i := 0; i < 2; i++ {
			select {
			case <-pings:
			case <-time.After(2 * time.Second):
				t.Fatalf("ping %d never reached the server", i+1)
			}
		}
	})
}

// Test_FluxSpeakFinishAfterReconnect verifies the graceful-close milestones are
// per-session: after a completed Finish and an AttemptReconnect, a second Finish
// waits for the NEW session's drain instead of returning instantly on the old
// session's already-fired milestones.
func Test_FluxSpeakFinishAfterReconnect(t *testing.T) {
	host, _, _, shutdown := newFluxSpeakTestServer(t, drainAndClose(150*time.Millisecond))
	defer shutdown()

	handler := newSpeakWireHandler()
	dgClient, err := speakv2.NewWSUsingCallback(
		context.Background(), MockAPIKey, speakWireClientOptions(host), speakWireOptions(), handler)
	if err != nil {
		t.Fatalf("NewWSUsingCallback failed: %s", err)
	}
	if !dgClient.Connect() {
		t.Fatal("Connect to the test server failed")
	}

	finishCtx, cancel := context.WithTimeout(context.Background(), wireWaitTimeout)
	defer cancel()
	if err := dgClient.Finish(finishCtx); err != nil {
		t.Fatalf("first Finish failed: %s", err)
	}
	<-handler.sessionMeta // drain the first session's metadata record

	if !dgClient.AttemptReconnect(context.Background(), 3) {
		t.Fatal("AttemptReconnect failed")
	}

	start := time.Now()
	finishCtx2, cancel2 := context.WithTimeout(context.Background(), wireWaitTimeout)
	defer cancel2()
	if err := dgClient.Finish(finishCtx2); err != nil {
		t.Fatalf("Finish after reconnect failed: %s", err)
	}
	// the server delays the second session's drain 150ms; an instant return would
	// mean Finish consumed the FIRST session's milestones
	if elapsed := time.Since(start); elapsed < 150*time.Millisecond {
		t.Errorf("Finish returned in %v — it did not wait for the new session's drain", elapsed)
	}
	select {
	case <-handler.sessionMeta:
	default:
		t.Fatal("the second session's SessionMetadata was not delivered before Finish returned")
	}
}

// Test_FluxSpeakChannelFinishAfterContextReplacement ensures the channel
// observer follows the context installed by ConnectWithCancel, rather than
// stopping with the construction context.
func Test_FluxSpeakChannelFinishAfterContextReplacement(t *testing.T) {
	host, _, _, shutdown := newFluxSpeakTestServer(t, drainAndClose(100*time.Millisecond))
	defer shutdown()

	constructionCtx, constructionCancel := context.WithCancel(context.Background())
	defer constructionCancel()
	handler := newSpeakWireChanHandler()
	dgClient, err := speakv2.NewWSUsingChan(
		constructionCtx, MockAPIKey, speakWireClientOptions(host), speakWireOptions(), handler)
	if err != nil {
		t.Fatalf("NewWSUsingChan failed: %s", err)
	}

	activeCtx, activeCancel := context.WithCancel(context.Background())
	defer activeCancel()
	if !dgClient.ConnectWithCancel(activeCtx, activeCancel, 3) {
		t.Fatal("ConnectWithCancel to the test server failed")
	}
	constructionCancel()

	finishCtx, cancel := context.WithTimeout(context.Background(), wireWaitTimeout)
	defer cancel()
	if err := dgClient.Finish(finishCtx); err != nil {
		t.Fatalf("Finish after context replacement failed: %s", err)
	}
	select {
	case <-handler.sessionMetaCh:
	default:
		t.Fatal("SessionMetadata was not dispatched before Finish returned")
	}
}

// Test_FluxSpeakChannelFinishAfterReconnectContextReplacement ensures a
// reconnect starts an observer for its replacement context and ignores the
// completed session's queued observer events.
func Test_FluxSpeakChannelFinishAfterReconnectContextReplacement(t *testing.T) {
	host, _, _, shutdown := newFluxSpeakTestServer(t, drainAndClose(100*time.Millisecond))
	defer shutdown()

	constructionCtx, constructionCancel := context.WithCancel(context.Background())
	defer constructionCancel()
	handler := newSpeakWireChanHandler()
	dgClient, err := speakv2.NewWSUsingChan(
		constructionCtx, MockAPIKey, speakWireClientOptions(host), speakWireOptions(), handler)
	if err != nil {
		t.Fatalf("NewWSUsingChan failed: %s", err)
	}
	if !dgClient.Connect() {
		t.Fatal("Connect failed")
	}

	finishCtx, cancel := context.WithTimeout(context.Background(), wireWaitTimeout)
	defer cancel()
	if err := dgClient.Finish(finishCtx); err != nil {
		t.Fatalf("first Finish failed: %s", err)
	}
	<-handler.sessionMetaCh
	constructionCancel()

	activeCtx, activeCancel := context.WithCancel(context.Background())
	defer activeCancel()
	if !dgClient.AttemptReconnectWithCancel(activeCtx, activeCancel, 3) {
		t.Fatal("AttemptReconnectWithCancel failed")
	}

	finishCtx2, cancel2 := context.WithTimeout(context.Background(), wireWaitTimeout)
	defer cancel2()
	if err := dgClient.Finish(finishCtx2); err != nil {
		t.Fatalf("Finish after reconnect context replacement failed: %s", err)
	}
	select {
	case <-handler.sessionMetaCh:
	default:
		t.Fatal("second SessionMetadata was not dispatched before Finish returned")
	}
}
