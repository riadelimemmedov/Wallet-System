package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	logger "github.com/riad/banksystemendtoend/pkg/log"
	"go.uber.org/zap"
)

const (
	reconnectDelay = 3 * time.Second
	maxRetries     = 5
)

var (
	instance *Client
	once     sync.Once
	mu       sync.RWMutex
)

// Client provides a wrapper for the RabbitMQ (connection) and (channel)
type Client struct {
	conn         *amqp.Connection
	ch           *amqp.Channel
	connNotify   chan *amqp.Error
	chNotify     chan *amqp.Error
	config       Config
	mu           sync.Mutex
	isConnected  bool
	reconnecting bool
	quit         chan struct{}
}

// Config contains the connection parameters for RabbitMQ
type Config struct {
	Host              string
	Port              string
	Username          string
	Password          string
	VHost             string
	ConnectionTimeout time.Duration
	HeartbeatInterval time.Duration
}

// Message represents a message to be published to RabbitMQ
type Message struct {
	ContentType string
	Body        []byte
	Priority    uint8
	Expiration  string
}

// GetClient returns the singleton instance of the RabbitMQ client
func GetClient() *Client {
	mu.RLock()
	if instance != nil {
		defer mu.RUnlock()
		return instance
	}
	mu.RUnlock()
	return nil
}

// InitClient initializes the singleton instance of RabbitMQ client
// This should be called once at the start of your application
func InitClient(cfg Config) error {
	var initErr error
	once.Do(func() {
		mu.Lock()
		defer mu.Unlock()

		client := &Client{
			config: cfg,
			quit:   make(chan struct{}),
		}

		if err := client.connect(); err != nil {
			initErr = err
			return
		}

		go client.handleReconnect()
		instance = client
	})
	return initErr
}

// connect establishes a connection to RabbitMQ and creates a channel
func (c *Client) connect() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.isConnected {
		logger.Info("RabbitMQ client is already connected")
		return nil
	}

	// Build connection URL
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		c.config.Username, c.config.Password, c.config.Host, c.config.Port, c.config.VHost)

	config := amqp.Config{
		Heartbeat: c.config.HeartbeatInterval,
		Dial:      amqp.DefaultDial(c.config.ConnectionTimeout),
	}

	conn, err := amqp.DialConfig(connStr, config)
	if err != nil {
		logger.Error("Failed to connect to RabbitMQ", zap.Error(err))
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	c.conn = conn

	// Get connection notification channel
	c.connNotify = conn.NotifyClose(make(chan *amqp.Error, 1))

	ch, err := conn.Channel()
	if err != nil {
		c.conn.Close()
		logger.Error("Failed to open a channel", zap.Error(err))
		return fmt.Errorf("failed to open a channel: %w", err)
	}
	c.ch = ch

	// Get channel notification channel
	c.chNotify = ch.NotifyClose(make(chan *amqp.Error, 1))

	c.isConnected = true

	logger.GetLogger().Info("Connected to RabbitMQ",
		zap.String("host", c.config.Host),
		zap.String("vhost", c.config.VHost))

	return nil
}

// handleReconnect monitors connection and reconnects if necessary
func (c *Client) handleReconnect() {
	for {
		select {
		case <-c.quit:
			return
		case err := <-c.connNotify:
			if err != nil {
				c.reconnect("connection", err.Error())
			}
		case err := <-c.chNotify:
			if err != nil {
				c.reconnect("channel", err.Error())
			}
		}
	}
}

// reconnect handles reconnection logic with exponential backoff
func (c *Client) reconnect(reason, errMsg string) {
	c.mu.Lock()
	if c.reconnecting {
		c.mu.Unlock()
		return
	}
	c.reconnecting = true
	c.isConnected = false
	c.mu.Unlock()

	logger.GetLogger().Warn("RabbitMQ disconnected",
		zap.String("reason", reason),
		zap.String("error", errMsg))

	for i := range maxRetries {
		time.Sleep(reconnectDelay * time.Duration(i+1))

		if err := c.connect(); err != nil {
			logger.GetLogger().Error("Failed to reconnect to RabbitMQ",
				zap.Int("attempt", i+1),
				zap.Error(err))
		} else {
			c.mu.Lock()
			c.reconnecting = false
			c.mu.Unlock()
			logger.GetLogger().Info("Successfully reconnected to RabbitMQ",
				zap.Int("attempt", i+1))
			return
		}
	}
	logger.GetLogger().Error("Failed to reconnect to RabbitMQ after multiple attempts",
		zap.Int("maxRetries", maxRetries))

	c.mu.Lock()
	c.reconnecting = false
	c.mu.Unlock()
}

