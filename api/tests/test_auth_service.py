from types import SimpleNamespace
from unittest.mock import Mock

import pytest
from fastapi import HTTPException

from app.repositories.user_repository import UserRepository
from app.schemas.user import UserCreate, UserLogin
from app.services import auth_service


@pytest.fixture
def repo():
    repository = Mock(spec=UserRepository)
    repository.get_by_username.return_value = None
    repository.get_by_email.return_value = None
    return repository


def test_register_hashes_password_and_returns_public_user(repo, monkeypatch):
    hasher = Mock()
    hasher.hash.return_value = "hashed"
    monkeypatch.setattr(auth_service, "pwd_context", hasher)
    repo.create_user.return_value = SimpleNamespace(
        id=7, username="alice", email="alice@example.com", password_hash="hashed"
    )
    result = auth_service.AuthService(repo).register(
        UserCreate(username="alice", email="alice@example.com", password="secret")
    )
    hasher.hash.assert_called_once_with("secret")
    repo.create_user.assert_called_once_with(
        username="alice", email="alice@example.com", password_hash="hashed"
    )
    assert result.model_dump() == {
        "id": 7,
        "username": "alice",
        "email": "alice@example.com",
    }


@pytest.mark.parametrize(
    "duplicate,detail",
    [
        ("get_by_username", "Username already exists"),
        ("get_by_email", "Email already exists"),
    ],
)
def test_register_rejects_duplicates(repo, duplicate, detail):
    getattr(repo, duplicate).return_value = object()
    with pytest.raises(HTTPException) as error:
        auth_service.AuthService(repo).register(
            UserCreate(username="alice", email="alice@example.com", password="secret")
        )
    assert (error.value.status_code, error.value.detail) == (400, detail)
    repo.create_user.assert_not_called()


@pytest.mark.parametrize("exists,valid", [(False, False), (True, False), (True, True)])
def test_login(repo, monkeypatch, exists, valid):
    repo.get_by_username.return_value = (
        SimpleNamespace(id=7, password_hash="hashed") if exists else None
    )
    hasher = Mock()
    hasher.verify.return_value = valid
    token = Mock(return_value="signed-token")
    monkeypatch.setattr(auth_service, "pwd_context", hasher)
    monkeypatch.setattr(auth_service, "create_access_token", token)
    service = auth_service.AuthService(repo)
    payload = UserLogin(username="alice", password="secret")
    if exists and valid:
        assert service.login(payload).model_dump() == {
            "access_token": "signed-token",
            "token_type": "bearer",
        }
        token.assert_called_once_with(7)
    else:
        with pytest.raises(HTTPException) as error:
            service.login(payload)
        assert (error.value.status_code, error.value.detail) == (
            401,
            "Invalid credentials",
        )
        token.assert_not_called()
    if exists:
        hasher.verify.assert_called_once_with("secret", "hashed")
    else:
        hasher.verify.assert_not_called()
