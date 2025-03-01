import { Command, Option } from "commander";

import {
  CmdLongDescription,
  CmdName,
  CmdShortDescription,
  CmdVersion,
} from "../constants.ts";

import { listContainersCmd } from "./lsCmd.ts";

const cmd = new Command();

cmd
  .name(CmdName)
  .summary(CmdShortDescription)
  .description(CmdLongDescription)
  .version(CmdVersion);

cmd
  .command("create")
  .description("Create a new database container")
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("delete")
  .description("Stop and delete a database container")
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("backup")
  .description("Backup a database container")
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("ls")
  .description("List all running database containers")
  .addOption(
    new Option("-t, --format <json>", "format to print")
      .choices(["json", "table"])
      .default("table")
  )
  .action(listContainersCmd);

cmd
  .command("inspect")
  .description("Inspect a running database container")
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("logs")
  .description("Stream live logs of a database container")
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("shell")
  .description("Connect to a running database container")
  .action(async () => {
    console.info("cmd");
  });

cmd
  .command("http")
  .description("Start the API and serve it")
  .action(async () => {
    console.info("cmd");
  });

export { cmd };
