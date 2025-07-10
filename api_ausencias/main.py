from fastapi import FastAPI, HTTPException, Path
from typing import Dict
from fastapi.middleware.cors import CORSMiddleware
import pandas as pd
from datetime import datetime
from models.schemas import Horario, Tarea, PrediccionInput, AsignacionInteligenteInput, AsignacionOutput
from services.task_service import gestionar_asignacion_tarea

app = FastAPI(
    title="API de Predicción de Ausentismo Laboral",
    description="Una API para gestionar horarios y predecir la ausencia de trabajadores.",
    version="1.0.0"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# --- Base de datos en memoria  ---
db_horarios: Dict[str, Horario] = {}
db_tareas: Dict[str, list[Tarea]] = {}

# --- Cargar db_horarios.csv ---
df_horarios = pd.read_csv("db_horarios.csv")

for _, row in df_horarios.iterrows():
    worker_id = row["worker_id"]
    dias_laborales = row["dias_laborales"].split(",")
    hora_inicio = datetime.strptime(row["hora_inicio"], "%H:%M:%S").time()
    hora_fin = datetime.strptime(row["hora_fin"], "%H:%M:%S").time()
    
    db_horarios[worker_id] = Horario(
        dias_laborales=dias_laborales,
        hora_inicio=hora_inicio,
        hora_fin=hora_fin
    )

# --- Cargar db_tareas.csv ---
df_tareas = pd.read_csv("db_tareas.csv")

for _, row in df_tareas.iterrows():
    worker_id = row["worker_id"]
    tarea = Tarea(
        descripcion=row["descripcion"],
        prioridad=row["prioridad"],
        timestamp=datetime.fromisoformat(row["timestamp"])
    )
    if worker_id not in db_tareas:
        db_tareas[worker_id] = []
    db_tareas[worker_id].append(tarea)


# --- Endpoints ---

@app.post("/horarios/{worker_id}", status_code=201)
def crear_o_actualizar_horario(worker_id: str, horario: Horario):
    """
    Recibe y almacena el horario de un trabajador.
    Si el worker_id ya existe, actualiza su horario.
    """
    db_horarios[worker_id] = horario
    return {"mensaje": f"Horario para el trabajador '{worker_id}' almacenado correctamente."}

@app.get("/horarios/{worker_id}", response_model=Horario)
def obtener_horario(worker_id: str):
    """
    Obtiene el horario de un trabajador específico.
    """
    if worker_id not in db_horarios:
        raise HTTPException(status_code=404, detail=f"Trabajador con id '{worker_id}' no encontrado.")
    return db_horarios[worker_id]

@app.post(
    "/asignar-tarea/{worker_id}",
    response_model=AsignacionOutput,
    tags=["Asignación de Tareas"]
)
def asignar_tarea_inteligente(
    worker_id: str = Path(..., description="El ID único del trabajador"),
    datos_entrada: AsignacionInteligenteInput = ...
):
    """
    Punto de entrada único para asignar una tarea a un trabajador.
    
    Este endpoint realiza las siguientes acciones:
    1. Predice la probabilidad de ausencia del trabajador usando un modelo de IA.
    2. Verifica si el trabajador está dentro de su horario laboral.
    3. Aplica reglas de negocio para decidir si la tarea debe ser asignada.
    4. Asigna la tarea si todas las condiciones se cumplen.
    """
    resultado = gestionar_asignacion_tarea(
        worker_id=worker_id,
        datos_entrada=datos_entrada,
        db_horarios=db_horarios,
        db_tareas=db_tareas
    )
    
    # Si el servicio devuelve un error, lanzamos una excepción HTTP
    if "error" in resultado:
        raise HTTPException(
            status_code=resultado["status_code"],
            detail=resultado["error"]
        )
        
    # Si todo fue bien, devolvemos la respuesta de éxito
    return AsignacionOutput(
        mensaje=resultado["mensaje"],
        tarea_asignada=resultado["tarea_asignada"]
    )