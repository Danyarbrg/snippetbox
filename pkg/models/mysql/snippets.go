package mysql

import (
	"database/sql"
	"errors"

	"github.com/Danyarbrg/snippetbox/pkg/models"
)

// обьявили тип соедниения
type SnippetModel struct {
	DB *sql.DB
}

// вставляет новые сниппеты в БД
func (sm *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	// используй метод Exec() во встроенном пуле подключений для выполнения инструкции.
	// первый параметр это SQL оператор за которым следуют значения
	// Этот метод возвращает объект sql.Result, который содержит некоторую базовую 
	// информацию о том, что произошло при выполнении инструкции.
	result, err := sm.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	// с поомщью методы .LastInsertId() получаем ID нашей последней вставленной записи
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	// возвращаемый ID имеет тип int64, мы его конвертируем в просто int
	return int(id), nil
}

// возвращает сниппет по id
func (sm *SnippetModel) Get(id int) (*models.Snippet, error) {
	s := &models.Snippet{}

	err := sm.DB.QueryRow("SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() AND id = ?", id,
	).Scan(&s.ID,&s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

// возвращает 10 последних созданных снипетов
func (sm *SnippetModel) Latest() ([]*models.Snippet, error) {
	rows, err := sm.DB.Query("SELECT id, title, content, created, expires FROM snippets WHERE expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10")
	if err != nil {
		return nil, err
	}

	snippets := []*models.Snippet{}

	defer rows.Close()

	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.ID,&s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}