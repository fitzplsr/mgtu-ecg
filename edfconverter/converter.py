import json
import pyedflib
from vecg import make_vecg

def edf_to_json_with_vecg(edf_path: str) -> json:
    # Открываем EDF
    reader = pyedflib.EdfReader(edf_path)
    n = reader.signals_in_file

    # Список всех каналов и словарь для быстрого доступа
    channels = []
    sigs = {}
    for i in range(n):
        label = reader.getLabel(i)
        data = reader.readSignal(i)
        meta = {
            "label": label,
            "dimension": reader.getPhysicalDimension(i),
            "sample_frequency": reader.getSampleFrequency(i),
            "physical_max": reader.getPhysicalMaximum(i),
            "physical_min": reader.getPhysicalMinimum(i),
            "digital_max": reader.getDigitalMaximum(i),
            "digital_min": reader.getDigitalMinimum(i),
            "prefilter": reader.getPrefilter(i),
            "transducer": reader.getTransducer(i)
        }
        channels.append({**meta, "signal": data.tolist()})
        sigs[label] = data
    reader.close()

    result = {"channels": channels}

    # Если есть минимум I, II и V1–V6 — считаем x,y,z
    needed = ['ECG I','ECG II','ECG V1','ECG V2','ECG V3','ECG V4','ECG V5','ECG V6']
    if all(ch in sigs for ch in needed):
        order = ['ECG I','ECG II','ECG V1','ECG V2','ECG V3','ECG V4','ECG V5','ECG V6']
        vecs = [sigs[ch] for ch in order]
        x, y, z = make_vecg(vecs)
        result["vector_ecg_xyz"] = {
            "x": x.tolist(),
            "y": y.tolist(),
            "z": z.tolist()
        }

    return result
