from fastapi import FastAPI, Query, HTTPException
from pathlib import Path
from pydantic import BaseModel
import logging
from model import AdvancedECGModel
import torch
from process import predict_ecg


SEQ_LENGTH = 2500
MODEL_PATH = 'epoch_14_auc_0.859V1.pth'
DEVICE = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)

app = FastAPI()

class FileRequest(BaseModel):
    input_file: str

@app.post("/process")
def process(request: FileRequest):
    path = Path(request.input_file)
    if not path.exists():
        logger.error(f"Файл не найден: {request.input_file}")
        raise HTTPException(status_code=404, detail=f'Файл не найден: {path}')

    try:
        result = predict_ecg(model, path.as_posix())
        logger.info(f'Get result: {result}')
        return result
    except Exception as e:
        logger.error(f"Ошибка при обработке файла {request.input_file}: {e}")
        raise HTTPException(status_code=500, detail=str(e))

# ==== Загрузка модели ====
def load_model(model_path):
    model = AdvancedECGModel(seq_length=SEQ_LENGTH)
    model.load_state_dict(torch.load(model_path, map_location=DEVICE))
    model.to(DEVICE)
    model.eval()
    return model

model = load_model(MODEL_PATH)


