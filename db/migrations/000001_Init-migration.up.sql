CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "firstname" varchar(30) NOT NULL,
  "lastname" varchar(30) NOT NULL,
  "username" varchar(20) UNIQUE NOT NULL,
  "email" varchar(30) UNIQUE NOT NULL,
  "password" varchar(100) NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "transactions" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint,
  "amount" bigint NOT NULL,
  "transaction_type" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("id");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "transactions" ("user_id");

CREATE INDEX ON "transactions" ("transaction_type");

CREATE INDEX ON "transactions" ("id");

CREATE INDEX ON "transactions" ("user_id", "transaction_type");

ALTER TABLE "transactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
