CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "firstname" varchar(30) NOT NULL,
  "lastname" varchar(30) NOT NULL,
  "username" varchar(20) UNIQUE NOT NULL,
  "email" varchar(30) UNIQUE NOT NULL,
  "password" varchar(100) NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "username" varchar(20) NOT NULL,
  "amount" bigint NOT NULL,
  "transaction_type" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "wallets" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint UNIQUE NOT NULL,
  "username" varchar(20) UNIQUE NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("id");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "transactions" ("user_id");

CREATE INDEX ON "transactions" ("username");

CREATE INDEX ON "transactions" ("transaction_type");

CREATE INDEX ON "transactions" ("id");

CREATE INDEX ON "transactions" ("user_id", "transaction_type");

CREATE INDEX ON "wallets" ("user_id");

CREATE INDEX ON "wallets" ("username");

ALTER TABLE "transactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "transactions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "wallets" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "wallets" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
