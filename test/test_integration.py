import os
import unittest
import http.server
import subprocess
from threading import Thread
from typing import Final
from tempfile import TemporaryDirectory


PORT: Final[int] = 8080
DIRECTORY: Final[str] = os.path.join(
    os.path.dirname(os.path.realpath(__file__)), "testdata/"
)


def does_file_exist(path: str) -> bool:
    """
    Checks if a file exists in the file system
    """
    return os.path.isfile(path)


class Handler(http.server.SimpleHTTPRequestHandler):
    """
    A simple HTTP request handler that serves files from the given directory
    """

    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=DIRECTORY, **kwargs)


class BigIntegrationTest(unittest.TestCase):
    server: http.server.HTTPServer

    def setUp(self) -> None:
        self.server = http.server.HTTPServer(("", PORT), Handler)
        thread = Thread(target=self.server.serve_forever)
        thread.start()

    def test_download_images_from_html_and_save_to_directory(self):
        with TemporaryDirectory() as temp_dir:
            # when the server is running, run the program to download images
            subprocess.run(
                [
                    "go",
                    "run",
                    "main.go",
                    "html",
                    f"http://localhost:{PORT}",
                    "-d",
                    temp_dir,
                ],
                check=True,
            )

            # then check if the images were downloaded to the file system
            self.assertTrue(does_file_exist(os.path.join(temp_dir, "300.jpeg")))
            self.assertTrue(does_file_exist(os.path.join(temp_dir, "800.jpeg")))
            self.assertTrue(does_file_exist(os.path.join(temp_dir, "1350.jpeg")))

    def tearDown(self) -> None:
        thread = Thread(target=self.server.shutdown)
        thread.start()


if __name__ == "__main__":
    unittest.main()
