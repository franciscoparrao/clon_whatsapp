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

//Mantener registro del flujo del usuario para la simulacion
type Session struct {
	ID           string
	CurrentFlow  string //Donde esta
	LastQuestion string //Que se le pregunto?
	PendingInput string //Si tiene que responder algo
}

var botSessions = map[string]*Session{}

var ChatBots = map[string]*ChatBot{
	"bot_support": {
		ID:          "bot_support",
		Name:        "Soporte Técnico",
		Description: "Bot de ayuda y soporte",
		Avatar:      "🤖",
		DefaultMsg:  "Lo siento, no entiendo tu pregunta. ¿Podrías reformularla?",
		Responses: map[string]string{
			"hola":       "¡Hola! Soy el bot de soporte. ¿En qué puedo ayudarte?",
			"ayuda":      "Puedo ayudarte con:\n• Información de productos\n• Estado de pedidos\n• Problemas técnicos\n• Preguntas frecuentes",
			"precio":     "Los precios varían según el producto. ¿Sobre qué producto específico necesitas información?",
			"horario":    "Nuestro horario de atención es de Lunes a Viernes de 9:00 a 18:00",
			"contacto":   "Puedes contactarnos en:\n📧 soporte@empresa.com\n📱 +54 11 1234-5678",
			"problema":   "Lamento que tengas problemas. ¿Podrías describir tu situación con más detalle?",
			"gracias":    "¡De nada! Estoy aquí para ayudarte cuando lo necesites.",
			"pedido":     "Para consultar el estado de tu pedido, por favor proporciona tu número de orden.",
			"envio":      "Los envíos tardan entre 3-5 días hábiles. ¿Necesitas información sobre un envío específico?",
			"devolucion": "Las devoluciones se aceptan dentro de los 30 días. ¿Necesitas iniciar una devolución?",
		},
	},
	"bot_sales": {
		ID:          "bot_sales",
		Name:        "Ventas Bot",
		Description: "Bot de ventas y promociones",
		Avatar:      "💼",
		DefaultMsg:  "No tengo información sobre eso, pero puedo contarte sobre nuestras ofertas actuales.",
		Responses: map[string]string{
			"hola":      "¡Hola! Soy el bot de ventas. ¿Buscas algún producto en particular?",
			"oferta":    "🎉 Ofertas de la semana:\n• 20% OFF en electrónica\n• 2x1 en productos seleccionados\n• Envío gratis en compras mayores a $5000",
			"descuento": "Los descuentos actuales son:\n• Primera compra: 15% OFF\n• Cliente frecuente: 10% OFF\n• Compras mayores a $10000: 25% OFF",
			"catalogo":  "Nuestro catálogo incluye:\n📱 Electrónica\n👕 Ropa\n🏠 Hogar\n🎮 Gaming\n¿Qué categoría te interesa?",
			"pago":      "Aceptamos:\n💳 Tarjetas de crédito/débito\n💰 Efectivo\n📱 Mercado Pago\n🏦 Transferencia bancaria",
			"cuotas":    "¡Sí! Ofrecemos:\n• 3 cuotas sin interés\n• 6 cuotas con 10% de interés\n• 12 cuotas con 15% de interés",
		},
	},
	"bot_template": {
		ID:          "bot_template",
		Name:        "Template Bot",
		Description: "Bot para gestión de plantillas de mensajes",
		Avatar:      "📋",
		DefaultMsg:  "No reconozco ese comando. Escribe 'ayuda' para ver las opciones disponibles.",
		Responses: map[string]string{
			"hola":      "¡Hola! Soy el bot de plantillas. Puedo ayudarte a crear y gestionar plantillas de mensajes.",
			"ayuda":     "Comandos disponibles:\n• /nueva - Crear nueva plantilla\n• /listar - Ver plantillas\n• /usar [nombre] - Usar plantilla\n• /variables - Ver variables disponibles",
			"nueva":     "Para crear una plantilla, usa el formato:\n/nueva [nombre] [mensaje]\n\nEjemplo: /nueva bienvenida Hola {{nombre}}, bienvenido a {{empresa}}",
			"listar":    "Plantillas disponibles:\n1. bienvenida\n2. confirmacion_pedido\n3. recordatorio\n4. promocion\n5. seguimiento",
			"variables": "Variables disponibles:\n• {{nombre}} - Nombre del cliente\n• {{empresa}} - Nombre de la empresa\n• {{fecha}} - Fecha actual\n• {{pedido}} - Número de pedido\n• {{producto}} - Nombre del producto",
			"ejemplo":   "Ejemplo de plantilla:\n\nHola {{nombre}},\n\nTu pedido #{{pedido}} ha sido confirmado.\nFecha de entrega estimada: {{fecha}}\n\nGracias por tu compra!",
		},
	},
	"bot_test": {
		ID:          "bot_test",
		Name:        "Tester Bot",
		Description: "Bot para pruebas de respuestas de plantillas de mensajes",
		Avatar:      "🤓",
		DefaultMsg:  "No reconozco ese comando. Falta hacerlo.",
		Responses: map[string]string{
			"hola":        "¡Hola! Soy el bot de simulacion de flujos. Simulare los flujos de fletzy!.",
			"ayuda":       "Comandos disponibles:\n• /nueva - Crear nueva plantilla\n• /listar - Ver plantillas\n• /usar [nombre] - Usar plantilla\n• /variables - Ver variables disponibles",
			"nueva":       "Para crear una plantilla, usa el formato:\n/nueva [nombre] [mensaje]\n\nEjemplo: /nueva bienvenida Hola {{nombre}}, bienvenido a {{empresa}}",
			"listar":      "Plantillas disponibles:\n1. bienvenida\n2. confirmacion_pedido\n3. recordatorio\n4. promocion\n5. seguimiento",
			"variables":   "Variables disponibles:\n• {{nombre}} - Nombre del cliente\n• {{empresa}} - Nombre de la empresa\n• {{fecha}} - Fecha actual\n• {{pedido}} - Número de pedido\n• {{producto}} - Nombre del producto",
			"ejemplo":     "Ejemplo de plantilla:\n\nHola {{nombre}},\n\nTu pedido #{{pedido}} ha sido confirmado.\nFecha de entrega estimada: {{fecha}}\n\nGracias por tu compra!",
			"simulacion1": `Confirmación de Asistencia 🎫
					Buenos días, {{Nombre}} {{Apellido}}. 🧑‍⚕️ Esperamos que esté bien. Le escribimos para confirmar su asistencia a la operación programada para hoy. 
					Por favor, responda con uno de los siguientes botones:\n

					1. 'Confirmo' ✅ si asistirá.
					2. 'No asistiré' ❌ en caso contrario.\n

					Agradecemos su pronta respuesta. Que tenga un excelente día. 🚀\n

					Powered by Fletzy`,
		},
	},
}


