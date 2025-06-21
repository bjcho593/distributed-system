# Imagen base oficial de Python
FROM python:3.11-slim

# Establecer el directorio de trabajo en el contenedor
WORKDIR /app

# Copiar los archivos del proyecto al contenedor
COPY . .

# Instalar las dependencias
RUN pip install --no-cache-dir -r requirements.txt

# Puerto en el que corre FastAPI (por defecto 8081 en este microservicio)
EXPOSE 8081

# Comando para ejecutar la aplicación con Uvicorn
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8081"]
