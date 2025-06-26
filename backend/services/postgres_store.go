package services

import (
	"database/sql"
	"fmt"
	"time"
	"whatsapp-clone/database"
	
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// PostgresStore implementa el almacenamiento usando PostgreSQL
type PostgresStore struct {
	db *sql.DB
}

// NewPostgresStore crea una nueva instancia del store PostgreSQL
func NewPostgresStore() *PostgresStore {
	return &PostgresStore{
		db: database.GetDB(),
	}
}

// CreateUser crea un nuevo usuario en la base de datos
func (s *PostgresStore) CreateUser(user *User) error {
	if user.ID == "" {
		user.ID = uuid.New().String() // Just use the UUID without prefix
	}
	user.CreatedAt = time.Now()
	
	// Usar el username como phone_number por ahora
	phoneNumber := user.Username
	if phoneNumber == "" {
		phoneNumber = "+1" + fmt.Sprintf("%010d", time.Now().Unix()%10000000000)
	}
	
	query := `
		INSERT INTO users (id, phone_number, email, name, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	
	_, err := s.db.Exec(query, user.ID, phoneNumber, user.Email, user.Name, user.Password, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}
	
	return nil
}

// GetUserByUsername obtiene un usuario por su username (phone_number en la DB)
func (s *PostgresStore) GetUserByUsername(username string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, phone_number, email, name, password_hash, created_at
		FROM users
		WHERE phone_number = $1
	`
	
	err := s.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.Password, &user.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetUserByID obtiene un usuario por su ID
func (s *PostgresStore) GetUserByID(id string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, phone_number, email, name, password_hash, created_at
		FROM users
		WHERE id = $1
	`
	
	err := s.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.Password, &user.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetUser es un alias de GetUserByID para compatibilidad
func (s *PostgresStore) GetUser(id string) (*User, error) {
	return s.GetUserByID(id)
}

// GetUserByEmail obtiene un usuario por su email
func (s *PostgresStore) GetUserByEmail(email string) (*User, error) {
	user := &User{}
	query := `
		SELECT id, phone_number, email, name, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	
	err := s.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Name, &user.Password, &user.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

// GetUsers obtiene todos los usuarios
func (s *PostgresStore) GetUsers() ([]*User, error) {
	query := `
		SELECT id, phone_number, email, name, created_at
		FROM users
		ORDER BY created_at DESC
	`
	
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var users []*User
	for rows.Next() {
		user := &User{}
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Name, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	
	return users, nil
}

// CreateChat crea un nuevo chat
func (s *PostgresStore) CreateChat(chat *Chat) error {
	if chat.ID == "" {
		chat.ID = uuid.New().String() // Just use the UUID without prefix
	}
	chat.CreatedAt = time.Now()
	
	// Iniciar transacción
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// Insertar chat
	query := `
		INSERT INTO chats (id, name, type, created_at)
		VALUES ($1, $2, $3, $4)
	`
	
	chatType := "individual"
	if len(chat.Participants) > 2 {
		chatType = "group"
	}
	
	_, err = tx.Exec(query, chat.ID, chat.Name, chatType, chat.CreatedAt)
	if err != nil {
		return err
	}
	
	// Insertar participantes
	for _, participantID := range chat.Participants {
		participantQuery := `
			INSERT INTO chat_participants (id, chat_id, user_id, joined_at)
			VALUES ($1, $2, $3, $4)
		`
		_, err = tx.Exec(participantQuery, uuid.New().String(), chat.ID, participantID, chat.CreatedAt)
		if err != nil {
			return err
		}
	}
	
	return tx.Commit()
}

// GetUserChats obtiene los chats de un usuario
func (s *PostgresStore) GetUserChats(userID string) ([]*Chat, error) {
	query := `
		SELECT DISTINCT c.id, c.name, c.created_at
		FROM chats c
		INNER JOIN chat_participants cp ON c.id = cp.chat_id
		WHERE cp.user_id = $1 AND cp.left_at IS NULL
		ORDER BY c.created_at DESC
	`
	
	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var chats []*Chat
	for rows.Next() {
		chat := &Chat{}
		err := rows.Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
		if err != nil {
			return nil, err
		}
		
		// Obtener participantes
		participants, err := s.getChatParticipants(chat.ID)
		if err != nil {
			return nil, err
		}
		chat.Participants = participants
		
		chats = append(chats, chat)
	}
	
	return chats, nil
}

// getChatParticipants obtiene los participantes de un chat
func (s *PostgresStore) getChatParticipants(chatID string) ([]string, error) {
	query := `
		SELECT user_id
		FROM chat_participants
		WHERE chat_id = $1 AND left_at IS NULL
	`
	
	rows, err := s.db.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var participants []string
	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		participants = append(participants, userID)
	}
	
	return participants, nil
}

// GetChat obtiene un chat por ID
func (s *PostgresStore) GetChat(chatID string) (*Chat, error) {
	chat := &Chat{}
	query := `
		SELECT id, name, created_at
		FROM chats
		WHERE id = $1
	`
	
	err := s.db.QueryRow(query, chatID).Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	// Obtener participantes
	participants, err := s.getChatParticipants(chat.ID)
	if err != nil {
		return nil, err
	}
	chat.Participants = participants
	
	return chat, nil
}

// CreateMessage crea un nuevo mensaje
func (s *PostgresStore) CreateMessage(message *Message) error {
	if message.ID == "" {
		message.ID = uuid.New().String()
	}
	message.CreatedAt = time.Now()
	
	query := `
		INSERT INTO messages (id, chat_id, sender_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	
	_, err := s.db.Exec(query, message.ID, message.ChatID, message.SenderID, message.Content, message.CreatedAt)
	return err
}

// GetChatMessages obtiene los mensajes de un chat
func (s *PostgresStore) GetChatMessages(chatID string) ([]*Message, error) {
	query := `
		SELECT id, chat_id, sender_id, content, created_at
		FROM messages
		WHERE chat_id = $1
		ORDER BY created_at ASC
	`
	
	rows, err := s.db.Query(query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []*Message
	for rows.Next() {
		msg := &Message{}
		err := rows.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.CreatedAt)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	
	return messages, nil
}

// IsUserInChat verifica si un usuario pertenece a un chat
func (s *PostgresStore) IsUserInChat(userID, chatID string) bool {
	query := `
		SELECT COUNT(*)
		FROM chat_participants
		WHERE user_id = $1 AND chat_id = $2 AND left_at IS NULL
	`
	
	var count int
	err := s.db.QueryRow(query, userID, chatID).Scan(&count)
	if err != nil {
		return false
	}
	
	return count > 0
}

// FindOrCreateDirectChat busca o crea un chat directo entre dos usuarios
func (s *PostgresStore) FindOrCreateDirectChat(user1ID, user2ID string) (*Chat, error) {
	// Buscar chat existente
	query := `
		SELECT DISTINCT c.id, c.name, c.created_at
		FROM chats c
		WHERE c.type = 'individual'
		AND EXISTS (
			SELECT 1 FROM chat_participants cp1
			WHERE cp1.chat_id = c.id AND cp1.user_id = $1 AND cp1.left_at IS NULL
		)
		AND EXISTS (
			SELECT 1 FROM chat_participants cp2
			WHERE cp2.chat_id = c.id AND cp2.user_id = $2 AND cp2.left_at IS NULL
		)
		AND (
			SELECT COUNT(*) FROM chat_participants cp3
			WHERE cp3.chat_id = c.id AND cp3.left_at IS NULL
		) = 2
		LIMIT 1
	`
	
	var chat Chat
	err := s.db.QueryRow(query, user1ID, user2ID).Scan(&chat.ID, &chat.Name, &chat.CreatedAt)
	if err == nil {
		// Chat encontrado, obtener participantes
		participants, err := s.getChatParticipants(chat.ID)
		if err != nil {
			return nil, err
		}
		chat.Participants = participants
		return &chat, nil
	}
	
	// Si no existe, crear uno nuevo
	newChat := &Chat{
		ID:           uuid.New().String(), // Just use the UUID without prefix
		Name:         "Direct Chat",
		Participants: []string{user1ID, user2ID},
		CreatedAt:    time.Now(),
	}
	
	if err := s.CreateChat(newChat); err != nil {
		return nil, err
	}
	
	return newChat, nil
}