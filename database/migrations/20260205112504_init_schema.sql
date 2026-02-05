-- +goose Up
-- create "categories" table
CREATE TABLE "categories" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "name" text NOT NULL, "description" text NULL, "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "idx_categories_deleted_at" to table: "categories"
CREATE INDEX "idx_categories_deleted_at" ON "categories" ("deleted_at");
-- create index "uni_categories_name" to table: "categories"
CREATE UNIQUE INDEX "uni_categories_name" ON "categories" ("name");
-- create "users" table
CREATE TABLE "users" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "name" text NOT NULL, "email" text NOT NULL, "password" text NOT NULL, "role" text NOT NULL, "provider" text NOT NULL, "picture" text NULL, "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"));
-- create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "users" ("deleted_at");
-- create index "uni_users_email" to table: "users"
CREATE UNIQUE INDEX "uni_users_email" ON "users" ("email");
-- create "products" table
CREATE TABLE "products" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "name" text NOT NULL, "description" text NULL, "sku" text NOT NULL, "price" numeric(15,2) NOT NULL, "stock" bigint NOT NULL DEFAULT 0, "category_id" uuid NULL, "is_active" boolean NOT NULL DEFAULT true, "image" text NULL, "created_at" timestamptz NULL, "updated_at" timestamptz NULL, "deleted_at" timestamptz NULL, PRIMARY KEY ("id"), CONSTRAINT "fk_categories_products" FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- create index "idx_products_deleted_at" to table: "products"
CREATE INDEX "idx_products_deleted_at" ON "products" ("deleted_at");
-- create index "uni_products_sku" to table: "products"
CREATE UNIQUE INDEX "uni_products_sku" ON "products" ("sku");

-- +goose Down
-- reverse: create index "uni_products_sku" to table: "products"
DROP INDEX "uni_products_sku";
-- reverse: create index "idx_products_deleted_at" to table: "products"
DROP INDEX "idx_products_deleted_at";
-- reverse: create "products" table
DROP TABLE "products";
-- reverse: create index "uni_users_email" to table: "users"
DROP INDEX "uni_users_email";
-- reverse: create index "idx_users_deleted_at" to table: "users"
DROP INDEX "idx_users_deleted_at";
-- reverse: create "users" table
DROP TABLE "users";
-- reverse: create index "uni_categories_name" to table: "categories"
DROP INDEX "uni_categories_name";
-- reverse: create index "idx_categories_deleted_at" to table: "categories"
DROP INDEX "idx_categories_deleted_at";
-- reverse: create "categories" table
DROP TABLE "categories";
