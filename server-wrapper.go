package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
)

// Переменная для записи в stdin Minecraft-сервера
var mcStdin *bufio.Writer

// Функция для отображения формы ввода команд
func handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Отдаем HTML-страницу с формой
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<html>
			<head><title>Введите команду Minecraft</title></head>
			<body>
				<h1>Введите команду для Minecraft</h1>
				<form action="/command" method="POST">
					<label for="command">Команда:</label>
					<input type="text" id="command" name="command" required>
					<button type="submit">Отправить команду</button>
				</form>
			</body>
			</html>
		`)
	} else if r.Method == http.MethodPost {
		// Получаем команду из формы
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
			return
		}

		command := r.FormValue("command")
		if command == "" {
			http.Error(w, "Команда не указана", http.StatusBadRequest)
			return
		}

		// Отправляем команду в stdin Minecraft-сервера
		_, err = mcStdin.WriteString(command + "\n")
		if err != nil {
			http.Error(w, "Ошибка отправки команды", http.StatusInternalServerError)
			return
		}

		mcStdin.Flush()

		// Отправляем ответ пользователю
		w.Write([]byte(fmt.Sprintf("Команда отправлена: %s", command)))
	}
}

// Функция для получения последних логов Minecraft-сервера
func handleLogs(w http.ResponseWriter, r *http.Request) {
	logFile := "/minecraft/logs/latest.log"
	logFileContent, err := os.ReadFile(logFile)
	if err != nil {
		http.Error(w, "Не удалось прочитать логи", http.StatusInternalServerError)
		return
	}

	// Отдаем логи в формате HTML
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h1>Последние логи:</h1><pre>%s</pre>", logFileContent)
}

func main() {
	// Запускаем Minecraft сервер с помощью Java
	cmd := exec.Command("java", "-Xmx10G", "-Xms2G", "-jar", "forge-1.12.2-14.23.5.2860.jar", "nogui")

	// Перенаправляем stdout и stderr для отображения в консоли
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Создаём канал для отправки команд в stdin
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatalf("Ошибка создания stdin: %v", err)
	}

	// Создаем буферизованный writer для stdin
	mcStdin = bufio.NewWriter(stdinPipe)

	// Запускаем сервер Minecraft
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Ошибка запуска Minecraft сервера: %v", err)
	}

	// Регистрируем обработчики HTTP запросов
	http.HandleFunc("/command", handleCommand)  // Отображение формы и обработка команд
	http.HandleFunc("/logs", handleLogs)        // Отображение логов

	// Запускаем HTTP сервер на порту 8080
	log.Println("HTTP сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	// Ожидаем завершения работы Minecraft сервера
	cmd.Wait()
}