// ProcessBotMessage procesa el mensaje y devuelve la respuesta del bot
func ProcessBotMessage(botID, message string) string {
	bot, exists := ChatBots[botID]
	if !exists {
		return "Bot no encontrado"
	}

	session := botSessions[botID]
	
	// Paso: esperando hora tras confirmación
	if session != nil && session.CurrentFlow == "simulacion_1" {
		//fmt.Printf("DEBUG: sesión encontrada para %s: %+v\n", botID, session)
		//fmt.Println("Procesando mensaje en flujo de simulación 1 para el bot:", botID)
		if session.PendingInput == "hora_asistencia" {
			//fmt.Printf("Recibiendo hora: %s | Sesión actual: %+v\n", message, session)
			botSessions[botID] = nil // resetear sesión
			return fmt.Sprintf("🕒 Gracias. Registramos su hora de asistencia: %s. ¡Nos vemos pronto!", message)
		}

		switch strings.ToLower(message) {
		case "confirmo":
			session.PendingInput = "hora_asistencia"
			//fmt.Println("Se detectó CONFIRMO. Pendiente de hora.")
			//fmt.Println("DEBUG: sesión actualizada para", botID, "con PendingInput:", session.PendingInput)
			return "✅ ¡Gracias por confirmar! Por favor, indique la hora a la que asistirá (formato HH:MM)."
		case "no asistiré", "no asistire":
			fmt.Println("Se detectó NO ASISTIRE. Se borra la sesión.")
			botSessions[botID] = nil
			return "❌ Entendido, lamentamos que no pueda asistir. Puede reprogramar su operación llamando al 800-000-000."
		}
	}

	// Activar simulación 1
	if strings.ToLower(message) == "simulacion1" {
		//fmt.Println("Iniciando simulación 1 para el bot:", botID)
		botSessions[botID] = &Session{
			ID:          botID, // ya que ahora solo usamos uno por bot
			CurrentFlow: "simulacion_1",
		}
		//fmt.Printf("DEBUG: sesión creada para %s: %+v\n", botID, botSessions[botID])
		//fmt.Println("flow" + botSessions[botID].CurrentFlow)
		return bot.Responses["simulacion1"]
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
