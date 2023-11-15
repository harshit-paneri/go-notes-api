package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Note struct {
	ID   uint32 `json:"id"`
	Note string `json:"note"`
}

type CreateNoteRequest struct {
	SID  string `json:"sid"`
	Note string `json:"note"`
}

type CreateNoteResponse struct {
	ID uint32 `json:"id"`
}

type DeleteNoteRequest struct {
	SID string `json:"sid"`
	ID  uint32 `json:"id"`
}

func handleSignup(c *gin.Context) {
	var request SignupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	c.JSON(http.StatusOK, request)
}

func handleLogin(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	sessionID := "unique_session_id"

	c.JSON(http.StatusOK, gin.H{"sid": sessionID})
}

func handleListNotes(c *gin.Context) {
	sid := c.Query("sid") // Assuming the sid is passed as a query parameter

	if sid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session ID"})
		return
	}

	notes := []Note{
		{ID: 1, Note: "First note"},
		{ID: 2, Note: "Second note"},
		{ID: 3, Note: "Third note"},
	}

	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func handleCreateNote(c *gin.Context) {
	var request CreateNoteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	response := CreateNoteResponse{ID: 123}

	c.JSON(http.StatusOK, response)
}

func handleDeleteNote(c *gin.Context) {
	var request DeleteNoteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if request.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid note ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

func main() {
	router := gin.Default()

	router.POST("/signup", handleSignup)
	router.POST("/login", handleLogin)
	router.GET("/notes", handleListNotes)
	router.POST("/notes", handleCreateNote)
	router.DELETE("/notes", handleDeleteNote)

	router.Run("localhost:8080")
}