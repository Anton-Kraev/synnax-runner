#!/usr/bin/env python3
"""
Пример Python скрипта с ошибкой для тестирования обработки ошибок
"""

import datetime

def main():
    print("🚀 Python скрипт с ошибкой запущен!")
    print(f"⏰ Время запуска: {datetime.datetime.now()}")
    
    # Имитация работы скрипта
    print("📊 Выполняю вычисления...")
    
    # Намеренно вызываем ошибку
    print("❌ Вызываю ошибку...")
    undefined_variable = some_undefined_function()
    
    print("✅ Этот код не должен выполниться!")

if __name__ == "__main__":
    main()
