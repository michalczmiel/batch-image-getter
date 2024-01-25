import pytest
from utils import Server, run_cli


@pytest.fixture(scope="module")
def server() -> Server:
    server = Server()
    server.start()
    yield server
    server.stop()


def test_download_images_from_html_and_save_to_directory(server, tmpdir):
    # run the program to download images from hosted html page
    run_cli(["html", server.url, "-d", tmpdir.strpath])

    # then check if the images were downloaded to the file system
    assert (tmpdir / "300.jpeg").exists()
    assert (tmpdir / "800.jpeg").exists()
    assert (tmpdir / "1350.jpeg").exists()


def test_download_images_from_txt_file_and_save_to_directory(server, tmpdir):
    # create a text file with the URLs of the images
    text_file = tmpdir / "urls.txt"
    text_file.write(
        "\n".join(
            [
                "https://picsum.photos/600/600",
                "https://picsum.photos/1200/1200",
                "https://picsum.photos/2400/2400",
            ]
        )
    )

    # run the program to download images from the text file
    run_cli(["file", text_file.strpath, "-d", tmpdir.strpath])

    # then check if the images were downloaded to the file system
    assert (tmpdir / "600.jpeg").exists()
    assert (tmpdir / "1200.jpeg").exists()
    assert (tmpdir / "2400.jpeg").exists()
