CREATE TYPE "friendship_status" AS ENUM (
  'pending',
  'accepted',
  'rejected',
  'blocked'
);

CREATE TYPE "post_visibility" AS ENUM (
  'friends',
  'public'
);

CREATE TYPE "reaction_target_type" AS ENUM (
  'post',
  'comment'
);

CREATE TYPE "reaction_type" AS ENUM (
  'like',
  'love',
  'haha',
  'wow',
  'sad',
  'angry'
);

CREATE TYPE "notification_type" AS ENUM (
  'friend_request',
  'friend_accept',
  'post_comment',
  'post_reaction',
  'comment_reply',
  'comment_reaction'
);

CREATE TABLE "users" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "username" varchar UNIQUE NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_hash" varchar NOT NULL,
  "display_name" varchar NOT NULL,
  "bio" text,
  "avatar_url" varchar,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE TABLE "friendships" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "requester_id" uuid NOT NULL,
  "addressee_id" uuid NOT NULL,
  "status" friendship_status NOT NULL DEFAULT 'pending',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);  

CREATE TABLE "posts" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "content" text NOT NULL,
  "visibility" post_visibility NOT NULL DEFAULT 'friends',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "deleted_at" timestamp,
  "post_parent_id" uuid,
  "author_id" uuid NOT NULL
);

CREATE TABLE "comments" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "content" text NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "deleted_at" timestamp,
  "author_id" uuid NOT NULL,
  "post_id" uuid NOT NULL,
  "parent_comment_id" uuid
);

CREATE TABLE "reactions" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "target_type" reaction_target_type NOT NULL,
  "type" reaction_type NOT NULL DEFAULT 'like',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "user_id" uuid NOT NULL,
  "target_id" uuid NOT NULL
);

CREATE TABLE "notifications" (
  "id" uuid PRIMARY KEY DEFAULT (gen_random_uuid()),
  "user_id" uuid NOT NULL,
  "actor_id" uuid NOT NULL,
  "type" notification_type NOT NULL,
  "target_type" varchar,
  "target_id" uuid,
  "read_at" timestamp,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("email");

CREATE UNIQUE INDEX ON "friendships" ("requester_id", "addressee_id");

CREATE INDEX ON "friendships" ("addressee_id");

CREATE INDEX ON "posts" ("author_id");

CREATE INDEX ON "posts" ("created_at");

CREATE INDEX ON "comments" ("post_id");

CREATE INDEX ON "comments" ("parent_comment_id");

CREATE UNIQUE INDEX ON "reactions" ("user_id", "target_type", "target_id");

CREATE INDEX ON "reactions" ("target_type", "target_id");

CREATE INDEX ON "notifications" ("user_id", "read_at");

CREATE INDEX ON "notifications" ("created_at");

COMMENT ON TABLE "friendships" IS 'Una fila por par de usuarios. requester_id = quien envía la solicitud.';

COMMENT ON TABLE "comments" IS 'parent_comment_id nulo = comentario directo al post. No nulo = respuesta a otro comentario (un solo nivel de anidación en el MVP).';

COMMENT ON TABLE "reactions" IS 'target_id apunta a posts.id o comments.id según target_type (polimórfico, sin FK real por diseño). El índice único evita reacciones duplicadas del mismo usuario.';

COMMENT ON TABLE "notifications" IS 'user_id = quien recibe la notificación. actor_id = quien la originó (quien reaccionó, comentó, etc).';

ALTER TABLE "friendships" ADD FOREIGN KEY ("requester_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "friendships" ADD FOREIGN KEY ("addressee_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "posts" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "comments" ADD FOREIGN KEY ("author_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "comments" ADD FOREIGN KEY ("parent_comment_id") REFERENCES "comments" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "reactions" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "notifications" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;

ALTER TABLE "notifications" ADD FOREIGN KEY ("actor_id") REFERENCES "users" ("id") DEFERRABLE INITIALLY IMMEDIATE;
