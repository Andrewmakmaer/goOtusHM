package memorystorage

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

var events = []storage.Event{
	{
		ID:           "evt-001",
		Title:        "Встреча с клиентом",
		Description:  "Обсуждение нового проекта с ООО 'Рога и копыта'",
		StartTime:    time.Date(2024, time.August, 1, 10, 0, 0, 0, time.UTC),
		EndTime:      time.Date(2024, time.August, 1, 11, 30, 0, 0, time.UTC),
		UserID:       "usr-123",
		CallDuration: 90 * time.Minute,
	},
	{
		ID:           "evt-002",
		Title:        "Командный брифинг",
		Description:  "Еженедельное собрание отдела разработки",
		StartTime:    time.Date(2024, time.August, 2, 14, 0, 0, 0, time.UTC),
		EndTime:      time.Date(2024, time.August, 2, 15, 0, 0, 0, time.UTC),
		UserID:       "usr-456",
		CallDuration: 60 * time.Minute,
	},
	{
		ID:           "evt-004",
		Title:        "1-2-1",
		Description:  "Раздавить по пивасику с лидом разработки, на часик максимум",
		StartTime:    time.Date(2024, time.July, 31, 19, 0, 0, 0, time.UTC),
		EndTime:      time.Date(2024, time.August, 1, 6, 30, 0, 0, time.UTC),
		UserID:       "usr-456",
		CallDuration: 60 * time.Minute,
	},
	{
		ID:           "evt-003",
		Title:        "Интервью кандидата",
		Description:  "Собеседование на позицию младшего разработчика",
		StartTime:    time.Date(2024, time.August, 3, 11, 0, 0, 0, time.UTC),
		EndTime:      time.Date(2024, time.August, 3, 12, 0, 0, 0, time.UTC),
		UserID:       "usr-789",
		CallDuration: 60 * time.Minute,
	},
	{
		ID:           "evt-005",
		Title:        "Корпоратив",
		Description:  "Отмечаем окончание лета",
		StartTime:    time.Date(2024, time.August, 31, 14, 0, 0, 0, time.UTC),
		EndTime:      time.Date(2024, time.August, 31, 21, 0, 0, 0, time.UTC),
		UserID:       "usr-789",
		CallDuration: 60 * time.Minute,
	},
}

func TestStorage(t *testing.T) {
	c := New()
	t.Run("simple tests", func(t *testing.T) {
		var ok error
		for _, item := range events {
			ok = c.AddEvent(item)
		}
		require.NoError(t, nil, ok)

		result, err := c.ListEventsDay("usr-123", time.Date(2024, time.August, 2, 10, 0, 0, 0, time.UTC))
		require.NoError(t, err)
		require.True(t, (result[0] == events[0]))

		result, err = c.ListEventsWeek("usr-456", time.Date(2024, time.July, 30, 10, 0, 0, 0, time.UTC))
		require.NoError(t, err)
		expected := events[1:3]
		require.True(t, (reflect.DeepEqual(result, expected)))

		result, err = c.ListEventsMonth("usr-789", time.Date(2024, time.August, 30, 10, 0, 0, 0, time.UTC))
		require.NoError(t, err)
		expected = events[3:]
		require.True(t, (reflect.DeepEqual(result, expected)))
	})

	t.Run("cross events", func(t *testing.T) {
		newEvent := storage.NewEvent(
			"evt-123",
			"Обед",
			"ням-ням",
			time.Date(2024, time.August, 3, 11, 30, 0, 0, time.UTC),
			time.Date(2024, time.August, 3, 12, 30, 0, 0, time.UTC),
			"usr-789",
			60*time.Minute,
		)
		ok := c.AddEvent(*newEvent)

		require.ErrorIs(t, storage.ErrDateBusy, ok)

		newEvent = storage.NewEvent(
			"evt-1234",
			"Сон",
			"Мощнецки отосплюсь",
			time.Date(2024, time.August, 2, 23, 50, 0, 0, time.UTC),
			time.Date(2024, time.August, 3, 13, 0o0, 0, 0, time.UTC),
			"usr-789",
			60*time.Minute,
		)

		ok = c.AddEvent(*newEvent)

		require.ErrorIs(t, storage.ErrDateBusy, ok)

		ok = c.AddEvent(
			storage.Event{
				ID:           "evt-009",
				Title:        "Сон",
				Description:  "Чуть менее мощнецки отосплюсь",
				StartTime:    time.Date(2024, time.August, 2, 23, 50, 0, 0, time.UTC),
				EndTime:      time.Date(2024, time.August, 3, 11, 0o0, 0, 0, time.UTC),
				UserID:       "usr-789",
				CallDuration: 60 * time.Minute,
			},
		)

		require.ErrorIs(t, nil, ok)
	})

	t.Run("incorrect data", func(t *testing.T) {
		ok := c.AddEvent(
			storage.Event{
				ID:           "evt-010",
				Title:        "event",
				Description:  "event",
				StartTime:    time.Date(2024, time.August, 3, 11, 0o0, 0, 0, time.UTC),
				EndTime:      time.Date(2024, time.August, 2, 23, 50, 0, 0, time.UTC),
				UserID:       "usr-789",
				CallDuration: 60 * time.Minute,
			},
		)

		require.ErrorIs(t, storage.ErrTimeShift, ok)
	})

	t.Run("Update Event", func(t *testing.T) {
		updatedEvent := storage.Event{
			ID:           "evt-009",
			Title:        "Сон",
			Description:  "Не, все таки совесть надо иметь",
			StartTime:    time.Date(2024, time.August, 2, 23, 50, 0, 0, time.UTC),
			EndTime:      time.Date(2024, time.August, 3, 11, 0o0, 0, 0, time.UTC),
			UserID:       "usr-789",
			CallDuration: 60 * time.Minute,
		}

		c.UpdateEvent(updatedEvent)

		var resultDescrip string
		listEvents, _ := c.ListEventsDay("usr-789", time.Date(2024, time.August, 2, 23, 50, 0, 0, time.UTC))
		for _, item := range listEvents {
			if item.ID == "evt-009" {
				resultDescrip = item.Description
			}
		}
		require.Equal(t, resultDescrip, "Не, все таки совесть надо иметь")
	})

	t.Run("Delete", func(t *testing.T) {
		c.DeleteEvent("evt-009", "usr-789")
		listEvents, _ := c.ListEventsDay("usr-789", time.Date(2024, time.August, 2, 23, 50, 0, 0, time.UTC))
		require.Len(t, listEvents, 2)
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := New()
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, item := range events {
			ok := c.AddEvent(item)
			require.NoError(t, nil, ok)
		}
	}()

	go func() {
		defer wg.Done()
		for _, item := range events {
			ok := c.AddEvent(item)
			require.NoError(t, nil, ok)
		}
	}()

	wg.Wait()
}
