package segments

import (
	"database/sql"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List() (Segments, error) {
	segments := make([]*Segment, 0)

	rows, err := r.db.Query("SELECT id, nome FROM public.eixo;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var segment Segment
		rows.Scan(&segment.ID, &segment.Name)

		segments = append(segments, &segment)
	}

	return segments, nil
}

func (r *Repository) Create(segment *Segment) (*Segment, error) {
	err := r.db.QueryRow("INSERT INTO public.eixo (nome) VALUES ($1) RETURNING id;", segment.Name).
		Scan(&segment.ID)
	if err != nil {
		return nil, err
	}

	return segment, nil
}

func (r *Repository) Read(id uint8) (*Segment, error) {
	var segment Segment
	err := r.db.QueryRow("SELECT id, nome FROM public.eixo r WHERE r.id = $1;", id).
		Scan(&segment.ID, &segment.Name)
	if err != nil {
		return nil, err
	}

	return &segment, nil
}

func (r *Repository) Update(segment *Segment) (int64, error) {
	result, err := r.db.Exec(
		"UPDATE public.eixo SET nome = $1 WHERE id= $2;",
		segment.Name,
		segment.ID,
	)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}

func (r *Repository) Delete(id uint8) (int64, error) {
	result, err := r.db.Exec("DELETE FROM public.eixo WHERE id = $1;", id)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}