// DeclareQueue declares a queue to the broker
func (c *Client) DeclareQueue(name string, durable, autoDelete, exclusive bool) (amqp.Queue, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		logger.Error("RabbitMQ client is not connected when declaring queue")
		return amqp.Queue{}, fmt.Errorf("RabbitMQ client is not connected when declaring queue")
	}

	return c.ch.QueueDeclare(
		name,
		durable,
		autoDelete,
		exclusive,
		false,
		nil,
	)
}

// DeclareExchange declares an exchange to the broker
func (c *Client) DeclareExchange(name, kind string, durable, autoDelete, internal bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		logger.Error("RabbitMQ client is not connected when declaring exchange")
		return fmt.Errorf("RabbitMQ client is not connected when declaring exchange")
	}

	return c.ch.ExchangeDeclare(
		name,
		kind,
		durable,
		autoDelete,
		internal,
		false,
		nil,
	)
}

// BindQueue binds a queue to an exchange
func (c *Client) BindQueue(name, key, exchange string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		logger.Error("RabbitMQ client is not connected when binding queue")
		return fmt.Errorf("RabbitMQ client is not connected when binding queue")
	}

	return c.ch.QueueBind(
		name,
		key,
		exchange,
		false,
		nil,
	)
}

// Publish publishes a message to an exchange
func (c *Client) Publish(ctx context.Context, exchange, routingKey string, msg Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		logger.Error("RabbitMQ client is not connected when publishing message")
		return fmt.Errorf("RabbitMQ client is not connected when publishing message")
	}

	return c.ch.PublishWithContext(
		ctx,
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			DeliveryMode: amqp.Persistent,
			Priority:     msg.Priority,
			Expiration:   msg.Expiration,
			Timestamp:    time.Now(),
		},
	)
}

// PublishJSON publishes a JSON message to an exchange
func (c *Client) PublishJSON(ctx context.Context, exchange, routingKey string, obj interface{}) error {
	data, err := json.Marshal(obj)
	if err != nil {
		logger.Error("Failed to marshal JSON", zap.Error(err))
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	msg := Message{
		ContentType: "application/json",
		Body:        data,
	}
	return c.Publish(ctx, exchange, routingKey, msg)
}

// Consume starts consuming messages from a queue
func (c *Client) Consume(queue, consumer string, autoAck, exclusive bool) (<-chan amqp.Delivery, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		logger.Error("RabbitMQ client is not connected when consuming messages")
		return nil, fmt.Errorf("RabbitMQ client is not connected when consuming messages")
	}

	return c.ch.Consume(
		queue,
		consumer,
		autoAck,
		exclusive,
		false,
		false,
		nil,
	)
}

// QoS sets the quality of service for the channel
func (c *Client) QoS(prefetchCount, prefetchSize int, global bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.isConnected {
		logger.Error("RabbitMQ client is not connected when setting QoS")
		return fmt.Errorf("RabbitMQ client is not connected when setting QoS")
	}

	return c.ch.Qos(prefetchCount, prefetchSize, global)
}

// Close closes the channel and connection
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	close(c.quit)

	if !c.isConnected {
		return nil
	}

	if err := c.ch.Close(); err != nil {
		logger.Error("Failed to close channel", zap.Error(err))
		return fmt.Errorf("failed to close channel: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		logger.Error("Failed to close connection", zap.Error(err))
		return fmt.Errorf("failed to close connection: %w", err)
	}

	c.isConnected = false
	logger.GetLogger().Info("RabbitMQ client closed Successfully")
	return nil
}
