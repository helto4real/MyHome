package cast

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	"golang.org/x/net/context"

	"github.com/helto4real/MyHome/core/net/cast/event"
	"github.com/ninjasphere/go-castv2/api"
)

const interval = time.Second * 5
const maxBacklog = 3

type HeartbeatController struct {
	pongs    int64
	ticker   *time.Ticker
	channel  *Channel
	eventsCh chan Event
}

var ping = PayloadHeaders{Type: "PING"}
var pong = PayloadHeaders{Type: "PONG"}

func NewHeartbeatController(conn *TlsConnection, eventsCh chan Event, sourceId, destinationId string) *HeartbeatController {
	controller := &HeartbeatController{
		channel:  conn.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.cast.tp.heartbeat"),
		eventsCh: eventsCh,
	}

	controller.channel.OnMessage("PING", controller.onPing)
	controller.channel.OnMessage("PONG", controller.onPong)

	return controller
}

func (c *HeartbeatController) onPing(_ *api.CastMessage) {
	err := c.channel.Send(pong)
	if err != nil {
		log.Printf("Error sending pong: %s", err)
	}
}

func (c *HeartbeatController) sendEvent(event Event) {
	select {
	case c.eventsCh <- event:
	default:
		log.Printf("Dropped event: %#v", event)
	}
}

func (c *HeartbeatController) onPong(_ *api.CastMessage) {
	atomic.StoreInt64(&c.pongs, 0)
}

func (c *HeartbeatController) Start(ctx context.Context) error {
	if c.ticker != nil {
		c.Stop()
	}

	c.ticker = time.NewTicker(interval)
	go func() {
	LOOP:
		for {
			select {
			case <-c.ticker.C:
				if atomic.LoadInt64(&c.pongs) >= maxBacklog {
					log.Printf("Missed %d pongs", c.pongs)
					c.sendEvent(event.Disconnected{errors.New("Ping timeout")})
					break LOOP
				}
				err := c.channel.Send(ping)
				atomic.AddInt64(&c.pongs, 1)
				if err != nil {
					log.Printf("Error sending ping: %s", err)
					c.sendEvent(event.Disconnected{err})
					break LOOP
				}
			case <-ctx.Done():
				log.Println("Heartbeat stopped")
				break LOOP
			}
		}
	}()

	log.Println("Heartbeat started")
	return nil
}

func (c *HeartbeatController) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
		c.ticker = nil
	}
}

type ConnectionController struct {
	channel *Channel
}

var connect = PayloadHeaders{Type: "CONNECT"}
var close = PayloadHeaders{Type: "CLOSE"}

func NewConnectionController(conn *TlsConnection, eventsCh chan Event, sourceId, destinationId string) *ConnectionController {
	controller := &ConnectionController{
		channel: conn.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.cast.tp.connection"),
	}

	return controller
}

func (c *ConnectionController) Start(ctx context.Context) error {
	return c.channel.Send(connect)
}

func (c *ConnectionController) Close() error {
	return c.channel.Send(close)
}

type Event interface {
}

type MediaController struct {
	interval       time.Duration
	channel        *Channel
	eventsCh       chan Event
	DestinationID  string
	MediaSessionID int
}

const NamespaceMedia = "urn:x-cast:com.google.cast.media"

var getMediaStatus = PayloadHeaders{Type: "GET_STATUS"}

var commandMediaPlay = PayloadHeaders{Type: "PLAY"}
var commandMediaPause = PayloadHeaders{Type: "PAUSE"}
var commandMediaStop = PayloadHeaders{Type: "STOP"}
var commandMediaLoad = PayloadHeaders{Type: "LOAD"}

type MediaCommand struct {
	PayloadHeaders
	MediaSessionID int `json:"mediaSessionId"`
}

type LoadMediaCommand struct {
	PayloadHeaders
	Media       MediaItem   `json:"media"`
	CurrentTime int         `json:"currentTime"`
	Autoplay    bool        `json:"autoplay"`
	CustomData  interface{} `json:"customData"`
}

type MediaItem struct {
	ContentId   string `json:"contentId"`
	StreamType  string `json:"streamType"`
	ContentType string `json:"contentType"`
}

type MediaStatusMedia struct {
	ContentId   string  `json:"contentId"`
	StreamType  string  `json:"streamType"`
	ContentType string  `json:"contentType"`
	Duration    float64 `json:"duration"`
}

