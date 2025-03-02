export const POSTGRES_DEFAULT_IMAGE_TAG = "postgres";
export const POSTGRES_DEFAULT_IMAGE_VERSION = "17.0-alpine3.19";
export const POSTGRES_VOLUME_NAME = "psql";
export const POSTGRES_MOUNT_TARGET = "/var/lib/postgresql/data";
export const POSTGRES_PORT = 5432;

export function getPostgresConnectionURL(pw: string, ip: string, port: number) {
  return `postgres://postgres:${pw}@${ip}:${port}/postgres`;
}

export const POSTGRES_REQUIRED_ENV = [
  "POSTGRES_PASSWORD",
  "POSTGRES_DB",
  "POSTGRES_USER",
];

// for future
export const postgresJsonTemplate = {
  "@type": "template_planner",
  image: {
    tag: POSTGRES_DEFAULT_IMAGE_TAG,
    version: POSTGRES_DEFAULT_IMAGE_VERSION,
  },
  envs: POSTGRES_DEFAULT_IMAGE_TAG,
};
