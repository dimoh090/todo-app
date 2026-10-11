package storage

import (
	"encoding/json"
	"fmt"
	"os"

	"todo-app/internal/todo"
)

// SaveJSON сохраняет список задач в файл формата JSON.
func SaveJSON(path string, tasks []todo.Task) error {
	// Преобразуем список задач в красиво оформленный JSON.
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("не удалось преобразовать задачи в JSON: %w", err)
	}

	// Записываем полученные данные в указанный файл.
	err = os.WriteFile(path, data, 0644)
	if err != nil {
		return fmt.Errorf("не удалось записать JSON-файл: %w", err)
	}

	// nil означает, что сохранение прошло без ошибки.
	return nil
}

// LoadJSON загружаем список задач из JSON-файла
func LoadJSON(path string) ([]todo.Task, error) {

	// Проверяем, существует ли файл по указанному пути.
	_, err := os.Stat(path)

	// Если файл не существует, подготавливаем пустой список задач.
	if os.IsNotExist(err) {
		emptyTasks := []todo.Task{}

		// Пытаемся создать отсутствующий файл с пустым списком задач.
		saveErr := SaveJSON(path, emptyTasks)
		if saveErr != nil {

			// Возвращаем пустой результат и ошибку создания файла.
			return nil, fmt.Errorf("не удалось создать файл %s: %w", path, saveErr)
		}

		// Файл успешно создан, возвращаем пустой список без ошибки.
		return emptyTasks, nil
	}

	// Если файл существует, но проверить его не удалось, возвращаем ошибку.
	if err != nil {
		return nil, fmt.Errorf("не удалось проверить файл %s: %w", path, err)
	}
	// Читаем всё содержимое существующего JSON-файла.
	data, err := os.ReadFile(path)

	// Если чтение файла не удалось, возвращаем ошибку.
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать файл %s: %w", path, err)
	}

	// Объявляем список, который позже заполним задачами их JSON
	var tasks []todo.Task

	// Разбираем JSON из data и записываем получение задачи в tasks.
	err = json.Unmarshal(data, &tasks)

	// Если JSON нельзя преобразовать в задачи, возвращаем ошибку.
	if err != nil {
		return nil, fmt.Errorf("не удалось разобрать JSON-файл %s: %w:", path, err)
	}

	// Возвращамем загруженный список задач и отсутствие ошибки.
	return tasks, nil
}
