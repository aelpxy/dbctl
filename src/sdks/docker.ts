import Docker from "dockerode";
import process from "node:process";

const docker = new Docker({ socketPath: "/var/run/docker.sock" });

try {
  await docker.ping();
} catch (error: any) {
  if (error.code === "FailedToOpenSocket") {
    console.error("Are you sure Docker is properly set up and running?");
    process.exit(1);
  }
}

export { docker };