func NewMediaController(conn *TlsConnection, eventsCh chan Event, sourceId, destinationID string) *MediaController {
	controller := &MediaController{
		channel:       conn.NewChannel(sourceId, destinationID, NamespaceMedia),
		eventsCh:      eventsCh,
		DestinationID: destinationID,
	}

	controller.channel.OnMessage("MEDIA_STATUS", controller.onStatus)

	return controller
}

func (c *MediaController) SetDestinationID(id string) {
	c.channel.DestinationId = id
	c.DestinationID = id
}

func (c *MediaController) sendEvent(event Event) {
	select {
	case c.eventsCh <- event:
	default:
		log.Printf("Dropped event: %#v", event)
	}
}

func (c *MediaController) onStatus(message *api.CastMessage) {
	response, err := c.parseStatus(message)
	if err != nil {
		log.Printf("Error parsing status: %s", err)
	}

	for _, status := range response.Status {
		c.sendEvent(*status)
	}
}

func (c *MediaController) parseStatus(message *api.CastMessage) (*MediaStatusResponse, error) {
	response := &MediaStatusResponse{}

	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
	}

	for _, status := range response.Status {
		c.MediaSessionID = status.MediaSessionID
	}

	return response, nil
}

type MediaStatusResponse struct {
	PayloadHeaders
	Status []*MediaStatus `json:"status,omitempty"`
}

type MediaStatus struct {
	PayloadHeaders
	MediaSessionID         int                    `json:"mediaSessionId"`
	PlaybackRate           float64                `json:"playbackRate"`
	PlayerState            string                 `json:"playerState"`
	CurrentTime            float64                `json:"currentTime"`
	SupportedMediaCommands int                    `json:"supportedMediaCommands"`
	Volume                 *Volume                `json:"volume,omitempty"`
	Media                  *MediaStatusMedia      `json:"media"`
	CustomData             map[string]interface{} `json:"customData"`
	RepeatMode             string                 `json:"repeatMode"`
	IdleReason             string                 `json:"idleReason"`
}

func (c *MediaController) Start(ctx context.Context) error {
	_, err := c.GetStatus(ctx)
	return err
}

func (c *MediaController) GetStatus(ctx context.Context) (*MediaStatusResponse, error) {
	message, err := c.channel.Request(ctx, &getMediaStatus)
	if err != nil {
		return nil, fmt.Errorf("Failed to get receiver status: %s", err)
	}

	return c.parseStatus(message)
}

func (c *MediaController) Play(ctx context.Context) (*api.CastMessage, error) {
	message, err := c.channel.Request(ctx, &MediaCommand{commandMediaPlay, c.MediaSessionID})
	if err != nil {
		return nil, fmt.Errorf("Failed to send play command: %s", err)
	}
	return message, nil
}

func (c *MediaController) Pause(ctx context.Context) (*api.CastMessage, error) {
	message, err := c.channel.Request(ctx, &MediaCommand{commandMediaPause, c.MediaSessionID})
	if err != nil {
		return nil, fmt.Errorf("Failed to send pause command: %s", err)
	}
	return message, nil
}

func (c *MediaController) Stop(ctx context.Context) (*api.CastMessage, error) {
	if c.MediaSessionID == 0 {
		// no current session to stop
		return nil, nil
	}
	message, err := c.channel.Request(ctx, &MediaCommand{commandMediaStop, c.MediaSessionID})
	if err != nil {
		return nil, fmt.Errorf("Failed to send stop command: %s", err)
	}
	return message, nil
}

