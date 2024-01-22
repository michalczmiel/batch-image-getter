import os
import unittest
import http.server
from threading import Thread
from typing import Final


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

    def test_download_images_from_html(self):
        # when the server is running, run the program to download images
        os.system(f"go run main.go html http://localhost:{PORT} -d images/")

        # then check if the images were downloaded to the file system
        self.assertTrue(does_file_exist("images/300.jpeg"))
        self.assertTrue(does_file_exist("images/800.jpeg"))
        self.assertTrue(does_file_exist("images/1350.jpeg"))

    def tearDown(self) -> None:
        thread = Thread(target=self.server.shutdown)
        thread.start()


if __name__ == "__main__":
    unittest.main()
