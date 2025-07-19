import aio_pika
import asyncio
from app.models.station import placeCreate
from app.repository.station import UserRepository
class DefectAlgo:
    rabbitmq:aio_pika.Connection
    @staticmethod
    async def initialize()-> None:
        DefectAlgo.rabbitmq = await aio_pika.connect("amqp://guest:guest@localhost/")
        return None
    
    async def checkconflict(data:placeCreate)->bool:
        await UserRepository
        
        
       
    async def  getchange():
       pass