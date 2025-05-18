Запуск всего
```
    docker compose up
```

Запуск только бд
```
    docker compose -f "local-docker-compose.yaml" up -d
```

Схема сигнала
```js
{
  "channels": [
    {
      "label": "ECG I",
      "dimension": "mV",
      "sample_frequency": 500.0,
      "physical_max": 7.4,
      "physical_min": -7.4,
      "digital_max": 32766,
      "digital_min": -32766,
      "prefilter": "NF:Off",
      "transducer": "",
      "signal": []
    }
  ],
  "vector_ecg_xyz": {
    "x": [],
    "y": [],
    "z": []
  }
}
```