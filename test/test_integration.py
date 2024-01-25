import pytest
import subprocess

from utils import Server


@pytest.fixture
def server() -> Server:
    server = Server()
    server.start()
    yield server
    server.stop()


def test_download_images_from_html_and_save_to_directory(server, tmpdir):
    # when the server is running, run the program to download images
    subprocess.run(
        [
            "go",
            "run",
            "main.go",
            "html",
            server.url,
            "-d",
            tmpdir.strpath,
        ],
        check=True,
    )

    # then check if the images were downloaded to the file system
    assert (tmpdir / "300.jpeg").exists()
    assert (tmpdir / "800.jpeg").exists()
    assert (tmpdir / "1350.jpeg").exists()
