from typing import Optional
import sys
import os
import signal
import argparse
import asyncio
from watchdog.observers import Observer
from watchdog.events import FileSystemEventHandler, FileSystemEvent, FileModifiedEvent


class EventTs(asyncio.Event):

    def set(self):
        self._loop.call_soon_threadsafe(super().set)

    def clear(self):
        self._loop.call_soon_threadsafe(super().clear)


class ProcessRunner:

    def __init__(self, command: str) -> None:
        self.commad = command
        self.process = None
        self.running = False

    async def respawn_process(self) -> None:

        if self.is_process_running:
            await self.terminate_process()

        self.process = await asyncio.create_subprocess_shell(
            self.commad,
            stdout=asyncio.subprocess.PIPE,
            stderr=asyncio.subprocess.PIPE,
            start_new_session=True
        )

    def is_process_running(self) -> bool:
        if self.process:
            if self.process.returncode is None:
                return True
        return False

    async def terminate_process(self) -> None:
        if self.process:
            if self.is_process_running():
                os.killpg(os.getpgid(self.process.pid), signal.SIGKILL)

    async def start(self, restart_ev: EventTs) -> None:

        await self.respawn_process()

        while True:

            if self.is_process_running():

                stdout_task = asyncio.create_task(
                    self.process.stdout.readline())
                stderr_task = asyncio.create_task(
                    self.process.stderr.readline())
                event_task = asyncio.create_task(restart_ev.wait())

                done, pending = await asyncio.wait(
                    [stdout_task, stderr_task, event_task,],
                    return_when=asyncio.FIRST_COMPLETED,
                )

                if event_task in done:
                    await self.respawn_process()
                    restart_ev.clear()
                    continue

                if stdout_task in done:
                    if stdout_task.result() == b'':
                        continue
                    print(stdout_task.result().decode(),
                          file=sys.stdout, end='')

                if stderr_task in done:
                    if stderr_task.result() == b'':
                        continue
                    print(stderr_task.result().decode(),
                          file=sys.stderr, end='')

                for ptask in pending:
                    ptask.cancel()

            else:
                await restart_ev.wait()
                await self.respawn_process()
                restart_ev.clear()


class FileChangedHandler(FileSystemEventHandler):

    def __init__(self, restart_ev: EventTs, file_extensions: Optional[list[str]], exclude_dirs: Optional[list[str]]):
        self.restart_ev = restart_ev
        self.file_extensions = file_extensions
        self.exclude_dirs = exclude_dirs

    def on_modified(self, file_ev: FileSystemEvent):

        if file_ev.is_directory:
            return

        if file_ev.event_type == FileModifiedEvent.event_type:

            trigger_restart = False

            if self.exclude_dirs:
                for exclude_dir in self.exclude_dirs:
                    if file_ev.src_path.startswith(exclude_dir):
                        return

            if self.file_extensions:
                for file_extension in self.file_extensions:
                    if file_ev.src_path.endswith(file_extension):
                        trigger_restart = True
            else:
                trigger_restart = True

            if trigger_restart:
                print(f'WATCHIT: Detected file change "{file_ev.src_path}". Restarting ...',
                      file=sys.stderr)
                self.restart_ev.set()

        return


async def main(config: dict, process_runner: ProcessRunner):

    restart_ev = EventTs()
    file_changed_handler = FileChangedHandler(
        restart_ev,
        config["file_extensions"],
        config["exclude_directories"])
    observer = Observer()
    observer.schedule(
        file_changed_handler,
        path=config["directory"],
        recursive=True)
    observer.start()

    try:
        await process_runner.start(restart_ev)

    except KeyboardInterrupt:
        pass
    finally:
        observer.stop()
        observer.join()
    print("Interrupted")


if __name__ == "__main__":

    parser = argparse.ArgumentParser(
        description="""Watchit runs a command in a subprocess,
                watches for file changes in a directory,
                and restarts the subprocess when changes occur.""",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument("-d", "--directory",
                        help="directory to watch", default="./")
    parser.add_argument("-e", "--file-extensions",
                        help="watch only files with these extensions", nargs="*")
    parser.add_argument("-x", "--exclude-directories",
                        help="Directories to exclude from watch", nargs="*")
    parser.add_argument("-c", "--command",
                        help="Command to run", required=True)
    args = parser.parse_args()
    config = vars(args)

    process_runner = ProcessRunner(config["command"])

    try:
        asyncio.run(main(config, process_runner))
    except KeyboardInterrupt:
        pass
    finally:

        async def cleanup_process():
            await process_runner.terminate_process()

        asyncio.run(cleanup_process())