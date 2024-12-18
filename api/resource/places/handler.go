package places

import (
	"database/sql"
	"encoding/json"
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
//	@summary		List places
//	@description	List places
//	@tags			place
//	@accept			json
//	@produce		json
//	@success		200	{array}		DTO
//	@failure		500	{object}	err.Error
//	@router			/places [get]
func (a *API) List(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 8)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespInvalidQueryParamPage)
		return
	}

	places, err := a.repository.List(uint8(page))
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	paginationMetadata, err := a.repository.PaginationMetadata()
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	paginationMetadata.Places = places.ToDto()

	if err := json.NewEncoder(w).Encode(paginationMetadata); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Create godoc
//
//	@summary		Create places
//	@description	Create places
//	@tags			place
//	@accept			json
//	@produce		json
//	@param			body	body	Form	true	"Place form"
//	@success		201
//	@failure		400	{object}	err.Error
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/places [post]
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

	newPlace := form.ToModel()

	referenceWay, err := a.repository.Create(r.Context(), &newPlace)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataInsertFailure)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Uint8("id", referenceWay.ID).Msg("new place created")
	w.WriteHeader(http.StatusCreated)
}

// Read godoc
//
//	@summary		Read place
//	@description	Read place
//	@tags			place
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"Reference Way ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/places/{id} [get]
func (a *API) Read(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 8)
	if err != nil {
		e.BadRequest(w, e.RespInvalidURLParamID)
		return
	}

	place, err := a.repository.Read(uint8(id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	dto := place.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

// Read godoc
//
//	@summary		Read place
//	@description	Read place
//	@tags			place
//	@accept			json
//	@produce		json
//	@param			id	path		string	true	"Place ID"
//	@success		200	{object}	DTO
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/places/{id} [get]
func (a *API) Filter(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	page, err := strconv.ParseUint(r.URL.Query().Get("page"), 10, 8)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespInvalidQueryParamPage)
		return
	}

	filters := Filters{
		ServiceTypes: parseUint8SliceQuery(r, "service-type"),
		Segments:     parseUint8SliceQuery(r, "segment"),
		Regionals:    parseUint8SliceQuery(r, "regional"),
		Name:         parseStringSliceQuery(r, "name"),
	}

	places, conditionals, err := a.repository.Filter(filters, uint8(page))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	paginationMetadata, err := a.repository.FilterPaginationMetadata(conditionals)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataAccessFailure)
		return
	}

	paginationMetadata.Places = places.ToDto()

	if err := json.NewEncoder(w).Encode(paginationMetadata); err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespJSONEncodeFailure)
		return
	}
}

func parseUint8SliceQuery(r *http.Request, param string) []uint8 {
	values := r.URL.Query()[param]
	var result []uint8
	for _, val := range values {
		parsedVal, err := strconv.ParseUint(val, 10, 8)
		if err == nil {
			result = append(result, uint8(parsedVal))
		}
	}
	return result
}

func parseStringSliceQuery(r *http.Request, param string) string {
	values := r.URL.Query()[param]
	var result string
	for _, val := range values {
		result = val
		break
	}
	return result
}

// Update godoc
//
//	@summary		Update place
//	@description	Update place
//	@tags			places
//	@accept			json
//	@produce		json
//	@param			id		path	string	true	"Place ID"
//	@param			body	body	Form	true	"Place form"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		422	{object}	err.Errors
//	@failure		500	{object}	err.Error
//	@router			/places/{id} [put]
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

	place := form.ToModel()
	place.ID = uint8(id)

	rows, err := a.repository.Update(r.Context(), &place)
	if err != nil {
		a.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		e.ServerError(w, e.RespDBDataUpdateFailure)
		return
	}
	if rows == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	a.logger.Info().Str(l.KeyReqID, reqID).Uint8("id", place.ID).Msg("place updated")
}

// Delete godoc
//
//	@summary		Delete places
//	@description	Delete places
//	@tags			places
//	@accept			json
//	@produce		json
//	@param			id	path	string	true	"Place ID"
//	@success		200
//	@failure		400	{object}	err.Error
//	@failure		404
//	@failure		500	{object}	err.Error
//	@router			/places/{id} [delete]
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

	a.logger.Info().Str(l.KeyReqID, reqID).Uint8("id", uint8(id)).Msg("place deleted")
}
