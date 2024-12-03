package places

import (
	"context"
	txUtil "cuide/util/db-tx"
	"database/sql"
	"encoding/json"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) List(page uint8) (Places, error) {
	places := make([]*Place, 0)

	rows, err := r.db.Query(`
	SELECT
    s.id AS servico_id,
    s.nome AS servico_nome,
    s.endereco AS servico_endereco,
    s.contato AS servico_contato,
    s.site AS servico_site,
    s.observacoes AS servico_observacoes,
    s.maps_link AS servico_maps_link,
    jsonb_build_object('id', ts.id, 'name', ts.nome) AS tipo_servico,
    jsonb_build_object('id', e.id, 'name', e.nome) AS eixo,
    jsonb_agg(DISTINCT jsonb_build_object('id', ta.id, 'name', ta.nome)) AS tipos_atendimento,
    jsonb_agg(DISTINCT jsonb_build_object('id', ca.id, 'name', ca.nome)) AS criterios_admissao,
    jsonb_agg(DISTINCT jsonb_build_object('id', fe.id, 'name', fe.nome)) AS formas_encaminhamento,
    jsonb_agg(DISTINCT jsonb_build_object('id', r.id, 'name', r.nome)) AS regionais
	FROM
			public.servico s
	LEFT JOIN public.tipo_servico ts
			ON s.tipo_servico_id = ts.id
	LEFT JOIN public.eixo e
			ON s.eixo_id = e.id
	LEFT JOIN public.tipo_atendimento_servico tas
			ON s.id = tas.servico_id
	LEFT JOIN public.tipo_atendimento ta
			ON tas.tipo_atendimento_id = ta.id
	LEFT JOIN public.criterios_admissao_servico cas
			ON s.id = cas.servico_id
	LEFT JOIN public.criterios_admissao ca
			ON cas.criterio_admissao_id = ca.id
	LEFT JOIN public.forma_encaminhamento_servico fes
			ON s.id = fes.servico_id
	LEFT JOIN public.forma_encaminhamento fe
			ON fes.forma_encaminhamento_id = fe.id
	LEFT JOIN public.regionais_servico rs
			ON s.id = rs.servico_id
	LEFT JOIN public.regionais r
			ON rs.regional_id = r.id
	GROUP BY
			s.id, ts.id, ts.nome, e.id, e.nome
	ORDER BY
			s.id
	LIMIT 20 OFFSET $1;`, (page-1)*10)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var (
			place                                                                                                      Place
			serviceTypeJson, segmentJson, attendanceTypesJson, admissionCriteriaJson, referenceWaysJson, regionalsJson string
		)

		err := rows.Scan(
			&place.ID,
			&place.Name,
			&place.Address,
			&place.PhoneNumber,
			&place.Website,
			&place.Observations,
			&place.GoogleMapsLink,
			&serviceTypeJson,
			&segmentJson,
			&attendanceTypesJson,
			&admissionCriteriaJson,
			&referenceWaysJson,
			&regionalsJson,
		)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(serviceTypeJson), &place.ServiceType)
		json.Unmarshal([]byte(segmentJson), &place.Segment)
		json.Unmarshal([]byte(attendanceTypesJson), &place.AttendanceType)
		json.Unmarshal([]byte(admissionCriteriaJson), &place.AdmissionCriteria)
		json.Unmarshal([]byte(referenceWaysJson), &place.ReferenceWays)
		json.Unmarshal([]byte(regionalsJson), &place.Regionals)

		places = append(places, &place)
	}

	return places, nil
}

func (r *Repository) Create(ctx context.Context, place *Place) (*Place, error) {
	err := txUtil.CallTx(ctx, r.db, func(tx *sql.Tx) error {
		var placeID uint8

		err := r.db.QueryRow(
			`INSERT INTO
			public.servico (
				tipo_servico_id,
				nome,
				endereco,
				contato,
				site,
				observacoes,
				eixo_id,
				maps_link
			)
		VALUES
		($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id;`,
			place.ServiceType.ID,
			place.Name,
			place.Address,
			place.PhoneNumber,
			place.Website,
			place.Observations,
			place.Segment.ID,
			place.GoogleMapsLink,
		).Scan(&placeID)
		if err != nil {
			return err
		}

		stmt, err := r.db.Prepare(
			`INSERT INTO public.criterios_admissao_servico (servico_id, criterio_admissao_id) VALUES ($1, $2)`,
		)
		if err != nil {
			return err
		}
		for _, ac := range place.AdmissionCriteria {
			_, err := stmt.Exec(placeID, ac.ID)
			if err != nil {
				return err
			}
		}

		stmt, err = r.db.Prepare(
			`INSERT INTO public.forma_encaminhamento_servico (servico_id, forma_encaminhamento_id) VALUES ($1, $2)`,
		)
		if err != nil {
			return err
		}
		for _, rws := range place.ReferenceWays {
			_, err := stmt.Exec(placeID, rws.ID)
			if err != nil {
				return err
			}
		}

		stmt, err = r.db.Prepare(
			`INSERT INTO public.regionais_servico (servico_id, regional_id) VALUES ($1, $2)`,
		)
		if err != nil {
			return err
		}
		for _, rs := range place.Regionals {
			_, err := stmt.Exec(placeID, rs.ID)
			if err != nil {
				return err
			}
		}

		stmt, err = r.db.Prepare(
			`INSERT INTO public.tipo_atendimento_servico (servico_id, tipo_atendimento_id) VALUES ($1, $2)`,
		)
		if err != nil {
			return err
		}
		for _, atts := range place.AttendanceType {
			_, err := stmt.Exec(placeID, atts.ID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return place, nil
}

func (r *Repository) Read(id uint8) (*Place, error) {
	var place Place
	err := r.db.QueryRow("SELECT id, nome FROM public.eixo r WHERE r.id = $1;", id).
		Scan(&place.ID, &place.Name)
	if err != nil {
		return nil, err
	}

	return &place, nil
}

func (r *Repository) Update(place *Place) (int64, error) {
	result, err := r.db.Exec(
		"UPDATE public.eixo SET nome = $1 WHERE id= $2;",
		place.Name,
		place.ID,
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
