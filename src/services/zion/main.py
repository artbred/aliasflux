import os
import asyncio
import nats
import json
import pickle

from dotenv import load_dotenv
from chat import send_chat_message
from prompts import business_name_chat_human_prompt_template
load_dotenv()


async def nats_listener():
    nc = await nats.connect(f"nats://{os.getenv('NATS_HOST')}:{os.getenv('NATS_PORT')}")

    async def message_handler(msg: nats.NATS.msg_class):
        try:
            json_mes = json.loads(msg.data.decode())
            res = send_chat_message(business_name_chat_human_prompt_template, json_mes["all_messages"], user_request=json_mes["message"])
            await nc.publish(msg.reply, json.dumps(json.loads(res)).encode("utf-8"))
        except Exception as e:
            print(e)

    channel = 'zion.create_business_names'
    await nc.subscribe(channel, cb=message_handler)

    while True:
        await asyncio.sleep(1)

if __name__ == '__main__':
    asyncio.run(nats_listener())
