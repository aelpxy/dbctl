import { Command, Option } from "commander";

import {
  CmdLongDescription,
  CmdName,
  CmdShortDescription,
  CmdVersion,
} from "../constants.ts";

import { createContainerCmd } from "./createContainerCmd.ts";
import { deleteContainerCmd } from "./deleteContainerCmd.ts";
import { listContainersCmd } from "./listContainersCmd.ts";

const cmd = new Command();

cmd
  .name(CmdName)
  .summary(CmdShortDescription)
  .description(CmdLongDescription)
  .version(CmdVersion);

cmd
  .command("create")
  .description("Create a new database container")
  .argument("<type>", "type of database to create")
  .addOption(new Option("-p, --port <port>", "port to use"))
  .addOption(new Option("-w --password <password>", "password to use"))
  .addOption(new Option("-n, --name <name>", "name to use"))
  .addOption(
    new Option("-f, --format <json>", "format to print")
      .choices(["json", "table"])
      .default("table")
  )
  .action(createContainerCmd);

cmd
  .command("delete")
  .alias("rm")
  .description("Stop and delete a database container")
  .argument("<id>", "id of the container to delete")
  .addOption(
    new Option("-f, --format <json>", "format to print")
      .choices(["json", "table"])
      .default("table")
  )
  .action(deleteContainerCmd);

cmd
  .command("backup")
  .alias("bk")
  .description("Backup a database container")
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("ls")
  .description("List all running database containers")
  .addOption(
    new Option("-f, --format <json>", "format to print")
      .choices(["json", "table"])
      .default("table")
  )
  .action(listContainersCmd);

cmd
  .command("inspect")
  .description("Inspect a running database container")
  .addOption(
    new Option("-f, --format <json>", "format to print")
      .choices(["json", "table"])
      .default("table")
  )
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("logs")
  .description("Stream live logs of a database container")
  .addOption(
    new Option("-f, --format <json>", "format to print")
      .choices(["json", "table"])
      .default("table")
  )
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("shell")
  .description("Connect to a running database container")
  .addOption(
    new Option("-f, --format <json>", "format to print")
      .choices(["json", "table"])
      .default("table")
  )
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("http")
  .description("Start the API and serve it")
  .addOption(
    new Option("-f, --format <json>", "format to print")
      .choices(["json", "table"])
      .default("table")
  )
  .action(async () => {
    console.info("cmd");
  });

export { cmd };
