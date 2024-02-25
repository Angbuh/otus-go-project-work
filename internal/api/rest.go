package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"my_notes_project/internal/database"
	"my_notes_project/internal/entities"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sirupsen/logrus"
)

type RestAPI struct {
	app    *fiber.App
	db     database.DBRepository
	logger *logrus.Logger
}

func NewRestAPI(db database.DBRepository, logger *logrus.Logger) *RestAPI {
	engine := html.New("./web/templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static/", "./web/static")

	return &RestAPI{
		app:    app,
		db:     db,
		logger: logger,
	}
}

func (r *RestAPI) HandlersInit() error {
	r.app.Get("favicon.ico", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("web/static/favicon.ico")
	}).Get("/", func(ctx *fiber.Ctx) error {
		m := fiber.Map{
			"IsAuthed": false,
			"Title": "Notes",
		}

		if username := ctx.Cookies("username"); username != "" {
			m["IsAuthed"] = true
			user, err := r.db.GetUserByName(username)
			if err != nil {
				return err
			}

			notes, err := r.db.GetNotesByUserId(user.ID)
			if err != nil {
				return err
			}

			m["Notes"] = notes
		}

		r.logger.Debug(m)
		return ctx.Render("index", m)
	}).Get("/note/get/:id", func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		r.logger.Debug(id)
		var note *entities.Note
		if note, err = r.db.GetNoteById(id); err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		noteBytes, err := json.Marshal(note)
		if err != nil {
			r.logger.Error(err)
			return err
		}

		return ctx.SendStream(bytes.NewReader(noteBytes))
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

		if password1 != password2 {
			return fmt.Errorf("")
		}

		user := &entities.User{
			Name:     name,
			Password: password1,
		}

		id, err := r.db.AddUser(user)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		user.ID = id

		r.logger.Debug(user)

		return nil
	}).Post("/auth", func(ctx *fiber.Ctx) error {
		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var name, password string
		if vals, exists := form.Value["username"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no username")
		} else {
			name = vals[0]
		}

		if vals, exists := form.Value["password"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no password")
		} else {
			password = vals[0]
		}

		var user *entities.User
		if user, err = r.db.GetUserByName(name); err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if user.Password != password {
			r.logger.Error("invalid data")
			return fiber.NewError(fiber.StatusBadRequest, "invalid data")
		}

		ctx.Cookie(&fiber.Cookie{
			Name:  "username",
			Value: name,
			Path:  "/",
		})

		return nil
	}).Get("/logout", func(ctx *fiber.Ctx) error {
		username := ctx.Cookies("username")
		if username == "" {
			return fmt.Errorf("not authed")
		}

		ctx.ClearCookie("username")

		return nil
	}).Get("/note/get", func(ctx *fiber.Ctx) error {
		notes, err := r.db.GetAllNotes()
		if err != nil {
			r.logger.Error(err)
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		noteBytes, err := json.Marshal(notes)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return ctx.SendStream(bytes.NewReader(noteBytes))
	}).Post("/note/add", func(ctx *fiber.Ctx) error {
		username := ctx.Cookies("username")
		if username == "" {
			return fiber.NewError(fiber.StatusBadRequest, "not authed")
		}

		user, err := r.db.GetUserByName(username)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		var title, content string
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

		noteID, err := r.db.AddNote(&entities.Note{
			Title:   title,
			Content: content,
			UserID:  user.ID,
		})
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		r.logger.Debug(noteID)
		return nil
	}).Get("/note/remove/:id", func(ctx *fiber.Ctx) error {
		id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		r.logger.Debug(id)
		if err := r.db.RemoveNoteByID(id); err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return nil
	}).Patch("/update/:id", func(ctx *fiber.Ctx) error {
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

		notes, err := r.db.GetNotesByUserName(username)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		note, exists := notes[id]
		if !exists {
			r.logger.Error("not found")
			return ctx.Status(fiber.StatusBadRequest).SendString("not found")
		}

		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		if vals, exists := form.Value["title"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no title")
		} else {
			note.Title = vals[0]
		}

		if vals, exists := form.Value["content"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no content")
		} else {
			note.Content = vals[0]
		}

		err = r.db.UpdateNote(note)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return nil
	})

	return nil
}

func (r *RestAPI) Listen(addr string) error {
	return r.app.Listen(addr)
}
