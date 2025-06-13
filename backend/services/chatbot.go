package services

import (
	"fmt"
	"strings"
	"time"
)

type ChatBot struct {
	ID          string
	Name        string
	Description string
	Avatar      string
	Responses   map[string]string
	DefaultMsg  string
}

var ChatBots = map[string]*ChatBot{
	"bot_support": {
		ID:          "bot_support",
		Name:        "Soporte Técnico",
		Description: "Bot de ayuda y soporte",
		Avatar:      "🤖",
		DefaultMsg:  "Lo siento, no entiendo tu pregunta. ¿Podrías reformularla?",
		Responses: map[string]string{
			"hola":           "¡Hola! Soy el bot de soporte. ¿En qué puedo ayudarte?",
			"ayuda":          "Puedo ayudarte con:\n• Información de productos\n• Estado de pedidos\n• Problemas técnicos\n• Preguntas frecuentes",
			"precio":         "Los precios varían según el producto. ¿Sobre qué producto específico necesitas información?",
			"horario":        "Nuestro horario de atención es de Lunes a Viernes de 9:00 a 18:00",
			"contacto":       "Puedes contactarnos en:\n📧 soporte@empresa.com\n📱 +54 11 1234-5678",
			"problema":       "Lamento que tengas problemas. ¿Podrías describir tu situación con más detalle?",
			"gracias":        "¡De nada! Estoy aquí para ayudarte cuando lo necesites.",
			"pedido":         "Para consultar el estado de tu pedido, por favor proporciona tu número de orden.",
			"envio":          "Los envíos tardan entre 3-5 días hábiles. ¿Necesitas información sobre un envío específico?",
			"devolucion":     "Las devoluciones se aceptan dentro de los 30 días. ¿Necesitas iniciar una devolución?",
		},
	},
	"bot_sales": {
		ID:          "bot_sales",
		Name:        "Ventas Bot",
		Description: "Bot de ventas y promociones",
		Avatar:      "💼",
		DefaultMsg:  "No tengo información sobre eso, pero puedo contarte sobre nuestras ofertas actuales.",
		Responses: map[string]string{
			"hola":           "¡Hola! Soy el bot de ventas. ¿Buscas algún producto en particular?",
			"oferta":         "🎉 Ofertas de la semana:\n• 20% OFF en electrónica\n• 2x1 en productos seleccionados\n• Envío gratis en compras mayores a $5000",
			"descuento":      "Los descuentos actuales son:\n• Primera compra: 15% OFF\n• Cliente frecuente: 10% OFF\n• Compras mayores a $10000: 25% OFF",
			"catalogo":       "Nuestro catálogo incluye:\n📱 Electrónica\n👕 Ropa\n🏠 Hogar\n🎮 Gaming\n¿Qué categoría te interesa?",
			"pago":           "Aceptamos:\n💳 Tarjetas de crédito/débito\n💰 Efectivo\n📱 Mercado Pago\n🏦 Transferencia bancaria",
			"cuotas":         "¡Sí! Ofrecemos:\n• 3 cuotas sin interés\n• 6 cuotas con 10% de interés\n• 12 cuotas con 15% de interés",
		},
	},
	"bot_template": {
		ID:          "bot_template",
		Name:        "Template Bot",
		Description: "Bot para gestión de plantillas de mensajes",
		Avatar:      "📋",
		DefaultMsg:  "No reconozco ese comando. Escribe 'ayuda' para ver las opciones disponibles.",
		Responses: map[string]string{
			"hola":           "¡Hola! Soy el bot de plantillas. Puedo ayudarte a crear y gestionar plantillas de mensajes.",
			"ayuda":          "Comandos disponibles:\n• /nueva - Crear nueva plantilla\n• /listar - Ver plantillas\n• /usar [nombre] - Usar plantilla\n• /variables - Ver variables disponibles",
			"nueva":          "Para crear una plantilla, usa el formato:\n/nueva [nombre] [mensaje]\n\nEjemplo: /nueva bienvenida Hola {{nombre}}, bienvenido a {{empresa}}",
			"listar":         "Plantillas disponibles:\n1. bienvenida\n2. confirmacion_pedido\n3. recordatorio\n4. promocion\n5. seguimiento",
			"variables":      "Variables disponibles:\n• {{nombre}} - Nombre del cliente\n• {{empresa}} - Nombre de la empresa\n• {{fecha}} - Fecha actual\n• {{pedido}} - Número de pedido\n• {{producto}} - Nombre del producto",
			"ejemplo":        "Ejemplo de plantilla:\n\nHola {{nombre}},\n\nTu pedido #{{pedido}} ha sido confirmado.\nFecha de entrega estimada: {{fecha}}\n\nGracias por tu compra!",
		},
	},
	"bot_test": {
		ID:          "bot_test",
		Name:        "Tester Bot",
		Description: "Bot para pruebas de respuestas de plantillas de mensajes",
		Avatar:      "🤓",
		DefaultMsg:  "No reconozco ese comando. Falta hacerlo.",
		Responses: map[string]string{
			"hola":           "¡Hola! Soy el bot de plantillas. Puedo ayudarte a crear y gestionar plantillas de mensajes.",
			"ayuda":          "Comandos disponibles:\n• /nueva - Crear nueva plantilla\n• /listar - Ver plantillas\n• /usar [nombre] - Usar plantilla\n• /variables - Ver variables disponibles",
			"nueva":          "Para crear una plantilla, usa el formato:\n/nueva [nombre] [mensaje]\n\nEjemplo: /nueva bienvenida Hola {{nombre}}, bienvenido a {{empresa}}",
			"listar":         "Plantillas disponibles:\n1. bienvenida\n2. confirmacion_pedido\n3. recordatorio\n4. promocion\n5. seguimiento",
			"variables":      "Variables disponibles:\n• {{nombre}} - Nombre del cliente\n• {{empresa}} - Nombre de la empresa\n• {{fecha}} - Fecha actual\n• {{pedido}} - Número de pedido\n• {{producto}} - Nombre del producto",
			"ejemplo":        "Ejemplo de plantilla:\n\nHola {{nombre}},\n\nTu pedido #{{pedido}} ha sido confirmado.\nFecha de entrega estimada: {{fecha}}\n\nGracias por tu compra!",
		},
	},
}

