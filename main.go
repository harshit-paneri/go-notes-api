package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
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

type User struct {
	ID       uint32 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	users      = make(map[string]User) 
	usersMutex sync.Mutex
)

var (
	notes      = make(map[uint32]Note)
	notesMutex sync.Mutex
)

func handleSignup(c *gin.Context) {
	var user User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	if _, exists := users[user.Email]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email already registered"})
		return
	}

	user.ID = uint32(len(users) + 1)

	users[user.Email] = user

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func handleLogin(c *gin.Context) {
	var request LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	usersMutex.Lock()
	defer usersMutex.Unlock()

	user, exists := users[request.Email]
	if !exists || user.Password != request.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	sessionID := fmt.Sprintf("%d", user.ID)

	c.JSON(http.StatusOK, gin.H{"sid": sessionID})
}

func handleListNotes(c *gin.Context) {
	sid := c.Query("sid")

	if sid == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing session ID"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notes": getAllNotes()})
}

func handleCreateNote(c *gin.Context) {
	var request CreateNoteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	noteID := createNote(request.Note)

	c.JSON(http.StatusOK, CreateNoteResponse{ID: noteID})
}

func handleDeleteNote(c *gin.Context) {
	var request DeleteNoteRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	if err := deleteNote(request.ID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}

// getAllNotes returns a slice of all notes.
func getAllNotes() []Note {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	result := make([]Note, 0, len(notes))
	for _, note := range notes {
		result = append(result, note)
	}

	return result
}

// createNote creates a new note and returns its ID.
func createNote(noteText string) uint32 {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	noteID := uint32(len(notes) + 1)
	newNote := Note{ID: noteID, Note: noteText}
	notes[noteID] = newNote

	return noteID
}

// deleteNote deletes a note with the given ID.
func deleteNote(noteID uint32) error {
	notesMutex.Lock()
	defer notesMutex.Unlock()

	if _, exists := notes[noteID]; !exists {
		return fmt.Errorf("Note with ID %d not found", noteID)
	}

	delete(notes, noteID)
	return nil
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