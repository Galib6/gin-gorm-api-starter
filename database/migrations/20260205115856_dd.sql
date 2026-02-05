-- +goose Up
-- modify "categories" table
ALTER TABLE "categories" ADD CONSTRAINT "uni_categories_name" UNIQUE USING INDEX "uni_categories_name";
-- modify "products" table
ALTER TABLE "products" ADD CONSTRAINT "uni_products_sku" UNIQUE USING INDEX "uni_products_sku";
-- modify "users" table
ALTER TABLE "users" ADD CONSTRAINT "uni_users_email" UNIQUE USING INDEX "uni_users_email";

-- +goose Down
-- reverse: modify "users" table
ALTER TABLE "users" DROP CONSTRAINT "uni_users_email";
-- reverse: modify "products" table
ALTER TABLE "products" DROP CONSTRAINT "uni_products_sku";
-- reverse: modify "categories" table
ALTER TABLE "categories" DROP CONSTRAINT "uni_categories_name";
