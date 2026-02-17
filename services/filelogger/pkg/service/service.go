package service

import (
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"

	"github.com/opencloud-eu/opencloud/pkg/log"
	"github.com/opencloud-eu/opencloud/services/filelogger/pkg/config"
	"github.com/opencloud-eu/reva/v2/pkg/events"
)

// FileLoggerService is the service responsible for logging file events to a JSON file
type FileLoggerService struct {
	log     log.Logger
	cfg     *config.Config
	ch      <-chan events.Event
	stopCh  chan struct{}
	stopped atomic.Bool
	logFile *os.File
}

// NewFileLoggerService returns a filelogger service
func NewFileLoggerService(logger log.Logger, cfg *config.Config, stream events.Stream) (*FileLoggerService, error) {
	if stream == nil {
		return nil, fmt.Errorf("need non nil stream to work properly")
	}

	logger.Info().
		Str("endpoint", cfg.Events.Endpoint).
		Str("cluster", cfg.Events.Cluster).
		Msg("Connecting to event bus")

	// We subscribe to common file-related event types
	registeredEvents := []events.Unmarshaller{
		events.UploadReady{},
		events.ItemTrashed{},
		events.ItemRestored{},
		events.ItemMoved{},
		events.ContainerCreated{},
	}

	ch, err := events.Consume(stream, "filelogger", registeredEvents...)
	if err != nil {
		return nil, err
	}

	logger.Info().Str("path", cfg.LogFilePath).Msg("Opening log file")
	f, err := os.OpenFile(cfg.LogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("could not open log file %s: %w", cfg.LogFilePath, err)
	}

	return &FileLoggerService{
		log:     logger,
		cfg:     cfg,
		ch:      ch,
		stopCh:  make(chan struct{}, 1),
		logFile: f,
	}, nil
}

// Run runs the service
func (fs *FileLoggerService) Run() error {
	fs.log.Info().Str("file", fs.cfg.LogFilePath).Msg("Starting filelogger service")
EventLoop:
	for {
		select {
		case event, ok := <-fs.ch:
			if !ok {
				break EventLoop
			}
			fs.processEvent(event)

			if fs.stopped.Load() {
				break EventLoop
			}
		case <-fs.stopCh:
			break EventLoop
		}
	}

	return fs.logFile.Close()
}

func (fs *FileLoggerService) Close() {
	if fs.stopped.CompareAndSwap(false, true) {
		close(fs.stopCh)
	}
}

func (fs *FileLoggerService) processEvent(event events.Event) {
	// Simple structure for logging
	logEntry := struct {
		ID          string      `json:"id"`
		Type        string      `json:"type"`
		InitiatorID string      `json:"initiator_id"`
		Event       interface{} `json:"event"`
	}{
		ID:          event.ID,
		Type:        fmt.Sprintf("%T", event.Event),
		InitiatorID: event.InitiatorID,
		Event:       event.Event,
	}

	b, err := json.Marshal(logEntry)
	if err != nil {
		fs.log.Error().Err(err).Msg("failed to marshal event")
		return
	}

	if _, err := fs.logFile.Write(append(b, '\n')); err != nil {
		fs.log.Error().Err(err).Msg("failed to write to log file")
	}
}
