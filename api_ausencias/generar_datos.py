import numpy as np
import pandas as pd
import os # Importar el módulo os para manejar rutas de archivos

def generar_datos_crudos_a_csv(num_muestras=1000, random_state=42, nombre_archivo='datos_simulados_crudos.csv'):
    """
    Genera un DataFrame con datos simulados que muestran una relación logística
    entre varias características y la probabilidad de ausencia, en su formato 'crudo'
    (variables categóricas como texto, numéricas sin escalar), y los guarda en un archivo CSV.

    Args:
        num_muestras (int): El número de filas de datos a generar.
        random_state (int): Semilla para la reproducibilidad de los resultados.
        nombre_archivo (str): El nombre del archivo CSV donde se guardarán los datos.
    """
    np.random.seed(random_state) # Para reproducibilidad

    # Definir rangos y opciones para las variables
    ausencias_ultimos_30d = np.random.randint(0, 15, num_muestras)
    patron_whatsapp = np.random.choice(['normal', 'lento', 'rapido'], num_muestras, p=[0.5, 0.3, 0.2])
    hora_peticion_minutos = np.random.randint(0, 1440, num_muestras) # Minutos en un día (24 * 60)
    dia_semana = np.random.choice(['lunes', 'martes', 'miércoles', 'jueves', 'viernes', 'sábado', 'domingo'], num_muestras)
    clima = np.random.choice(['soleado', 'nublado', 'lluvioso', 'tormenta'], num_muestras, p=[0.4, 0.3, 0.2, 0.1])
    distancia_km = np.random.uniform(1, 30, num_muestras)

    # Convertir variables categóricas a numéricas para calcular la probabilidad logística
    # IMPORTANTE: Estos mapeos son SOLO para la generación interna de 'se_ausento',
    # no afectarán el formato final de las columnas del CSV.
    patron_whatsapp_map = {'normal': 0, 'lento': 1, 'rapido': -1}
    clima_map = {'soleado': -0.5, 'nublado': 0, 'lluvioso': 0.5, 'tormenta': 1}
    dia_semana_map = {'lunes': 0, 'martes': 0.1, 'miércoles': 0.2, 'jueves': 0.3, 'viernes': 0.8, 'sábado': 0.6, 'domingo': 0.7}

    patron_whatsapp_num = np.array([patron_whatsapp_map[p] for p in patron_whatsapp])
    clima_num = np.array([clima_map[c] for c in clima])
    dia_semana_num = np.array([dia_semana_map[d] for d in dia_semana])

    # Definir los coeficientes para la función logística
    beta_0 = -3.0 # Intercepto
    beta_ausencias = 0.4
    beta_whatsapp = 1.5
    beta_hora = 0.001
    beta_dia = 1.0
    beta_clima = 2.0
    beta_distancia = 0.1

    # Calcular la parte lineal de la función logística (logit)
    z = (beta_0 +
         beta_ausencias * ausencias_ultimos_30d +
         beta_whatsapp * patron_whatsapp_num +
         beta_hora * hora_peticion_minutos +
         beta_dia * dia_semana_num +
         beta_clima * clima_num +
         beta_distancia * distancia_km)

    # Aplicar la función sigmoide para obtener la probabilidad de ausencia
    prob_ausencia = 1 / (1 + np.exp(-z))

    # Generar el resultado binario (se_ausento) basándose en la probabilidad
    se_ausento = (np.random.rand(num_muestras) < prob_ausencia).astype(int)

    # Crear el DataFrame con el formato de salida deseado (datos 'crudos')
    data = pd.DataFrame({
        'ausencias_ultimos_30d': ausencias_ultimos_30d,
        'patron_whatsapp': patron_whatsapp,
        'hora_peticion_minutos': hora_peticion_minutos,
        'dia_semana': dia_semana,
        'clima': clima,
        'distancia_km': distancia_km,
        'se_ausento': se_ausento
    })

    # Guardar el DataFrame en un archivo CSV
    data.to_csv(nombre_archivo, index=False)
    print(f"Datos generados en formato 'crudo' y guardados exitosamente en '{nombre_archivo}'")
    print(f"Número de filas: {data.shape[0]}, Número de columnas: {data.shape[1]}")
    print(f"Ejemplo de las primeras 5 filas del CSV:\n{data.head()}")

# --- Uso de la función ---
if __name__ == "__main__":
    # Generar 1000 muestras y guardarlas en 'datos_ausencias_crudos.csv'
    generar_datos_crudos_a_csv(num_muestras=1000, nombre_archivo='datos_ausencia.csv')

    # Para verificar el archivo generado
    # try:
    #     df_cargado = pd.read_csv('datos_ausencias_crudos.csv')
    #     print("\nDatos cargados desde CSV para verificación:")
    #     print(df_cargado.head())
    # except FileNotFoundError:
    #     print("El archivo CSV no se encontró. Asegúrate de que la función se haya ejecutado correctamente.")