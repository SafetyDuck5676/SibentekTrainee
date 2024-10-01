package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv" // Пакет для работы с переменными окружения
)

const (
	defaultLogFilePath = "app.log" // Путь к лог-файлу по умолчанию
	defaultConfigFile  = ".env"    // Файл конфигурации по умолчанию
)

// Читает данные из файла или stdin
func readNumbers(source string) ([]int, error) {
	var data []byte
	var err error

	if source == "stdin" {
		// Чтение данных из стандартного ввода
		fmt.Println("Введите массив чисел в формате JSON и нажмите Enter:")
		reader := bufio.NewReader(os.Stdin)
		data, err = reader.ReadBytes('\n')
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных из stdin: %v", err)
		}
	} else {
		// Чтение данных из файла
		data, err = os.ReadFile(source)
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении файла: %v", err)
		}
	}

	var numbers []int
	if err := json.Unmarshal(data, &numbers); err != nil {
		return nil, fmt.Errorf("ошибка при парсинге JSON: %v", err)
	}

	return numbers, nil
}

// Считает сумму массива чисел
func sumNumbers(numbers []int) int {
	sum := 0
	for _, num := range numbers {
		sum += num
	}
	return sum
}

// Выполняет HTTP GET запрос и проверяет статус ответа
func checkURL(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("ошибка при выполнении HTTP запроса: %v", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

// Инициализирует логгер с выводом в файл
func initLogger(logFile string) (*log.Logger, error) {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("ошибка при создании файла логов: %v", err)
	}
	logger := log.New(file, "LOG: ", log.Ldate|log.Ltime|log.Lshortfile)
	return logger, nil
}

// Загрузка переменных окружения из файла конфигурации
func loadConfig(configFile string) error {
	err := godotenv.Load(configFile)
	if err != nil {
		return fmt.Errorf("ошибка при загрузке файла конфигурации: %v", err)
	}
	return nil
}

func main() {
	// Параметры командной строки
	sourcePtr := flag.String("source", "stdin", "Источник данных: 'stdin' или путь к файлу")
	outputPtr := flag.String("output", "output.txt", "Файл для сохранения результата")
	logPtr := flag.String("log", defaultLogFilePath, "Файл для логов")
	configPtr := flag.String("config", defaultConfigFile, "Файл конфигурации")
	flag.Parse()

	// Загрузка файла конфигурации
	if err := loadConfig(*configPtr); err != nil {
		fmt.Printf("Ошибка при загрузке конфигурации: %v\n", err)
		return
	}

	// Получение URL из переменных окружения
	url := os.Getenv("TARGET_URL")
	if url == "" {
		fmt.Println("URL для HTTP запроса не задан в конфигурации")
		return
	}

	// Инициализация логгера
	logger, err := initLogger(*logPtr)
	if err != nil {
		fmt.Printf("Не удалось инициализировать логгер: %v\n", err)
		return
	}
	logger.Println("Программа запущена")

	// Чтение данных
	numbers, err := readNumbers(*sourcePtr)
	if err != nil {
		logger.Printf("Ошибка при чтении данных: %v\n", err)
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	logger.Printf("Массив чисел: %v\n", numbers)

	// Суммирование чисел
	sum := sumNumbers(numbers)
	logger.Printf("Сумма чисел: %d\n", sum)

	// Выполнение HTTP GET запроса
	statusCode, err := checkURL(url)
	if err != nil {
		logger.Printf("Ошибка при выполнении HTTP запроса: %v\n", err)
		fmt.Printf("Ошибка HTTP запроса: %v\n", err)
		return
	}
	logger.Printf("HTTP GET запрос на URL '%s', статус ответа: %d\n", url, statusCode)

	// Проверка успешного статуса
	if statusCode == http.StatusOK {
		logger.Println("Запрос завершился успешно (статус 200)")
	} else {
		logger.Printf("Неверный статус ответа: %d\n", statusCode)
	}

	// Сохранение результата в указанный файл
	outputFile, err := os.Create(*outputPtr)
	if err != nil {
		logger.Printf("Ошибка при создании файла для вывода: %v\n", err)
		fmt.Printf("Ошибка создания файла: %v\n", err)
		return
	}
	defer outputFile.Close()

	result := fmt.Sprintf("Массив чисел: %v\nСумма чисел: %d\nСтатус HTTP запроса: %d\n", numbers, sum, statusCode)
	outputFile.WriteString(result)

	logger.Println("Результат сохранён в файл:", *outputPtr)
	logger.Println("Программа завершена успешно")
	fmt.Println("Программа завершена успешно. Результат сохранён в файл:", *outputPtr)
}
