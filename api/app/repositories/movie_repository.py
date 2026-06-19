from sqlalchemy.orm import Session

from app.models.movie import Movie


class MovieRepository:
    def __init__(self, db: Session):
        self.db = db

    def list_movies(self):
        return self.db.query(Movie).all()

    def get_movie(self, movie_id: int):
        return self.db.query(Movie).filter(Movie.id == movie_id).first()

    def create_movie(self, movie: Movie):
        self.db.add(movie)
        self.db.commit()
        self.db.refresh(movie)
        return movie
