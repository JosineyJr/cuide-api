package admission_criteria

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"

	e "cuide/api/resource/common/err"
	l "cuide/api/resource/common/log"
	ctxUtil "cuide/util/ctx"
	validatorUtil "cuide/util/validator"
)

type API struct {
	logger     *zerolog.Logger
	validator  *validator.Validate
	repository *Repository
}

func New(logger *zerolog.Logger, validator *validator.Validate, db *sql.DB) *API {
	return &API{
		logger:     logger,
		validator:  validator,
		repository: NewRepository(db),
	}
}

// List godoc
//
//	@summary		List admission criteria
//	@description	List admission criteria
//	@tags			admission-criteria
//	@accept			json
//	@produce		json
//	@success		200	{array}		DTO
//	@failure		500	{object}	err.Error
//	@router			/admission-criteria [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	admissionCriteriaList, err := a.repository.List()
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	if len(admissionCriteriaList) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	if err := json.NewEncoder(w).Encode(admissionCriteriaList.ToDto()); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Create godoc
//
//	@summary		Create admission criteria
//	@description	Create admission criteria
//	@tags			admission-criteria
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"AdmissionCriteria form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/admission-criteria [post]
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.BadRequest(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	newAdmissionCriteria := form.ToModel()

	admissionCriteria, err := a.repository.Create(&newAdmissionCriteria)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Uint8("id", admissionCriteria.ID).Msg("new book created")
	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read admission criteria
//	@description	Read admission criteria
//	@tags			admission-criteria
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"AdmissionCriteria ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/admission-criteria/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	admissionCriteria, err := a.repository.Read(uint8(id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	dto := admissionCriteria.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Update godoc
//
//	@summary		Update admission criteria
//	@description	Update admission criteria
//	@tags			admission-criteria
//	@accept			json
//	@produce		json
//	@param			id		path	string	true	"AdmissionCriteria ID"
//	@param			body	body	Form	true	"AdmissionCriteria form"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/admission-criteria/{id} [put]
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	form := &Form{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.BadRequest(w, e.RespJSONDecodeFailure)
		return
	}

	if err := a.validator.Struct(form); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
			e.ServerError(w, e.RespJSONEncodeFailure)
			return
		}

		e.ValidationErrors(w, respBody)
		return
	}

	admissionCriteria := form.ToModel()
	admissionCriteria.ID = uint8(id)

	rows, err := a.repository.Update(&admissionCriteria)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataUpdateFailure)
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Uint8("id", admissionCriteria.ID).Msg("book updated")
}

// Delete godoc
//
//	@summary		Delete admission criteria
//	@description	Delete admission criteria
//	@tags			admission-criteria
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"AdmissionCriteria ID"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/admission-criteria/{id} [delete]
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	rows, err := a.repository.Delete(uint8(id))
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataRemoveFailure)
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Uint8("id", uint8(id)).Msg("book deleted")
}
