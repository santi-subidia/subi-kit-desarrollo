---
name: clean-code
title: Principios de Clean Code y Calidad de Software
category: core
always_on: true
description: Principios fundamentales de legibilidad, arquitectura limpia, SOLID, DRY, KISS y manejo defensivo de errores.
tags: [core, clean-code, architecture]
---

# Principios de Clean Code y Calidad de Software

## 1. Nombres y Expresividad
- **Nombres con significado**: Variables, funciones y clases deben revelar su intención sin necesidad de comentarios explicativos obvios.
- **Funciones pequeñas y de responsabilidad única (SRP)**: Una función debe hacer una sola cosa y hacerla bien.
- **Evitar números mágicos**: Usar constantes descriptivas con nombres en mayúsculas o enums.

## 2. Simplicidad y Diseño
- **KISS (Keep It Simple, Stupid)**: Preferir la solución más simple que cumpla con los requisitos. Evitar sobreingeniería o abstracciones prematuras.
- **DRY (Don't Repeat Yourself)**: Reutilizar lógica común extrayéndola en helpers, hooks o servicios cuando haya duplicación evidente (regla de 3).
- **YAGNI (You Aren't Gonna Need It)**: No construir código para casos de uso hipotéticos futuros.

## 3. Manejo de Errores y Programación Defensiva
- **Fail Fast**: Validar entradas y argumentos en las fronteras de las funciones y APIs.
- **Manejo explícito de excepciones**: Evitar bloques `catch` vacíos o que silencien errores sin registrar contexto.
- **Inmutabilidad por defecto**: Preferir estructuras de datos inmutables y funciones puras cuando sea posible.
