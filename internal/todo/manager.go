package todo

import "fmt"

// Add создаёт новую задачу и добавляет её в список.
func Add(tasks []Task, desc string) []Task {
	// Первая задача должна получить ID, равный 1.
	nextID := 1

	// Перебираем существующие задачи, чтобы найти самый большой ID.
	for _, task := range tasks {
		// Если текущий ID больше или равен следующему ID,
		// увеличиваем его на единицу.
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}

	// Создаём и заполняем новую задачу.
	newTask := Task{
		ID:          nextID,
		Description: desc,
		Done:        false,
	}

	// Добавляем новую задачу в конец списка и возвращаем обновлённый список.
	return append(tasks, newTask)
}

// List возвращает задачи, соответствующие выбранному фильтру.
func List(tasks []Task, filter string) []Task {
	// Если выбран фильтр all, возвращаем все задачи без отбора.
	if filter == "all" {
		return tasks
	}

	// Создаём пустой список для подходящих задач.
	filteredTasks := []Task{}

	// По очереди проверяем каждую задачу.
	for _, task := range tasks {
		// Фильтр done добавляет только выполненные задачи.
		if filter == "done" && task.Done {
			filteredTasks = append(filteredTasks, task)
		}

		// Фильтр pending добавляет только невыполненные задачи.
		if filter == "pending" && !task.Done {
			filteredTasks = append(filteredTasks, task)
		}
	}

	// Возвращаем отфильтрованный список.
	return filteredTasks
}

// CompLete отмечает задачу с указанным ID как обычную выполненую.
func Complete(tasks []Task, id int) ([]Task, error) {
	//Перебираем индексы задач, что бы изменить задачу внутри списка.
	for i := range tasks {
		// Проверяем, совпадает ли ID текущей задачи с указаным ID.
		if tasks[i].ID == id {
			tasks[i].Done = true

			//nil означает, что функция выполнилась без ошибки.
			return tasks, nil
		}
	}

	//Если задача не найдена, возращаем список и понятную ошибку.
	return tasks, fmt.Errorf("задача с ID %d не найдена", id)
}

//Delete удаляет задачу с указанным ID из списка.
func Delete(tasks []Task, id int) ([]Task, error) {
	// Перебераем задачи вместе с их индексами.
	for i, task := range tasks {
		// Проверяем, совпадает ли ID текущей задачи с указаным ID.
		if task.ID == id {
			// Соеденяем задачи до найденого элемента и после него.
			// Сама найденая задача в новыйй список не совпадает.
			tasks = append(tasks[:i], tasks[i+1:]...)

			// Возвращаем обновленный список без ошибоки.
			return tasks, nil
		}

	}

	// Если задача не найдена, возвращаем исходный список и ошибку.
	return tasks, fmt.Errorf("задача с ID не найдена", id)
}
