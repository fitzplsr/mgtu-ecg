from fastapi import FastAPI, Query, HTTPException
from converter import edf_to_json_with_vecg
from pathlib import Path
from pydantic import BaseModel
import logging

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

app = FastAPI()

class FileRequest(BaseModel):
    input_file: str

@app.post("/convert")
def process(request: FileRequest):
    path = Path(request.input_file)
    if not path.exists():
        logger.error(f"Файл не найден: {request.input_file}")
        raise HTTPException(status_code=404, detail=f'Файл не найден: {path}')

    try:
        result = edf_to_json_with_vecg(path.as_posix())
        return result
    except Exception as e:
        logger.error(f"Ошибка при обработке файла {request.input_file}: {e}")
        raise HTTPException(status_code=500, detail=str(e))


