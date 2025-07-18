from prisma import Prisma
from app.models.station import placeCreate

class UserRepository:
    @staticmethod
    async def create_place(db: Prisma, user_data: placeCreate):
        return await db.stations.create(
            data={
                "Name": user_data.Name,
                "Latitude": user_data.Latitude,
                "Longitude": user_data.Longitude,
                "type": "hybrid_bus"  
            }
        )

    # @staticmethod
    # async def get_user(db: Prisma, user_id: int):
    #     return await db.station.find_unique(where={"id": user_id})
