from datetime import date
from io import BytesIO
from types import SimpleNamespace
from unittest.mock import Mock

import pytest
from fastapi import HTTPException, UploadFile

from app.constants.movie import MovieStatus
from app.repositories.movie_repository import MovieRepository
from app.schemas.movie import MovieCreate
from app.services import movie_service


@pytest.fixture
def repo():
    return Mock(spec=MovieRepository)


def test_list_movies(repo):
    repo.list_movies.return_value = []
    assert movie_service.MovieService(repo).list_movies() == []
    repo.list_movies.assert_called_once_with()


@pytest.mark.parametrize("status", list(MovieStatus))
def test_get_movie_builds_stream_urls(repo, monkeypatch, status):
    monkeypatch.setattr(movie_service.settings, "STREAM_BASE_URL", "http://stream.test")
    movie = SimpleNamespace(
        filename="video.mp4", status=status, hls_playlist_path="7/index.m3u8"
    )
    repo.get_movie.return_value = movie
    assert movie_service.MovieService(repo).get_movie(7) is movie
    repo.get_movie.assert_called_once_with(7)
    assert movie.stream_url == "http://stream.test/stream/video.mp4"
    assert movie.hls_url == (
        "/hls/7/index.m3u8" if status == MovieStatus.READY else None
    )


def test_get_movie_not_found(repo):
    repo.get_movie.return_value = None
    with pytest.raises(HTTPException) as error:
        movie_service.MovieService(repo).get_movie(7)
    assert (error.value.status_code, error.value.detail) == (404, "Movie not found")


def test_create_movie(repo):
    payload = MovieCreate(
        title="Film",
        description="Story",
        filename="film.mp4",
        release_date=None,
        language=None,
        genres=None,
        age_rating=None,
        director=None,
        cast=None,
    )
    assert (
        movie_service.MovieService(repo).create_movie(payload)
        is repo.create_movie.return_value
    )
    movie = repo.create_movie.call_args.args[0]
    assert (movie.title, movie.description, movie.filename) == (
        "Film",
        "Story",
        "film.mp4",
    )


@pytest.mark.parametrize("with_poster", [False, True])
def test_upload_stores_files_and_metadata(repo, monkeypatch, tmp_path, with_poster):
    media, posters = tmp_path / "movies", tmp_path / "posters"
    posters.mkdir()
    monkeypatch.setattr(movie_service, "MEDIA_DIR", media)
    monkeypatch.setattr(movie_service, "POSTER_DIR", posters)
    ids = iter(["video-id", "poster-id"])
    monkeypatch.setattr(movie_service.uuid, "uuid4", lambda: next(ids))
    video = UploadFile(filename="original.mp4", file=BytesIO(b"video content"))
    poster = (
        UploadFile(filename="cover.jpg", file=BytesIO(b"poster content"))
        if with_poster
        else None
    )
    result = movie_service.MovieService(repo).upload_movie(
        "Film",
        "Story",
        date(2024, 1, 2),
        "English",
        ["Drama"],
        "PG",
        "Director",
        ["Actor"],
        video,
        poster,
    )
    assert result is repo.create_movie.return_value
    assert (media / "video-id.mp4").read_bytes() == b"video content"
    movie = repo.create_movie.call_args.args[0]
    assert movie.filename == "video-id.mp4"
    assert (movie.title, movie.description, movie.release_date, movie.language) == (
        "Film",
        "Story",
        date(2024, 1, 2),
        "English",
    )
    assert (movie.genres, movie.age_rating, movie.director, movie.cast) == (
        ["Drama"],
        "PG",
        "Director",
        ["Actor"],
    )
    assert movie.poster_filename == ("poster-id.jpg" if with_poster else None)
    if with_poster:
        assert (posters / movie.poster_filename).read_bytes() == b"poster content"


def test_upload_rejects_missing_filename(repo, monkeypatch, tmp_path):
    monkeypatch.setattr(movie_service, "MEDIA_DIR", tmp_path / "movies")
    with pytest.raises(HTTPException) as error:
        movie_service.MovieService(repo).upload_movie(
            "Film",
            "",
            None,
            None,
            None,
            None,
            None,
            None,
            UploadFile(filename=None, file=BytesIO(b"video")),
        )
    assert (error.value.status_code, error.value.detail) == (400, "No file provided")
    repo.create_movie.assert_not_called()
    assert not (tmp_path / "movies").exists()
