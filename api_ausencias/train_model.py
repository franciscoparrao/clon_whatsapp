import os
import random
import time
import numpy as np
import pandas as pd
import tensorflow as tf
import matplotlib.pyplot as plt

from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
import joblib

from tensorflow.keras.models import Sequential
from tensorflow.keras.layers import Dense, Dropout, BatchNormalization, Input, Activation
from tensorflow.keras.callbacks import EarlyStopping
from tensorflow.keras import metrics
from tensorflow.keras import regularizers

# --- SEMILLA PARA REPRODUCIBILIDAD ---
SEED = 42
random.seed(SEED)
np.random.seed(SEED)
tf.random.set_seed(SEED)
os.environ['PYTHONHASHSEED'] = str(SEED)

print("Iniciando el entrenamiento del modelo de predicción de ausentismo...")

# --- 1. Cargar y Preparar los Datos ---
data = pd.read_csv('datos_ausencia.csv')
df = pd.DataFrame(data)

# --- 2. Preprocesamiento de Datos ---

# a) Convertir variables categóricas a numéricas (One-Hot Encoding)
df_processed = pd.get_dummies(df, columns=['patron_whatsapp', 'dia_semana', 'clima'], drop_first=True)

# b) Separar características (X) y objetivo (y)
X = df_processed.drop('se_ausento', axis=1)
y = df_processed['se_ausento']

# Guardar las columnas del modelo
model_columns = X.columns.tolist()
joblib.dump(model_columns, 'model_columns.pkl')
print(f"Columnas del modelo guardadas en 'model_columns.pkl'. Total: {len(model_columns)}")

# c) Dividir en conjuntos de entrenamiento y prueba (train/val/test)
X_trainval, X_test, y_trainval, y_test = train_test_split(X, y, test_size=0.2, random_state=SEED, stratify=y)
X_train, X_val, y_train, y_val = train_test_split(X_trainval, y_trainval, test_size=0.25, random_state=SEED, stratify=y_trainval)

# d) Escalar características numéricas
scaler = StandardScaler()
X_train_scaled = scaler.fit_transform(X_train)
X_val_scaled = scaler.transform(X_val)
X_test_scaled = scaler.transform(X_test)

# Guardar el scaler
joblib.dump(scaler, 'scaler.pkl')
print("Scaler guardado en 'scaler.pkl'")

# --- 3. Definición y Compilación de la Red Neuronal Mejorada ---
model = Sequential([
    Input(shape=(X_train_scaled.shape[1],)),

    Dense(128, use_bias=False, kernel_regularizer=regularizers.l2(0.001)),
    BatchNormalization(),
    Activation('relu'),
    Dropout(0.3),

    Dense(64, use_bias=False, kernel_regularizer=regularizers.l2(0.001)),
    BatchNormalization(),
    Activation('relu'),
    Dropout(0.3),

    Dense(32, activation='relu', kernel_regularizer=regularizers.l2(0.001)),
    Dropout(0.3),

    Dense(1, activation='sigmoid')
])

model.compile(
    optimizer=tf.keras.optimizers.Adam(learning_rate=0.0005),
    loss='binary_crossentropy',
    metrics=[
        'accuracy',
        metrics.Precision(name='precision'),
        metrics.Recall(name='recall'),
        metrics.AUC(name='auc')
    ]
)

model.summary()

# --- 4. Entrenamiento del Modelo con EarlyStopping ---
early_stopping = EarlyStopping(monitor='val_loss', patience=15, restore_best_weights=True)

print("\nEntrenando el modelo...")
history = model.fit(
    X_train_scaled,
    y_train,
    epochs=100,
    batch_size=16,
    validation_data=(X_val_scaled, y_val),
    callbacks=[early_stopping],
    verbose=1
)

# --- 5. Evaluación del Modelo ---
loss, accuracy, precision, recall, auc = model.evaluate(X_test_scaled, y_test)
f1_score = 2 * (precision * recall) / (precision + recall + 1e-7)

print(f"\nEvaluación en datos de prueba:")
print(f"  - Pérdida: {loss:.4f}")
print(f"  - Accuracy: {accuracy*100:.2f}%")
print(f"  - Precision: {precision:.4f}")
print(f"  - Recall: {recall:.4f}")
print(f"  - AUC: {auc:.4f}")
print(f"  - F1-Score: {f1_score:.4f}")

# --- 6. Guardado del Modelo Entrenado ---
model.save('modelo_ausentismo.h5')
print("\n¡Modelo entrenado y guardado como 'modelo_ausentismo.h5'!")

# --- 7. Graficación de métricas ---
def plot_metric(history, metric):
    plt.plot(history.history[metric], label='Entrenamiento')
    plt.plot(history.history[f'val_{metric}'], label='Validación', linestyle='--')
    plt.title(f'{metric.upper()} vs Épocas')
    plt.xlabel('Épocas')
    plt.ylabel(metric.capitalize())
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.show()

for metric in ['loss', 'accuracy', 'precision', 'recall', 'auc']:
    if f'val_{metric}' in history.history:
        plot_metric(history, metric)