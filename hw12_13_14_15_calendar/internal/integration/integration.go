package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	calendarServiceURL = os.Getenv("INTEGRATION_CALENDAR_URL") // "http://localhost:8888"
	rabbitMQURL        = os.Getenv("INTEGRATION_RABBIT_URL")   // "amqp://guest:guest@localhost:5672/"
	queueName          = os.Getenv("INTEGRATION_RABBIT_QUEUE") // "sender-log"
)

type Event struct {
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	StartTime    time.Time `json:"starttime"`
	EndTime      time.Time `json:"endtime"`
	CallDuration string    `json:"callDuration"`
}

func RunIntegrationTests() error {
	fmt.Println("Запуск интеграционных тестов...")

	// Ожидание запуска сервисов
	if err := waitForServices(); err != nil {
		return fmt.Errorf("ошибка при ожидании запуска сервисов: %w", err)
	}

	event := Event{
		Title:        "Test Event",
		Description:  "This is a test event",
		StartTime:    time.Now().Local().Add(5 * time.Minute),
		EndTime:      time.Now().Local().Add(6 * time.Minute),
		CallDuration: "5m",
	}
	eventID := uuid.New().String()
	userID := uuid.New().String()

	defer deleteEvent(userID, eventID)

	res, err := createEvent(userID, eventID, event)
	if err != nil {
		return fmt.Errorf("ошибка при создании события: %w", err)
	}
	fmt.Println(res)

	event.Title = "Overlap Title"
	eventID = uuid.New().String()

	res, err = createEvent(userID, eventID, event)
	if err == nil && bytes.Contains([]byte(res), []byte("event overlaps with existing events")) {
		fmt.Println(res)
		fmt.Println("Успех! Событие с наложением времени не создалось")
	} else {
		return fmt.Errorf("некоректный ответ на бизнес-ошибку")
	}

	defer deleteEvent(userID, eventID)

	time.Sleep(6 * time.Second)
	message, err := getMessageFromQueue()
	if err != nil {
		return fmt.Errorf("ошибка при получении сообщения из очереди: %w", err)
	}
	if !bytes.Contains(message, []byte(userID)) {
		return fmt.Errorf("сообщение в очереди не содержит ID пользователя")
	}
	fmt.Println("Событие успешно обнаружено в очереди")

	event.StartTime = event.StartTime.Add(24 * time.Hour)
	event.EndTime = event.EndTime.Add(24 * time.Hour)
	event.Title = "Test event 2"
	eventID = uuid.New().String()

	createEvent(userID, eventID, event)
	defer deleteEvent(userID, eventID)

	events, err := listEvents(userID, "day")
	if err != nil {
		return fmt.Errorf("ошибка при получении событий: %w", err)
	}
	if !strings.Contains(events, userID) && !strings.Contains(events, event.Title) && strings.Contains(events, eventID) {
		return fmt.Errorf("ошибка, в списке нет ожидаемого события или есть лишнее: %s", events)
	}
	fmt.Println("Найдены все ожидаемые события за день")

	events, _ = listEvents(userID, "week")
	if !strings.Contains(events, userID) && !strings.Contains(events, event.Title) && !strings.Contains(events, eventID) {
		return fmt.Errorf("ошибка, в списке нет ожидаемого события или есть лишнее: %s", events)
	}
	fmt.Println("Найдены все ожидаемые события за неделю")

	events, _ = listEvents(userID, "month")
	if !strings.Contains(events, userID) && !strings.Contains(events, event.Title) && !strings.Contains(events, eventID) {
		return fmt.Errorf("ошибка, в списке нет ожидаемого события или есть лишнее: %s", events)
	}
	fmt.Println("Найдены все ожидаемые события за месяц")

	fmt.Println("Все интеграционные тесты успешно пройдены!")
	return nil
}

func waitForServices() error {
	for i := 0; i < 30; i++ {
		_, err := http.NewRequest(http.MethodGet, calendarServiceURL, nil)
		if err == nil {
			break
		}
		fmt.Println(err)
		time.Sleep(time.Second)
	}

	for i := 0; i < 30; i++ {
		conn, err := amqp.Dial(rabbitMQURL)
		if err == nil {
			conn.Close()
			return nil
		}
		fmt.Println(err)
		time.Sleep(time.Second)
	}

	return fmt.Errorf("сервисы не запустились в отведенное время")
}

func createEvent(userID, eventID string, event Event) (string, error) {
	url := fmt.Sprintf("%s/events/%s/%s", calendarServiceURL, userID, eventID)
	method := "POST"
	b, _ := json.Marshal(event)
	payload := bytes.NewReader(b)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload) //nolint:noctx
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

func deleteEvent(userID, eventID string) error {
	url := fmt.Sprintf("%s/events/%s/%s", calendarServiceURL, userID, eventID)
	req, _ := http.NewRequest(http.MethodDelete, url, nil) //nolint:noctx
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неожиданный статус код при удалении: %d", resp.StatusCode)
	}

	return nil
}

func listEvents(userID, time string) (string, error) {
	url := fmt.Sprintf("%s/events/%s/%s", calendarServiceURL, userID, time)
	req, _ := http.NewRequest(http.MethodGet, url, nil) //nolint:noctx
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный статус код при получении событий: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return string(body), nil
}

func getMessageFromQueue() ([]byte, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	defer ch.Close()

	msg, ok, err := ch.Get(queueName, true)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("очередь пуста")
	}

	return msg.Body, nil
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
