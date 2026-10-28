from datetime import UTC, datetime
from types import SimpleNamespace
from unittest.mock import Mock

import pytest

from app.constants.watch_event import WatchEventType
from app.repositories.watch_event_repository import WatchEventRepository
from app.repositories.watch_progress_repository import WatchProgressRepository
from app.schemas.watch_event import WatchEventCreate
from app.services.watch_event_service import WatchEventService
from app.services.watch_progress_service import WatchProgressService


@pytest.mark.parametrize("event_type", list(WatchEventType))
def test_create_event_updates_progress_according_to_event(event_type):
    repo = Mock(spec=WatchEventRepository)
    progress = Mock(spec=WatchProgressRepository)
    repo.create.return_value = SimpleNamespace(
        id=1, user_id=3, movie_id=7, event_type=event_type, position_seconds=42
    )
    result = WatchEventService(repo, progress).create_event(
        WatchEventCreate(movie_id=7, event_type=event_type, position_seconds=42),
        SimpleNamespace(id=3),
    )
    repo.create.assert_called_once_with(
        user_id=3, movie_id=7, event_type=event_type, position_seconds=42
    )
    assert result.model_dump() == repo.create.return_value.__dict__
    if event_type == WatchEventType.PLAY:
        progress.upsert.assert_not_called()
    else:
        expected = {
            "user_id": 3,
            "movie_id": 7,
            "position_seconds": 42,
            "is_completed": event_type == WatchEventType.COMPLETE,
        }
        if event_type == WatchEventType.HEARTBEAT:
            expected["update_watch_time"] = True
        progress.upsert.assert_called_once_with(**expected)


def test_progress_lookup_delegates_to_repository():
    repo = Mock(spec=WatchProgressRepository)
    service = WatchProgressService(repo)
    user = SimpleNamespace(id=3)
    assert service.get_user_progress(user) is repo.get_by_user.return_value
    assert service.get_movie_progress(3, 7) is repo.get_by_user_and_movie.return_value
    assert (
        WatchEventService(Mock(), repo).get_user_progress(user)
        is repo.get_by_user.return_value
    )
    assert repo.get_by_user.call_count == 2
    repo.get_by_user.assert_called_with(3)
    repo.get_by_user_and_movie.assert_called_once_with(3, 7)


def test_continue_watching_maps_repository_rows():
    repo = Mock(spec=WatchProgressRepository)
    repo.get_continue_watching.return_value = [
        (
            SimpleNamespace(last_position_seconds=42),
            SimpleNamespace(id=7, title="Film", poster_url=None),
        )
    ]
    result = WatchProgressService(repo).get_continue_watching(3)
    assert [row.model_dump() for row in result] == [
        {
            "movie_id": 7,
            "title": "Film",
            "poster_url": None,
            "last_position_seconds": 42,
        }
    ]
    repo.get_continue_watching.assert_called_once_with(3)


@pytest.mark.parametrize(
    "duration,expected", [(120, 35.0), (127, 33.07), (0, 0), (-1, 0)]
)
def test_watch_history_maps_metadata_and_percentage(duration, expected):
    repo = Mock(spec=WatchProgressRepository)
    timestamp = datetime(2024, 1, 2, tzinfo=UTC)
    repo.get_watch_history.return_value = [
        (
            SimpleNamespace(
                last_position_seconds=42,
                started_at=timestamp,
                last_watched_at=timestamp,
                is_completed=False,
            ),
            SimpleNamespace(
                id=7,
                title="Film",
                poster_url="/posters/7.jpg",
                duration_seconds=duration,
            ),
        )
    ]
    result = WatchProgressService(repo).get_watch_history(SimpleNamespace(id=3))
    assert result[0].model_dump() == {
        "movie": {
            "id": 7,
            "title": "Film",
            "poster_url": "/posters/7.jpg",
            "duration_seconds": duration,
        },
        "progress": {
            "last_position_seconds": 42,
            "watch_percentage": expected,
            "started_at": timestamp,
            "last_watched_at": timestamp,
            "is_completed": False,
        },
    }
    repo.get_watch_history.assert_called_once_with(3)


def test_empty_watch_lists():
    repo = Mock(spec=WatchProgressRepository)
    repo.get_continue_watching.return_value = []
    repo.get_watch_history.return_value = []
    service = WatchProgressService(repo)
    assert service.get_continue_watching(3) == []
    assert service.get_watch_history(SimpleNamespace(id=3)) == []
