# 📩 Email Collector with Fiber and SQLite

This is a simple web app built with the [Fiber](https://github.com/gofiber/fiber) web framework and the [GORM](https://gorm.io) ORM that collects user emails and client information (IP address and user agent) and stores it in a SQLite database.

## 🚀 Quick Start

1. Clone the repository:
   ```
   git clone https://github.com/xeronith/email-collector.git
   ```

2. Navigate to the project directory:
   ```
   cd email-collector
   ```

3. Build the Docker image:
   ```
   docker build -t email-collector .
   ```

4. Run the Docker container and mount a volume for the database:
   ```
   docker run \
      -p 8080:8080 \
      -e POSTMARK_TOKEN="your-token" \
      -e POSTMARK_FROM="from@yourdomain.com" \
      -e POSTMARK_TEMPLATE_ALIAS="template-alias" \
      -v /path/to/local/db:/app/db \
      email-collector
   ```

5. Visit `http://localhost:8080` in your web browser to see the app in action.

## 📝 API Endpoints

The app has a single endpoint that collects user emails and client information and returns a JSON response:

- `POST /subscribe`

  Collects the user's email address and any additional data provided in the request body, as well as the client IP address information and user agent. Stores this information in the SQLite database and returns a JSON response containing a message.

  **Request Body:**
  ```
  {
    "email": "user@example.com",
    "data": "additional data to log"
  }
  ```

  **Response Body:**
  ```
  {
    "message": "your message goes here"
  }
  ```

## 📚 Dependencies

The app uses the following third-party dependencies:

- [Fiber](https://github.com/gofiber/fiber) - Web framework for Go.
- [GORM](https://gorm.io) - ORM library for Go.
- [SQLite3](https://github.com/mattn/go-sqlite3) - SQLite3 driver for Go.
- [Postmark](https://postmarkapp.com) - Email Delivery Service

## 📝 License

This project is licensed under the [MIT License](LICENSE).
