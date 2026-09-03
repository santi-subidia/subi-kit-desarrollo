---
name: supabase-sql
title: Base de Datos PostgreSQL, SQL y Supabase
category: backend
always_on: false
description: Buenas prácticas de PostgreSQL, Supabase, Row Level Security (RLS), migraciones, índices y generación de tipos.
tags: [supabase, postgres, sql, database, backend]
---

# PostgreSQL y Supabase: Estándares y Seguridad

## 1. Seguridad y Row Level Security (RLS)
- **RLS Obligatorio**: Todas las tablas públicas deben tener `ENABLE ROW LEVEL SECURITY`.
- **Políticas explícitas**: Definir políticas separadas para `SELECT`, `INSERT`, `UPDATE` y `DELETE` usando `auth.uid()`.
- **Service Role vs Anon Key**: Usar `supabaseAdmin` (service role) únicamente en entornos de servidor seguros (Server Actions / API Routes autenticadas) para operaciones administrativas que deben evadir RLS de forma controlada.

## 2. Modelado de Datos y Migraciones
- **Claves Primarias y Foráneas**: Usar `UUID` (`gen_random_uuid()`) o `BIGINT GENERATED ALWAYS AS IDENTITY` para PKs. Definir constraints de integridad referencial (`ON DELETE CASCADE` o `RESTRICT`).
- **Nombres en snake_case**: Tablas y columnas deben nombrarse en minúsculas con guiones bajos (`user_profiles`, `created_at`).
- **Migraciones Versionadas**: No alterar schemas en producción manualmente; utilizar scripts SQL de migración reproducibles (`supabase/migrations/`).
- **Índices estratégicos**: Crear índices en columnas frecuentemente filtradas o utilizadas en `JOIN` (ej. `user_id`, `status`, `created_at`).

## 3. Tipado con TypeScript
- Generar y sincronizar los tipos de TypeScript con la base de datos (`supabase gen types typescript --project-id ... > src/types/database.types.ts`).
- Usar el cliente tipado `createClient<Database>()` para autocompletado y type safety total en queries.
