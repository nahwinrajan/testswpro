-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

CREATE TABLE "estates" (
  "estate_id" text PRIMARY KEY,
  "width" int DEFAULT 0,
  "length" int DEFAULT 0,
  "count" int DEFAULT 0,
  "min" int DEFAULT 0,
  "max" int DEFAULT 0,
  "median" int DEFAULT 0,
  "patrol_distance" int DEFAULT 0,
  "patrol_route" text DEFAULT '',
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
);

CREATE TABLE "trees" (
  "tree_id" text PRIMARY KEY,
  "estate_id" text,
  "x" int NOT NULL,
  "y" int NOT NULL,
  "height" int NOT NULL,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  UNIQUE ("estate_id", "x", "y")
);

CREATE INDEX "idx_trees_estate_x_y" ON "trees" ("estate_id", "x", "y");

ALTER TABLE "trees" ADD FOREIGN KEY ("estate_id") REFERENCES "estates" ("estate_id");