func (c *MediaController) LoadMedia(ctx context.Context, media MediaItem, currentTime int, autoplay bool, customData interface{}) (*api.CastMessage, error) {
	message, err := c.channel.Request(ctx, &LoadMediaCommand{
		PayloadHeaders: commandMediaLoad,
		Media:          media,
		CurrentTime:    currentTime,
		Autoplay:       autoplay,
		CustomData:     customData,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to send load command: %s", err)
	}

	response := &PayloadHeaders{}
	err = json.Unmarshal([]byte(*message.PayloadUtf8), response)
	if err != nil {
		return nil, err
	}
	if response.Type == "LOAD_FAILED" {
		return nil, errors.New("Load media failed")
	}

	return message, nil
}

type Volume struct {
	Level *float64 `json:"level,omitempty"`
	Muted *bool    `json:"muted,omitempty"`
}

type ReceiverController struct {
	interval time.Duration
	channel  *Channel
	eventsCh chan Event
	status   *ReceiverStatus
}

var getStatus = PayloadHeaders{Type: "GET_STATUS"}
var commandLaunch = PayloadHeaders{Type: "LAUNCH"}
var commandStop = PayloadHeaders{Type: "STOP"}

func NewReceiverController(conn *TlsConnection, eventsCh chan Event, sourceId, destinationId string) *ReceiverController {
	controller := &ReceiverController{
		channel:  conn.NewChannel(sourceId, destinationId, "urn:x-cast:com.google.cast.receiver"),
		eventsCh: eventsCh,
	}

	controller.channel.OnMessage("RECEIVER_STATUS", controller.onStatus)

	return controller
}

func (c *ReceiverController) sendEvent(event Event) {
	select {
	case c.eventsCh <- event:
	default:
		log.Printf("Dropped event: %#v", event)
	}
}

func (c *ReceiverController) onStatus(message *api.CastMessage) {
	response := &StatusResponse{}
	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)
	if err != nil {
		log.Printf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
		return
	}

	previous := map[string]*ApplicationSession{}
	if c.status != nil {
		for _, app := range c.status.Applications {
			previous[*app.AppID] = app
		}
	}

	c.status = response.Status
	vol := response.Status.Volume
	c.sendEvent(event.StatusUpdated{Level: *vol.Level, Muted: *vol.Muted})

	for _, app := range response.Status.Applications {
		if _, ok := previous[*app.AppID]; ok {
			// Already running
			delete(previous, *app.AppID)
			continue
		}
		event := event.AppStarted{
			AppID:       *app.AppID,
			DisplayName: *app.DisplayName,
			StatusText:  *app.StatusText,
		}
		c.sendEvent(event)
	}

	// Stopped apps
	for _, app := range previous {
		event := event.AppStopped{
			AppID:       *app.AppID,
			DisplayName: *app.DisplayName,
			StatusText:  *app.StatusText,
		}
		c.sendEvent(event)
	}
}

type StatusResponse struct {
	PayloadHeaders
	Status *ReceiverStatus `json:"status,omitempty"`
}

type ReceiverStatus struct {
	PayloadHeaders
	Applications []*ApplicationSession `json:"applications"`
	Volume       *Volume               `json:"volume,omitempty"`
}

type LaunchRequest struct {
	PayloadHeaders
	AppId string `json:"appId"`
}

func (s *ReceiverStatus) GetSessionByNamespace(namespace string) *ApplicationSession {
	for _, app := range s.Applications {
		for _, ns := range app.Namespaces {
			if ns.Name == namespace {
				return app
			}
		}
	}
	return nil
}

func (s *ReceiverStatus) GetSessionByAppId(appId string) *ApplicationSession {
	for _, app := range s.Applications {
		if *app.AppID == appId {
			return app
		}
	}
	return nil
}

type ApplicationSession struct {
	AppID       *string      `json:"appId,omitempty"`
	DisplayName *string      `json:"displayName,omitempty"`
	Namespaces  []*Namespace `json:"namespaces"`
	SessionID   *string      `json:"sessionId,omitempty"`
	StatusText  *string      `json:"statusText,omitempty"`
	TransportId *string      `json:"transportId,omitempty"`
}

type Namespace struct {
	Name string `json:"name"`
}

func (c *ReceiverController) Start(ctx context.Context) error {
	// noop
	return nil
}

func (c *ReceiverController) GetStatus(ctx context.Context) (*ReceiverStatus, error) {
	message, err := c.channel.Request(ctx, &getStatus)
	if err != nil {
		return nil, fmt.Errorf("Failed to get receiver status: %s", err)
	}

	response := &StatusResponse{}
	err = json.Unmarshal([]byte(*message.PayloadUtf8), response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message: %s - %s", err, *message.PayloadUtf8)
	}

	return response.Status, nil
}

func (c *ReceiverController) SetVolume(ctx context.Context, volume *Volume) (*api.CastMessage, error) {
	return c.channel.Request(ctx, &ReceiverStatus{
		PayloadHeaders: PayloadHeaders{Type: "SET_VOLUME"},
		Volume:         volume,
	})
}

func (c *ReceiverController) GetVolume(ctx context.Context) (*Volume, error) {
	status, err := c.GetStatus(ctx)
	if err != nil {
		return nil, err
	}
	return status.Volume, err
}

