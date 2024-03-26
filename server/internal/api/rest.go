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
	// механизм,для связи с серверами через HTTP,логгером...
	app    *fiber.App
	logger *logrus.Logger
	core   core.ServiceCore
}

func NewRestAPI(core core.ServiceCore, logger *logrus.Logger) *RestAPI {
	//Новый экземпляр для шаблонизатора
	//Новый экземпляр  для файбер,в которую передаем дополнительные параметры конфигурации
	engine := html.New("./web/templates", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	//Для улучшения внешнего вид домашней страницы, добавила в проект статические файлы
	app.Static("/static/", "./web/static")

	return &RestAPI{
		app:    app,
		logger: logger,
		core:   core,
	}
}

func (r *RestAPI) HandlersInit() error {
	// Добавляем иконку на сайт, обращаемся к статическому файлу
	r.app.Get("favicon.ico", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("web/static/favicon.ico")
	}).Get("/", func(ctx *fiber.Ctx) error {
		// Создаем гет запрос, который у нас будет главной страницей...
		m := fiber.Map{
			"IsAuthed": false,
			"Title":    "Notes",
		}

		// Создаем ключ, проверяем на пустую строку
		// Получаем заметки пользователя,если они у него были
		if username := ctx.Cookies("username"); username != "" {
			m["IsAuthed"] = true
			notes, err := r.core.GetNotesByUserName(username)
			if err != nil {
				r.logger.Error(err)
				return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
			}

			m["Notes"] = notes
		}

		r.logger.Debug(m)
		return ctx.Render("index", m)
	}).Post("/reg", func(ctx *fiber.Ctx) error {
		// Создаем пост запрос для регистрации пользователя
		// Создаем форму для анализа, поступивших данных
		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		//Проверяем на валидацию имя пользователя
		var name, password1, password2 string
		if vals, exists := form.Value["username"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no username")
		} else {
			name = vals[0]
		}

		//Проверяем на валидацию полученный пароль
		if vals, exists := form.Value["password1"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no password")
		} else {
			password1 = vals[0]
		}

		//Проверяем на соответствие этому паролю уже ранее полученный пароль
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
		//Создаем обработчик для аутентификации пользователя
		//Создаем форму
		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		//Проверяем имя пользователя
		var name, password string
		if vals, exists := form.Value["username"]; !exists || len(vals) == 0 {
			r.logger.Error("no username")
			return ctx.Status(fiber.StatusBadRequest).SendString("no username")
		} else {
			name = vals[0]
		}

		//Проверяем пароль пользователя
		if vals, exists := form.Value["password"]; !exists || len(vals) == 0 {
			r.logger.Error("no password")
			return ctx.Status(fiber.StatusBadRequest).SendString("no password")
		} else {
			password = vals[0]
		}

		//Проверяем действительно ли сопадает с данными пользователя
		isValid, err := r.core.IsValidUserCredentials(name, password)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		} else if !isValid {
			return fmt.Errorf("invalid credentials")
		}

		//Читаем поля формы и позволяем пользователю зайти
		ctx.Cookie(&fiber.Cookie{
			Name:  "username",
			Value: name,
			Path:  "/",
		})

		return ctx.RedirectBack("/")
	}).Get("/logout", func(ctx *fiber.Ctx) error {
		//Проверяем авторизован ли пользователь, если да, то удаляем его куки
		username := ctx.Cookies("username")
		if username == "" {
			return fmt.Errorf("not authed")
		}

		ctx.ClearCookie("username")

		return ctx.RedirectBack("/")
	}).Post("/note/add", func(ctx *fiber.Ctx) error {
		//Создаем обработчик для добавления заметок
		//Для начала проверяем на авторизацию
		username := ctx.Cookies("username")
		if username == "" {
			r.logger.Error("not authed")
			return fiber.NewError(fiber.StatusBadRequest, "not authed")
		}

		//Создаем форму
		form, err := ctx.MultipartForm()
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		//Проверяем наличие заголовка заметки
		var title, content string
		if vals, exists := form.Value["title"]; !exists || len(vals) == 0 {
			r.logger.Error("no title")
			return ctx.Status(fiber.StatusBadRequest).SendString("no title")
		} else {
			title = vals[0]
		}

		//Проверяем содержание заметки
		if vals, exists := form.Value["content"]; !exists || len(vals) == 0 {
			r.logger.Error("no content")
			return ctx.Status(fiber.StatusBadRequest).SendString("no content")
		} else {
			content = vals[0]
		}

		//Читаем поля формы и добавляем заметку
		if err = r.core.AddNoteToUserByName(username, &entities.Note{
			Title:   title,
			Content: content,
		}); err != nil {
			return err
		}

		return ctx.RedirectBack("/")
	}).Get("/note/remove/:id", func(ctx *fiber.Ctx) error {
		//Создаем обработчик гет для удаления статьи
		//Получаем id, парсим его, получаем значение
		id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		//По полученному id удаляем заметку
		r.logger.Debug(id)
		if err := r.core.RemoveNoteByID(id); err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		return ctx.RedirectBack("/")
	}).Post("/note/update/:id", func(ctx *fiber.Ctx) error {
		//Создаем обработчик гет для обновления статьи
		//Проверяем на авторизацию пользователя
		username := ctx.Cookies("username")
		if username == "" {
			return fmt.Errorf("not authed")
		}

		r.logger.Debug(username)

		//Получаем id от пользователя, парсим его, получаем значение, по которому нужно обновить заметку
		id, err := strconv.ParseUint(ctx.Params("id"), 10, 64)
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		r.logger.Debug(id)

		//Создаем форму для будущей заметки
		form, err := ctx.MultipartForm()
		var title, content string
		if err != nil {
			r.logger.Error(err)
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		//Проверяем  наличие заголовка заметки
		if vals, exists := form.Value["title"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no title")
		} else {
			title = vals[0]

		}

		//Проверяем содержание заметки
		if vals, exists := form.Value["content"]; !exists || len(vals) == 0 {
			return ctx.Status(fiber.StatusBadRequest).SendString("no content")
		} else {
			content = vals[0]
		}

		//Читаем поля формы и обновляем соответствующие поля у заметки
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
	//
	return r.app.Listen(addr)
}
