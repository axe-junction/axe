import aio_pika
from typing import List
from app.models.station import placeCreate
from app.repository.station import UserRepository
import json

class DefectAlgo:
    rabbitmq: aio_pika.Connection = None

    @staticmethod
    async def initialize() -> None:
        if DefectAlgo.rabbitmq is None or DefectAlgo.rabbitmq.is_closed:
            DefectAlgo.rabbitmq = await aio_pika.connect("amqp://guest:guest@localhost/")
        print("✅ RabbitMQ connection established")

    @staticmethod
    def is_conflicting(local: placeCreate, remote: placeCreate) -> bool:
        return local.id == remote.id and local != remote

    @staticmethod
    async def checkconflict(data: List[placeCreate]) -> bool:
        places: List[placeCreate] = await UserRepository.get_places()
        db_map = {place.id: place for place in places}

        conflict_found = False
        DefectAlgo._conflicts = []

        for incoming in data:
            if incoming.id in db_map:
                if DefectAlgo.is_conflicting(incoming, db_map[incoming.id]):
                    conflict_found = True
                    DefectAlgo._conflicts.append({
                        "id": incoming.id,
                        "local": incoming.dict(),
                        "remote": db_map[incoming.id].dict()
                    })

        return conflict_found

    @staticmethod
    async def getchange() -> None:
        if not hasattr(DefectAlgo, "_conflicts") or not DefectAlgo._conflicts:
            print("ℹ️ No conflicts to publish")
            return

        channel = await DefectAlgo.rabbitmq.channel()
        queue = await channel.declare_queue("placetopic", durable=True)

        message_body = json.dumps({
            "type": "conflict",
            "payload": DefectAlgo._conflicts
        })

        await channel.default_exchange.publish(
            aio_pika.Message(body=message_body.encode(), content_type="application/json"),
            routing_key=queue.name,
        )

