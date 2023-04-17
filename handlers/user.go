package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func GetUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		user := User{}

		err := db.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
		}

		return c.JSON(http.StatusOK, user)
	}
}

func CreateUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.FormValue("name")
		email := c.FormValue("email")

		var userID int
		err := db.QueryRow("INSERT INTO users(name, email) VALUES($1, $2) RETURNING id", name, email).Scan(&userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
		}

		user := User{
			ID:    userID,
			Name:  name,
			Email: email,
		}

		return c.JSON(http.StatusCreated, user)
	}
}

func UpdateUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		name := c.FormValue("name")
		email := c.FormValue("email")

		var userID int
		err := db.QueryRow("UPDATE users SET name=$1, email=$2 WHERE id=$3 RETURNING id", name, email, id).Scan(&userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user"})
		}

		user := User{
			ID:    userID,
			Name:  name,
			Email: email,
		}

		return c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		var userID int
		err := db.QueryRow("DELETE FROM users WHERE id=$1 RETURNING id", id).Scan(&userID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to delete user"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
	}
}
