import torch
import numpy as np
import pyedflib

SEQ_LENGTH = 2500
DEVICE = torch.device('cuda' if torch.cuda.is_available() else 'cpu')

# ==== Загрузка и предобработка сигнала ====
def load_and_preprocess_signal(edf_path, seq_length=SEQ_LENGTH):
    with pyedflib.EdfReader(edf_path) as f:
        ch = f.getSignalLabels()
        try:
            idx = ch.index('ECG V1')
        except ValueError:
            raise ValueError("Канал 'ECG V1' не найден в файле!")
        sig = f.readSignal(idx)
    sig = sig[:seq_length].astype(np.float32)
    sig = (sig - sig.mean()) / (sig.std() + 1e-6)
    sig = torch.from_numpy(sig).unsqueeze(0).unsqueeze(0)  # (1, 1, seq_length)
    return sig

# ==== Предсказание ====
def predict_ecg(model, edf_path):
    signal = load_and_preprocess_signal(edf_path).to(DEVICE)
    with torch.no_grad():
        output = model(signal).item()
        prob = torch.sigmoid(torch.tensor(output)).item()
        if prob > 0.5:
            prediction = True
        else:
            prediction = False
    return {'predict': str(prob), 'result': prediction}
