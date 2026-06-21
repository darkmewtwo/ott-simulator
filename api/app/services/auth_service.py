from passlib.context import CryptContext
from fastapi import HTTPException

from app.core.security import create_access_token
from app.repositories.user_repository import UserRepository
from app.schemas.user import UserCreate, UserResponse, TokenResponse, UserLogin


pwd_context = CryptContext(
    schemes=["bcrypt"],
    deprecated="auto",
)


class AuthService:
    def __init__(
        self,
        user_repo: UserRepository,
    ):
        self.user_repo = user_repo

    def register(
        self,
        payload: UserCreate,
    ) -> UserResponse:

        existing_user = self.user_repo.get_by_username(payload.username)

        if existing_user:
            raise HTTPException(
                status_code=400,
                detail="Username already exists",
            )

        existing_email = self.user_repo.get_by_email(payload.email)

        if existing_email:
            raise HTTPException(
                status_code=400,
                detail="Email already exists",
            )

        print(payload.password)
        print(type(payload.password))
        print(len(payload.password))
        password_hash = pwd_context.hash(payload.password)

        user = self.user_repo.create_user(
            username=payload.username,
            email=payload.email,
            password_hash=password_hash,
        )

        return UserResponse.model_validate(user)

    def verify_password(self, plain_password: str, password_hash: str):
        return pwd_context.verify(
            plain_password,
            password_hash,
        )

    def login(
        self,
        payload: UserLogin,
    ) -> TokenResponse:

        user = self.user_repo.get_by_username(payload.username)

        if not user:
            raise HTTPException(
                status_code=401,
                detail="Invalid credentials",
            )

        if not self.verify_password(
            payload.password,
            user.password_hash,
        ):
            raise HTTPException(
                status_code=401,
                detail="Invalid credentials",
            )

        token = create_access_token(user.id)

        return TokenResponse(
            access_token=token,
            token_type="bearer",
        )
