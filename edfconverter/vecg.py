def make_vecg(signals):
    """
    Преобразует сигналы ЭКГ в ВЭКГ и добавляет координаты x, y, z.
    signals: список из 8 numpy-массивов [DI, DII, V1…V6]
    Возвращает три numpy-массива (x, y, z).
    """
    DI, DII, V1, V2, V3, V4, V5, V6 = signals

    # Коэффициенты из вашего примера
    x = -(-0.172*V1 - 0.074*V2 + 0.122*V3 + 0.231*V4 + 0.239*V5 +
          0.194*V6 + 0.156*DI - 0.01*DII)
    y = (0.057*V1 - 0.019*V2 - 0.106*V3 - 0.022*V4 + 0.041*V5 +
         0.048*V6 - 0.227*DI + 0.887*DII)
    z = -(-0.229*V1 - 0.31*V2 - 0.246*V3 - 0.063*V4 + 0.055*V5 +
          0.108*V6 + 0.022*DI + 0.102*DII)

    return x, y, z
