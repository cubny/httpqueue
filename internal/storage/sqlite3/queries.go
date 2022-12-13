package sqlite3

const queryInsertCart = `
INSERT INTO carts(user_id, created_at, updated_at) values (?,?,?)
`
const queryCartsByIDAndUserID = `
SELECT id, user_id, created_at, updated_at FROM carts 
WHERE id = ? AND user_id = ?
`

const queryInsertItem = `
INSERT INTO line_items (cart_id, product_id, quantity, price, created_at, updated_at) 
values (?,?,?,?,?,?)
`
const queryItemsByCartIDAndProductID = `
SELECT id, cart_id, product_id, quantity, price, created_at, updated_at FROM line_items
WHERE cart_id = ? and product_id = ?
`
const queryItemByID = `
SELECT id, cart_id, product_id, quantity, price, created_at, updated_at FROM line_items
WHERE id = ? 
`
const queryRemoveItem = `
DELETE FROM line_items where id = ?;
`

const queryRemoveItemsByCartID = `
DELETE FROM line_items where cart_id = ?;
`

// Migrations -----------------------

const migration01MigrationCreateCartsTable = `
CREATE TABLE IF NOT EXISTS "carts" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  "user_id" integer,
  "created_at" datetime NOT NULL,
  "updated_at" datetime NOT NULL,
  CONSTRAINT "fk_32e17d5e33" FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
`

const migration02AddCartIndex = `
CREATE INDEX "index_cart_on_user_id" ON "carts" ("user_id");
`

const migration03CreateLineItemsTable = `
CREATE TABLE IF NOT EXISTS "line_items" (
  "id" INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
  "cart_id" integer,
  "product_id" integer,
  "quantity" integer DEFAULT 1,
  "price" decimal,
  "created_at" datetime NOT NULL,
  "updated_at" datetime NOT NULL,
  CONSTRAINT "fk_af645e8e5f" FOREIGN KEY ("cart_id") REFERENCES "carts" ("id")
);
`

const migration04AddLineItemsIndex = `
CREATE INDEX "index_line_items_on_cart_id" ON "line_items" ("cart_id");
`
const migration05AddLineItemsIndex2 = `
CREATE INDEX "index_line_items_on_product_id" ON "line_items" ("product_id");
`

const truncateCartsTable = `DELETE FROM carts;`
const truncateLineItemsTable = `DELETE FROM line_items;`
