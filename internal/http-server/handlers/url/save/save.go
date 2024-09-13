package save

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	rand "url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(alias string, urlToSave string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("Failed to decode request body", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to decode request body"))
			return
		}

		log.Info("Request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			log.Error("Invalid request", sl.Err(err))
			render.JSON(w, r, resp.Error("Invalid request"))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = rand.NewRandomString(resp.AliasLength)
		}

		id, err := urlSaver.SaveURL(alias, req.URL)
		if errors.Is(err, storage.ErrExists) {
			log.Info("url alreay exists", slog.String("url", req.URL))
			render.JSON(w, r, resp.Error("url alreay exists"))
			return
		}
		if err != nil {
			log.Error("Failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to save url"))
			return
		}

		log.Info("URL saved", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: resp.OK(),
			Alias:    alias,
		})
	}
}
