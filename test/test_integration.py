import subprocess
import unittest
from tempfile import TemporaryDirectory

from utils import Server, does_directory_contains_file


class BigIntegrationTest(unittest.TestCase):
    server: Server
    directory: TemporaryDirectory

    def setUp(self) -> None:
        self.server = Server()
        self.server.start()
        self.directory = TemporaryDirectory()

    def test_download_images_from_html_and_save_to_directory(self):
        # when the server is running, run the program to download images
        subprocess.run(
            [
                "go",
                "run",
                "main.go",
                "html",
                self.server.url,
                "-d",
                self.directory.name,
            ],
            check=True,
        )

        # then check if the images were downloaded to the file system
        self.assertTrue(does_directory_contains_file(self.directory.name, "300.jpeg"))
        self.assertTrue(does_directory_contains_file(self.directory.name, "800.jpeg"))
        self.assertTrue(does_directory_contains_file(self.directory.name, "1350.jpeg"))


    def tearDown(self) -> None:
        self.server.stop()
        self.directory.cleanup()


if __name__ == "__main__":
    unittest.main()
