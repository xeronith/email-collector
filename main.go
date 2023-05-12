package main

import (
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/jinzhu/configor"
)

// Define a struct for the user model
type User struct {
	gorm.Model
	Email      string `gorm:"unique"`
	IP         string
	UserAgent  string
	Referer    string
	RemoteAddr string
	Data       string
}

// Define a struct for the configuration
var Config = struct {
	PostmarkToken         string `env:"POSTMARK_TOKEN"`
	PostmarkFrom          string `env:"POSTMARK_FROM"`
	PostmarkTemplateAlias string `env:"POSTMARK_TEMPLATE_ALIAS"`
}{}

func main() {
	// Load the configuration
	if err := configor.Load(&Config); err != nil {
		panic(err)
	}

	// Connect to the SQLite database
	db, err := gorm.Open(sqlite.Open("./db/users.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the database schema
	err = db.AutoMigrate(&User{})
	if err != nil {
		panic(err)
	}

	// Create a new Fiber app
	app := fiber.New()

	// Allow all origins
	app.Use(cors.New())

	// Define a rate limiter middleware to limit requests to 10 requests per minute per IP address
	limiter := limiter.New(limiter.Config{
		Max:        10,
		Expiration: 1 * time.Minute,
	})

	// Apply the rate limiter middleware to the entire app
	app.Use(limiter)

	// Apply the logger middleware to the entire app
	app.Use(logger.New())

	// Define a route for handling POST requests to /subscribe
	app.Post("/subscribe", func(c *fiber.Ctx) error {
		// Parse the email from the request body
		body := new(struct {
			Email string `json:"email"`
			Data  string `json:"data"`
		})
		if err := c.BodyParser(body); err != nil {
			return err
		}

		// Get the client's IP address, user agent, referer, and remote address
		ip := c.IP()
		userAgent := c.Get("User-Agent")
		referer := c.Get("Referer")
		remoteAddr := c.Get("Remote-Addr")

		// Create a new user record in the database
		user := &User{
			Email:      body.Email,
			IP:         ip,
			UserAgent:  userAgent,
			Referer:    referer,
			RemoteAddr: remoteAddr,
			Data:       body.Data,
		}
		if err := db.Create(user).Error; err != nil {
			if err.Error() == "UNIQUE constraint failed: users.email" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Email already exists",
				})
			}
			return err
		}

		// Send confirmation email
		if err := SendEmail(
			Config.PostmarkToken,
			Config.PostmarkFrom,
			body.Email,
			Config.PostmarkTemplateAlias,
		); err != nil {
			log.Println(err)
		}

		// Return a success response
		return c.JSON(fiber.Map{
			"message": "Email collected successfully",
		})
	})

	// Start the HTTP server
	app.Listen(":8080")
}
