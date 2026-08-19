package storage

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"tasks-api/internal/models"
)

// MemoryStorage — потокобезопасная реализация Storage, хранящая задачи в памяти.
// Доступ к карте задач защищён sync.RWMutex: чтение (List, Get) использует
// RLock и может выполняться параллельно из разных горутин, запись
// (Create, Update, Delete) использует эксклюзивный Lock.
type MemoryStorage struct {
	mu     sync.RWMutex
	tasks  map[int]models.Task
	nextID int
}

// NewMemoryStorage создаёт новое пустое in-memory хранилище.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tasks:  make(map[int]models.Task),
		nextID: 1,
	}
}

// List возвращает все задачи, отсортированные по ID для стабильного порядка вывода.
func (s *MemoryStorage) List() []models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]models.Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ID < tasks[j].ID
	})

	return tasks
}

// Create создаёт новую задачу: ID генерируется на сервере (игнорируя
// любой ID, присланный клиентом), CreatedAt проставляется автоматически,
// если не передан клиентом.
func (s *MemoryStorage) Create(task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task.ID = s.nextID
	s.nextID++

	if task.CreatedAt == "" {
		task.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}

	s.tasks[task.ID] = task
	return task, nil
}

// Get возвращает задачу по ID и флаг, найдена ли она.
func (s *MemoryStorage) Get(id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	return task, ok
}

// Update полностью заменяет задачу с данным ID. ID и CreatedAt сохраняются
// от оригинальной записи независимо от того, что передал клиент в теле запроса.
func (s *MemoryStorage) Update(id int, task models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.tasks[id]
	if !ok {
		return models.Task{}, fmt.Errorf("task with id %d not found", id)
	}

	task.ID = id
	task.CreatedAt = existing.CreatedAt
	s.tasks[id] = task

	return task, nil
}

// Delete удаляет задачу по ID.
func (s *MemoryStorage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task with id %d not found", id)
	}

	delete(s.tasks, id)
	return nil
}
