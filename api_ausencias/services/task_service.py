from typing import Dict, Any
from models.schemas import AsignacionInteligenteInput, Horario, PrediccionInput
from services.prediction_service import predecir_probabilidad_ausencia
from datetime import timezone, timedelta

# Aplicar zona horaria local
CHILE_TZ = timezone(timedelta(hours=-4))

# Mapeo de día de la semana de Python a nuestro formato
DIAS_SEMANA_MAP = {
    0: "lunes", 1: "martes", 2: "miercoles", 3: "jueves", 
    4: "viernes", 5: "sabado", 6: "domingo"
}

def gestionar_asignacion_tarea(
    worker_id: str,
    datos_entrada: AsignacionInteligenteInput,
    db_horarios: Dict[str, Horario],
    db_tareas: Dict[str, list]
) -> Dict[str, Any]:
    """
    Función única que centraliza la lógica para asignar una tarea de forma inteligente.
    Retorna un diccionario con el resultado y el código de estado HTTP sugerido.
    """
    # --- 1. Validar existencia de horario ---
    horario_trabajador = db_horarios.get(worker_id)
    if not horario_trabajador:
        return {
            "error": f"No se encontró horario para el trabajador '{worker_id}'.",
            "status_code": 404
        }

    # --- 2. Preparar datos y realizar predicción ---
    tarea = datos_entrada.tarea
    datos_prediccion = datos_entrada.datos_prediccion
    
    timestamp_local = tarea.timestamp.astimezone(CHILE_TZ)

    # Extraemos el día y la hora del timestamp de la tarea para consistencia
    dia_semana_peticion = DIAS_SEMANA_MAP[timestamp_local.weekday()]
    hora_peticion = timestamp_local.time().replace(tzinfo=None)
    
    # Creamos el objeto PrediccionInput completo para el modelo
    input_para_modelo = PrediccionInput(
        historial_asistencia_30d=datos_prediccion.historial_asistencia_30d,
        patron_respuesta_whatsapp=datos_prediccion.patron_respuesta_whatsapp,
        condiciones_climaticas=datos_prediccion.condiciones_climaticas,
        distancia_trabajo_km=datos_prediccion.distancia_trabajo_km,
        hora_peticion=hora_peticion,
        dia_semana=dia_semana_peticion
    )

    try:
        # Llamamos al servicio de predicción con los datos completos
        probabilidad_ausencia = predecir_probabilidad_ausencia(input_para_modelo)
    except RuntimeError as e:
         return {"error": str(e), "status_code": 503} # Service Unavailable

    # --- 3. Validar si el trabajador está en horario laboral ---
    if not (dia_semana_peticion in horario_trabajador.dias_laborales and
            horario_trabajador.hora_inicio <= hora_peticion <= horario_trabajador.hora_fin):
        return {
            "error": f"No se puede asignar la tarea. El trabajador '{worker_id}' está fuera de su horario laboral.",
            "status_code": 400 # Bad Request
        }

    # --- 4. Aplicar lógica de negocio para la asignación ---

# Si la probabilidad es alta, solo se permite asignar tareas de prioridad baja

    if probabilidad_ausencia >= 0.7 and tarea.prioridad in ["alta", "media"]:
        return {
        "error": f"No se asignan tareas de prioridad '{tarea.prioridad}' con una probabilidad de ausencia alta ({probabilidad_ausencia:.2f}). Solo se permiten tareas de prioridad 'baja'.",
        "status_code": 400
    }

# Si la probabilidad es baja, no se asignan tareas de baja prioridad
    if probabilidad_ausencia < 0.5 and tarea.prioridad == "baja":
        return {
        "error": f"Se priorizan tareas importantes cuando la probabilidad de ausencia es baja ({probabilidad_ausencia:.2f}). No se asignan tareas de prioridad 'baja'.",
        "status_code": 400
    }


    # --- 5. Si todas las validaciones pasan, asignar la tarea ---
    if worker_id not in db_tareas:
        db_tareas[worker_id] = []
    
    db_tareas[worker_id].append(tarea)
    
    return {
        "mensaje": f"Tarea '{tarea.descripcion}' asignada exitosamente a '{worker_id}'. Probabilidad de ausencia: {probabilidad_ausencia:.2f}",
        "tarea_asignada": True,
        "status_code": 201 # Created
    }