package admission_criteria

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

func (r *Repository) List() (AdmissionCriteriaList, error) {
	admissionCriteriaList := make([]*AdmissionCriteria, 0)

	rows, err := r.db.Query("SELECT id, nome FROM public.criterios_admissao;")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var admissionCriteria AdmissionCriteria
		rows.Scan(&admissionCriteria.ID, &admissionCriteria.Name)

		admissionCriteriaList = append(admissionCriteriaList, &admissionCriteria)
	}

	return admissionCriteriaList, nil
}

func (r *Repository) Create(admissionCriteria *AdmissionCriteria) (*AdmissionCriteria, error) {
	err := r.db.QueryRow("INSERT INTO public.criterios_admissao (nome) VALUES ($1) RETURNING id;", admissionCriteria.Name).
		Scan(&admissionCriteria.ID)
	if err != nil {
		return nil, err
	}

	return admissionCriteria, nil
}

func (r *Repository) Read(id uint8) (*AdmissionCriteria, error) {
	var admissionCriteria AdmissionCriteria
	err := r.db.QueryRow("SELECT id, nome FROM public.criterios_admissao r WHERE r.id = $1;", id).
		Scan(&admissionCriteria.ID, &admissionCriteria.Name)
	if err != nil {
		return nil, err
	}

	return &admissionCriteria, nil
}

func (r *Repository) Update(admissionCriteria *AdmissionCriteria) (int64, error) {
	result, err := r.db.Exec(
		"UPDATE public.criterios_admissao SET nome = $1 WHERE id= $2;",
		admissionCriteria.Name,
		admissionCriteria.ID,
	)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}

func (r *Repository) Delete(id uint8) (int64, error) {
	result, err := r.db.Exec("DELETE FROM public.criterios_admissao WHERE id = $1;", id)
	rows, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rows, err
}