// ProcessBotMessage procesa el mensaje y devuelve la respuesta del bot
func ProcessBotMessage(botID, message string) string {
	bot, exists := ChatBots[botID]
	if !exists {
		return "Bot no encontrado"
	}

	// Convertir mensaje a minúsculas para comparación
	lowerMessage := strings.ToLower(message)

	// Buscar coincidencias en las respuestas
	for keyword, response := range bot.Responses {
		if strings.Contains(lowerMessage, keyword) {
			return response
		}
	}

	// Si no hay coincidencia, devolver mensaje por defecto
	return bot.DefaultMsg
}

// GetBotResponse genera una respuesta de bot con delay simulado
func GetBotResponse(botID, userMessage string, chatID string, hub interface{}) {
	// Simular tiempo de procesamiento
	time.Sleep(1 * time.Second)

	// Obtener respuesta del bot
	response := ProcessBotMessage(botID, userMessage)

	// Crear mensaje de respuesta
	botMessage := &Message{
		ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
		ChatID:    chatID,
		SenderID:  botID,
		Content:   response,
		CreatedAt: time.Now(),
	}

	// Guardar mensaje en el store
	Store.CreateMessage(botMessage)

	// Enviar mensaje a través del WebSocket
	if hub != nil {
		messageData := map[string]interface{}{
			"type":      "new_message",
			"message":   botMessage,
			"chatId":    chatID,
			"timestamp": time.Now().Unix(),
		}
		
		// Broadcast to all clients
		
		// Use type assertion to call broadcast method
		if h, ok := hub.(interface{ BroadcastToChat(string, interface{}) }); ok {
			h.BroadcastToChat(chatID, messageData)
		}
	}
}

// InitializeBotChats crea chats iniciales con los bots para cada usuario
func InitializeBotChats(userID string) {
	fmt.Printf("Initializing bot chats for user: %s\n", userID)
	for _, bot := range ChatBots {
		// Crear usuario bot si no existe
		botUser := &User{
			ID:        bot.ID,
			Username:  bot.ID,
			Email:     bot.ID + "@bot.local",
			Name:      bot.Name,
			CreatedAt: time.Now(),
		}
		Store.CreateUser(botUser)
		fmt.Printf("Created bot user: %s\n", bot.Name)

		// Crear chat con el bot
		chat := &Chat{
			ID:           fmt.Sprintf("chat_%s_%s", userID, bot.ID),
			Name:         bot.Name,
			Participants: []string{userID, bot.ID},
			CreatedAt:    time.Now(),
		}
		Store.CreateChat(chat)
		fmt.Printf("Created chat with bot: %s (ID: %s)\n", bot.Name, chat.ID)

		// Mensaje de bienvenida del bot
		welcomeMsg := &Message{
			ID:        fmt.Sprintf("msg_%d", time.Now().UnixNano()),
			ChatID:    chat.ID,
			SenderID:  bot.ID,
			Content:   bot.Responses["hola"],
			CreatedAt: time.Now(),
		}
		Store.CreateMessage(welcomeMsg)
	}
}