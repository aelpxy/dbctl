import { Defaults } from "../constants";
import { docker } from "../sdks/docker";

export async function listContainersCmd({
  format,
}: {
  format: "json" | "table";
}) {
  const filteredContainers = [];

  try {
    const containers = await docker.listContainers();

    for (const container of containers) {
      if (
        container.Names[0]
          .replace(/^\/+/, "")
          .startsWith(Defaults.DockerContainerPrefix)
      ) {
        const containerInstance = docker.getContainer(container.Id);
        const containerInfo = await containerInstance.inspect();

        const { HostIp, HostPort } =
          containerInfo.NetworkSettings.Ports[
            Object.keys(containerInfo.NetworkSettings.Ports)[0]
          ][0];

        filteredContainers.push({
          ID: container.Id.slice(0, 12),
          Name: container.Names[0].slice(7),
          State: container.State,
          Uptime: container.Status,
          IP: `${HostIp}:${HostPort}`,
          Image: container.Image.split(":")[0],
        });
      }
    }

    if (format !== "json") {
      console.table(filteredContainers);
      return;
    }

    console.info(JSON.stringify(filteredContainers));
  } catch (error) {
    process.exit(1);
  }
}
