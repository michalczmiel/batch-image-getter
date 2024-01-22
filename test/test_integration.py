import os
import unittest
import subprocess
from tempfile import TemporaryDirectory

from utils import does_file_exist, Server


class BigIntegrationTest(unittest.TestCase):
    server: Server

    def setUp(self) -> None:
        self.server = Server()
        self.server.start()

    def test_download_images_from_html_and_save_to_directory(self):
        with TemporaryDirectory() as temp_dir:
            # when the server is running, run the program to download images
            subprocess.run(
                [
                    "go",
                    "run",
                    "main.go",
                    "html",
                    self.server.url,
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
        self.server.stop()


if __name__ == "__main__":
    unittest.main()
