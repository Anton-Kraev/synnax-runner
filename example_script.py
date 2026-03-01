#!/usr/bin/env python3
"""
Пример Python скрипта для тестирования Telegram бота
"""

import datetime
import json
import random

def main():
    print("🚀 Python скрипт запущен!")
    print(f"⏰ Время запуска: {datetime.datetime.now()}")
    
    # Имитация работы скрипта
    print("📊 Выполняю вычисления...")
    
    # Генерируем случайные данные
    data = {
        "timestamp": datetime.datetime.now().isoformat(),
        "random_number": random.randint(1, 100),
        "status": "success",
        "message": "Скрипт выполнен успешно!"
    }
    
    print("📋 Результат работы:")
    print(json.dumps(data, indent=2, ensure_ascii=False))
    
    print("✅ Скрипт завершен успешно!")

if __name__ == "__main__":
    main()
