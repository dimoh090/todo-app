package todo

// Task описывает одну учебную задачу по изучению Go.
type Task struct {
	ID          int    `json:"id"`          // Уникальный номер задачи.
	Description string `json:"description"` // Текст с описанием задачи.
	Done        bool   `json:"done"`        // Статус выполнения: true или false.
}
