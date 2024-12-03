package regionals

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

func (r *Repository) List() (Regionals, error) {
	regionals := make([]*Regional, 0)

	rows, err := r.db.Query("SELECT id, nome FROM public.regionais;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var regional Regional
		rows.Scan(&regional.ID, &regional.Name)

		regionals = append(regionals, &regional)
	}

	return regionals, nil
}

func (r *Repository) Create(regional *Regional) (*Regional, error) {
	err := r.db.QueryRow("INSERT INTO public.regionais (nome) VALUES ($1) RETURNING id;", regional.Name).
		Scan(&regional.ID)
	if err != nil {
		return nil, err
	}

	return regional, nil
}

func (r *Repository) Read(id uint8) (*Regional, error) {
	var regional Regional
	err := r.db.QueryRow("SELECT id, nome FROM public.regionais r WHERE r.id = $1;", id).
		Scan(&regional.ID, &regional.Name)
	if err != nil {
		return nil, err
	}

	return &regional, nil
}

func (r *Repository) Update(regional *Regional) (int64, error) {
	result, err := r.db.Exec(
		"UPDATE public.regionais SET nome = $1 WHERE id= $2;",
		regional.Name,
		regional.ID,
	)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}

func (r *Repository) Delete(id uint8) (int64, error) {
	result, err := r.db.Exec("DELETE FROM public.regionais WHERE id = $1;", id)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}
