// import type { ContainerCreateOptions } from "dockerode";
// import { Defaults, SupportedDatabases } from "../constants";
// import {
//   POSTGRES_DEFAULT_IMAGE_TAG,
//   POSTGRES_DEFAULT_IMAGE_VERSION,
//   POSTGRES_MOUNT_TARGET,
//   POSTGRES_PORT,
// } from "../dbs/postgres";

// import { docker } from "../sdks/docker";

type FlagOptions = {
  port: string | null;
  password: string | null;
  name: string | null;
  format: "json" | "table";
};

export async function createContainerCmd(type: string, options: FlagOptions) {
  try {
    if (type === "postgres") {
      //   const volume = await docker.createVolume({ Name: "dbctl.volume.1" });

      //   const containerConfig: ContainerCreateOptions = {
      //     Image: `${POSTGRES_DEFAULT_IMAGE_TAG}:${POSTGRES_DEFAULT_IMAGE_VERSION}`,
      //     name: options.name || "hello-world",
      //     Env: [
      //       `POSTGRES_PASSWORD=${options.password || "password"}`,
      //       `POSTGRES_DB=postgres`,
      //       "POSTGRES_USER=postgres",
      //     ],
      //     HostConfig: {
      //       PortBindings: {
      //         "5432/tcp": [{ HostPort: options.port }],
      //       },
      //       Binds: [`${volume.Name}:${POSTGRES_MOUNT_TARGET}`],
      //     },
      //   };

      //   const container = await docker.createContainer(containerConfig);

      //   await container.start();

      if (options.format === "json") {
        process.exit(1);
      }
    }
  } catch (error) {
    process.exit(1);
  }
}
