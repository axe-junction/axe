from fastapi import APIRouter,  Request
from app.models.station import placeCreate
from app.repository.station import UserRepository
router = APIRouter()


@router.post("/")
async def create_station(station:placeCreate , request: Request):
    db = request.app.state.db
    
    created = await UserRepository.create_place(db, station)
    return created

