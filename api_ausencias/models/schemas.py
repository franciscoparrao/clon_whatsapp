from pydantic import BaseModel, Field
from typing import List, Literal
from datetime import datetime, time

class Horario(BaseModel):
    dias_laborales: List[Literal["lunes", "martes", "miercoles", "jueves", "viernes", "sabado", "domingo"]]
    hora_inicio: time
    hora_fin: time

class Tarea(BaseModel):
    descripcion: str
    prioridad: Literal["alta", "media", "baja"]
    timestamp: datetime = Field(default_factory=datetime.now, description="Timestamp de cuando se crea la tarea")


class PrediccionInput(BaseModel):
    historial_asistencia_30d: List[int] = Field(..., min_items=30, max_items=30)
    patron_respuesta_whatsapp: Literal["rapido", "normal", "lento", "no_responde"]
    condiciones_climaticas: Literal["soleado", "nublado", "lluvioso", "tormenta"]
    distancia_trabajo_km: float = Field(..., gt=0)
    hora_peticion: time
    dia_semana: Literal["lunes", "martes", "miercoles", "jueves", "viernes", "sabado", "domingo"]

class AsignacionInteligenteInput(BaseModel):
    tarea: Tarea
    datos_prediccion: PrediccionInput

class AsignacionOutput(BaseModel):
    mensaje: str
    tarea_asignada: bool