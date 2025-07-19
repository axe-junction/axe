import aio_pika
import asyncio
from app.models.station import placeCreate
class DefectAlgo:
    rabbitmq:aio_pika.Connection
    @staticmethod
    async def initialize()-> None:
        DefectAlgo.rabbitmq = await aio_pika.connect("amqp://guest:guest@localhost/")
        return None
    
    async def checkconflict(data:placeCreate)->bool:
        
        
       pass
    async def  getchange():
       pass