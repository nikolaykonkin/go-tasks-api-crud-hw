package storage

import "tasks-api/internal/models"

// Storage описывает контракт хранилища задач.
// Любая реализация (in-memory, БД и т.д.) должна удовлетворять этому интерфейсу,
// что позволяет заменить хранение данных без изменения обработчиков.
type Storage interface {
	List() []models.Task
	Create(models.Task) (models.Task, error)
	Get(id int) (models.Task, bool)
	Update(id int, task models.Task) (models.Task, error)
	Delete(id int) error
}
