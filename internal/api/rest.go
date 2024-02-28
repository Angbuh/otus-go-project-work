package api

import (
	"fmt"
	"my_notes_project/internal/core"
	"my_notes_project/internal/entities"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
)

type RestAPI struct {
	app    *fiber.App
	logger *logrus.Logger
	core   core.ServiceCore
}

func NewRestAPI(core core.ServiceCore, logger *logrus.Logger) *RestAPI {
	engine := html.New("./web/templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static/", "./web/static")

	return &RestAPI{
		app:    app,
		core:   core,
		logger: logger,
	}
}

func (r *RestAPI) HandlersInit() error {
	r.app.Get("favicon.ico", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("web/static/favicon.ico")
	}).Get("/", func(ctx *fiber.Ctx) error {
		m := fiber.Map{
			"IsAuthed": false,
			"Title":    "Notes",
		}

		if username := ctx.Cookies("username"); username != "" {
			m["IsAuthed"] = true
			notes, err := r.core.GetNotesByUserName(username)
			if err != nil {
				// TODO: ctx.Status...
				return err
			}

			m["Notes"] = notes
		}

		r.logger.Debug(m)
		return ctx.Render("index", m)
	}).Post("/reg", func(ctx *fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var name, password1, password2 string
		if vals, exists := form.Value["username"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no username")
		} else {
			name = vals[0]
		}

		if vals, exists := form.Value["password1"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no password")
		} else {
			password1 = vals[0]
		}

		if vals, exists := form.Value["password2"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no password")
		} else {
			password2 = vals[0]
		}

		if err = r.core.RegisterUser(name, password1, password2); err != nil {
			return err
		}

		return ctx.RedirectBack("/")
	}).Post("/auth", func(ctx *fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var name, password string
		if vals, exists := form.Value["username"]; !exists || len(vals) == 0 {
			r.logger.Error("no username")
			return ctx.Status(fiber.StatusBadRequest).SendString("no username")
		} else {
			name = vals[0]
		}

		if vals, exists := form.Value["password"]; !exists || len(vals) == 0 {
			r.logger.Error("no password")
			return ctx.Status(fiber.StatusBadRequest).SendString("no password")
		} else {
			password = vals[0]
		}

		isValid, err := r.core.IsValidUserCredentials(name, password)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		} else if !isValid {
			return fmt.Errorf("invalid credentials")
		}

		ctx.Cookie(&fiber.Cookie{
			Name:  "username",
			Value: name,
			Path:  "/",
		})

		return ctx.RedirectBack("/")
	}).Get("/logout", func(ctx *fiber.Ctx) error {
		username := ctx.Cookies("username")
		if username == "" {
			return fmt.Errorf("not authed")
		}

		ctx.ClearCookie("username")

		return ctx.RedirectBack("/")
	}).Post("/note/add", func(ctx *fiber.Ctx) error {
		username := ctx.Cookies("username")
		if username == "" {
			r.logger.Error("not authed")
			return fiber.NewError(fiber.StatusBadRequest, "not authed")
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var title, content string
		if vals, exists := form.Value["title"]; !exists || len(vals) == 0 {
			r.logger.Error("no title")
			return ctx.Status(fiber.StatusBadRequest).SendString("no title")
		} else {
			title = vals[0]
		}

		if vals, exists := form.Value["content"]; !exists || len(vals) == 0 {
			r.logger.Error("no content")
			return ctx.Status(fiber.StatusBadRequest).SendString("no content")
		} else {
			content = vals[0]
		}

		if err = r.core.AddNoteToUserByName(username, &entities.Note{
			Title:   title,
			Content: content,
		}); err != nil {
			return err
		}

		return ctx.RedirectBack("/")
	}).Get("/note/remove/:id", func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		r.logger.Debug(id)
		if err := r.core.RemoveNoteByID(id); err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return ctx.RedirectBack("/")
	}).Post("/note/update/:id", func(ctx *fiber.Ctx) error {
		username := ctx.Cookies("username")
		if username == "" {
			return fmt.Errorf("not authed")
		}

		r.logger.Debug(username)

		id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		r.logger.Debug(id)

		form, err := ctx.MultipartForm()
		var title, content string
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if vals, exists := form.Value["title"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no title")
		} else {
			title = vals[0]

		}

		if vals, exists := form.Value["content"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no content")
		} else {
			content = vals[0]
		}

		err = r.core.UpdateNoteByUserName(username, &entities.Note{
			ID:      id,
			Title:   title,
			Content: content,
		})
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return ctx.RedirectBack("/")
	})

	return nil
}

func (r *RestAPI) Listen(addr string) error {
	return r.app.Listen(addr)
}
