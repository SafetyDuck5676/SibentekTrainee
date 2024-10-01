package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Тест для функции readNumbers (чтение из файла)
func TestReadNumbersFromFile(t *testing.T) {
	// Создаём временный файл с JSON содержимым
	testFileName := "test_numbers.json"
	fileContent := `[1, 2, 3, 4, 5]`
	err := os.WriteFile(testFileName, []byte(fileContent), 0644)
	if err != nil {
		t.Fatalf("Ошибка при создании временного файла: %v", err)
	}
	defer os.Remove(testFileName) // Удаляем файл после теста

	// Чтение чисел из файла
	numbers, err := readNumbers(testFileName)
	if err != nil {
		t.Fatalf("Ошибка при чтении чисел из файла: %v", err)
	}

	// Проверка, что числа прочитаны корректно
	expected := []int{1, 2, 3, 4, 5}
	for i, num := range numbers {
		if num != expected[i] {
			t.Errorf("Ожидалось %d, получено %d на позиции %d", expected[i], num, i)
		}
	}
}

// Тест для функции sumNumbers (суммирование чисел)
func TestSumNumbers(t *testing.T) {
	// Входные данные
	numbers := []int{1, 2, 3, 4, 5}
	// Ожидаемый результат
	expected := 15

	// Вызываем функцию
	result := sumNumbers(numbers)

	// Проверка
	if result != expected {
		t.Errorf("Ожидалось %d, получено %d", expected, result)
	}
}

// Тест для функции checkURL (проверка HTTP статуса)
func TestCheckURL(t *testing.T) {
	// Создаем тестовый HTTP сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK) // Отвечаем статусом 200
	}))
	defer server.Close()

	// Вызываем функцию
	statusCode, err := checkURL(server.URL)
	if err != nil {
		t.Fatalf("Ошибка при вызове checkURL: %v", err)
	}

	// Проверка статуса
	if statusCode != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d", statusCode)
	}
}

// Тест для функции checkURL при ошибке
func TestCheckURL_Error(t *testing.T) {
	// Неправильный URL для проверки
	statusCode, err := checkURL("http://invalid-url")
	if err == nil {
		t.Fatal("Ожидалась ошибка при выполнении HTTP запроса на неправильный URL, но ошибка не возникла")
	}

	// Проверка, что статус не 200
	if statusCode != 0 {
		t.Errorf("Ожидался статус 0 при ошибке, получен %d", statusCode)
	}
}
