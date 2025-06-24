package main

import (
	"fmt"
	"math/rand"
	"time"
)

func RunEventsLoop() {
	fmt.Println("Организация круга событий в Go")

	userEvents := make(chan string, 10)
	systemEvents := make(chan string, 10)
	timerEvents := make(chan string, 10)
	shutdown := make(chan struct{})

	go userEventProducer(userEvents)
	go systemEventProducer(systemEvents)
	go timerEventProducer(timerEvents)

	go eventLoop(userEvents, systemEvents, timerEvents, shutdown)

	time.Sleep(5 * time.Second)

	fmt.Println("Закрытие круга событий")
	close(shutdown)

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Круг событий завершён")
}

func eventLoop(userEvents, systemEvents, timerEvents <-chan string, shutdown <-chan struct{}) {
	fmt.Println("Цикл событий начался...")
	for {
		select {
		case event := <-userEvents:
			fmt.Printf("Обработка событий пользователя: %s\n", event)
			processUserEvent(event)
		case event := <-systemEvents:
			fmt.Printf("Обработка системного события: %s\n", event)
			processSystemEvent(event)
		}
	}
}

func userEventProducer(events chan<- string) {
	userActions := []string{"login", "logout", "click", "scroll", "submit"}
	for i := 0; i < 8; i++ {
		time.Sleep(time.Duration(rand.Intn(800)+200) * time.Millisecond)
		action := userActions[rand.Intn(len(userActions))]
		events <- fmt.Sprintf("%s (user_%d)", action, i+1)
	}
}

func systemEventProducer(events chan<- string) {
	systemEvents := []string{"backup", "update", "maintenance", "alert", "sync"}
	for i := 0; i < 6; i++ {
		time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)
		event := systemEvents[rand.Intn(len(systemEvents))]
		events <- fmt.Sprintf("%s (system_%d)", event, i+1)
	}
}

func timerEventProducer(events chan<- string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	count := 0
	for range ticker.C {
		if count >= 5 {
			break
		}
		events <- fmt.Sprintf("Процесс тайминга (timer_%d)", count+1)
		count++
	}
}
