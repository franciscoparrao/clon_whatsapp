package services

import (
	"log"
	"os"
	"whatsapp-clone/database"
)

// StoreInterface define los métodos que debe implementar cualquier store
type StoreInterface interface {
	// User methods
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserByID(id string) (*User, error)
	GetUsers() ([]*User, error)
	GetUser(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	
	// Chat methods
	CreateChat(chat *Chat) error
	GetUserChats(userID string) ([]*Chat, error)
	GetChat(chatID string) (*Chat, error)
	IsUserInChat(userID, chatID string) bool
	FindOrCreateDirectChat(user1ID, user2ID string) (*Chat, error)
	
	// Message methods
	CreateMessage(message *Message) error
	GetChatMessages(chatID string) ([]*Message, error)
}

// Store es la instancia global del almacenamiento
var Store StoreInterface

// InitializeStore inicializa el store según la configuración
func InitializeStore() {
	usePostgres := os.Getenv("USE_POSTGRES") != "false"
	
	if usePostgres && database.GetDB() != nil {
		log.Println("Usando PostgreSQL como almacenamiento")
		Store = NewPostgresStore()
		
		// Crear bots automáticamente si no existen
		createDefaultBots()
	} else {
		log.Println("Usando almacenamiento en memoria")
		Store = &MemoryStore{
			users:    make(map[string]*User),
			chats:    make(map[string]*Chat),
			messages: make(map[string][]*Message),
		}
	}
}

// createDefaultBots crea los bots por defecto si no existen
func createDefaultBots() {
	bots := []struct {
		id       string
		username string
		name     string
		email    string
	}{
		{
			id:       "", // Let the store generate a UUID
			username: "support_bot",
			name:     "Soporte Técnico",
			email:    "support@whatsapp-clone.com",
		},
		{
			id:       "", // Let the store generate a UUID
			username: "sales_bot",
			name:     "Ventas Bot",
			email:    "sales@whatsapp-clone.com",
		},
		{
			id:       "", // Let the store generate a UUID
			username: "template_bot",
			name:     "Template Bot",
			email:    "template@whatsapp-clone.com",
		},
	}
	
	for _, bot := range bots {
		// Verificar si el bot ya existe por email
		existingBot, _ := Store.GetUserByEmail(bot.email)
		if existingBot == nil {
			// Crear el bot
			botUser := &User{
				ID:       "", // Let the store generate the UUID
				Username: bot.username,
				Email:    bot.email,
				Name:     bot.name,
				Password: "bot_password_hash", // Los bots no necesitan contraseña real
			}
			
			if err := Store.CreateUser(botUser); err != nil {
				log.Printf("Error creando bot %s: %v", bot.name, err)
			} else {
				log.Printf("Bot %s creado exitosamente", bot.name)
			}
		}
	}
}