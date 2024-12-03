package reference_ways

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

func (r *Repository) List() (ReferenceWays, error) {
	referenceWays := make([]*ReferenceWay, 0)

	rows, err := r.db.Query("SELECT id, nome FROM public.forma_encaminhamento;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var referenceWay ReferenceWay
		rows.Scan(&referenceWay.ID, &referenceWay.Name)

		referenceWays = append(referenceWays, &referenceWay)
	}

	return referenceWays, nil
}

func (r *Repository) Create(referenceWay *ReferenceWay) (*ReferenceWay, error) {
	err := r.db.QueryRow("INSERT INTO public.forma_encaminhamento (nome) VALUES ($1) RETURNING id;", referenceWay.Name).
		Scan(&referenceWay.ID)
	if err != nil {
		return nil, err
	}

	return referenceWay, nil
}

func (r *Repository) Read(id uint8) (*ReferenceWay, error) {
	var referenceWay ReferenceWay
	err := r.db.QueryRow("SELECT id, nome FROM public.forma_encaminhamento r WHERE r.id = $1;", id).
		Scan(&referenceWay.ID, &referenceWay.Name)
	if err != nil {
		return nil, err
	}

	return &referenceWay, nil
}

func (r *Repository) Update(referenceWay *ReferenceWay) (int64, error) {
	result, err := r.db.Exec(
		"UPDATE public.forma_encaminhamento SET nome = $1 WHERE id= $2;",
		referenceWay.Name,
		referenceWay.ID,
	)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}

func (r *Repository) Delete(id uint8) (int64, error) {
	result, err := r.db.Exec("DELETE FROM public.forma_encaminhamento WHERE id = $1;", id)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}
