import pandas as pd
import tensorflow as tf
import joblib
from models.schemas import PrediccionInput

# --- Carga de artefactos del modelo (se hace una sola vez al iniciar la API) ---
try:
    MODEL_PATH = 'modelo_ausentismo.h5'
    SCALER_PATH = 'scaler.pkl'
    COLUMNS_PATH = 'model_columns.pkl'
    
    # Cargar el modelo de red neuronal
    model = tf.keras.models.load_model(MODEL_PATH)
    
    # Cargar el scaler
    scaler = joblib.load(SCALER_PATH)
    
    # Cargar las columnas con las que el modelo fue entrenado
    model_columns = joblib.load(COLUMNS_PATH)
    
    print("Modelo y artefactos cargados correctamente.")

except FileNotFoundError as e:
    print(f"Error: No se encontró un archivo del modelo. Asegúrate de ejecutar train_model.py primero. {e}")
    model = None
    scaler = None
    model_columns = None


def predecir_probabilidad_ausencia(datos: PrediccionInput) -> float:
    """
    Usa la red neuronal entrenada para predecir la probabilidad de ausencia.
    """
    if not all([model, scaler, model_columns]):
        # Si el modelo no se cargó, devuelve un valor por defecto o lanza un error
        raise RuntimeError("El modelo de predicción no está disponible. Revisa los logs del servidor.")

    # 1. Convertir el JSON de entrada a un DataFrame de pandas
    input_data = {
        'ausencias_ultimos_30d': [sum(1 for dia in datos.historial_asistencia_30d if dia == 0)],
        'patron_whatsapp': [datos.patron_respuesta_whatsapp],
        # Convertir la hora a un valor numérico (ej: 8:30 -> 8.5)
        'hora_peticion_num': [datos.hora_peticion.hour + datos.hora_peticion.minute / 60],
        'dia_semana': [datos.dia_semana],
        'clima': [datos.condiciones_climaticas],
        'distancia_km': [datos.distancia_trabajo_km],
        'hora_peticion': [datos.hora_peticion.strftime('%H:%M')]
    }
    df = pd.DataFrame.from_dict(input_data)

    # 2. Preprocesar los datos EXACTAMENTE como en el entrenamiento
    
    # a) One-Hot Encoding
    df_processed = pd.get_dummies(df)
    
    # b) Alineación de columnas: Asegurar que el DataFrame tenga las mismas columnas
    # que el modelo espera. Las columnas faltantes (ej: 'clima_soleado') se añaden con valor 0.
    df_aligned = df_processed.reindex(columns=model_columns, fill_value=0)
    
    # c) Escalar los datos con el scaler cargado
    scaled_data = scaler.transform(df_aligned)

    # 3. Realizar la predicción
    probabilidad = model.predict(scaled_data)
    
    # El resultado de model.predict es un array de numpy, ej: [[0.753]]
    # Extraemos el valor flotante.
    resultado_final = float(probabilidad[0][0])
    
    return round(resultado_final, 4)