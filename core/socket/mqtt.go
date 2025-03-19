package socket

import (
	"fmt"
	"log"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
)

// MQTTClient 封装MQTT客户端功能
type MQTTClient struct {
	client      paho.Client
	brokerURL   string
	clientID    string
	username    string
	password    string
	isConnected bool
}

// NewMQTTClient 创建新的MQTT客户端
func NewMQTTClient(brokerURL, clientID, username, password string) *MQTTClient {
	return &MQTTClient{
		brokerURL: brokerURL,
		clientID:  clientID,
		username:  username,
		password:  password,
	}
}

// Connect 连接到MQTT代理
func (m *MQTTClient) Connect() error {
	opts := paho.NewClientOptions()
	opts.AddBroker(m.brokerURL)
	opts.SetClientID(m.clientID)

	if m.username != "" {
		opts.SetUsername(m.username)
		opts.SetPassword(m.password)
	}

	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(5 * time.Second)
	opts.SetConnectionLostHandler(m.onConnectionLost)
	opts.SetOnConnectHandler(m.onConnect)

	m.client = paho.NewClient(opts)
	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("连接MQTT代理失败: %v", token.Error())
	}

	return nil
}

// Disconnect 断开MQTT连接
func (m *MQTTClient) Disconnect() {
	if m.client != nil && m.client.IsConnected() {
		m.client.Disconnect(250)
	}
	m.isConnected = false
}

// Publish 发布消息到指定主题
func (m *MQTTClient) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT客户端未连接")
	}

	token := m.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("发布消息失败: %v", token.Error())
	}

	return nil
}

// Subscribe 订阅主题
func (m *MQTTClient) Subscribe(topic string, qos byte, callback paho.MessageHandler) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT客户端未连接")
	}

	token := m.client.Subscribe(topic, qos, callback)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("订阅主题失败: %v", token.Error())
	}

	return nil
}

// Unsubscribe 取消订阅主题
func (m *MQTTClient) Unsubscribe(topics ...string) error {
	if !m.client.IsConnected() {
		return fmt.Errorf("MQTT客户端未连接")
	}

	token := m.client.Unsubscribe(topics...)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("取消订阅失败: %v", token.Error())
	}

	return nil
}

// IsConnected 返回连接状态
func (m *MQTTClient) IsConnected() bool {
	return m.client != nil && m.client.IsConnected()
}

// 连接丢失处理函数
func (m *MQTTClient) onConnectionLost(client paho.Client, err error) {
	m.isConnected = false
	log.Printf("MQTT连接丢失: %v", err)
}

// 连接成功处理函数
func (m *MQTTClient) onConnect(client paho.Client) {
	m.isConnected = true
	log.Println("MQTT连接成功")
}
