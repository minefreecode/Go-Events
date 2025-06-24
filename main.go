package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Организация круга событий в Go")

	usrEvents := make(chan string, 10) // Создание канала пользовательских событий
	stmEvents := make(chan string, 10) //Создание канала системных событий
	tmrEvents := make(chan string, 10) //Канал событий тайминга
	destroy := make(chan struct{})     //Канал для событий закртия программы

	go userEventProducer(usrEvents)   //Создать пользовательские события
	go systemEventProducer(stmEvents) //Создать системные события
	go timerEventProducer(tmrEvents)  //Создать события по таймингу

	go eventCatch(usrEvents, stmEvents, tmrEvents, destroy) //Поймать ошибку

	time.Sleep(5 * time.Second)

	fmt.Println("Закрытие круга событий")
	close(destroy)

	time.Sleep(500 * time.Millisecond)
	fmt.Println("Круг событий завершён")
}

func eventCatch(userEvents, systemEvents, timerEvents <-chan string, shutdown <-chan struct{}) {
	fmt.Println("Цикл событий начался...")
	for {
		select {
		case event := <-userEvents:
			fmt.Printf("Обработка событий пользователя: %s\n", event)
			processUserEvent(event)
		case event := <-systemEvents:
			fmt.Printf("Обработка системного события: %s\n", event)
			processSystemEvent(event)
		case event := <-timerEvents:
			fmt.Printf("Обработка событий тайминга: %s\n", event)
			processTimerEvent(event)
		case <-shutdown:
			fmt.Println("Событие закртытия получено")
			return
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

func processUserEvent(event string) {
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("   -> Пользовательское событие обработано: %s\n", event)
}

func processSystemEvent(event string) {
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("   -> Системное событие обработано: %s\n", event)
}

func processTimerEvent(event string) {
	time.Sleep(50 * time.Millisecond)
	fmt.Printf("   -> Событие таймера обработано: %s\n", event)
}
