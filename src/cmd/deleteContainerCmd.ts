import ora from "ora";

// import { Defaults } from "../constants";
import { docker } from "../sdks/docker";

export async function deleteContainerCmd(
  containerId: string,
  options: { format: "json" | "table" }
) {
  const spinner = ora("Starting the process... \n").start();

  try {
    const container = docker.getContainer(containerId);

    spinner.text = "Stopping container...";
    await container.stop();

    spinner.text = "Deleting container...";
    await container.remove();

    spinner.succeed("Container successfully deleted!");
  } catch (error: any) {
    spinner.clear();

    if (error.statusCode === 404) {
      console.error("The container id you've provided is invalid.");
    }

    process.exit(1);
  }
}
