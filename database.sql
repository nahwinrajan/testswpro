-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "estates" (
  "estate_id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "width" int DEFAULT 0,
  "length" int DEFAULT 0,
  "count" int DEFAULT 0,
  "min" int DEFAULT 0,
  "max" int DEFAULT 0,
  "median" int DEFAULT 0,
  "patrol_distance" int DEFAULT 0,
  "patrol_route" text DEFAULT ''
);

CREATE TABLE "trees" (
  "tree_id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "estate_id" UUID,
  "x" int NOT NULL,
  "y" int NOT NULL,
  "height" int NOT NULL,
  CONSTRAINT "unique_trees_estate_x_y" UNIQUE ("estate_id", "x", "y")
);

CREATE INDEX "idx_trees_estate_x_y" ON "trees" ("estate_id", "x", "y");

ALTER TABLE "trees" ADD FOREIGN KEY ("estate_id") REFERENCES "estates" ("estate_id");
