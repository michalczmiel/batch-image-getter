import os
import http.server
from threading import Thread
from typing import Final

PORT: Final[int] = 8080
DIRECTORY: Final[str] = os.path.join(
    os.path.dirname(os.path.realpath(__file__)), "testdata/"
)


def does_directory_contains_file(path: str, file: str) -> bool:
    """
    Checks if the given directory contains the given file
    """
    return os.path.isfile(os.path.join(path, file))


class Server:
    """
    A simple HTTP server that serves files from the given directory
    """

    class Handler(http.server.SimpleHTTPRequestHandler):
        def __init__(self, *args, **kwargs):
            super().__init__(*args, directory=DIRECTORY, **kwargs)

    def __init__(self):
        self.server = http.server.HTTPServer(("", PORT), self.Handler)

    def start(self) -> None:
        thread = Thread(target=self.server.serve_forever)
        thread.start()

    def stop(self) -> None:
        thread = Thread(target=self.server.shutdown)
        thread.start()

    @property
    def url(self) -> str:
        return f"http://localhost:{self.server.server_port}"