func (c *ReceiverController) LaunchApp(ctx context.Context, appId string) (*ReceiverStatus, error) {
	message, err := c.channel.Request(ctx, &LaunchRequest{
		PayloadHeaders: commandLaunch,
		AppId:          appId,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed sending request: %s", err)
	}

	response := &StatusResponse{}
	err = json.Unmarshal([]byte(*message.PayloadUtf8), response)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message: %s - %s", err, *message.PayloadUtf8)
	}
	return response.Status, nil
}

func (c *ReceiverController) QuitApp(ctx context.Context) (*api.CastMessage, error) {
	return c.channel.Request(ctx, &commandStop)
}

type URLController struct {
	interval      time.Duration
	channel       *Channel
	eventsCh      chan Event
	DestinationID string
	URLSessionID  int
}

const NamespaceURL = "urn:x-cast:com.url.cast"

var getURLStatus = PayloadHeaders{Type: "GET_STATUS"}

var commandURLLoad = PayloadHeaders{Type: "LOAD"}

type LoadURLCommand struct {
	PayloadHeaders
	URL  string `json:"url"`
	Type string `json:"type"`
}

type URLStatusURL struct {
	ContentId   string  `json:"contentId"`
	StreamType  string  `json:"streamType"`
	ContentType string  `json:"contentType"`
	Duration    float64 `json:"duration"`
}

func NewURLController(conn *TlsConnection, eventsCh chan Event, sourceId, destinationID string) *URLController {
	controller := &URLController{
		channel:       conn.NewChannel(sourceId, destinationID, NamespaceURL),
		eventsCh:      eventsCh,
		DestinationID: destinationID,
	}

	controller.channel.OnMessage("URL_STATUS", controller.onStatus)

	return controller
}

func (c *URLController) SetDestinationID(id string) {
	c.channel.DestinationId = id
	c.DestinationID = id
}

func (c *URLController) sendEvent(event Event) {
	select {
	case c.eventsCh <- event:
	default:
		log.Printf("Dropped event: %#v", event)
	}
}

func (c *URLController) onStatus(message *api.CastMessage) {
	response, err := c.parseStatus(message)
	if err != nil {
		log.Printf("Error parsing status: %s", err)
	}

	for _, status := range response.Status {
		c.sendEvent(*status)
	}
}

func (c *URLController) parseStatus(message *api.CastMessage) (*URLStatusResponse, error) {
	response := &URLStatusResponse{}

	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
	}

	for _, status := range response.Status {
		c.URLSessionID = status.URLSessionID
	}

	return response, nil
}

type URLStatusResponse struct {
	PayloadHeaders
	Status []*URLStatus `json:"status,omitempty"`
}

type URLStatus struct {
	PayloadHeaders
	URLSessionID         int                    `json:"mediaSessionId"`
	PlaybackRate         float64                `json:"playbackRate"`
	PlayerState          string                 `json:"playerState"`
	CurrentTime          float64                `json:"currentTime"`
	SupportedURLCommands int                    `json:"supportedURLCommands"`
	Volume               *Volume                `json:"volume,omitempty"`
	URL                  *URLStatusURL          `json:"media"`
	CustomData           map[string]interface{} `json:"customData"`
	RepeatMode           string                 `json:"repeatMode"`
	IdleReason           string                 `json:"idleReason"`
}

func (c *URLController) Start(ctx context.Context) error {
	_, err := c.GetStatus(ctx)
	return err
}

func (c *URLController) GetStatus(ctx context.Context) (*URLStatusResponse, error) {
	message, err := c.channel.Request(ctx, &getURLStatus)
	if err != nil {
		return nil, fmt.Errorf("Failed to get receiver status: %s", err)
	}

	return c.parseStatus(message)
}

func (c *URLController) LoadURL(ctx context.Context, url string) (*api.CastMessage, error) {
	message, err := c.channel.Request(ctx, &LoadURLCommand{
		PayloadHeaders: commandURLLoad,
		URL:            url,
		Type:           "loc",
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to send load command: %s", err)
	}

	response := &PayloadHeaders{}
	err = json.Unmarshal([]byte(*message.PayloadUtf8), response)
	if err != nil {
		return nil, err
	}
	if response.Type == "LOAD_FAILED" {
		return nil, errors.New("Load URL failed")
	}

	return message, nil
}
